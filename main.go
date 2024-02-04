package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

const VERSION string = "Unstable 2.0.0"
const FILENAME string = "image.png"
const EncryptedFilename string = "encrypt.png"
const DecryptedFilename string = "decrypt.png"

func main() {
	fmt.Println("GoPixEnc v" + VERSION + "!")
	fmt.Println("PixEnc implementation in Go.")
	fmt.Print("\nI want to encode(e)/decode(d): ")

	var choice string
	_, err := fmt.Scan(&choice)
	if err != nil {
		panic(err)
	}

	fmt.Print("Enter password: ")
	var password string
	_, err = fmt.Scan(&password)
	if err != nil {
		panic(err)
	}

	if choice == "e" {
		img, err := openAndDecodeImage(FILENAME)
		if err != nil {
			panic(err)
		}

		multiThreadOperation(img, password, true)
	} else if choice == "d" {
		img, err := openAndDecodeImage(EncryptedFilename)
		if err != nil {
			panic(err)
		}

		multiThreadOperation(img, password, false)
	}

}

func openAndDecodeImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	var img image.Image
	if format == "png" {
		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
	} else if format == "jpeg" || format == "jpg" {
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return img, nil
}
