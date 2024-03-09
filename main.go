package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"

	"strings"

	"github.com/gen2brain/go-fitz"
)

func main() {
	var pdfPaths []string

	err := filepath.Walk("./files", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if strings.HasSuffix(strings.ToLower(path), ".pdf") {
			fmt.Printf("Processing: %s\n", path)
			pdfPaths = append(pdfPaths, path)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, pdfPath := range pdfPaths {
		err := extractFirstPage(pdfPath)
		if err != nil {
			fmt.Printf("Error extracting first page from %s: %v\n", pdfPath, err)
		}
	}
}
func extractFirstPage(pdfPath string) error {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		fmt.Println("Error opening PDF:", err)
		return err
	}
	defer doc.Close()

	outputDir := filepath.Join(".", "output_images")

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err := os.Mkdir(outputDir, 0755)
		if err != nil {
			fmt.Println("Error creating output directory:", err)
			return err
		}
	}

	img, err := doc.Image(0)
	if err != nil {
		fmt.Println("Error extracting image from page:", err)
		return err
	}

	baseName := filepath.Base(pdfPath)

	outputFileName := fmt.Sprintf("first_page_%s.jpg", baseName)
	outputFilePath := filepath.Join(outputDir, outputFileName)

	f, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println("Error creating image file:", err)
		return err
	}

	err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		fmt.Println("Error encoding image to JPEG:", err)
		return err
	}

	f.Close()

	fmt.Printf("PDF first page converted to image successfully: %s\n", pdfPath)
	return nil
}
