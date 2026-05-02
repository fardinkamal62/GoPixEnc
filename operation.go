package main

import (
	"crypto/sha256"
    "encoding/binary"
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
// quad: RGBA xor values
// history: history of the pixels
// historyMutex: mutex for the history
// rgbaMutex: mutex for the image
func operation(img *image.Image, randomNumbers *[]int, start int, end int, wg *sync.WaitGroup, rgba **image.NRGBA, quad *[]int, history *map[int]bool, historyMutex *sync.Mutex, rgbaMutex *sync.Mutex) {
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
		currentColor := color.NRGBAModel.Convert((*img).At(currentX, currentY)).(color.NRGBA)

		newX := (*randomNumbers)[i] % width
		newY := (*randomNumbers)[i] / width
		newColor := color.NRGBAModel.Convert((*img).At(newX, newY)).(color.NRGBA)

		xorR := uint8((*quad)[0])
		xorG := uint8((*quad)[1])
		xorB := uint8((*quad)[2])
		xorA := uint8((*quad)[3])

		rgbaMutex.Lock()
		(*rgba).SetNRGBA(newX, newY, color.NRGBA{
			R: currentColor.R ^ xorR,
			G: currentColor.G ^ xorG,
			B: currentColor.B ^ xorB,
			A: currentColor.A ^ xorA,
		})
		(*rgba).SetNRGBA(currentX, currentY, color.NRGBA{
			R: newColor.R ^ xorR,
			G: newColor.G ^ xorG,
			B: newColor.B ^ xorB,
			A: newColor.A ^ xorA,
		})
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
	rgba := image.NewNRGBA(bounds)
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

	seedValue := convertToAscii(password)
	randomNumbers, err := generateUniqueRandomArray(pixels, seedValue)
	if err != nil {
		panic(err)
	}

	quad, err := generateQuadrupleUniqueNumbers(seedValue)
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
		go operation(&img, &randomNumbers, start, end, &wg, &rgba, &quad, &history, &historyMutex, &rgbaMutex)
	}

	// Last thread
	// The last thread will take the remaining pixels
	start := (usableThreads - 1) * (pixels / usableThreads)
	end := pixels - 1
	wg.Add(1)
	go operation(&img, &randomNumbers, start, end, &wg, &rgba, &quad, &history, &historyMutex, &rgbaMutex)
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
	if encrypt {
		defer fmt.Println("\nDone creating image: " + EncryptedFilename)
	} else {
		defer fmt.Println("\nDone creating image: " + DecryptedFilename)
	}

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
func generateUniqueRandomArray(length int, seed int64) ([]int, error) {
	var start = 0
	var end = length - 1

	rng := rand.New(rand.NewSource(seed))

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

// convertToAscii converts a string into a stable 64-bit seed.
// str: string from which to generate the seed
// returns the seed as int64
func convertToAscii(str string) int64 {
    sum := sha256.Sum256([]byte(str))
    seed := binary.LittleEndian.Uint64(sum[:8])
    return int64(seed)
}

// generateQuadrupleUniqueNumbers generates RGBA values.
// seed: the seed to use for the random number generator.
// returns 4 numbers: R, G, B, A (0-255).
func generateQuadrupleUniqueNumbers(seed int64) ([]int, error) {
	rng := rand.New(rand.NewSource(seed))

	numbers := make([]int, 4)

	for i := 0; i < 4; i++ {
		numbers[i] = rng.Intn(256)
	}

	return numbers, nil
}
