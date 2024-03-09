package main

import (
	"bytes"
	"context"
	j "encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
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
	Enabled bool
	Name    string
	Value   string
}

type Query struct {
	Enabled bool
	Name    string
	Value   string
}

type Response struct {
	Status string
	Header http.Header
	Body   string
}

func (a *App) MakeRequest(method string, url string, headers []Header, query []Query, body string) Response {
	data := []byte(body)
	for i := 0; i < len(query); i++ {
		if query[i].Enabled && strings.TrimSpace(query[i].Name) != "" && strings.TrimSpace(query[i].Value) != "" {
			if strings.Contains(url, "?") {
				url += fmt.Sprintf("&%s=%s", query[i].Name, query[i].Value)
			} else {
				url += fmt.Sprintf("?%s=%s", query[i].Name, query[i].Value)
			}
		}
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return Response{
			"", http.Header{}, "",
		}
	}
	for i := 0; i < len(headers); i++ {
		regexp, _ := regexp.Compile(`^[A-Za-z\d[\]{}()<>\/@?=:";,-]*$`)
		if headers[i].Enabled && strings.TrimSpace(headers[i].Name) != "" && regexp.MatchString(headers[i].Name) && strings.TrimSpace(headers[i].Value) != "" {
			req.Header.Add(headers[i].Name, headers[i].Value)
		}
	}
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return Response{
			"", http.Header{}, "",
		}
	}
	var resBody []byte
	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		var jsonBody map[string]interface{}
		bytes, _ := io.ReadAll(res.Body)
		j.Unmarshal(bytes, &jsonBody)
		resBody, _ = j.MarshalIndent(jsonBody, "", "  ")
	} else {
		bytes, _ := io.ReadAll(res.Body)
		resBody = bytes
	}
	return Response{
		res.Status, res.Header, string(resBody),
	}
}
