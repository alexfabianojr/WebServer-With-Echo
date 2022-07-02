package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func main() {

	fmt.Println("Starting server")

	webClient := echo.New()

	webClient.GET("/alive", func(context echo.Context) error {
		return context.String(http.StatusOK, "I am alive!")
	})

	webClient.Start(":8080")
}
