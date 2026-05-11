package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// SummarizeGoFile reads a Go source file and returns a summary containing only
// its package declaration, imports, types, interfaces, constants, and function
// signatures (both standalone and methods).
// This is highly useful for providing context to LLMs without consuming too many tokens.
func SummarizeGoFile(filepath string, src []byte) (string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filepath, src, 0)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("package %s\n\n", node.Name.Name))

	var imports []string
	var consts []string
	var vars []string
	var types []string
	var funcs []string

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			switch x.Tok {
			case token.IMPORT:
				for _, spec := range x.Specs {
					if is, ok := spec.(*ast.ImportSpec); ok {
						imports = append(imports, is.Path.Value)
					}
				}
			case token.CONST:
				for _, spec := range x.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						names := []string{}
						for _, name := range vs.Names {
							names = append(names, name.Name)
						}
						typeStr := ""
						if vs.Type != nil {
							typeStr = " " + formatNode(fset, vs.Type)
						}
						consts = append(consts, fmt.Sprintf("const %s%s", strings.Join(names, ", "), typeStr))
					}
				}
			case token.VAR:
				for _, spec := range x.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						names := []string{}
						for _, name := range vs.Names {
							names = append(names, name.Name)
						}
						typeStr := ""
						if vs.Type != nil {
							typeStr = " " + formatNode(fset, vs.Type)
						}
						vars = append(vars, fmt.Sprintf("var %s%s", strings.Join(names, ", "), typeStr))
					}
				}
			case token.TYPE:
				for _, spec := range x.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						typeStr := formatNode(fset, ts.Type)
						types = append(types, fmt.Sprintf("type %s %s", ts.Name.Name, typeStr))
					}
				}
			}
		case *ast.FuncDecl:
			recv := ""
			if x.Recv != nil && len(x.Recv.List) > 0 {
				recvStr := formatNode(fset, x.Recv.List[0].Type)
				// Clean up pointers if present simply
				recvStr = strings.ReplaceAll(recvStr, "*", "*") 
				recv = fmt.Sprintf("(%s) ", recvStr)
			}
			
			params := formatFieldList(fset, x.Type.Params)
			results := formatFieldList(fset, x.Type.Results)
			
			retStr := ""
			if results != "" {
				if strings.Contains(results, ",") {
					retStr = fmt.Sprintf(" (%s)", results)
				} else {
					retStr = fmt.Sprintf(" %s", results)
				}
			}

			funcs = append(funcs, fmt.Sprintf("func %s%s(%s)%s { ... }", recv, x.Name.Name, params, retStr))
		}
		return true
	})

	if len(imports) > 0 {
		sb.WriteString(fmt.Sprintf("import (\n\t%s\n)\n\n", strings.Join(imports, "\n\t")))
	}

	if len(consts) > 0 {
		for _, c := range consts {
			sb.WriteString(c + "\n")
		}
		sb.WriteString("\n")
	}

	if len(vars) > 0 {
		for _, v := range vars {
			sb.WriteString(v + "\n")
		}
		sb.WriteString("\n")
	}

	if len(types) > 0 {
		for _, t := range types {
			sb.WriteString(t + "\n\n")
		}
	}

	if len(funcs) > 0 {
		for _, f := range funcs {
			sb.WriteString(f + "\n")
		}
	}

	return sb.String(), nil
}

func formatNode(fset *token.FileSet, n ast.Node) string {
	if n == nil {
		return ""
	}
	var buf bytes.Buffer
	// Basic rudimentary printing
	ast.Fprint(&buf, fset, n, nil)
	// We'd ideally use format.Node here but to keep things simple for the LLM let's do simple type extraction
	switch x := n.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.StarExpr:
		return "*" + formatNode(fset, x.X)
	case *ast.ArrayType:
		return "[]" + formatNode(fset, x.Elt)
	case *ast.SelectorExpr:
		return formatNode(fset, x.X) + "." + x.Sel.Name
	case *ast.InterfaceType:
		return "interface{ ... }"
	case *ast.StructType:
		return "struct{ ... }"
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", formatNode(fset, x.Key), formatNode(fset, x.Value))
	case *ast.FuncType:
		return "func(...)"
	case *ast.ChanType:
		return "chan " + formatNode(fset, x.Value)
	}
	return fmt.Sprintf("%T", n)
}

func formatFieldList(fset *token.FileSet, fl *ast.FieldList) string {
	if fl == nil {
		return ""
	}
	var parts []string
	for _, field := range fl.List {
		typeStr := formatNode(fset, field.Type)
		if len(field.Names) == 0 {
			parts = append(parts, typeStr)
		} else {
			for _, name := range field.Names {
				parts = append(parts, fmt.Sprintf("%s %s", name.Name, typeStr))
			}
		}
	}
	return strings.Join(parts, ", ")
}
