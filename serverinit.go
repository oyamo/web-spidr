package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"github.com/oyamoh-brian/spidr/downloader"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"path/filepath"
)

func  Init()  {
	// App
	engine := html.New("./template", ".html")
	app = fiber.New(fiber.Config{
		Views: engine,
	})

	dir, err := filepath.Abs(filepath.Dir(".."))
	if err != nil {
		log.Fatal(err)
	}

	if err := godotenv.Load(filepath.Join(dir, ".env")); err != nil {
		fmt.Print("No .env file found")
		panic(err)
	}
	app.Static("/static", "./static")

	// Downloader Settings
	var args = os.Args[1:]
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx  := context.Background()

	err = client.Connect(ctx)

	if err != nil {
		panic(err)
	}

	database = client.Database(os.Getenv("MONGO_DATABASE"))

	downloaderConfig = downloader.Config{
		UrlTXTPath: args[0],
		Database:   database,
	}
}
