package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/4lxprime/gitdl"
)

const (
	replaceID           = "AABBCCDD"
	bootstrapFolderName = "bootstrap"
	reponame            = "https://github.com/anthdm/superkit.git"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println()
		fmt.Println("install requires your project name as the first argument")
		fmt.Println()
		fmt.Println("\tgo run superkit/install.go [your_project_name]")
		fmt.Println()
		os.Exit(1)
	}

	projectName := args[0]

	secret := generateSecret()

	fmt.Println("-- setting up bootstrap")
	fmt.Println("-- generating secure secret")
  if err := gitdl.DownloadGit(
		"anthdm/superkit",
		"bootstrap",
		projectName,
		gitdl.WithBranch("master"),
		gitdl.WithReplace(gitdl.Map{
			replaceID:        projectName,
			"{{app_secret}}": secret,
		}),
		//gitdl.WithLogs,
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println("-- renaming .env.local -> .env")
	if err := os.Rename(
		path.Join(projectName, ".env.local"),
		path.Join(projectName, ".env"),
	); err != nil {
		log.Fatal(err)
	}

	if err := exec.Command("cd", projectName).Run(); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("go", "clean", "-modcache").Run(); err != nil {
		log.Fatal(err)
	}
	if err := exec.Command("go", "get", "-u", "./...").Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("-- project (%s) successfully installed!\n", projectName)
}

func generateSecret() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(bytes)
}
