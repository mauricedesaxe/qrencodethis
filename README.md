# QR Encode This 

This is a web application built in Go that allows users to encode any data into a QR code. 

## Features

1. Generate a QR code from any given data
2. Download the generated QR code as an image
3. Copy a shareable link that encodes the given data into a QR code

## Usage

1. Enter the data to be encoded in the text field.
2. Click on 'QR encode this' to generate the QR code.
3. Download the QR code by clicking on 'Download image'.
4. To share a link that will generate the same QR code, click on 'Copy link'. 

## Setup & Installation

Make sure you have installed Go, SQLite and Redis.

1. Clone this repository: `git clone https://github.com/<yourusername>/qrencodethis.git`
2. Navigate to the project directory: `cd qrencodethis`
3. Copy `.env.example` to `.env` and set your Redis URL: `cp .env.example .env`
4. Install dependencies: `go mod download`
5. Run the project: `go run main.go`

The server will start on `http://localhost:3000`, or on the port you've set in the PORT environment variable.

## Technologies

- [Go](https://golang.org)
- [Fiber](https://github.com/gofiber/fiber/v2)
- [GORM](https://gorm.io)
- [Redis](https://redis.io)
- [SQLite](https://www.sqlite.org/index.html)
- [go-qrcode](https://github.com/skip2/go-qrcode)

## License

GPL-3
