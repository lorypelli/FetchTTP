package main

import (
	"bytes"
	"context"
	j "encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	r "runtime"
	"strconv"
	"strings"
	"time"

	"github.com/474420502/gcurl"
	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const APP_VERSION = "1.3.8"

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
	bin, _ := os.Executable()
	path := filepath.Dir(bin)
	goos := r.GOOS
	switch goos {
	case "windows": {
		bat := filepath.Join(path, "temp.bat")
		_, err := os.Stat(bat)
		if !os.IsNotExist(err) {
			os.Remove(bat)
		}
		break
	}
	case "linux": {
		sh := filepath.Join(path, "temp.sh")
		_, err := os.Stat(sh)
		if !os.IsNotExist(err) {
			os.Remove(sh)
		}
		break
	}
	}
	bat := filepath.Join(path, "temp.bat")
	_, err := os.Stat(bat)
	if !os.IsNotExist(err) {
		os.Remove(bat)
	}
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

type HTTPResponse struct {
	URL    string
	Status string
	Header http.Header
	Body   string
	Error  string
}

type WSResponse struct {
	Ws      *websocket.Conn
	Status  string
	Header  http.Header
	Message string
}

type Update struct {
	IsLatest    bool
	Version     string
	Description string
	Error       string
}

type Github struct {
	Tag_name string
	Body     string
	Assets   []Download
}

type Download struct {
	Browser_download_url string
}

func (a *App) HTTP(method string, url string, headers []Header, query []Query, body string) HTTPResponse {
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
		return HTTPResponse{
			url, "", http.Header{}, "", err.Error(),
		}
	}
	for i := 0; i < len(headers); i++ {
		regexp, _ := regexp.Compile(`^[A-Za-z\d[\]{}()<>\/@?=:";,-]*$`)
		if headers[i].Enabled && strings.TrimSpace(headers[i].Name) != "" && regexp.MatchString(headers[i].Name) && strings.TrimSpace(headers[i].Value) != "" {
			req.Header.Add(headers[i].Name, headers[i].Value)
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return HTTPResponse{
			url, "", http.Header{}, "", err.Error(),
		}
	}
	var resBody []byte
	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		var jsonBody interface{}
		bytes, _ := io.ReadAll(res.Body)
		j.Unmarshal(bytes, &jsonBody)
		resBody, _ = j.MarshalIndent(jsonBody, "", "\t")
	} else {
		bytes, _ := io.ReadAll(res.Body)
		resBody = bytes
	}
	return HTTPResponse{
		url, res.Status, res.Header, string(resBody), "",
	}
}

var currentConnection *websocket.Conn
var currentResponse *http.Response
var currentError error

func (a *App) WS(url string, headers []Header, query []Query, connected bool) string {
	for i := 0; i < len(query); i++ {
		if query[i].Enabled && strings.TrimSpace(query[i].Name) != "" && strings.TrimSpace(query[i].Value) != "" {
			if strings.Contains(url, "?") {
				url += fmt.Sprintf("&%s=%s", query[i].Name, query[i].Value)
			} else {
				url += fmt.Sprintf("?%s=%s", query[i].Name, query[i].Value)
			}
		}
	}
	var header http.Header
	for i := 0; i < len(headers); i++ {
		regexp, _ := regexp.Compile(`^[A-Za-z\d[\]{}()<>\/@?=:";,-]*$`)
		if headers[i].Enabled && strings.TrimSpace(headers[i].Name) != "" && regexp.MatchString(headers[i].Name) && strings.TrimSpace(headers[i].Value) != "" {
			header.Add(headers[i].Name, headers[i].Value)
		}
	}
	if connected {
		ws, res, err := websocket.DefaultDialer.Dial(url, header)
		currentConnection = ws
		currentResponse = res
		currentError = err
	}
	if currentError == nil {
		go Connect(currentResponse, currentConnection, connected, a)
	} else {
		return currentError.Error()
	}
	return ""
}

