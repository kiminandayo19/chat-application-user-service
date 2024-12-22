package middlewares

import (
	"chat-application/users/config"
	"chat-application/users/internal/domain"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HeadersCheck struct {
	Name          string
	RequiredValue string
	StatusCode    int
	Message       string
}

type SecurityHeaders struct {
	XContentTypeOptions     string
	XXSSProtection          string
	StrictTransportSecurity string
	XFrameOptions           string
	APIKey                  string
}

const (
	DefaultXContentTypeOptions     = "nosniff"
	DefaultXXSSProtection          = "1; mode=block"
	DefaultStrictTransportSecurity = "max-age=31536000; includeSubDomains; preload"
	DefaultXFrameOptions           = "SAMEORIGIN"
)

func NewHeaderMiddleware() *SecurityHeaders {
	env, err := config.NewEnvDev()

	if err != nil {
		log.Print("[middleware] - Failed to load env")
	}

	return &SecurityHeaders{
		XContentTypeOptions:     DefaultXContentTypeOptions,
		XXSSProtection:          DefaultXXSSProtection,
		StrictTransportSecurity: DefaultStrictTransportSecurity,
		XFrameOptions:           DefaultXFrameOptions,
		APIKey:                  env.APIKEY,
	}
}

func HeaderMiddleware(config *SecurityHeaders) gin.HandlerFunc {
	if config == nil {
		config = NewHeaderMiddleware()
	}

	headersCheck := []HeadersCheck{
		{
			Name:          "X-Content-Type-Options",
			RequiredValue: config.XContentTypeOptions,
			StatusCode:    http.StatusBadRequest,
			Message:       "Bad request. Missing or invalid x-content-type-options header",
		},
		{
			Name:          "X-XSS-Protection",
			RequiredValue: config.XXSSProtection,
			StatusCode:    http.StatusBadRequest,
			Message:       "Bad request. Missing or invalid x-xss-protection header",
		},
		{
			Name:          "Strict-Transport-Security",
			RequiredValue: config.StrictTransportSecurity,
			StatusCode:    http.StatusBadRequest,
			Message:       "Bad request. Missing or invalid strict-transport-security header",
		},
		{
			Name:          "X-Frame-Options",
			RequiredValue: config.XFrameOptions,
			StatusCode:    http.StatusBadRequest,
			Message:       "Bad request. Missing or invalid x-frame-options headers",
		},
		{
			Name:          "APIKey",
			RequiredValue: config.APIKey,
			StatusCode:    http.StatusUnauthorized,
			Message:       "Unautorized. Missing or invalid api key",
		},
	}

	return func(ctx *gin.Context) {
		setHeaders(ctx, config)

		for _, header := range headersCheck {
			rHeader := ctx.GetHeader(header.Name)
			if rHeader == "" || rHeader != header.RequiredValue {
				response := domain.NewResponse(
					header.StatusCode,
					false,
					header.Message,
					nil,
				)
				ctx.AbortWithStatusJSON(header.StatusCode, response)
				return
			}
		}
		ctx.Next()
	}
}

func setHeaders(ctx *gin.Context, config *SecurityHeaders) {
	ctx.Header("X-Content-Type-Options", config.XContentTypeOptions)
	ctx.Header("X-XSS-Protection", config.XXSSProtection)
	ctx.Header("Strict-Transport-Security", config.StrictTransportSecurity)
	ctx.Header("X-Frame-Options", config.XFrameOptions)
	ctx.Header("Content-Type", "application/json")
}
