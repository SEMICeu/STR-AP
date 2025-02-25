package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping godoc
//
//	@Summary		ping
//	@Schemes		https
//	@Description	ping test to check application health
//	@Tags			health
//	@Produce		json
//	@Success		200	{object} Status
//	@Router			/ping [get]
func Ping(ctx *gin.Context) {
	statusOk := Status{Status: "ok"}
	ctx.JSON(http.StatusOK, statusOk)
}
