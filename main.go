package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/skip2/go-qrcode"
)

var ctx = context.Background()

func main() {
	// Create a new template engine
	engine := html.New("./templates", ".html")
	// Or from an embedded system
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retreive the custom status code if it's an fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Send custom error page
			err = c.Render("error", fiber.Map{
				"Title": "Something bad happened 😳",
				"Code":  code,
				"Error": err.Error(),
			}, "layouts/main")

			// Return from handler
			return nil
		},
	})
	app.Use(logger.New())

	// Cache middleware for all routes
	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	// Serve static files
	app.Static("/", "./public")

	// Render index template
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":       "QR Encode This",
			"Description": "This site allows you to encode any data into a QR code. You can then scan the QR code with your phone to get the data back. Or you can download the QR code as an image. Or you can copy the URL of the page and share it with someone else.",
			"Url":         "https://qrencodethis.com",
		}, "layouts/app")
	})

	// Render a QR code
	app.Get("/qr", func(c *fiber.Ctx) error {
		data := c.Query("data")
		// if no data is provided, render the form
		if data == "" {
			return c.Render("form", fiber.Map{
				"Title":       "QR Encode This",
				"Description": "This site allows you to encode any data into a QR code. You can then scan the QR code with your phone to get the data back. Or you can download the QR code as an image. Or you can copy the URL of the page and share it with someone else.",
				"Url":         "https://qrencodethis.com",
				"Data":        data,
			})
		}
		// if data is provided, render the QR code
		var png []byte
		png, err := qrcode.Encode(data, qrcode.Medium, 256)
		if err != nil {
			log.Println("❌ Error generating QR code")
			log.Println(err)
			return c.Render("form", fiber.Map{
				"Title":       "QR Encode This",
				"Description": "This site allows you to encode any data into a QR code. You can then scan the QR code with your phone to get the data back. Or you can download the QR code as an image. Or you can copy the URL of the page and share it with someone else.",
				"Url":         "https://qrencodethis.com",
				"Data":        data,
				"Error":       "❌ Error generating QR code, please try again.",
			})
		}

		// Convert to base64
		b64 := base64.StdEncoding.EncodeToString(png)

		return c.Render("qr", fiber.Map{
			"Title":       "QR Encode This",
			"Description": "This site allows you to encode any data into a QR code. You can then scan the QR code with your phone to get the data back. Or you can download the QR code as an image. Or you can copy the URL of the page and share it with someone else.",
			"Url":         "https://qrencodethis.com",
			"Data":        data,
			"Image":       b64,
		})
	})

	app.Get("/share", func(c *fiber.Ctx) error {
		data := c.Query("data")
		if data == "" {
			log.Println("❌ No data provided")
			return c.Render("qr_image", fiber.Map{
				"Title":       "QR Encode This",
				"Description": "This site allows you to encode any data into a QR code. You can then scan the QR code with your phone to get the data back. Or you can download the QR code as an image. Or you can copy the URL of the page and share it with someone else.",
				"Url":         "https://qrencodethis.com",
				"Data":        data,
				"Error":       "❌ No data provided.",
			}, "layouts/app")
		}
		// if data is provided, render the QR code
		var png []byte
		png, err := qrcode.Encode(data, qrcode.Medium, 256)
		if err != nil {
			log.Println("❌ Error generating QR code")
			log.Println(err)
			return c.Render("qr_image", fiber.Map{
				"Title":       "QR Encode This",
				"Description": "This site allows you to encode any data into a QR code. You can then scan the QR code with your phone to get the data back. Or you can download the QR code as an image. Or you can copy the URL of the page and share it with someone else.",
				"Url":         "https://qrencodethis.com",
				"Data":        data,
				"Error":       "❌ Error generating QR code",
			}, "layouts/app")
		}

		// Convert to base64
		b64 := base64.StdEncoding.EncodeToString(png)

		return c.Render("qr_image", fiber.Map{
			"Title":       "Someone shared this QR code with you",
			"Description": "This site allows you to encode any data into a QR code. You can then scan the QR code with your phone to get the data back. Or you can download the QR code as an image. Or you can copy the URL of the page and share it with someone else.",
			"Url":         "https://qrencodethis.com",
			"Data":        data,
			"Image":       b64,
		}, "layouts/app")
	})

	// Handle page not found; must be the last route
	app.Use(func(c *fiber.Ctx) error {
		return c.Render("error", fiber.Map{
			"Title":       "Page not found 😭",
			"Description": "We tried our best... looked everywhere, but we couldn't find this page.",
			"Url":         "https://qrencodethis.com",
			"Code":        404,
			"Error":       "We tried our best... looked everywhere, but we couldn't find this page.",
		}, "layouts/main")
	})

	// Start server
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000" // default port if environment variable is not set
	}

	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
