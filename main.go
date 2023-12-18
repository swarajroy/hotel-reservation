package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/swarajroy/hotel-reservation/api"
)

func main() {

	listenAddr := flag.String("listenAddr", ":3000", "The API Servers port")
	flag.Parse()
	app := fiber.New()

	apiv1 := app.Group("/api/v1")

	apiv1.Get("/users", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)

	app.Listen(*listenAddr)
}
