package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oyamoh-brian/spidr/downloader"
	"log"
	"net/http"
)


func Success(c *fiber.Ctx) error {

	return c.Render("success", nil)
}

func Fetch(c *fiber.Ctx) error  {
	var d = downloader.New(downloaderConfig)
	_, fileReadingErr := d.Download()
	if fileReadingErr != nil {
		log.Fatal(fileReadingErr)
	}
	return c.Redirect("/success", http.StatusFound)
}