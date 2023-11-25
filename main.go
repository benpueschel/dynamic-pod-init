package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/plus3it/gorecurcopy"
	"golift.io/xtractr"
)


func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func downloadFile(url string, target string) error {
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetFlags(0)

	targetDir := getenv("TARGET_DIR", "/")
	templateDir := os.Getenv("TEMPLATE_DIR")
	patches := os.Getenv("PATCHES")

	if len(templateDir) != 0 {
		fmt.Printf("copying files from %s to %s\n", templateDir, targetDir)
		err := gorecurcopy.CopyDirectory(templateDir, targetDir)
		if err != nil {
			fmt.Printf("fatal: %s\n", err)
			os.Exit(1)
		}
	}

	if len(patches) != 0 {
		log.Println("applying patches")
		patches := strings.Split(patches, ",")
		for i, patch := range patches {
			fmt.Printf("applying patch %s\n", patch)

			archivePath := strconv.Itoa(i) + ".patch"
			err := downloadFile(patch, archivePath)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				continue
			}

			x := &xtractr.XFile {
				FilePath:	archivePath,
				OutputDir: 	targetDir,
			}
			size, files, archives, err := xtractr.ExtractFile(x)
			if err != nil || files == nil {
				fmt.Printf("error: %d, %s, %s, %s\n", size, files, archives, err)
				continue
			}
		}
	}

}
