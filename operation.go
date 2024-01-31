package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var history = make([]int, 0)

func operation(img image.Image, password string, encrypt bool) {
	// Initialize variables
	bounds := img.Bounds() // get image bounds (width and height)
	width, height := bounds.Max.X, bounds.Max.Y
	pixels := width * height
	var outputFile *os.File
	var err error

	if encrypt {
		outputFile, err = os.Create(EncryptedFilename)
		if err != nil {
			panic(err)
		}
	} else {
		outputFile, err = os.Create(DecryptedFilename)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("\nStarting operation...")
	defer fmt.Println("\nOperation complete.")

	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

	randomNums, err := generateUniqueRandomArray(pixels, password)
	if err != nil {
		panic(err)
	}

	startTime := time.Now()                   // operation start timer
	bar := progressbar.Default(int64(pixels)) // progress bar
	for i := 0; i < pixels; i++ {
		err := bar.Add(1)
		if err != nil {
			return
		}

		if containsElement(i) || containsElement(randomNums[i]) {
			continue
		}

		history = append(history, i)
		history = append(history, randomNums[i])

		currentX := i % width
		currentY := i / width
		currentR, currentG, currentB, currentA := img.At(currentX, currentY).RGBA()

		newX := randomNums[i] % width
		newY := randomNums[i] / width
		newR, newG, newB, newA := img.At(newX, newY).RGBA()

		rgba.SetRGBA(newX, newY, color.RGBA{R: uint8(currentR / 257), G: uint8(currentG / 257), B: uint8(currentB / 257), A: uint8(currentA / 257)})
		rgba.SetRGBA(currentX, currentY, color.RGBA{R: uint8(newR / 257), G: uint8(newG / 257), B: uint8(newB / 257), A: uint8(newA / 257)})
	}
	endTime := time.Now()
	fmt.Println("\nOperation took", endTime.Sub(startTime))

	fmt.Println("\nCreating image file...")
	defer fmt.Println("\nDone creating image", DecryptedFilename)

	startTime = time.Now() // file write start timer

	err = png.Encode(outputFile, rgba)
	if err != nil {
		return
	}
	err = outputFile.Close()
	endTime = time.Now()

	fmt.Println("\nWriting to file took", endTime.Sub(startTime))
	if err != nil {
		return
	}
}

// generateUniqueRandomArray generates a random array of numbers from 0 to length-1.
// length: the length of the array.
// seed: the seed to use for the random number generator.
// returns the array of numbers and an error if there is one.
func generateUniqueRandomArray(length int, seed string) ([]int, error) {
	var start = 0
	var end = length - 1

	seedValue, err := strconv.Atoi(seed)
	if err != nil {
		return nil, err
	}
	rand.Seed(int64(seedValue))

	numbers := make([]int, length)
	used := make(map[int]bool)

	for i := 0; i < length; {
		randomNumber := rand.Intn(end-start+1) + start
		if !used[randomNumber] {
			numbers[i] = randomNumber
			used[randomNumber] = true
			i++
		}
	}

	return numbers, nil
}

// containsElement checks if an element is in the history array.
// target: the element to check for.
// returns true if it is, false if it isn't
func containsElement(target int) bool {
	for _, element := range history {
		if element == target {
			return true
		}
	}
	return false
}
