<div align="center"><h1>GoPixEnc</h1></div>

# Abstract
This is a Go implementation of my [PixEnc](https://github.com/fardinkamal62/PixEnc) project, which was originally written in Python.
It's different from the original project in terms of programming language and approach of encrypting image pixels.

I made this project to learn Go and I tried a different approach of encrypting image pixels.

#### How it works is
It gets pixel count, generates random numbers the amount by which pixels are calculated; ranging from 0 to pixel count. Then it swaps the RGBA values of the pixels with the random numbered pixels.

# Procedure
## Encryption
1. Read the image with Image library. It will return an Image object.
2. From the Image object, we can get the pixel count, width and height of the image
3. Generate random numbers the amount by which pixels are calculated; ranging from 0 to pixel count.
      
   **Example:** If the image has 1024 pixels, then generate 1024 random numbers ranging from 0 to 1024 and each value will be unique. With the same password, we will get the same 1024 random values everytime. 
4. Swap RGBA values with corresponding random numbered pixels.
      
   **Example:** If the 10-th random number is 3, then swap the 10-th pixel's RGBA value with the 3-rd pixel's RGBA value. So the 10-th pixel's RGBA value will be the 3-rd pixel's RGBA value and the 3-rd pixel's RGBA value will be the 10-th pixel's RGBA value.
5. Save the image.

### Example

I have a 32x32 image. It has 1024 pixels. First generate 1024 random numbers with password as seed, so everytime with the correct password, we will get the same 1024 random values.

Let’s say we’ll change the 10-th pixel. Our 10-th random number is 3 and the 10-th pixel value is R: 81 G: 123 B:21 A:100. Pixel value of 3-rd pixel is R: 82 G:120 B: 129 A: 210.

Now swap the 10-th pixel’s RGBA value with the 3-rd pixel’s RGBA value. So the 10-th pixel’s RGBA value will be R: 82 G:120 B: 129 A: 210 and the 3-rd pixel’s RGBA value will be R: 81 G: 123 B:21 A:100.

## Decryption
Decryption is same as encryption.

As our 10-th pixel's RGBA value is in the 3-rd pixel and the 3-rd pixel's RGBA value is in the 10-th pixel, we can swap the RGBA values again to get the original image back. 
