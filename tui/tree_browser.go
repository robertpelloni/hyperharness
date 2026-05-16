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
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	Path       string              `json:"path"`
	IsDir      bool                `json:"isDir"`
	ChildCount int                 `json:"childCount"`
	Children   []TreeBrowserItem   `json:"children,omitempty"`
	Selected   bool                `json:"selected"`
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
		return fmt.Sprintf("[Tree] Entered directory: %s", item.Name), nil
	}
	return fmt.Sprintf("[Tree] Selected file: %s", item.Name), nil
}

// refreshPinnedFoundationTreeBrowser refreshes the pinned browser.
func refreshPinnedFoundationTreeBrowser(m *model) {
	if m == nil || m.foundationSessionID == "" {
		return
	}
	// Refresh browser items from working directory
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
		if strings.HasPrefix(entry.Name(), ".") && entry.Name() != ".git" {
			continue
		}
		info, _ := entry.Info()
		isDir := entry.IsDir()
		childCount := 0
		if isDir && info != nil {
			subEntries, err := os.ReadDir(filepath.Join(rootPath, entry.Name()))
			if err == nil {
				childCount = len(subEntries)
			}
		}
		items = append(items, TreeBrowserItem{
			ID:         filepath.Join(rootPath, entry.Name()),
			Name:       entry.Name(),
			Path:       filepath.Join(rootPath, entry.Name()),
			IsDir:      isDir,
			ChildCount: childCount,
		})
	}
	return items, nil
}