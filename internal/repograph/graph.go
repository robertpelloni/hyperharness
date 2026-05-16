package repograph

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type NodeType string

const (
	NodeFile      NodeType = "file"
	NodeFunction  NodeType = "function"
	NodeTypeDef   NodeType = "type"
	NodeInterface NodeType = "interface"
	NodeImport    NodeType = "import"
)

type Node struct {
	ID         string   `json:"id"`
	Type       NodeType `json:"type"`
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	Package    string   `json:"package,omitempty"`
	Language   string   `json:"language,omitempty"`
	LineStart  int      `json:"lineStart,omitempty"`
	IsExported bool     `json:"isExported"`
}

type Edge struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Type   string `json:"type"`
	Weight int    `json:"weight"`
}

type Graph struct {
	Nodes    map[string]*Node `json:"nodes"`
	Edges    []Edge           `json:"edges"`
	RootPath string           `json:"rootPath"`
	BuiltAt  time.Time        `json:"builtAt"`
	Stats    GraphStats       `json:"stats"`
}

type GraphStats struct {
	TotalFiles     int `json:"totalFiles"`
	TotalFunctions int `json:"totalFunctions"`
	TotalTypes     int `json:"totalTypes"`
	TotalImports   int `json:"totalImports"`
	TotalEdges     int `json:"totalEdges"`
}

type RepoGraphService struct {
	root  string
	mu    sync.RWMutex
	graph *Graph
}

func NewRepoGraphService(root string) *RepoGraphService {
	return &RepoGraphService{root: root}
}

