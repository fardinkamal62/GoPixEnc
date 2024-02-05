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
	"runtime"
	"sync"
	"time"
)

// operation performs the operation on the image.
// img: the image
// randomNumbers: random numbers to use for the operation
// start: start index of the random numbers
// end: end index (inclusive) of the random numbers
// wg: wait group for the threads
// rgba: image to write to
// history: history of the pixels
// historyMutex: mutex for the history
// rgbaMutex: mutex for the image
func operation(img *image.Image, randomNumbers *[]int, start int, end int, wg *sync.WaitGroup, rgba **image.RGBA, history *map[int]bool, historyMutex *sync.Mutex, rgbaMutex *sync.Mutex) {
	// Initialize variables
	bounds := (*img).Bounds() // get image bounds (width and height)
	width, height := bounds.Max.X, bounds.Max.Y
	pixels := width * height

	var err error

	fmt.Println("\nStarting operation...")
	defer fmt.Println("\nOperation completed")

	startTime := time.Now()                   // operation start timer
	bar := progressbar.Default(int64(pixels)) // progress bar
	for i := start; i <= end; i++ {
		err := bar.Add(1)
		if err != nil {
			panic(err)
		}

		historyMutex.Lock()
		if containsElement(history, i) || containsElement(history, (*randomNumbers)[i]) {
			historyMutex.Unlock()
			continue
		}

		(*history)[i] = true
		(*history)[(*randomNumbers)[i]] = true
		historyMutex.Unlock()

		currentX := i % width
		currentY := i / width
		currentR, currentG, currentB, currentA := (*img).At(currentX, currentY).RGBA()

		newX := (*randomNumbers)[i] % width
		newY := (*randomNumbers)[i] / width
		newR, newG, newB, newA := (*img).At(newX, newY).RGBA()

		rgbaMutex.Lock()
		(*rgba).SetRGBA(newX, newY, color.RGBA{R: uint8(currentR >> 8), G: uint8(currentG >> 8), B: uint8(currentB >> 8), A: uint8(currentA >> 8)})
		(*rgba).SetRGBA(currentX, currentY, color.RGBA{R: uint8(newR >> 8), G: uint8(newG >> 8), B: uint8(newB >> 8), A: uint8(newA >> 8)})
		rgbaMutex.Unlock()
	}
	endTime := time.Now()
	fmt.Println("\nOperation took", endTime.Sub(startTime))

	if err != nil {
		panic(err)
	}
	wg.Done()
}

// multiThreadOperation performs the operation on the image.
// img: the image
// password: the password
// encrypt: whether to encrypt or decrypt
func multiThreadOperation(img image.Image, password string, encrypt bool) {
	var err error

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	pixels := width * height

	var rgbaMutex sync.Mutex
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

	randomNumbers, err := generateUniqueRandomArray(pixels, password)
	if err != nil {
		panic(err)
	}

	totalThreads := runtime.NumCPU()
	usableThreads := totalThreads / 2

	var wg sync.WaitGroup

	var historyMutex sync.Mutex
	history := make(map[int]bool)

	for i := 0; i < usableThreads-1; i++ {
		start := i * (pixels / usableThreads)
		end := (i+1)*(pixels/usableThreads) - 1
		wg.Add(1)
		go operation(&img, &randomNumbers, start, end, &wg, &rgba, &history, &historyMutex, &rgbaMutex)
	}

	// Last thread
	// The last thread will take the remaining pixels
	start := (usableThreads - 1) * (pixels / usableThreads)
	end := pixels - 1
	wg.Add(1)
	go operation(&img, &randomNumbers, start, end, &wg, &rgba, &history, &historyMutex, &rgbaMutex)
	//

	wg.Wait() // wait for all threads to finish

	var outputFile *os.File
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
	fmt.Println("\nCreating image file...")
	defer fmt.Println("\nDone creating image")

	startTime := time.Now() // file write start timer

	err = png.Encode(outputFile, rgba)
	if err != nil {
		return
	}
	err = outputFile.Close()
	endTime := time.Now()

	fmt.Println("\nWriting to file took", endTime.Sub(startTime))
}

// generateUniqueRandomArray generates a random array of numbers from 0 to length-1.
// length: the length of the array.
// seed: the seed to use for the random number generator.
// returns the array of numbers and an error if there is one.
func generateUniqueRandomArray(length int, seed string) ([]int, error) {
	var start = 0
	var end = length - 1

	seedValue := convertToAscii(seed)

	rng := rand.New(rand.NewSource(int64(seedValue)))

	numbers := make([]int, length)
	used := make(map[int]bool)

	for i := 0; i < length; {
		randomNumber := rng.Intn(end - start + 1)

		if used[i] && used[randomNumber] {
			i++
			continue
		}
		if !used[i] && !used[randomNumber] {
			numbers[i] = randomNumber
			numbers[randomNumber] = i

			used[randomNumber] = true
			used[i] = true
			i++
		}
	}

	return numbers, nil
}

// containsElement checks if an element is in the history array.
// arr: the map to check in.
// target: the element to check for.
// returns true if it is, false if it isn't
func containsElement(arr *map[int]bool, target int) bool {
	_, ok := (*arr)[target]
	return ok
}

// convertToAscii converts ASCII characters to decimal values
// str: string from which to generate the ASCII values
// returns the ASCII  values consolidated
func convertToAscii(str string) int {
	var asciiSlice []int
	result := 0

	for i := 0; i < len(str); i++ {
		asciiSlice = append(asciiSlice, int(str[i]))
		result = (result * 100) + int(str[i])
	}

	return result
}
