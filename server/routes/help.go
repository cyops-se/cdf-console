package routes

import (
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"path"
	"server/help"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type DocsStructureResponse struct {
	Success bool          `json:"success"`
	Errors  []string      `json:"errors"`
	Root    *help.DocNode `json:"root,omitempty"`
}

type DocContentResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
	Content string   `json:"content,omitempty"`
	Path    string   `json:"path,omitempty"`
}

func registerHelpRoutes(api fiber.Router) {
	api.Get("/help/structure", GetDocsStructure)
	api.Get("/help/content/*", GetDocContent)
	api.Get("/help/images/*", GetDocImage)
}

// GetDocsStructure returns the folder structure of all help documentation
func GetDocsStructure(c *fiber.Ctx) error {
	response := DocsStructureResponse{Success: false, Errors: []string{}}

	help.ListDocsFileSystem()
	root, err := help.BuildDocTree("docs", "docs")
	if err != nil {
		response.Errors = append(response.Errors, err.Error())
		return c.JSON(response)
	}

	response.Success = true
	response.Root = root
	return c.Status(http.StatusOK).JSON(response)
}

// GetDocContent retrieves the content of a specific markdown file
func GetDocContent(c *fiber.Ctx) error {
	response := DocContentResponse{Success: false, Errors: []string{}}

	// Get the path from the URL (everything after /help/content/)
	docPath := c.Params("*")
	if docPath == "" {
		response.Errors = append(response.Errors, "Document path is required")
		return c.JSON(response)
	}

	// Decode URL encoding (handles spaces and UTF-8 characters)
	decodedPath, err := url.PathUnescape(docPath)
	if err != nil {
		response.Errors = append(response.Errors, "Invalid path encoding: "+err.Error())
		return c.JSON(response)
	}

	// Read the file content using the help package
	content, err := help.GetDocContent(decodedPath)
	if err != nil {
		response.Errors = append(response.Errors, "Document not found: "+err.Error())
		return c.JSON(response)
	}

	response.Success = true
	response.Content = content
	response.Path = decodedPath
	return c.Status(http.StatusOK).JSON(response)
}

// GetDocImage serves image files from the embedded docs filesystem
func GetDocImage(c *fiber.Ctx) error {
	// Get the image path from the URL (everything after /help/images/)
	imagePath := c.Params("*")
	if imagePath == "" {
		return c.Status(http.StatusBadRequest).SendString("Image path is required")
	}

	// Decode URL encoding (handles spaces and UTF-8 characters)
	decodedPath, err := url.PathUnescape(imagePath)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid path encoding: " + err.Error())
	}

	// Clean the path to prevent directory traversal
	decodedPath = path.Clean(decodedPath)
	if strings.Contains(decodedPath, "../") {
		// return c.Status(http.StatusBadRequest).SendString("Invalid image path")
		decodedPath = strings.ReplaceAll(decodedPath, "../", "")
	}

	// Construct the full path within the docs folder
	fullPath := path.Join("docs", "images", decodedPath)

	// Read the image file from embedded filesystem
	imageData, err := fs.ReadFile(help.DocsFS, fullPath)
	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Image not found: " + err.Error())
	}

	// Determine content type based on file extension
	ext := path.Ext(decodedPath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Set content type and return the image
	c.Set("Content-Type", contentType)
	return c.Status(http.StatusOK).Send(imageData)
}
