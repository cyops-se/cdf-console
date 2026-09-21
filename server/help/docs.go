package help

import (
	"embed"
	"io/fs"
	"path"
	"regexp"
	"server/logger"
	"strings"
)

//go:embed all:docs/*
var DocsFS embed.FS

type DocNode struct {
	Name     string     `json:"name"`
	Title    string     `json:"title"`
	Path     string     `json:"path"`
	IsDir    bool       `json:"isDir"`
	Children []*DocNode `json:"children,omitempty"`
}

// extractFirstHeading extracts the first # heading from markdown content
func extractFirstHeading(content string) string {
	// Match first # heading (not ##, ###, etc.)
	// This regex matches a line that starts with exactly one # followed by space
	re := regexp.MustCompile(`(?m)^#\s+(.+?)(?:\r?\n|$)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		// Remove any trailing special characters
		title = strings.TrimRight(title, " \t\r\n")
		return title
	}
	return ""
}

// BuildDocTree recursively builds a tree structure of the docs folder
func BuildDocTree(name, docpath string) (*DocNode, error) {
	entries, err := DocsFS.ReadDir(docpath)
	if err != nil {
		return nil, err
	}

	node := &DocNode{
		Name:     name,
		Title:    name,
		Path:     strings.TrimPrefix(docpath, "docs/"),
		IsDir:    true,
		Children: []*DocNode{},
	}

	for _, entry := range entries {
		childPath := path.Join(docpath, entry.Name())

		if entry.IsDir() {
			childNode, err := BuildDocTree(entry.Name(), childPath)
			if err != nil {
				continue
			}
			node.Children = append(node.Children, childNode)
		} else if strings.HasSuffix(entry.Name(), ".md") {
			// Read file content to extract title
			content, err := fs.ReadFile(DocsFS, childPath)
			title := entry.Name()
			if err == nil {
				extractedTitle := extractFirstHeading(string(content))
				if extractedTitle != "" {
					title = extractedTitle
					logger.Trace("BuildDocTree", "Extracted title '%s' from file '%s'", title, entry.Name())
				} else {
					logger.Trace("BuildDocTree", "No title found in file '%s', using filename", entry.Name())
				}
			} else {
				logger.Trace("BuildDocTree", "Error reading file '%s': %v", entry.Name(), err)
			}

			node.Children = append(node.Children, &DocNode{
				Name:  entry.Name(),
				Title: title,
				Path:  strings.TrimPrefix(childPath, "docs/"),
				IsDir: false,
			})
		}
	}

	return node, nil
}

// GetDocContent reads the content of a markdown file from the embedded FS
func GetDocContent(docPath string) (string, error) {
	// Clean the path to prevent directory traversal
	docPath = path.Clean(docPath)
	if strings.Contains(docPath, "..") {
		return "", fs.ErrInvalid
	}

	// Construct the full path
	fullPath := path.Join("docs", docPath)

	// Read the file content
	logger.Trace("GetDocContent", "Reading file: %s", fullPath)
	content, err := fs.ReadFile(DocsFS, fullPath)
	if err != nil {
		logger.Error("GetDocContent", "Cannot read file: %s, error: %s", fullPath, err.Error())
		return "", err
	}

	return string(content), nil
}

func ListFolder(folder string) {
	if entries, err := DocsFS.ReadDir(folder); err == nil {
		for _, e := range entries {
			logger.Trace("GetDocContent", "Entry: %s", path.Join(folder, e.Name()))
			if e.IsDir() {
				subfolder := path.Join(folder, e.Name())
				ListFolder(subfolder)
			}
		}

	} else {
		logger.Trace("Error listing DocsFS:", err.Error())
	}
}

func ListDocsFileSystem() {
	logger.Trace("GetDocContent", "--== LISTING DOCSFS ==--")
	ListFolder("docs")
	logger.Trace("GetDocContent", "--== DONE ==--")
}
