package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log/slog"
	"net/http"
	"os"
	"str/docs"
	"str/internal"
	"str/pkg/handler"
	"str/pkg/routes"
)

//	@title			EU STR - Single Digital Entry Point
//	@description	A gateway for the electronic transmission of data between online short-term rental platforms and competent authorities, ensuring timely, reliable and efficient data sharing processes
//	@description	Priority of development is: 1. listings, 2. orders, 3. activity, 4. area
//	@description	<br>
//	@description	To obtain API credentials, please contact: wouter.travers@pwc.com via e-mail
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Wouter Travers
//	@contact.email	wouter.travers@pwc.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api/v0

//	@securitydefinitions.oauth2.accessCode	OAuth2AccessCode
//	@tokenUrl								https://tt-dp-dev.eu.auth0.com/oauth/token
//	@authorizationUrl						https://tt-dp-dev.eu.auth0.com/authorize?audience=https://str.eu
//	@scope.str								Grants read access

// @externalDocs.description	STR Application Profile (STR-AP)
// @externalDocs.url			https://semiceu.github.io/STR-AP/releases/1.0.1/

func main() {
	internal.Config()
	// programmatically set swagger info
	docs.SwaggerInfo.Version = handler.Version
	docs.SwaggerInfo.Host = viper.GetString("HOST")

	logger := internal.CallNewLogger()

	//internal.Infof("Short Term Rental Application Profile")
	logger.Info("Short Term Rental Application Profile", slog.String("version", handler.Version))

	ginMode := viper.GetString("GIN_MODE")
	logger.Info("GIN Mode", slog.String("mode", ginMode))
	//internal.Infof("Gin mode: %v", ginMode)
	gin.SetMode(ginMode)

	router := gin.Default()
	routes.SetupRoutes(router)

	addr := ":" + viper.GetString("INTERNAL_PORT")
	logger.Info("Webserver binding", slog.String("addr", addr))
	//internal.Infof("Webserver binding: %s", addr)

	err := http.ListenAndServeTLS(addr, viper.GetString("CERT_FILE"), viper.GetString("KEY_FILE"), router)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	//internal.Fatalf("Unable to start HTTP server: %s", err)
}
