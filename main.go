package main

import (
	"fmt"
	"log"
	"os"
)


func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	log.SetFlags(0)

	targetDir := getenv("TARGET_DIR", "/")
	templateDir := os.Getenv("TEMPLATE_DIR")
	patches := os.Getenv("PATCHES")

	if len(templateDir) != 0 {
		fmt.Printf("copying files from %s to %s\n", templateDir, targetDir)
		err := CopyDir(templateDir, targetDir, true)
		if err != nil {
			fmt.Printf("fatal: %s\n", err)
			os.Exit(1)
		}
	}

	if len(patches) != 0 {
		fmt.Println("applying patches")
		ApplyPatches(patches, targetDir)
	}

}
