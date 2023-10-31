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

	frequencies := make([][]int, 256)
	for i := 0; i < len(frequencies); i++ {
		frequencies[i] = make([]int, 256)
	}
	buffer := make([]byte, 2) // Create a buffer to hold each byte pair
	maxFreq := 0
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // Exit loop at end of file
			}
			fmt.Printf("Read error: %v\n", err)
			return // Exit if we encounter an error other than EOF
		}

		if bytesRead < 2 {
			break // Exit loop if we don't have a complete byte pair
		}

		// update max frequency
		frequencies[buffer[0]][buffer[1]]++
		if frequencies[buffer[0]][buffer[1]] > maxFreq {
			maxFreq = frequencies[buffer[0]][buffer[1]]
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
