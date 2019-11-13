package main

import (
	"github.com/labstack/echo"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	h := &handler{}
	e.POST("/login", h.login)
	e.Logger.Fatal(e.Start(":1323"))
}
