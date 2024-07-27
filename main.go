package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	extractor "thumb_pro/pkg"
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

	ch := make(chan string)
	var wg sync.WaitGroup

	for _, pdfPath := range pdfPaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			err := extractor.ExtractFirstPage(path, ch)
			if err != nil {
				fmt.Printf("Error extracting first page from %s: %v\n", path, err)
			}
		}(pdfPath)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		fmt.Println(result)
	}
}
