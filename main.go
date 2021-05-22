package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/oyamoh-brian/spidr/downloader"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var downloaderConfig downloader.Config
var client mongo.Client
var database *mongo.Database
var app *fiber.App


func main() {

	Init()
	InitRoutes()

	defer client.Disconnect(context.Background())
	log.Fatal(app.Listen(":3000"))
}





func InitRoutes()  {
	app.Get("/", HomePage)
	app.Get("/download", Fetch)
	app.Get("/success", Success)
}