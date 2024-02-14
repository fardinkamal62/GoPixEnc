package main

import (
	"fmt"
	"github.com/sqweek/dialog"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

const Version string = "2.2.1"
const ExampleImage string = "images/example.jpg"
const EncryptedFilename string = "images/encrypt.png"
const DecryptedFilename string = "images/decrypt.png"

func main() {
	fmt.Println("GoPixEnc v" + Version + "!")
	fmt.Println("PixEnc implementation in Go.")
	fmt.Print("\nI want to encrypt(e)/decrypt(d): ")

	var choice string
	var FILENAME string
	var err error

	_, err = fmt.Scan(&choice)
	if err != nil {
		panic(err)
	}

	for {
		FILENAME, err = dialog.File().Title("Select a file").SetStartDir("images").Filter("All image files (*.png;*.jpg;*.jpeg)", "jpg", "jpeg", "png").Load()
		if err != nil {
			if _, err := os.Stat(ExampleImage); os.IsNotExist(err) {
				fmt.Println("Image file not found: " + ExampleImage + " (or select a file manually)")
				continue // Continue the loop to prompt again
			}
			FILENAME = ExampleImage
		}
		fmt.Println("\nSelected file:", FILENAME+"\n")
		break // Break the loop if an image is selected
	}

	fmt.Print("Enter password: ")
	var password string
	_, err = fmt.Scan(&password)
	if err != nil {
		panic(err)
	}

	if choice == "e" {
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		img, err := openAndDecodeImage(FILENAME)
		if err != nil {
			panic(err)
		}

		fmt.Println("Encrypting...")
		multiThreadOperation(img, password, true)
	} else if choice == "d" {
		img, err := openAndDecodeImage(EncryptedFilename)
		if err != nil {
			panic(err)
		}

		fmt.Println("Decrypting...")
		multiThreadOperation(img, password, false)
	}

}

// openAndDecodeImage opens and decodes the image.
// filename: the filename
// returns: the image and an error
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
