package main

import (
	"errors"
	"log"
	"path"
	"strings"

	goadb "github.com/electricbubble/gadb"
)

type FileSystem struct {
	device *goadb.Device
}

func (f *FileSystem) List(path string) []goadb.DeviceFileInfo {
	dirpath, err := cleanPath(path)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := f.device.List(dirpath)
	if err != nil {
		log.Fatal(err)
	}

	return entries
}

func (f *FileSystem) Download(remote, local string) {}

func (f *FileSystem) Upload(local, remote string) {}

func (f *FileSystem) Delete(path string) {}

func (f *FileSystem) Move(old, new string) {}

func (f *FileSystem) MakeDir(path string) {}

func (f *FileSystem) init(device *goadb.Device) {
	f.device = device
}

func cleanPath(dirpath string) (string, error) {
	const STORAGE_DIR = "/storage/"

	if strings.ContainsAny(";&|", dirpath) {
		return "", errors.New("invalid characters")
	}

	cleanPath := path.Clean(dirpath)
	if !strings.HasSuffix(cleanPath, "/") {
		cleanPath += "/"
	}

	if !strings.HasPrefix(cleanPath, STORAGE_DIR) {
		return "", errors.New("path escapes " + STORAGE_DIR)
	}

	return cleanPath, nil
}
