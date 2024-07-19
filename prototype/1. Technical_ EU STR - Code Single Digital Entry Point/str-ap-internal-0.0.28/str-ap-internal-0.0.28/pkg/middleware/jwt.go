package middleware

import (
	"context"
	"errors"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/url"
	"str/internal"
	"time"
)

var (
	// We want this struct to be filled in with
	// our custom claims from the token.
	customClaims = func() validator.CustomClaims {
		return &CustomClaims{}
	}
)

// CheckJWT is a gin.HandlerFunc middleware
// that will check the validity of our JWT.
func CheckJWT() gin.HandlerFunc {
	issuer := viper.GetString("AUTH0_DOMAIN")
	audience := []string{viper.GetString("AUTH0_AUDIENCE")}

	issuerURL, err := url.Parse("https://" + issuer + "/")
	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	// Set up the validator.
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		audience,
		validator.WithCustomClaims(customClaims),
		validator.WithAllowedClockSkew(5*time.Minute),
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(ctx *gin.Context) {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.Request = r
			ctx.Next()
		}

		middleware.CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)

		if encounteredError {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				map[string]string{"message": "JWT is invalid."},
			)
		}
	}
}

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Name         string `json:"https://str-eu/name"`
	CA           string `json:"https://str-eu/ca-id"`
	ShouldReject bool   `json:"shouldReject,omitempty"`
}

// Validate errors out if `ShouldReject` is true.
func (c *CustomClaims) Validate(ctx context.Context) error {
	if c.ShouldReject {
		return errors.New("should reject was set to true")
	}
	return nil
}

func GetCustomClaims(ctx *gin.Context) *CustomClaims {
	// Get all claims
	claims, ok := ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{"message": "Failed to get validated JWT claims."},
		)
	}

	customClaims, ok := claims.CustomClaims.(*CustomClaims)
	if !ok {
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			map[string]string{"message": "Failed to cast custom JWT claims to specific type."},
		)
	}

	internal.Infof("Custom claim Name: %v", customClaims.Name)
	internal.Infof("Custom claim CA: %v", customClaims.CA)

	return customClaims
}
