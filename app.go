package main

import (
	"context"
	"io"
	j "encoding/json"
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

type Header struct {
	disabled bool
	name string
	value string
}

type Query struct {
	disabled bool
	name string
	value string
}

func (a *App) MakeRequest(method string, url string, headers []Header, query []Query, body string) (string, http.Header, []byte) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", http.Header{}, []byte{}
	}
	req.Header.Set("User-Agent", "FetchTTP")
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return "", http.Header{}, []byte{}
	}
	var resBody []byte
	if res.Header.Get("Content-Type") == "application/json" {
		var jsonBody map[string]interface{}
		bytes, _ := io.ReadAll(res.Body)
		j.Unmarshal(bytes, &jsonBody)
		resBody, _ = j.MarshalIndent(jsonBody, "", "  ")
	} else {
		bytes, _ := io.ReadAll(res.Body)
		resBody = bytes
	}
	return res.Status, res.Header, resBody
}
