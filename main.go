package main

import (
	"io"
	"log"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	// Create Fiber app
	app := fiber.New()
	app.Static("/", "./docs/")
	// Define file upload endpoint
	app.Post("/compare-documents", handleUpload)

	// Start server
	log.Fatal(app.Listen(":8080"))
}

func handleUpload(c *fiber.Ctx) error {
	// Get uploaded files
	file1, err := c.FormFile("file1")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing file1"})
	}

	file2, err := c.FormFile("file2")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing file2"})
	}

	// Read file contents into memory
	file1Content, err := readFormFile(file1)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read content of file1"})
	}

	file2Content, err := readFormFile(file2)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read content of file2"})
	}

	// Compare document contents
	diffs := compareDocuments(file1Content, file2Content)

	// Return comparison results
	return c.JSON(fiber.Map{"diffs": diffs})
}

func readFormFile(file *multipart.FileHeader) ([]byte, error) {
	// Open the uploaded file
	fileReader, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	// Read the file content into a byte slice
	fileContent, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}

func compareDocuments(file1Data []byte, file2Data []byte) []diffmatchpatch.Diff {
	// Convert file contents to strings
	text1 := string(file1Data)
	text2 := string(file2Data)

	// Perform diffing
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)

	return diffs
}
