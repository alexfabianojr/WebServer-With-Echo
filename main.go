package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func main() {

	fmt.Println("Starting server")

	webClient := echo.New()

	webClient.GET("/alive", handleAlive)

	webClient.GET("/cats/:data", handleSearchCats)

	webClient.Start(":8080")
}

func handleAlive(context echo.Context) error {
	return context.String(http.StatusOK, "I am alive!")
}

func handleSearchCats(context echo.Context) error {
	catName := context.QueryParam("name")
	catType := context.QueryParam("type")
	dataType := context.Param("data")

	if dataType == "string" {
		return context.String(http.StatusOK, fmt.Sprintf("Your cat name is %s and type is %s", catName, catType))
	}

	if dataType == "json" {
		return context.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return context.JSON(http.StatusBadRequest, map[string]string{
		"error": "You need to insert string or json data type",
	})
}
