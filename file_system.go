package main

import (
	"errors"
	"log"
	"os"
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

func (f *FileSystem) Download(idx int, remote, local string) {
	files, ok := f.cache.get(f.currentPath)
	if !ok {
		log.Fatal("current dir not found")
	}

	if l := len(files); l <= 0 || l <= idx {
		log.Fatal("invalid index")
	}

	remoteDir, err := cleanPath(path.Dir(remote))
	if err != nil {
		log.Fatal(err)
	}

	if remote != path.Join(remoteDir, files[idx].Name) {
		log.Fatal("path error ", path.Join(remoteDir, files[idx].Name), remote)
	}

	dest, err := os.OpenFile(local, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, files[idx].Mode)
	if err != nil {
		log.Fatal(err)
	}

	defer closeIO(dest)
	if err := f.device.Pull(remote, dest); err != nil {
		log.Fatal(err)
	}
}

// updates cache
func (f *FileSystem) Upload(local, remote string) {
	remoteDir, err := cleanPath(remote)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(local)
	if err != nil {
		log.Fatal(err)
	}

	defer closeIO(file)
	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	f.cache.invalidate(remoteDir)
	if err := f.device.Push(file, remoteDir, stat.ModTime(), stat.Mode()); err != nil {
		log.Fatal(err)
	}
}

func (f *FileSystem) Delete(path string) {}

func (f *FileSystem) Rename(old, new string) {}

func (f *FileSystem) MakeDir(path string) {}

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
