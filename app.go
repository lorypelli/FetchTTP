package main

import (
	"context"
	"net/http"
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

func (a *App) MakeRequest(method string, url string) string {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "FetchTTP")
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return ""
	}
	return res.Status
}
