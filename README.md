<div align="center"><h1>GoPixEnc</h1></div>
<div align="center"><a href="https://github.com/fardinkamal62/PixEnc">PixEnc</a> in Go</div>
<div align="center">Encrypt image by manipulating pixels</div>
<div align="center" style="color: grey"><sub>Version: 2.2.2</sub></div>
<div align="center">
  <strong>
    <a href="https://fardinkamal62.vercel.app/projects/pixenc">Website</a>
    •
    <a href="https://docs.google.com/document/d/173xWvlrEQd1esI3rtD1SmtqtZ1rmFFwKzwRIdWKSTQw/edit?usp=sharing">Docs</a>
    </strong>
</div>
<br>
GoPixEnc is an application written in Go, designed for encrypting and decrypting images securely. This app was developed as a project submission for the CSE Spring Fest 2024 at <a href="https://diu.ac">Dhaka International University</a>

# Technologies

- Go

# Example

**Original Image**

![Original Image](https://i.ibb.co/717YFZ3/image.png)
![Original Image](https://i.ibb.co/GPrdJjp/image.png)
<img src="https://i.ibb.co/wd6dLdg/barakah.jpg" alt="barakah-original" style="width:20%; height: auto"/>

**Encrypted Image**

![Encrypted Image](https://i.ibb.co/tQF5Pn7/encrypt.png)
![Encrypted Image](https://i.ibb.co/cCzGLgL/encrypt.png)
<img src="https://i.ibb.co/z7HHJMZ/encrypt.png" alt="barakah-encrypted" style="width:20%; height: auto"/>

**Decrypted Image**

![Decrypted Image](https://i.ibb.co/9rhKkgr/decrypt.png)
![Decrypted Image](https://i.ibb.co/HgSTFV5/decrypt.png)
<img src="https://i.ibb.co/5K0GgQ0/decrypt.png" alt="barakah-decrypted" style="width:20%; height: auto"/>

# Build

- Clone the repository
- Run `go build -o build/GoPixEnc` to build & generate an executable file at `build` folder
- Run `go run .` to run the program without building 

### Cross-Platform Build

**Important Variables**
- `GOOS` for target OS
- `GOARCH` for target architecture

**Example**

To build for current OS & architecture: `go build -o build/GoPixEnc_vx.y.z_os_arch`

To build for Windows: `GOOS=windows GOARCH=amd64 go build -o build/GoPixEnc_vx.y.z_windows_arch.exe`

To build for Linux: `GOOS=linux GOARCH=amd64 go build -o build/GoPixEnc_vx.y.z_linux_arch`

# Usage

- Download the appropriate latest release from GitHub Releases
- Run the executable file

### Encrypting Image

1. Run the executable file
2. Chose the image you want to encrypt from the file explorer or leave it empty to use the default image
3. In choice menu, select `e` to encrypt image
4. Enter the password you want to use to encrypt the image
5. It will generate a file named `encrypt.png` in the same directory as the executable file

### Decrypting Image

1. Run the executable file
2. Chose the image you want to decrypt from the file explorer
3. In choice menu, select `d` to decrypt image
4. Enter the password you used to encrypt the image
5. It will generate a file named `decrypt.png` in the same directory as the executable file

# Release Note

### 2.2.2 (Current)

- Create `images` folder if not exists

### 2.2.1

- File explorer UX improvement: Now, it will show the file explorer after encryption/decryption choice is made & file explorer will remain open until a file is selected
- Better error handling in file explorer
- Updated package versions
- Separate images folder

### 2.2.0

- Dialogue box to choose image
- Included example image in the repository

### Beta 2.1.1

- Image can be image.png or image.jpg
- Updated README with new example images

### Beta 2.1.0

- Fixed issue with multithreading of version 2.0.0
- New algorithm to generate random numbers for pixels

### 2.0.0 (Unstable)

- Multithreading

Check note for more info

### Note

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

### 1.1.1

- Password now supports any ASCII value. In v1.0.0 only numbers were allowed as password
- Using `rand.New(rand.NewSource(seed))` instead of deprecated `rand.seed(seed)`
- Fixed wrong output on decryption: `Done creating image decrypt.png` -> `Done creating image`

### 1.0.0

- Can encrypt & decrypt images

# Planned Features

- [x] Multi-threading
- [x] Better random number generation
- [x] File explorer
- [ ] GUI
