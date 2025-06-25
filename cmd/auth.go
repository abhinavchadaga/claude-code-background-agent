package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/google/go-github/v57/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/term"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Store GitHub token locally (optional)",
	Long: `Validates and provides guidance for setting up GitHub authentication.
The CLI primarily uses the GITHUB_TOKEN environment variable.`,
	RunE: runAuthLogin,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
}

func runAuthLogin(cmd *cobra.Command, args []string) error {
	fmt.Println("GitHub Authentication Setup")
	fmt.Println("============================")

	existingToken := os.Getenv("GITHUB_TOKEN")
	if existingToken != "" {
		fmt.Println("✓ GITHUB_TOKEN environment variable is already set")

		if err := validateToken(existingToken); err != nil {
			fmt.Printf("✗ Token validation failed: %v\n", err)
			fmt.Println("\nPlease check your token has the 'repo' scope and is not expired.")
			return nil
		}

		fmt.Println("✓ Token is valid")
		return nil
	}

	fmt.Println("No GITHUB_TOKEN found in environment.")
	fmt.Println("\nTo set up authentication:")
	fmt.Println("1. Go to https://github.com/settings/tokens/new")
	fmt.Println("2. Create a token with 'repo' scope")
	fmt.Println("3. Export it: export GITHUB_TOKEN=your_token_here")
	fmt.Println("\nAlternatively, enter your token now to validate it:")

	fmt.Print("GitHub token (optional): ")

	tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("failed to read token: %w", err)
	}

	token := strings.TrimSpace(string(tokenBytes))
	fmt.Println()

	if token == "" {
		fmt.Println("No token provided. Set GITHUB_TOKEN environment variable when ready.")
		return nil
	}

	if err := validateToken(token); err != nil {
		fmt.Printf("✗ Token validation failed: %v\n", err)
		return nil
	}

	fmt.Println("✓ Token is valid!")
	fmt.Printf("To use this token, run: export GITHUB_TOKEN=%s\n", token[:8]+"...")

	return nil
}

func validateToken(token string) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	fmt.Printf("Authenticated as: %s\n", user.GetLogin())
	return nil
}
