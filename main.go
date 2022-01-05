package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"

	"github.com/hittaito/go-practice/ioseek"
)

func main() {
	fp, err := os.Open("sample.png")
	if err != nil {
		os.Exit(1)
	}
	defer fp.Close()

	/* io.ReadAtLeast
	b, err := ioreadatleast.IsPNG(fp)
	if err != nil {
		os.Exit(0)
	}
	fmt.Printf("result: %v", b)
	*/

	// io.Seek
	b, err := ioseek.IsPNG(fp)
	if err != nil {
		os.Exit(0)
	}
	fmt.Printf("result: %v\n", b)

	_, str, err := image.DecodeConfig(fp)
	if err != nil {
		fmt.Println(fmt.Printf("error: %v", err.Error()))
		os.Exit(0)
	}
	fmt.Println(str)
}