func Connect(res *http.Response, ws *websocket.Conn, connected bool, a *App) {
	runtime.EventsOff(a.ctx, "connected")
	runtime.EventsOff(a.ctx, "message")
	runtime.EventsOn(a.ctx, "connected", func(data ...interface{}) {
		d := fmt.Sprint(data[0])
		connected, _ = strconv.ParseBool(d)
	})
	runtime.EventsOn(a.ctx, "message", func(data ...interface{}) {
		d := fmt.Sprint(data[0])
		if strings.TrimSpace(d) != "" {
			ws.WriteMessage(websocket.TextMessage, []byte(d))
		}
	})
	for {
		if connected {
			_, msg, err := ws.ReadMessage()
			if err == nil {
				runtime.EventsEmit(a.ctx, "websocket", WSResponse{
					ws, res.Status, res.Header, string(msg),
				})
			}
		} else {
			ws.Close()
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (a *App) CheckUpdates() Update {
	res, err := http.Get("https://api.github.com/repos/lorypelli/FetchTTP/releases/latest")
	if err != nil {
		return Update{
			true, "", "", err.Error(),
		}
	}
	jsonBody := Github{}
	bytes, _ := io.ReadAll(res.Body)
	j.Unmarshal(bytes, &jsonBody)
	tag := jsonBody.Tag_name
	description := jsonBody.Body
	version := strings.ReplaceAll(tag, ".", "")
	currentVersion := strings.ReplaceAll(APP_VERSION, ".", "")
	version_int, _ := strconv.Atoi(version)
	currentVersion_int, _ := strconv.Atoi(currentVersion)
	if version_int > currentVersion_int {
		return Update{
			false, tag, description, "",
		}
	}
	return Update{
		true, "", "", "",
	}
}

func (a *App) Update() {
	res, err := http.Get("https://api.github.com/repos/lorypelli/FetchTTP/releases/latest")
	if err != nil {
		return
	}
	jsonBody := Github{}
	bytes, _ := io.ReadAll(res.Body)
	j.Unmarshal(bytes, &jsonBody)
	bin, _ := os.Executable()
	path, _ := filepath.Abs(bin)
	dir := filepath.Dir(bin)
	var windows int
	var linux int
	for i := 0; i < len(jsonBody.Assets); i++ {
		switch jsonBody.Assets[i].Browser_download_url {
		case "FetchTTP": {
			linux = i
			break
		}
		case "FetchTTP.exe": {
			windows = i
			break
		}
		}
	}
	goos := r.GOOS
	switch goos {
	case "windows":
		{
			res, err = http.Get(jsonBody.Assets[windows].Browser_download_url)
			if err != nil {
				return
			}
			bytes, _ = io.ReadAll(res.Body)
			newPath := strings.Split(path, ".exe")[0] + "1.exe"
			file, _ := os.Create(newPath)
			file.Write(bytes)
			tempPath := filepath.Join(dir, "temp.bat")
			tempFile, _ := os.Create(tempPath)
			exe := filepath.Base(path)
			commands := `@echo off
cd ` + dir + `
taskkill /F /IM ` + exe + `
taskkill /F /IM ` + exe + `
del /F ` + exe + `
rename "` + filepath.Base(newPath) + `" "` + exe + `"
start ` + exe + `
exit`
			tempFile.WriteString(commands)
			cmd := exec.Command("cmd.exe", "/C", tempPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			break
		}
	case "linux":
		{
			res, err = http.Get(jsonBody.Assets[linux].Browser_download_url)
			if err != nil {
				return
			}
			bytes, _ = io.ReadAll(res.Body)
			newPath := path + "1"
			file, _ := os.Create(newPath)
			file.Write(bytes)
			tempPath := filepath.Join(dir, "temp.sh")
			tempFile, _ := os.Create(tempPath)
			exe := filepath.Base(path)
			commands := `cd ` + dir + `
pkill -9 ` + exe + `
rm -rf ` + exe + `
mv "` + filepath.Base(newPath) + `" "` + exe + `"
chmod +x ` + exe + `
./` + exe + `
exit`
			tempFile.WriteString(commands)
			cmd := exec.Command("bash", tempPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
			break
		}
	}
}

func (a *App) CURL(url string) HTTPResponse {
	curl := gcurl.Parse(url)
	req, err := http.NewRequest(curl.Method, url, curl.Body)
	req.Header.Add("User-Agent", "FetchTTP")
	for k, v := range curl.Header {
		req.Header.Add(k, strings.Join(v, ","))
	}
	if err != nil {
		return HTTPResponse{
			url, "", http.Header{}, "", err.Error(),
		}
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return HTTPResponse{
			url, "", http.Header{}, "", err.Error(),
		}
	}
	var resBody []byte
	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		var jsonBody interface{}
		bytes, _ := io.ReadAll(res.Body)
		j.Unmarshal(bytes, &jsonBody)
		resBody, _ = j.MarshalIndent(jsonBody, "", "\t")
	} else {
		bytes, _ := io.ReadAll(res.Body)
		resBody = bytes
	}
	return HTTPResponse{
		url, res.Status, res.Header, string(resBody), "",
	}
}
