// tree_browser.go provides tree browser support for the TUI.
package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TreeBrowserItem represents a file or directory in the tree browser.
type TreeBrowserItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	IsDir      bool   `json:"isDir"`
	ChildCount int    `json:"childCount"`
	Children   []TreeBrowserItem `json:"children,omitempty"`
	Selected   bool   `json:"selected"`
}

// visibleTreeBrowserItems returns items matching the filter.
func visibleTreeBrowserItems(items []TreeBrowserItem, filter string, collapsed map[string]bool) []TreeBrowserItem {
	if filter == "" {
		return items
	}
	filter = strings.ToLower(filter)
	var result []TreeBrowserItem
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), filter) {
			result = append(result, item)
		}
	}
	return result
}

// openSelectedTreeBrowser opens the selected tree browser item.
func openSelectedTreeBrowser(workingDir, sessionID string, items []TreeBrowserItem, index, limit int) (string, error) {
	if index < 0 || index >= len(items) {
		return "", fmt.Errorf("invalid index")
	}
	item := items[index]
	if item.IsDir {
		// List directory contents
		entries, err := os.ReadDir(item.Path)
		if err != nil {
			return fmt.Sprintf("[Tree] Entered directory: %s (error: %v)", item.Name, err), nil
		}
		var lines []string
		lines = append(lines, fmt.Sprintf("[Tree] Directory: %s (%d entries)", item.Name, len(entries)))
		for i, entry := range entries {
			if i >= limit { break }
			icon := "📄"
			if entry.IsDir() { icon = "📁" }
			lines = append(lines, fmt.Sprintf("  %s %s", icon, entry.Name()))
		}
		return strings.Join(lines, "\n"), nil
	}
	// Read file preview
	content, err := os.ReadFile(item.Path)
	if err != nil {
		return fmt.Sprintf("[Tree] Selected file: %s (error: %v)", item.Name, err), nil
	}
	preview := string(content)
	if len(preview) > limit*200 {
		preview = preview[:limit*200] + "\n  ... (truncated)"
	}
	return fmt.Sprintf("[Tree] File: %s\n%s", item.Name, preview), nil
}

// refreshPinnedFoundationTreeBrowser refreshes the pinned browser.
func refreshPinnedFoundationTreeBrowser(m *model) {
	if m == nil || m.foundationSessionID == "" {
		return
	}
	items, _ := buildTreeItems(m.director.WorkingDir)
	m.browserItems = items
	m.browserIndex = 0
}

// buildTreeItems builds tree browser items from a directory.
func buildTreeItems(rootPath string) ([]TreeBrowserItem, error) {
	if rootPath == "" {
		rootPath = "."
	}
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}
	var items []TreeBrowserItem
	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden files except .git
		if strings.HasPrefix(name, ".") && name != ".git" && name != ".github" {
			continue
		}
		// Skip node_modules, vendor, __pycache__
		if name == "node_modules" || name == "vendor" || name == "__pycache__" || name == "dist" || name == "bin" {
			continue
		}
		info, _ := entry.Info()
		isDir := entry.IsDir()
		childCount := 0
		if isDir && info != nil {
			subEntries, err := os.ReadDir(filepath.Join(rootPath, name))
			if err == nil {
				childCount = len(subEntries)
			}
		}
		items = append(items, TreeBrowserItem{
			ID:         filepath.Join(rootPath, name),
			Name:       name,
			Path:       filepath.Join(rootPath, name),
			IsDir:      isDir,
			ChildCount: childCount,
		})
	}
	return items, nil
}

// renderTreeBrowser renders the tree browser with the theme system.
func renderTreeBrowser(items []TreeBrowserItem, index int, filter string, confirmPending bool, collapsed map[string]bool, grouped bool, paneHeight int, title string, preview bool) string {
	t := DefaultTheme
	visible := visibleTreeBrowserItems(items, filter, collapsed)
	var lines []string

	if title != "" {
		lines = append(lines, t.Bold(t.AccentText(title)))
	}

	// Filter bar
	if filter != "" {
		lines = append(lines, t.Dim("  filter: ")+t.Fg(t.TextColor, filter))
	}

	// Items
	startIdx := 0
	if paneHeight > 0 && len(visible) > paneHeight {
		// Scroll to keep index visible
		if index > paneHeight-2 {
			startIdx = index - paneHeight + 3
		}
	}
	endIdx := len(visible)
	if paneHeight > 0 && startIdx+paneHeight < len(visible) {
		endIdx = startIdx + paneHeight
	}

	for i := startIdx; i < endIdx; i++ {
		item := visible[i]
		icon := "📄"
		if item.IsDir { icon = "📁" }

		prefix := "  "
		sel := ""
		if i == index {
			prefix = " ▶"
			sel = " " + t.Fg(t.TextColor, item.Name)
		}

		childInfo := ""
		if item.IsDir && item.ChildCount > 0 {
			childInfo = t.Dim(fmt.Sprintf(" (%d)", item.ChildCount))
			if collapsed != nil && collapsed[item.ID] {
				childInfo = t.Dim(" ▸") + childInfo
			} else {
				childInfo = t.Dim(" ▾") + childInfo
			}
		}

		if i == index {
			lines = append(lines, t.AccentText(prefix)+icon+" "+t.Bold(t.Fg(t.TextColor, item.Name))+childInfo)
		} else {
			lines = append(lines, t.Dim(prefix)+t.Dim(icon)+" "+t.Fg(t.Muted, item.Name)+childInfo)
		}
		_ = sel // selection highlight applied via prefix
	}

	// Preview
	if preview && index >= 0 && index < len(visible) {
		item := visible[index]
		if !item.IsDir {
			content, err := os.ReadFile(item.Path)
			if err == nil {
				preview := string(content)
				if len(preview) > 300 { preview = preview[:300] + "…" }
				lines = append(lines, t.Dim("  ┌─ preview ──────────────────"))
				for _, line := range strings.Split(preview, "\n") {
					if len(lines) > startIdx+paneHeight+5 { break }
					lines = append(lines, t.Dim("  │ ")+t.Fg(t.DimColor, line))
				}
				lines = append(lines, t.Dim("  └────────────────────────────"))
			}
		}
	}

	// Confirm prompt
	if confirmPending {
		lines = append(lines, t.WarningText("  Confirm? [Enter=Y, Esc=N]"))
	}

	// Scroll indicator
	if paneHeight > 0 && len(visible) > paneHeight {
		lines = append(lines, t.Dim(fmt.Sprintf("  ↑↓ scroll  (%d/%d)", index+1, len(visible))))
	}

	return strings.Join(lines, "\n")
}
