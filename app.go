package main

import (
	"context"
	"errors"
	"log"
	"path"
	"strings"

	goadb "github.com/electricbubble/gadb"
)

type App struct {
	ctx    context.Context
	client goadb.Client
	device goadb.Device
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) NewADBClient() []string {
	if err := startADBServer(); err != nil {
		log.Fatal(err)
	}

	client, err := goadb.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	a.client = client
	devices, err := client.DeviceList()
	if err != nil {
		log.Fatal(err)
	}

	labels := make([]string, 0, len(devices))
	for _, device := range devices {
		info := device.DeviceInfo()
		label := info["device"] + ":" + info["model"]
		labels = append(labels, label)
	}

	return labels
}

func (a *App) SelectDevice(idx int) {
	devices, err := a.client.DeviceList()
	if err != nil {
		log.Fatal(err)
	}

	if l := len(devices); l <= 0 || l <= idx {
		log.Fatal("Invalid device index")
	}

	a.device = devices[idx]
}

func (a *App) ListDirectory(path string) []goadb.DeviceFileInfo {
	dirpath, err := a.getCleanDirPath(path)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := a.device.List(dirpath)
	if err != nil {
		log.Fatal(err)
	}

	return entries
}

func (a *App) getCleanDirPath(dirpath string) (string, error) {
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
