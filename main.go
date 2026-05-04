package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/sqweek/dialog"
)

const Version string = "3.1.0"
const ExampleImage string = "images/example.jpg"
const EncryptedFilename string = "images/encrypt.png"
const DecryptedFilename string = "images/decrypt.png"
const divider string = "----------------------------------------------"

const colorReset string = "\033[0m"
const colorCyan string = "\033[36m"
const colorYellow string = "\033[33m"
const colorDim string = "\033[2m"
const colorBold string = "\033[1m"

func main() {
	printBanner()
	fmt.Println(colorDim + "PixEnc implementation in Go." + colorReset)
	printDivider()
	fmt.Println(colorCyan + "[1/3] Mode" + colorReset)
	fmt.Print(colorYellow + "  e = encrypt, d = decrypt > " + colorReset)

	var choice string
	var FILENAME string
	var err error

	_, err = fmt.Scan(&choice)
	if err != nil {
		panic(err)
	}

	fmt.Println(colorCyan + "\n[2/3] Source Image" + colorReset)
	for {
		FILENAME, err = dialog.File().Title("Select a file").SetStartDir("images").Filter("All image files (*.png;*.jpg;*.jpeg)", "jpg", "jpeg", "png").Load()
		if err != nil {
			if _, err := os.Stat(ExampleImage); os.IsNotExist(err) {
				fmt.Println("Image file not found: " + ExampleImage + " (or select a file manually)")
				continue // Continue the loop to prompt again
			}
			FILENAME = ExampleImage
		}
		fmt.Println("\n"+colorBold+"Selected file:"+colorReset, FILENAME+"\n")
		break // Break the loop if an image is selected
	}

	fmt.Println(colorCyan + "[3/3] Password" + colorReset)
	fmt.Print(colorYellow + "  Enter password > " + colorReset)
	var password string
	_, err = fmt.Scan(&password)
	if err != nil {
		panic(err)
	}

	checkFolderExistence()

	if choice == "e" {
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		img, err := openAndDecodeImage(FILENAME)
		if err != nil {
			panic(err)
		}

		fmt.Println("\nEncrypting...")
		multiThreadOperation(img, password, true)
	} else if choice == "d" {
		img, err := openAndDecodeImage(EncryptedFilename)
		if err != nil {
			panic(err)
		}

		fmt.Println("\nDecrypting from", EncryptedFilename+"...")
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
	switch format {
	case "png":
		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
	case "jpeg", "jpg":
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return img, nil
}

// checkFolderExistence checks if the images folder exists and creates it if it doesn't.
func checkFolderExistence() {
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		err := os.Mkdir("images", 0755)
		if err != nil {
			panic(err)
		}
	}
}

func printBanner() {
	fmt.Print("\033[H\033[2J") // ANSI escape code to clear the screen and move the cursor to the top-left corner
	banner, err := os.ReadFile("banner.txt")
	if err == nil {
		fmt.Print(colorCyan + string(banner) + colorReset)
	} else {
		fmt.Print(colorCyan + "  ____       ____ _      ______           \n" +
			" / ___| ___ |  _ \\ | ___|  ____|__  _ __  \n" +
			"| |  _ / _ \\| |_) | |/ _ \\  _| / _ \\| '_ \\\n" +
			"| |_| | (_) |  __/| |  __/ |__| (_) | | | |\n" +
			" \\____|\\___/|_|   |_|\\___|_____|\\___/|_| |_|\n" + colorReset)
	}
	fmt.Println(colorBold + "v" + Version + colorReset)
	printDivider()
}

func printDivider() {
	fmt.Println(colorCyan + divider + colorReset)
}
