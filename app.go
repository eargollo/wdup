package main

import (
	"changeme/pkg/dupfile"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct1
type App struct {
	ctx context.Context
	log logger.Logger
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{}
	app.log = logger.NewDefaultLogger()
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) PathSelect(name string) string {
	selection, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Cabrum",
	})
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	return selection
}

// Greet returns a greeting for the given name
func (a *App) DuplicateSearch(paths []string) [][]*dupfile.File {
	a.log.Info(fmt.Sprintf("paths: %v", paths))

	cachePath, err := a.cachePath()
	if err != nil {
		a.log.Error(err.Error())
		return [][]*dupfile.File{}
	}

	df, err := dupfile.New(dupfile.WithPaths(paths), dupfile.WithCache(cachePath))
	if err != nil {
		a.log.Error(err.Error())
		return [][]*dupfile.File{}
	}

	return df.Run()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) [][]*dupfile.File {
	return [][]*dupfile.File{
		{
			&dupfile.File{
				Path: "first/first",
				Name: "first",
				Size: 123,
				Hash: []byte("hash"),
			},
			&dupfile.File{
				Path: "first/second",
				Name: "second",
				Size: 123,
				Hash: []byte("hash"),
			},
		},
		{
			&dupfile.File{
				Path: "second/first",
				Name: "first",
				Size: 1234,
				Hash: []byte("hash"),
			},
			&dupfile.File{
				Path: "second/second",
				Name: "second",
				Size: 1234,
				Hash: []byte("hash"),
			},
			&dupfile.File{
				Path: "second/third",
				Name: "third",
				Size: 1234,
				Hash: []byte("hash"),
			},
			&dupfile.File{
				Path: "second/fourth",
				Name: "fourth",
				Size: 1234,
				Hash: []byte("hash"),
			},
		},
	}
}

// Returns the path for the cached file hashes
func (a *App) cachePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return homedir, err
	}

	cacheDir := homedir + "/.wdup"
	return cacheDir, nil
}
