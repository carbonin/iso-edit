package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	cpio "github.com/cavaliercoder/go-cpio"
	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
)

func main() {
	inPath := flag.String("in", "isos/rhcos-4.6.1-x86_64-live.x86_64.iso", "input ISO path")
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

	info, err := os.Stat(inPath)
	if err != nil {
		return err
	}

	return packISO(dir, outPath, info.Size())
}

// takes an iso path and returns a writable directory containing the iso contents
func unpackISO(isoPath string) (string, error) {
	dir := filepath.Join(os.TempDir(), "iso-test")
	if err := os.Mkdir(dir, 0755); err != nil {
		return "", err
	}

	disk, err := diskfs.OpenWithMode(isoPath, diskfs.ReadOnly)
	if err != nil {
		return "", err
	}

	fs, err := disk.GetFilesystem(0)
	if err != nil {
		return "", err
	}

	files, err := fs.ReadDir("/")
	if err != nil {
		return "", err
	}
	err = copyAll(fs, "/", files, dir)
	if err != nil {
		return "", err
	}

	return dir, nil
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

			fmt.Printf("wrote %d bytes to %s\n", wrote, osName)

			if err := osFile.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// adds all the files at filesPath to the unpacked iso at isoPath as an additional image
func addFiles(filesPath, isoPath string) error {
	f, err := os.Create(filepath.Join(isoPath, "images/my_image.img"))
	if err != nil {
		return fmt.Errorf("failed to open image file: %s", err)
	}

	w := cpio.NewWriter(f)
	addFileToArchive := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			hdr := &cpio.Header{
				Name: path,
				Mode: 040775,
				Size: 0,
			}
			if err := w.WriteHeader(hdr); err != nil {
				return err
			}
		} else {
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
		}

		return nil
	}

	if err := filepath.Walk(filesPath, addFileToArchive); err != nil {
		w.Close()
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	// edit config to add new image to initrd
	err = editFile(filepath.Join(isoPath, "EFI/redhat/grub.cfg"), `(?m)^(\s+initrd) (.+| )+$`, "$1 $2 /images/my_image.img")
	if err != nil {
		return err
	}
	return editFile(filepath.Join(isoPath, "isolinux/isolinux.cfg"), `(?m)^(\s+append.*initrd=\S+) (.*)$`, "${1},/images/my_image.img ${2}")
}

func editFile(fileName string, reString string, replacement string) error {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(reString)
	newContent := re.ReplaceAllString(string(content), replacement)

	if err := ioutil.WriteFile(fileName, []byte(newContent), 0644); err != nil {
		return err
	}

	return nil
}

// creates a new iso out of the directory structure at isoDir and writes it to outPath
func packISO(isoDir string, outPath string, size int64) error {
	d, err := diskfs.Create(outPath, size, diskfs.Raw)
	if err != nil {
		return err
	}

	d.LogicalBlocksize = 2048
	fspec := disk.FilesystemSpec{Partition: 0, FSType: filesystem.TypeISO9660, VolumeLabel: "rhcos-46.82.202010091720-0"}
	fs, err := d.CreateFilesystem(fspec)
	if err != nil {
		return err
	}

	addFileToISO := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		p, err := filepath.Rel(isoDir, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return fs.Mkdir(p)
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		rw, err := fs.OpenFile(p, os.O_CREATE|os.O_RDWR)
		if err != nil {
			return err
		}

		_, err = rw.Write(content)
		return err
	}
	if err := filepath.Walk(isoDir, addFileToISO); err != nil {
		return err
	}

	iso, ok := fs.(*iso9660.FileSystem)
	if !ok {
		fmt.Errorf("not an iso9660 filesystem")
	}

	options := iso9660.FinalizeOptions{
		VolumeIdentifier: "rhcos-46.82.202010091720-0",
		RockRidge:        true,
		ElTorito: &iso9660.ElTorito{
			BootCatalog: "isolinux/boot.cat",
			Entries: []*iso9660.ElToritoEntry{
				&iso9660.ElToritoEntry{
					Platform:  iso9660.BIOS,
					Emulation: iso9660.NoEmulation,
					BootFile:  "isolinux/isolinux.bin",
					BootTable: true,
					LoadSize:  4,
				},
				&iso9660.ElToritoEntry{
					Platform:  iso9660.EFI,
					Emulation: iso9660.NoEmulation,
					BootFile:  "images/efiboot.img",
				},
			},
		},
	}
	err = iso.Finalize(options)
	if err != nil {
		return err
	}

	return nil
}
