package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	middleware2 "github.com/labstack/echo/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {

	fmt.Println("Starting server")

	webClient := echo.New()

	webClient.Use(ServerHeader)

	webClient.GET("/alive", handleAlive)
	webClient.GET("/cats/:data", handleSearchCats)
	webClient.GET("/login", handleLogin)

	webClient.POST("/cats", handlerAddCat)
	webClient.POST("/dogs", handleAddDog)
	webClient.POST("/hamster", handleAddHamster) // cleaner way with pure echo

	groups := webClient.Group("/admin")

	groups.Use(middleware2.LoggerWithConfig(middleware2.LoggerConfig{
		Format: `[${time_rfc3339} ${status} ${method} ${path}]` + "/n",
	})) // best way to add middleware

	groups.Use(middleware2.BasicAuth(
		func(username string, password string, context echo.Context) (bool, error) {
			if username == "jack" && password == "1234" {
				return true, nil
			} else {
				return false, echo.ErrForbidden
			}
		}))

	groups.GET("/main", handleMainAdmin)

	cookieGroup := webClient.Group("/cookies")

	cookieGroup.Use(checkCookie)

	cookieGroup.GET("", handleLoginCookies)

	error := webClient.Start(":8080")

	if error != nil {
		panic(error)
	}
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

func handleMainAdmin(context echo.Context) error {
	return context.String(http.StatusOK, "You found an grouped link")
}

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		context.Response().Header().Set(echo.HeaderServer, "BlueBot 1.0")
		return next(context)
	}
}

func handleLoginCookies(context echo.Context) error {
	return context.String(http.StatusOK, "Cookies page")
}

func handleLogin(context echo.Context) error {
	username := context.QueryParam("username")
	password := context.QueryParam("password")

	if username == "jack" && password == "1234" {
		cookie := &http.Cookie{} // same than cookie := new(http.Cookie)

		cookie.Name = "sessionID"
		cookie.Value = "some_string"
		cookie.Expires = time.Now().Add(48 * time.Hour)

		context.SetCookie(cookie)

		return context.String(http.StatusOK, "You are logged in!")
	}

	return context.String(http.StatusUnauthorized, "Wrong username or password")
}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		cookie, error := context.Cookie("sessionID")

		if error != nil {
			log.Println(error.Error())

			if strings.Contains(error.Error(), "named cookie not present") {
				return context.String(http.StatusUnauthorized, "Needed cookie not present")
			}

			return error
		}

		if cookie.Value == "some_string" {
			return next(context)
		}

		return context.String(http.StatusUnauthorized, "You don't have the right cookie")
	}
}
