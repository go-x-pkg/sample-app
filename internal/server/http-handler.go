package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *App) HTTPHandler() http.Handler {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(
		gin.Recovery(),
	)

	r.HandleMethodNotAllowed = true

	return r
}
