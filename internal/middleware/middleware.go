package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// middleware который проверяет допустимые query параметры запроса(допустимые query параметры см. AllowedParams), если имеются недопустимые параметры выдает 400
func StrictQueryParamsMiddleware(allowedParams map[string]struct{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParams := c.Request.URL.Query()
		for param := range queryParams {
			if _, ok := allowedParams[param]; !ok {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("unknown request parameter: %s", param),
				})
				return
			}
		}
		c.Next()
	}
}

// все допустимые query параметры в запросе
var AllowedParams = map[string]struct{}{
	"name":        {},
	"surname":     {},
	"age":         {},
	"gender":      {},
	"nationality": {},
	"page":        {},
	"limit":       {},
}
