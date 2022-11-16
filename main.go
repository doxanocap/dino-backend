package main

import (
	"github.com/doxanocap/reactNative/dino-back/pkg/database"
	"github.com/doxanocap/reactNative/dino-back/pkg/routes"
)

func main() {
	database.Connect()
	routes.Router()
}
