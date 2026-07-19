package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	goadb "github.com/electricbubble/gadb"
)

type App struct {
	ctx    context.Context
	client goadb.Client
	FileSystem
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(_ context.Context) {
	if err := a.client.KillServer(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) NewADBClient(adbPath string, port int) []string {
	if err := startADBServer(adbPath, port); err != nil {
		log.Fatal(err)
	}

	client, err := goadb.NewClientWith("localhost", port)
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

	a.FileSystem.init(&devices[idx])
}

func (a *App) DownloadADB() string {
	tmp, err := downloadADB()
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmp)
	adbDir, err := getADBDir()
	if err != nil {
		log.Fatal(err)
	}

	if err := unzip(tmp, filepath.Dir(adbDir)); err != nil {
		log.Fatal(err)
	}

	adbPath, err := findADBExecutable(adbDir)
	if err != nil {
		log.Fatal(err)
	}

	return adbPath
}

func (a *App) KillServer(adbPath string, port int) {
	if err := killADBServer(adbPath, port); err != nil {
		log.Fatal(err)
	}
}
