package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	tokenURL := os.Getenv("TOKEN_URL")

	// cc := clientcredentials.Config{
	// 	ClientID:     "xangkham.apps.laoit.dev",
	// 	ClientSecret: "1rMllP5bJr9wG0lh3aF2v1KLasdDND8LKYFtzCmayWCAj86flMCANAOnKOJAi9qoiIivulX5yH8sRkQmpQ3vBjq1tVYtFTM6BJi31ZEPol6gzuci84pdJMAyTMN87aK2",
	// 	TokenURL:     "https://accounts.laoit.dev/application/o/token/",
	// }

	cc := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}

	client := cc.Client(context.Background())
	res, err := client.Get("http://localhost:8080/events")
	if err != nil {
		slog.Error("Failed to make request", "error", err)
	}
	fmt.Printf("res: %v\n", res)
	fmt.Printf("status: %v", res.StatusCode)
}
