package main

import (
  "log"
  "net/http"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  middlewarez "jwt-claims-extractor/middleware"
)

func main() {
  // Echo instance
  e := echo.New()

  // Middleware
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  // Usage example
  e.Use(middlewarez.JWTExtractor(middlewarez.JWTExtractorConfig{ DataFields: []string{"name", "email"} }))

  // Routes
  e.GET("/", hello)

  // Start server
  e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
  log.Println(c.Request().Header.Get("X-Consumer-Token-Name"))
  log.Println(c.Request().Header.Get("X-Consumer-Token-Email"))
  return c.String(http.StatusOK, "Hello, World!")
}
