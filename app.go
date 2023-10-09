package main

import (
	"changeme/pkg/dupfile"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

func (a *App) PathSelect() string {
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
	// Further event to update the UI on search progress
	// go func() {
	// 	counter := 0
	// 	for {
	// 		log.Printf("Emitting event %d", counter)
	// 		runtime.EventsEmit(a.ctx, "abacaxi", fmt.Sprintf("Hello from Go! Event %d", counter))

	// 		time.Sleep(time.Second)
	// 		counter++
	// 	}
	// }()

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

// Emulate returns a fake result
func (a *App) Emulate() [][]*dupfile.File {
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

func (a *App) OpenPath(file string) string {
	path := filepath.Dir(file)
	cmd := fmt.Sprintf("open '%s'", path)
	a.log.Info(cmd)
	err := exec.Command("/bin/bash", "-c", cmd).Start()
	if err != nil {
		a.log.Error(err.Error())
		return err.Error()
	}
	return ""
}

func (a *App) OpenFile(file string) string {
	cmd := fmt.Sprintf("open '%s'", file)
	a.log.Info(cmd)
	err := exec.Command("/bin/bash", "-c", cmd).Start()
	if err != nil {
		a.log.Error(err.Error())
		return err.Error()
	}
	return ""
}

type BoolReturn struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (a *App) DeleteFile(file string) BoolReturn {
	a.log.Info(fmt.Sprintf("Deleting file '%s'", file))
	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Comfirm Delete",
		Message:       fmt.Sprintf("Are you sure you want to delete '%s'?", file),
		Buttons:       []string{"Yes", "No"},
		DefaultButton: "No",
		CancelButton:  "No",
	})
	if err != nil {
		a.log.Error(err.Error())
		return BoolReturn{Success: false, Error: err.Error()}
	}

	if result == "Yes" {
		err := os.Remove(file)
		// err := fmt.Errorf("Not implemented")
		if err != nil {
			a.log.Error(err.Error())
			return BoolReturn{Success: false, Error: err.Error()}
		}
		return BoolReturn{Success: true, Error: ""}
	} else {
		return BoolReturn{Success: false, Error: ""}
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
