package main

import (
	"context"
	"fmt"
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

	ctx := context.Background()
	k, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		log.Fatalf("Failed to create a keyfunc.Keyfunc from the server's URL.\nError: %s", err)
	}
	defer func() {
		println("=> close ctx")
		ctx.Done()
	}()

	e.Use(Auth(k))

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
		fmt.Printf("=> client id: %v\n", claims["aud"])

		return c.String(http.StatusOK, "Protected resource accessed.")
	})

	slog.Info("Server started on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}

func Auth(k keyfunc.Keyfunc) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		KeyFunc: k.Keyfunc,
		ErrorHandler: func(c echo.Context, err error) error {
			fmt.Printf("=> error: %v\n", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		},
	})
}
