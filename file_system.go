package main

import (
	"errors"
	"log"
	"path"
	"strings"

	goadb "github.com/electricbubble/gadb"
)

type FileSystem struct {
	device      *goadb.Device
	currentPath string
	cache       *DirCache
}

func (f *FileSystem) List(path string) []goadb.DeviceFileInfo {
	dirpath, err := cleanPath(path)
	if err != nil {
		log.Fatal(err)
	}

	if entries, ok := f.cache.get(dirpath); ok {
		f.currentPath = dirpath
		return entries
	}

	entries, err := f.device.List(dirpath)
	if err != nil {
		log.Fatal(err)
	}

	f.cache.set(dirpath, entries)
	f.currentPath = dirpath
	return entries
}

func (f *FileSystem) Download(remote, local string) {}

// updates cache
func (f *FileSystem) Upload(local, remote string) {}

func (f *FileSystem) Delete(path string) {}

func (f *FileSystem) Rename(old, new string) {}

func (f *FileSystem) MakeDir(path string) {}

func (f *FileSystem) getCachedFile(device *goadb.Device) {

}

func (f *FileSystem) init(device *goadb.Device) {
	f.device = device
	f.cache = newDirCache(5)
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
