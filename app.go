package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path"
	"strings"

	goadb "github.com/electricbubble/gadb"
)

type App struct {
	ctx         context.Context
	client      goadb.Client
	device      goadb.Device
	rootAliases [2][2]string
}

func NewApp() *App {
	return &App{rootAliases: [2][2]string{{"internal", ""}, {"external", ""}}}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) NewADBClient() []goadb.Device {
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

	return devices
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
	rootDir := a.rootAliases[0][1]
	rest, ok := strings.CutPrefix(dirpath, "/internal")
	if !ok {
		if a.rootAliases[1][1] == "" {
			return "", errors.New("no external storage")
		}

		rest, ok = strings.CutPrefix(dirpath, "/external")
		rootDir = a.rootAliases[1][1]

		if !ok {
			return "", errors.New("unknown root")
		}
	}

	if strings.ContainsAny(";&|", rest) {
		return "", errors.New("invalid characters")
	}

	cleanPath := path.Clean(path.Join(rootDir, rest))
	if !strings.HasSuffix(cleanPath, "/") {
		cleanPath += "/"
	}

	if !strings.HasPrefix(cleanPath, rootDir) {
		return "", fmt.Errorf("path escapes: %s", rootDir)
	}

	return cleanPath, nil
}
