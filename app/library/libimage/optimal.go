package libimage

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

// OptimizePath ...
func OptimizePath(inPath string, outPath string) {
	quantization := 20
	var colorConversion ColorConversion

	// load image
	inFile, openErr := os.Open(inPath)
	if openErr != nil {
		fmt.Printf("couldn't open %v: %v\n", inPath, openErr)
		return
	}

	inFile.Stat()
	decoded, _, decodeErr := image.Decode(inFile)
	inFile.Close()
	if decodeErr != nil {
		fmt.Printf("couldn't decode %v: %v\n", inPath, decodeErr)
		return
	}

	optimized := Compress(decoded, colorConversion, quantization)

	// save optimized image
	outFile, createErr := os.Create(outPath)
	if createErr != nil {
		fmt.Printf("couldn't create %v: %v\n", outPath, createErr)
		return
	}

	encodeErr := png.Encode(outFile, optimized)
	outFile.Stat()
	outFile.Close()
	if encodeErr != nil {
		fmt.Printf("couldn't encode %v: %v\n", inPath, encodeErr)
		return
	}

	err := os.Rename(outPath, inPath)
	if err == nil {
		outPath = inPath
	} else {
		fmt.Printf("Cannot rewrite original file %s \n", err)
	}
}
