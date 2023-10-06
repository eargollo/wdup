package main

import (
	"changeme/pkg/dupfile"
	"context"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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
