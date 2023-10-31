package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("USAGE: bpfreq [FILENAME] ([OUTPUT DIRECTORY])")
		return
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println("name:", file.Name())
	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	// Count byte pair frequencies in file
	frequencies := make([][]int, 256)
	for i := 0; i < len(frequencies); i++ {
		frequencies[i] = make([]int, 256)
	}
	maxFreq := 0

	for i := 0; i+1 < len(fileBytes); i++ {
		b1 := fileBytes[i]
		b2 := fileBytes[i+1]

		frequencies[b1][b2]++
		if frequencies[b1][b2] > maxFreq {
			maxFreq = frequencies[b1][b2]
		}

	}

	// turn freqs into image
	img := image.NewGray(image.Rect(0, 0, 256, 256))
	for y, freqs := range frequencies {
		for x, val := range freqs {
			brightness := math.Log(1+float64(val)) / math.Log(1+float64(maxFreq)) // From 0 to 1
			img.SetGray(x, y, color.Gray{uint8(brightness * 255)})
		}
	}

	extension := filepath.Ext(file.Name())
	outputFileName := strings.TrimSuffix(filepath.Base(file.Name()), extension) + "-bpvis.png"
	outputDir := ""
	if len(os.Args) == 3 {
		info, err := os.Stat(os.Args[2])
		if err != nil {
			panic(err)
		}

		if !info.IsDir() {
			fmt.Println("USAGE: bpfreq [FILENAME] ([OUTPUT DIRECTORY])")
			return
		}
		outputDir = os.Args[2]
	}
	outFile, err := os.Create(filepath.Join(outputDir, outputFileName))

	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	png.Encode(outFile, img)

}
