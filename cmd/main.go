package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	cpio "github.com/cavaliercoder/go-cpio"
	isoutil "github.com/kdomanski/iso9660/util"
)

func main() {
	inPath := flag.String("in", "isos/rhcos-46.82.202009222340-0-live.x86_64.iso", "input ISO path")
	outPath := flag.String("out", "isos/my-rhcos.iso", "output ISO path")
	filesPath := flag.String("files", "files", "directory to add to the iso")

	err := patchISO(*inPath, *filesPath, *outPath)
	if err != nil {
		fmt.Printf("Failed to create patched iso: %s\n", err)
		os.Exit(1)
	}
}

func patchISO(inPath, filesPath, outPath string) error {
	dir, err := unpackISO(inPath)
	if err != nil {
		return err
	}

	if err := addFiles(filesPath, dir); err != nil {
		return err
	}

	return packISO(dir, outPath)
}

// takes an iso path and returns a writable directory containing the iso contents
func unpackISO(isoPath string) (string, error) {
	dir := filepath.Join(os.TempDir(), "iso-test")

	f, err := os.Open(isoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %s", err)
	}
	defer f.Close()

	if err = isoutil.ExtractImageToDirectory(f, dir); err != nil {
		return "", fmt.Errorf("failed to extract image: %s", err)
	}

	return dir, nil
}

// adds all the files at filesPath to the unpacked iso at isoPath as an additional image
func addFiles(filesPath, isoPath string) error {
	f, err := os.Create(filepath.Join(isoPath, "IMAGES/my_image"))
	if err != nil {
		return fmt.Errorf("failed to open image file: %s", err)
	}

	w := cpio.NewWriter(f)
	// find and read files
	err = filepath.Walk(filesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			hdr := &cpio.Header{
				Name: path,
				Mode: 0775,
				Size: 0,
			}
			if err := w.WriteHeader(hdr); err != nil {
				return err
			}
		}

		hdr := &cpio.Header{
			Name: path,
			Mode: 0664,
			Size: info.Size(),
		}
		if err := w.WriteHeader(hdr); err != nil {
			return err
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if _, err := w.Write(content); err != nil {
			return err
		}

		return nil
	})
	if err := w.Close(); err != nil {
		return err
	}
	// edit config
	return nil
}

// creates a new iso out of the directory structure at isoDir and writes it to outPath
func packISO(isoDir, outPath string) error {
	return nil
}
