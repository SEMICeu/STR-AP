package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "str/docs"
	"str/pkg/handler"
	"str/pkg/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// v0 routes
	v0 := r.Group("/api/v0")
	// str routes
	addStrRoutes(v0)
	// ca routes
	addCaRoutes(v0)

	// status routes
	addStatusRoutes(v0)

}

// @tag.name			str
// @tag.description	Endpoints for use by the Short Term Rental Platforms
func addStrRoutes(rg *gin.RouterGroup) {
	str := rg.Group("/str")
	str.POST("/activity-data", middleware.CheckJWT(), handler.ActivityDataPush)
	str.POST("/listings", middleware.CheckJWT(), handler.ListingsPush)

	str.GET("/area", middleware.CheckJWT(), handler.GetShapefiles)
	str.GET("/area/:id", middleware.CheckJWT(), handler.DownloadShapefile)
	
}

// @tag.name			ca
// @tag.description	Endpoints for use by the Competent Authorities
func addCaRoutes(rg *gin.RouterGroup) {
	ca := rg.Group("/ca")
	ca.GET("/activity-data", middleware.CheckJWT(), handler.ActivityDataPull)
	ca.GET("/listings", middleware.CheckJWT(), handler.ListingsPull)

	ca.POST("/area", middleware.CheckJWT(), handler.Areaupload)

	ca.DELETE("/area/:luid", middleware.CheckJWT(), handler.AreaDelete)
}

func addStatusRoutes(rg *gin.RouterGroup) {
	rg.GET("/ping", handler.Ping)
}
