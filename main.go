package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

func main() {
	// Define command-line flags
	provider := flag.String("provider", "github.com", "Git provider (e.g., github.com, gitlab.com)")
	gopath := flag.String("gopath", os.Getenv("GOPATH"), "GOPATH to use")
	username := flag.String("username", "", "Git username (defaults to global git config)")
	help := flag.Bool("h", false, "Display help")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Go Project Initializer\n\n")
		fmt.Fprintf(os.Stderr, "This program creates a new Go project directory structure, initializes a Git repository, sets up a Go module, and generates a basic main.go file.\n")
		fmt.Fprintf(os.Stderr, "It uses GOPATH and git global config, with options to override defaults.\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <project-name>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  project-name    Name of the project to create (required)\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Display help if -h flag is used
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Check for project name argument
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Error: Project name is required as an argument.")
		flag.Usage()
		os.Exit(1)
	}
	projectName := args[0]

	// Validate GOPATH
	if *gopath == "" {
		fmt.Println("Error: GOPATH is not set. Please set GOPATH environment variable or provide it using the -gopath flag.")
		os.Exit(1)
	}

	// Get username from git config if not provided
	if *username == "" {
		var err error
		*username, err = getGitUsername()
		if err != nil {
			fmt.Printf("Error getting git username: %v\n", err)
			os.Exit(1)
		}
	}

	// Create project path
	projectPath := filepath.Join(*gopath, "src", *provider, *username, projectName)

	// Create project directory
	err := os.MkdirAll(projectPath, 0755)
	if err != nil {
		fmt.Printf("Error creating project directory: %v\n", err)
		os.Exit(1)
	}

	// Change to project directory
	err = os.Chdir(projectPath)
	if err != nil {
		fmt.Printf("Error changing to project directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize Git repository
	_, err = git.PlainInit(projectPath, false)
	if err != nil {
		fmt.Printf("Error initializing Git repository: %v\n", err)
		os.Exit(1)
	}

	// Initialize Go module
	//modulePath := fmt.Sprintf("%s/%s/%s", *provider, *username, projectName)
	cmd := exec.Command("go", "mod", "init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error initializing Go module: %v\n%s\n", err, output)
		os.Exit(1)
	}

	// Generate main.go file
	mainContent := []byte(fmt.Sprintf(`package main

import "fmt"

func main() {
	fmt.Println("Hello from %s!")
}
`, projectName))

	err = os.WriteFile("main.go", mainContent, 0644)
	if err != nil {
		fmt.Printf("Error creating main.go file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created and set up Go project at %s\n", projectPath)
	fmt.Println("Created: Git repository, Go module, and main.go file")
}

func getGitUsername() (string, error) {
	cfg, err := config.LoadConfig(config.GlobalScope)
	if err != nil {
		return "", fmt.Errorf("failed to load git config: %w", err)
	}

	name := cfg.User.Name
	if name == "" {
		return "", fmt.Errorf("git user.name is not set in global config")
	}

	return name, nil
}
