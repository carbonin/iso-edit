package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/filesystem"
)

func main() {
	inPath := flag.String("in", "isos/rhcos-4.6.1-x86_64-live.x86_64.iso", "input ISO path")
	outPath := flag.String("out", "isos/my-rhcos", "output ISO path")

	err := unpackISO(*inPath, *outPath)
	if err != nil {
		fmt.Printf("Failed to unpack ISO: %s\n", err)
		os.Exit(1)
	}
}

func unpackISO(isoPath, dir string) error {
	if err := os.Mkdir(dir, 0755); err != nil {
		return err
	}

	disk, err := diskfs.OpenWithMode(isoPath, diskfs.ReadOnly)
	fmt.Printf("Disk: %+v\n", disk)
	if err != nil {
		return err
	}

	fs, err := disk.GetFilesystem(0)
	if err != nil {
		return err
	}
	fmt.Printf("FileSystem: %+v\n\n", fs)

	files, err := fs.ReadDir("/")
	if err != nil {
		return err
	}

	return copyAll(fs, "/", files, dir)
}

// recursive function for unpacking all files and directores from the given iso filesystem starting at fsDir
func copyAll(fs filesystem.FileSystem, fsDir string, infos []os.FileInfo, targetDir string) error {
	for _, info := range infos {
		osName := filepath.Join(targetDir, info.Name())
		fsName := filepath.Join(fsDir, info.Name())

		if info.IsDir() {
			if err := os.Mkdir(osName, info.Mode().Perm()); err != nil {
				return err
			}

			files, err := fs.ReadDir(fsName)
			if err != nil {
				return err
			}
			if err := copyAll(fs, fsName, files, osName); err != nil {
				return err
			}
		} else {
			fmt.Printf("Opening file: %s\n", fsName)
			fsFile, err := fs.OpenFile(fsName, os.O_RDONLY)
			if err != nil {
				return err
			}

			osFile, err := os.Create(osName)
			if err != nil {
				return err
			}

			wrote, err := io.Copy(osFile, fsFile)
			if err != nil {
				osFile.Close()
				return err
			}

			fmt.Printf("wrote %d bytes to %s, size: %d\n\n", wrote, osName, info.Size())

			if err := osFile.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
