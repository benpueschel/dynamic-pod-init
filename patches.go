package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/mholt/archiver/v4"
)

func ApplyPatches(patchUrls string, targetDir string) {
	patches := strings.Split(patchUrls, ",")

	for _, patch := range patches {
		fmt.Printf("applying patch %s\n", patch)

		// get archive from url
		response, err := http.Get(patch)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		defer response.Body.Close()

		// identify the archive format
		format, archive, err := archiver.Identify("", response.Body)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}

		ctx := context.WithValue(context.Background(), "targetDir", targetDir)

		// decompress if needed
		if decom, ok := format.(archiver.Decompressor); ok {
			rc, err := decom.OpenReader(archive)
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return
			}
			defer rc.Close()
			archive = rc // use decompressed reader instead
		}

		err = extractArchive(ctx, archive)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
	}
}

func extractArchive(ctx context.Context, reader io.Reader) error {
	format, archive, err := archiver.Identify("", reader)
	if err != nil {
		return err
	}

	if ex, ok := format.(archiver.Extractor); ok {
		bReader, err := byteReader(archive) // we need to convert the reader to a byte reader because zip is weird, see
											// https://pkg.go.dev/github.com/mholt/archiver/v4#Zip.Extract for more info
		if err != nil {
			return err
		}
		err = ex.Extract(ctx, bReader, nil, extractFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func extractFile(ctx context.Context, f archiver.File) error {
	archiveFile, err := f.Open()
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	// create parent dir if not exists
	targetDir := ctx.Value("targetDir").(string)
	err = os.MkdirAll(path.Dir(path.Join(targetDir, f.NameInArchive)), fs.ModePerm)
	if err != nil {
		return err
	}

	// create the file
	file, err := os.Create(path.Join(targetDir, f.NameInArchive))
	if err != nil {
		return err
	}
	defer file.Close()

	// copy archive file contents to file
	_, err = io.Copy(file, archiveFile)
	if err != nil {
		return err
	}

	fmt.Printf("extracted file %s\n", file.Name())
	return nil
}

func byteReader(reader io.Reader) (*bytes.Reader, error) {
	buff := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buff, reader)
	if err != nil {
		return bytes.NewReader(nil), err
	}
	return bytes.NewReader(buff.Bytes()), nil
}
