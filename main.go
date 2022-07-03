package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	middleware2 "github.com/labstack/echo/middleware"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Starting server")

	webClient := echo.New()

	webClient.GET("/alive", handleAlive)
	webClient.GET("/cats/:data", handleSearchCats)

	webClient.POST("/cats", handlerAddCat)
	webClient.POST("/dogs", handleAddDog)
	webClient.POST("/hamster", handleAddHamster) // cleaner way with pure echo

	webClient.Start(":8080")

	groups := webClient.Group("/admin", middleware2.Logger())

	groups.Use(middleware2.LoggerWithConfig(middleware2.LoggerConfig{
		Format: `[${time_rfc3339} ${status} ${method} ${path}]` + "/n",
	})) // best way to add middleware

	groups.GET("/admin", mainAdmin)
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

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func handlerAddCat(context echo.Context) error {
	cat := Cat{}

	defer context.Request().Body.Close()

	body, error := ioutil.ReadAll(context.Request().Body)

	if error != nil {
		log.Printf("Failed to read request body, error: %s", error)
		return context.String(http.StatusInternalServerError, "Error fetching request body")
	}

	error = json.Unmarshal(body, &cat)

	if error != nil {
		log.Printf("json body failed parse to data structure")
		return context.String(http.StatusInternalServerError, "Error parsing request body")
	}

	log.Printf("This is your cat: %v", cat)

	return context.String(http.StatusAccepted, "Created cat")
}

type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func handleAddDog(context echo.Context) error {
	dog := Dog{}

	defer context.Request().Body.Close()

	error := json.NewDecoder(context.Request().Body).Decode(&dog)

	if error != nil {
		log.Printf("Failed to read request body, error: %s", error)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("This is your dog: %v", dog)

	return context.String(http.StatusAccepted, "Created dog")
}

type Hamster struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func handleAddHamster(context echo.Context) error {
	hamster := Hamster{}

	error := context.Bind(&hamster)

	if error != nil {
		log.Printf("Failed to read request body, error: %s", error)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("This is your hamster: %v", hamster)

	return context.String(http.StatusAccepted, "Created dog")
}

func mainAdmin(context echo.Context) error {
	return context.String(http.StatusOK, "You found an grouped link")
}