func (rgs *RepoGraphService) Build(ctx context.Context) (*Graph, error) {
	graph := &Graph{Nodes: make(map[string]*Node), RootPath: rgs.root, BuiltAt: time.Now().UTC()}
	err := filepath.WalkDir(rgs.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := strings.ToLower(d.Name())
			switch name {
			case "node_modules", ".git", "dist", "build", "coverage", "vendor", "__pycache__", ".next", ".cache", "target":
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".go", ".ts", ".tsx", ".js", ".jsx", ".py", ".rs":
			rgs.indexFile(graph, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	rgs.mu.Lock()
	rgs.graph = graph
	rgs.mu.Unlock()
	graph.Stats = GraphStats{
		TotalFiles:     rgs.countByType(graph, NodeFile),
		TotalFunctions: rgs.countByType(graph, NodeFunction),
		TotalTypes:     rgs.countByType(graph, NodeTypeDef),
		TotalImports:   rgs.countByType(graph, NodeImport),
		TotalEdges:     len(graph.Edges),
	}
	return graph, nil
}

func (rgs *RepoGraphService) GetGraph() *Graph {
	rgs.mu.RLock()
	defer rgs.mu.RUnlock()
	return rgs.graph
}

func (rgs *RepoGraphService) SearchSymbols(query string, limit int) []*Node {
	rgs.mu.RLock()
	defer rgs.mu.RUnlock()
	if rgs.graph == nil {
		return nil
	}
	q := strings.ToLower(query)
	var results []*Node
	for _, node := range rgs.graph.Nodes {
		if node.Type == NodeFile || node.Type == NodeImport {
			continue
		}
		if strings.Contains(strings.ToLower(node.Name), q) {
			results = append(results, node)
		}
		if len(results) >= limit {
			break
		}
	}
	return results
}

func (rgs *RepoGraphService) indexFile(graph *Graph, filePath string) {
	relPath, _ := filepath.Rel(rgs.root, filePath)
	fileID := "file:" + relPath
	graph.Nodes[fileID] = &Node{ID: fileID, Type: NodeFile, Name: filepath.Base(filePath), Path: relPath, Language: languageFromExt(filepath.Ext(filePath))}
	data, err := os.ReadFile(filePath)
	if err != nil || len(data) > 500*1024 {
		return
	}
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		rgs.indexGoFile(graph, relPath, string(data))
	case ".ts", ".tsx", ".js", ".jsx":
		rgs.indexTSFile(graph, relPath, string(data))
	case ".py":
		rgs.indexPythonFile(graph, relPath, string(data))
	}
}

var (
	goFuncRe  = regexp.MustCompile(`^func\s+(?:\([^)]+\)\s+)?(\w+)`)
	goTypeRe  = regexp.MustCompile(`^type\s+(\w+)\s+(struct|interface)`)
	goPkgRe   = regexp.MustCompile(`^package\s+(\w+)`)
	goImpRe   = regexp.MustCompile(`import\s+"([^"]+)"`)
	pyFuncRe  = regexp.MustCompile(`^def\s+(\w+)`)
	pyClassRe = regexp.MustCompile(`^class\s+(\w+)`)
)

func (rgs *RepoGraphService) indexGoFile(graph *Graph, relPath, content string) {
	fileID := "file:" + relPath
	pkgName := ""
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if matches := goPkgRe.FindStringSubmatch(line); len(matches) > 1 {
			pkgName = matches[1]
			if f := graph.Nodes[fileID]; f != nil {
				f.Package = pkgName
			}
		}
		if matches := goImpRe.FindStringSubmatch(line); len(matches) > 1 {
			graph.Edges = append(graph.Edges, Edge{From: fileID, To: "import:" + matches[1], Type: "imports"})
		}
		if matches := goFuncRe.FindStringSubmatch(line); len(matches) > 1 {
			name := matches[1]
			graph.Nodes[relPath+"#"+name] = &Node{ID: relPath + "#" + name, Type: NodeFunction, Name: name, Path: relPath, Package: pkgName, LineStart: lineNum, IsExported: name[0] >= 'A' && name[0] <= 'Z', Language: "go"}
		}
		if matches := goTypeRe.FindStringSubmatch(line); len(matches) > 1 {
			name := matches[1]
			nt := NodeTypeDef
			if matches[2] == "interface" {
				nt = NodeInterface
			}
			graph.Nodes[relPath+"#"+name] = &Node{ID: relPath + "#" + name, Type: nt, Name: name, Path: relPath, Package: pkgName, LineStart: lineNum, IsExported: name[0] >= 'A' && name[0] <= 'Z', Language: "go"}
		}
	}
}

func (rgs *RepoGraphService) indexTSFile(graph *Graph, relPath, content string) {
	_ = "file:" + relPath // reserved for import edges
	tsFuncRe := regexp.MustCompile(`(?:export\s+)?(?:async\s+)?function\s+(\w+)`)
	tsClassRe := regexp.MustCompile(`(?:export\s+)?(?:abstract\s+)?class\s+(\w+)`)
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if matches := tsFuncRe.FindStringSubmatch(line); len(matches) > 1 {
			graph.Nodes[relPath+"#"+matches[1]] = &Node{ID: relPath + "#" + matches[1], Type: NodeFunction, Name: matches[1], Path: relPath, LineStart: lineNum, IsExported: strings.Contains(line, "export"), Language: "typescript"}
		}
		if matches := tsClassRe.FindStringSubmatch(line); len(matches) > 1 {
			graph.Nodes[relPath+"#"+matches[1]] = &Node{ID: relPath + "#" + matches[1], Type: NodeTypeDef, Name: matches[1], Path: relPath, LineStart: lineNum, IsExported: strings.Contains(line, "export"), Language: "typescript"}
		}
	}
}

func (rgs *RepoGraphService) indexPythonFile(graph *Graph, relPath, content string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if matches := pyFuncRe.FindStringSubmatch(line); len(matches) > 1 {
			name := matches[1]
			graph.Nodes[relPath+"#"+name] = &Node{ID: relPath + "#" + name, Type: NodeFunction, Name: name, Path: relPath, LineStart: lineNum, IsExported: !strings.HasPrefix(name, "_"), Language: "python"}
		}
		if matches := pyClassRe.FindStringSubmatch(line); len(matches) > 1 {
			graph.Nodes[relPath+"#"+matches[1]] = &Node{ID: relPath + "#" + matches[1], Type: NodeTypeDef, Name: matches[1], Path: relPath, LineStart: lineNum, IsExported: !strings.HasPrefix(matches[1], "_"), Language: "python"}
		}
	}
}

func languageFromExt(ext string) string {
	switch ext {
	case ".go": return "go"
	case ".ts", ".tsx": return "typescript"
	case ".js", ".jsx": return "javascript"
	case ".py": return "python"
	case ".rs": return "rust"
	default: return "unknown"
	}
}

func (rgs *RepoGraphService) countByType(graph *Graph, nodeType NodeType) int {
	count := 0
	for _, node := range graph.Nodes {
		if node.Type == nodeType {
			count++
		}
	}
	return count
}
