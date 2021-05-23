package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/oyamoh-brian/spidr/downloader"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"strings"
)

var downloaderConfig downloader.Config
var client mongo.Client
var database *mongo.Database
var app *fiber.App


func main() {

	Init()
	InitRoutes()

	defer client.Disconnect(context.Background())
	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}





func InitRoutes()  {

	app.Use(func(c *fiber.Ctx) error {
		scheme := c.Protocol()
		if scheme == "http" {
			url := strings.ReplaceAll(c.BaseURL(), "http://", "https://")
			return c.Redirect(url)
		}
		return c.Next()
	})

	app.Get("/", HomePage)
	app.Get("/download", Fetch)
	app.Get("/success", Success)
}