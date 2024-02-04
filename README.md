<div align="center"><h1>GoPixEnc</h1></div>
<div align="center"><a href="https://github.com/fardinkamal62/PixEnc">PixEnc</a> in Go</div>
<div align="center">Encrypt image by manipulating pixels</div>
<div align="center" style="color: grey"><sub>Version: Unstable 2.0.0</sub></div>
<div align="center">
  <strong>
    <a href="https://fardinkamal62.vercel.app/projects/pixenc">Website</a>
    •
    <a href="https://docs.google.com/document/d/173xWvlrEQd1esI3rtD1SmtqtZ1rmFFwKzwRIdWKSTQw/edit?usp=sharing">Docs</a>
    </strong>
</div>

<hr />

This is a Go implementation of my [PixEnc](https://github.com/fardinkamal62/PixEnc) project, which was originally
written in Python.

I made this project to learn Go and I tried a different approach of encrypting image pixels.

I am not a Go developer. So, if you find any issue, please create an issue or pull request.

# Note

This version is not stable & accurate. Due to multithreading, some pixels are not being encrypted/decrypted properly

### Problem

I'm generating random numbers for each pixel. Suppose if pixel 10 gets 5 as random number, 5 may/may not
get 10 as its random number. Even though I'm keeping a history operation on pixels, if 5 gets 2 as its random number,
and it's on different thread, it will swap 5 with 2, not 10 with 5. So, the image will be corrupted.

### Potential Solution

I'm thinking of using a different approach. Instead of generating random numbers for each pixel, I will generate n/2
random numbers for n pixels.

Suppose if I have 100 pixels, I will generate 50 random numbers from 51 to 100 and assign 0-50 to them. So that, if 55
gets 2 as its random number, 2 will get 55 as its random number. I'm still thinking about a better approach.

# Technologies

- Go

# Example

## Original Image

**Original Image**

![Original Image](https://i.ibb.co/717YFZ3/image.png)
![Original Image](https://i.ibb.co/GPrdJjp/image.png)

**Encrypted Image**

![Encrypted Image](https://i.ibb.co/tQF5Pn7/encrypt.png)
![Encrypted Image](https://i.ibb.co/cCzGLgL/encrypt.png)

**Decrypted Image**

![Decrypted Image](https://i.ibb.co/9rhKkgr/decrypt.png)
![Decrypted Image](https://i.ibb.co/HgSTFV5/decrypt.png)

# Build

- Clone the repository
- Run `go build` to build & generate an executable file
- Run `go run .` to run the program

# Usage

- Download the appropriate latest release from GitHub Releases
- Run the executable file

### Encrypting Image

1. Keep the image you want to encrypt in the same directory as the executable file
2. Rename the image to `image.png`(if PNG) or `image.jpg`(if JPG)
3. Run the executable file
4. In choice menu, select `e` to encrypt image
5. Enter the password you want to use to encrypt the image
6. It will generate a file named `encrypt.png`

### Decrypting Image

1. Keep the image you want to decrypt in the same directory as the executable file | If you have encrypted the image
   using this program, then you can skip this step
2. Run the executable file
3. In choice menu, select `d` to decrypt image
4. Enter the password you used to encrypt the image
5. It will generate a file named `decrypt.png`

# Release Note

### 2.0.0 (Unstable, Current)

- Multithreading

Check note for more info
### 1.1.1 (Current)

- Password now supports any ASCII value. In v1.0.0 only numbers were allowed as password
- Using `rand.New(rand.NewSource(seed))` instead of deprecated `rand.seed(seed)`
- Fixed wrong output on decryption: `Done creating image decrypt.png` -> `Done creating image`

### 1.0.0

- Can encrypt & decrypt images

# Planned Features

- [x] Multi-threading
- [ ] Better random number generation
- [ ] File explorer
- [ ] GUI
