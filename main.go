package main

import (
	"go-api/app"
)

func main() {
	a := app.App{}
	a.Initialize("admin", "123", "go_api")

	a.Run(":8080")
}
