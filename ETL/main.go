package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/mholt/archiver/v3"
)

func imageBW(img image.Image) image.Image {
	bwImg := imaging.Grayscale(img)
	return bwImg
}

func uniqueFileName(outputPath string) string {
	ext := filepath.Ext(outputPath)
	base := strings.TrimSuffix(outputPath, ext)
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s%s", base, timestamp, ext)
}

func processImages(input, output string, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	files, err := os.ReadDir(input)
	if err != nil {
		log.Fatal(err)
	}

	if len(files) == 0 {
		log.Println("No input images found.")
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		wg.Add(1)
		sem <- struct{}{} // acquire a token

		go func(file os.DirEntry) {
			defer wg.Done()
			defer func() { <-sem }() // release the token

			inputPath := filepath.Join(input, file.Name())
			outputPath := filepath.Join(output, file.Name())

			// To save as new  if already exists
			if _, err := os.Stat(outputPath); err == nil {
				outputPath = uniqueFileName(outputPath)
			}

			img, err := imaging.Open(inputPath)
			if err != nil {
				log.Println("Failed to open image:", err)
				return
			}

			bwImg := imageBW(img)

			outFile, err := os.Create(outputPath)
			if err != nil {
				log.Println("Failed to create output file:", err)
				return
			}
			defer outFile.Close()

			err = jpeg.Encode(outFile, bwImg, nil)
			if err != nil {
				log.Println("Failed to encode image:", err)
			}
		}(file)
	}
}

func createZip(output, zipFilePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	files, err := os.ReadDir(output)
	if err != nil {
		log.Fatal("Failed to read output directory:", err)
	}

	var filePaths []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePaths = append(filePaths, filepath.Join(output, file.Name()))
	}

	zip := archiver.NewZip()

	// To store new file if already exists
	if _, err := os.Stat(zipFilePath); err == nil {
		zipFilePath = uniqueFileName(zipFilePath)
	}
	err = zip.Archive(filePaths, zipFilePath)
	if err != nil {
		log.Println("Failed to create zip archive:", err)
	}
}

func main() {
	inputDir := "./input"
	outputDir := "./output"
	zipFilePath := "./output/images.zip"

	// Create output directory if it doesn't exist
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create output directory:", err)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	wg.Add(1) // process image
	go processImages(inputDir, outputDir, &wg, sem)

	wg.Wait() // image processing complete

	wg.Add(1)
	go createZip(outputDir, zipFilePath, &wg)

	wg.Wait()

	log.Println("ETL process completed successfully.")
}
