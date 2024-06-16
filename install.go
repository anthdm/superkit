package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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

	// check if superkit folder already exists, if so, delete
	_, err := os.Stat("superkit")
	if !os.IsNotExist(err) {
		fmt.Println("-- deleting superkit folder cause its already present")
		if err := os.RemoveAll("superkit"); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("-- cloning", reponame)
	clone := exec.Command("git", "clone", reponame)
	if err := clone.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("-- renaming bootstrap ->", projectName)
	if err := os.Rename(path.Join("superkit", bootstrapFolderName), projectName); err != nil {
		log.Fatal(err)
	}

	err = filepath.Walk(path.Join(projectName), func(fullPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		b, err := os.ReadFile(fullPath)
		if err != nil {
			return err
		}

		contentStr := string(b)
		if strings.Contains(contentStr, replaceID) {
			replacedContent := strings.ReplaceAll(contentStr, replaceID, projectName)
			file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = file.WriteString(replacedContent)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-- renaming .env.local -> .env")
	if err := os.Rename(
		path.Join(projectName, ".env.local"),
		path.Join(projectName, ".env"),
	); err != nil {
		log.Fatal(err)
	}

	fmt.Println("-- generating secure secret")
	pathToDotEnv := path.Join(projectName, ".env")
	b, err := os.ReadFile(pathToDotEnv)
	if err != nil {
		log.Fatal(err)
	}
	secret := generateSecret()
	replacedContent := strings.Replace(string(b), "{{app_secret}}", secret, -1)
	file, err := os.OpenFile(pathToDotEnv, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.WriteString(replacedContent)
	if err != nil {
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
