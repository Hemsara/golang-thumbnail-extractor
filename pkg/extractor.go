package extractor

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
)

func ExtractFirstPage(pdfPath string, ch chan<- string) error {

	doc, err := fitz.New(pdfPath)
	if err != nil {
		return fmt.Errorf("error opening PDF %s: %w", pdfPath, err)
	}
	defer doc.Close()

	outputDir := filepath.Join(".", "output_images")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %w", err)
		}
	}

	img, err := doc.Image(0)
	if err != nil {
		return fmt.Errorf("error extracting image from page of %s: %w", pdfPath, err)
	}

	imgImage, ok := img.(image.Image)
	if !ok {
		return fmt.Errorf("extracted image is not of type image.Image for %s", pdfPath)
	}

	baseName := filepath.Base(pdfPath)
	outputFileName := fmt.Sprintf("first_page_%s.jpg", baseName)
	outputFilePath := filepath.Join(outputDir, outputFileName)

	f, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("error creating image file %s: %w", outputFilePath, err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, imgImage, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
		return fmt.Errorf("error encoding image to JPEG: %w", err)
	}

	ch <- fmt.Sprintf("PDF first page converted to image successfully: %s", pdfPath)
	return nil
}
