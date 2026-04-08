package fs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// FileInfo provides detailed file information beyond os.FileInfo.
type FileInfo struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"is_dir"`
	Size     int64       `json:"size"`
	Mode     string      `json:"mode"`
	Modified time.Time   `json:"modified"`
	LineCount int        `json:"line_count,omitempty"`
	Language  string     `json:"language,omitempty"`
}

// FileTree represents a directory tree with metadata.
type FileTree struct {
	Root     string     `json:"root"`
	Files    []FileInfo `json:"files"`
	Dirs     []FileInfo `json:"dirs"`
	TotalSize int64     `json:"total_size"`
	FileCount int       `json:"file_count"`
	DirCount  int       `json:"dir_count"`
}

// EnsureDir creates a directory and all parent directories if they don't exist.
func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0o755)
	}
	return nil
}

// EnsureParentDir creates parent directories for a file path.
func EnsureParentDir(filePath string) error {
	return EnsureDir(filepath.Dir(filePath))
}

// FileExists checks if a file exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsGitIgnored checks if a path matches common gitignore patterns.
// This is a simplified version that checks common patterns.
func IsGitIgnored(path string) bool {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for _, part := range parts {
		switch part {
		case ".git", "node_modules", "__pycache__", ".DS_Store",
			".vscode", ".idea", "dist", "build", "target",
			".next", ".nuxt", "coverage", ".cache",
			"vendor", ".terraform", ".tox", "venv", ".venv",
			"bin", "obj", "Debug", "Release":
			return true
		}
		if strings.HasPrefix(part, ".") && len(part) > 1 {
			// Hidden directories except .config, .agents, .skills
			switch part {
			case ".config", ".agents", ".skills", ".github":
				// not ignored
			default:
				// other hidden dirs are typically ignored
			}
		}
	}
	return false
}

// DetectLanguage detects the programming language from a file extension.
func DetectLanguage(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".js", ".mjs", ".cjs":
		return "javascript"
	case ".ts", ".tsx", ".mts", ".cts":
		return "typescript"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".kt", ".kts":
		return "kotlin"
	case ".c", ".h":
		return "c"
	case ".cpp", ".cc", ".cxx", ".hpp", ".hxx":
		return "cpp"
	case ".cs":
		return "csharp"
	case ".rb":
		return "ruby"
	case ".php":
		return "php"
	case ".swift":
		return "swift"
	case ".zig":
		return "zig"
	case ".lua":
		return "lua"
	case ".sh", ".bash", ".zsh":
		return "shell"
	case ".ps1":
		return "powershell"
	case ".sql":
		return "sql"
	case ".html", ".htm":
		return "html"
	case ".css", ".scss", ".sass", ".less":
		return "css"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".toml":
		return "toml"
	case ".xml":
		return "xml"
	case ".md", ".mdx":
		return "markdown"
	case ".dockerfile":
		return "dockerfile"
	case ".proto":
		return "protobuf"
	case ".graphql", ".gql":
		return "graphql"
	default:
		// Check filename patterns
		base := strings.ToLower(filepath.Base(filename))
		switch base {
		case "dockerfile", "dockerfile.dev", "dockerfile.prod":
			return "dockerfile"
		case "makefile", "gnumakefile":
			return "makefile"
		case "cmakelists.txt":
			return "cmake"
		case "go.mod", "go.sum":
			return "go"
		case "cargo.toml", "cargo.lock":
			return "rust"
		case "package.json", "tsconfig.json":
			return "json"
		default:
			return ""
		}
	}
}

// GetFileInfo returns detailed information about a file.
func GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	fi := &FileInfo{
		Name:     info.Name(),
		Path:     path,
		IsDir:    info.IsDir(),
		Size:     info.Size(),
		Mode:     info.Mode().String(),
		Modified: info.ModTime(),
		Language: DetectLanguage(path),
	}

	// Count lines for text files
	if !fi.IsDir && fi.Language != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			fi.LineCount = strings.Count(string(data), "\n") + 1
		}
	}

	return fi, nil
}

// WalkDirectory walks a directory tree and returns a FileTree.
// Respects common ignore patterns.
func WalkDirectory(root string, maxDepth int) (*FileTree, error) {
	tree := &FileTree{
		Root: root,
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}

		// Skip ignored paths
		if IsGitIgnored(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Depth check
		if maxDepth > 0 {
			depth := len(strings.Split(rel, string(filepath.Separator)))
			if depth > maxDepth {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		fi := FileInfo{
			Name:     d.Name(),
			Path:     rel,
			IsDir:    d.IsDir(),
			Size:     info.Size(),
			Modified: info.ModTime(),
			Language: DetectLanguage(path),
		}

		if d.IsDir() {
			// Skip the root directory itself
			if rel == "." {
				return nil
			}
			tree.Dirs = append(tree.Dirs, fi)
			tree.DirCount++
			tree.TotalSize += info.Size()
		} else {
			tree.Files = append(tree.Files, fi)
			tree.FileCount++
			tree.TotalSize += info.Size()
		}

		return nil
	})

	return tree, err
}

// FindProjectRoot searches upward for a project root marker file.
func FindProjectRoot(startDir string) string {
	markers := []string{
		".git",
		"go.mod",
		"package.json",
		"Cargo.toml",
		"pyproject.toml",
		"pom.xml",
		"build.gradle",
		".hypercode",
		".agents",
	}

	dir := startDir
	for {
		for _, marker := range markers {
			if FileExists(filepath.Join(dir, marker)) {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return startDir
}

// HomeDir returns the user's home directory.
func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/"
	}
	return home
}

// ConfigDir returns the hyperharness config directory.
func ConfigDir() string {
	return filepath.Join(HomeDir(), ".hyperharness")
}

// TempDir returns a temporary directory for hyperharness.
func TempDir() string {
	return filepath.Join(os.TempDir(), "hyperharness")
}

// IsTextFile determines if a file is likely a text file based on extension.
func IsTextFile(path string) bool {
	lang := DetectLanguage(path)
	return lang != ""
}

// FormatSize formats bytes into human-readable size.
func FormatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// Platform returns the current platform string.
func Platform() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}
