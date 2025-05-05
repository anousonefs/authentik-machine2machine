package main

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func main() {
	jwksURL := os.Getenv("JWKS_URL")
	e := echo.New()

	// jwksURL := "https://accounts.laoit.dev/application/o/xangkham/jwks/"

	// Load KeyFunc from remote JWKS endpoint
	// Create the keyfunc.Keyfunc.
	k, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL}) // Context is used to end the refresh goroutine.
	if err != nil {
		log.Fatalf("Failed to create a keyfunc.Keyfunc from the server's URL.\nError: %s", err)
	}

	if err != nil {
		slog.Error("Failed to get JWKS", "error", err)
		os.Exit(1)
	}

	// JWT middleware with KeyFunc (not static SigningKey!)
	e.Use(echojwt.WithConfig(echojwt.Config{
		KeyFunc: k.Keyfunc,
	}))

	// Protected → requires Bearer token
	e.GET("/events", func(c echo.Context) error {
		ctx := c.Request().Context()

		// Read request body safely
		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to read body", "error", err)
			return c.String(http.StatusInternalServerError, "Internal Error")
		}

		// Log JWT claims (optional)
		user := c.Get("user").(*jwt.Token)
		claims, _ := user.Claims.(jwt.MapClaims)

		slog.InfoContext(ctx, "Received JWT-protected request", "claims", claims, "body", string(body))
		slog.Info("header", c.Request().Header)
		return c.String(http.StatusOK, "Protected resource accessed.")
	})

	slog.Info("Server started on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
