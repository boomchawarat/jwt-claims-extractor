package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JWTExtractorConfig struct {
	DataFields []string
}

func JWTExtractor(config JWTExtractorConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			val := c.Request().Header.Get("Authorization")
			if val == "" {
				return echo.NewHTTPError(http.StatusBadRequest)
			}

			splitToken := strings.Split(val, "Bearer ")
			if len(splitToken) != 2 {
				return echo.NewHTTPError(http.StatusBadRequest)
			}

			// Extract claims info
			token, _, err := new(jwt.Parser).ParseUnverified(splitToken[1], jwt.MapClaims{})
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				for _, key := range config.DataFields {
					if name, found := claims[key]; found {
						if key != "role" {
							c.Request().Header.Add("X-Consumer-Token-"+strings.Title(key), name.(string))
						} else {
							var roleString = "["
							var r string
							chkRoles := claims[key].([]interface{})
							for _, role := range chkRoles {
								r += ","
								r += "\""
								r += role.(string)
								r += "\""
							}

							r = r[1:len(r)]
							
							roleString += r
							roleString += "]"
							c.Request().Header.Add("X-Consumer-Token-"+strings.Title(key), roleString)
						}
					}
				}
			}
			// Proceed
			return next(c)
		}
	}
}
