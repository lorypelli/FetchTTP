package main

import (
	"bytes"
	j "encoding/json"
	"fmt"
	"io"
	"net/http"
	u "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/websocket"
)

func makeRequest(method string, url string, header_container *fyne.Container, body string) (string, http.Header, []byte) {
	data := []byte(body)
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return "", http.Header{}, []byte{}
	}
	req.Header.Set("User-Agent", "FetchTTP")
	for i := 0; i < len(header_container.Objects); i++ {
		header_border := header_container.Objects[i].(*fyne.Container)
		header_grid := header_border.Objects[0].(*fyne.Container)
		enabled := header_border.Objects[1].(*widget.Check)
		name := header_grid.Objects[0].(*widget.Entry).Text
		value := header_grid.Objects[1].(*widget.Entry).Text
		regexp, _ := regexp.Compile(`^[A-Za-z\d[\]{}()<>\/@?=:";,-]*$`)
		if enabled.Checked && name != "" && regexp.MatchString(name) && value != "" {
			req.Header.Add(name, value)
		}
	}
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

type Response struct {
	Headers http.Header
	Status string
	Msg     []byte
}

func ConnectWS(url string, header_container *fyne.Container, msg string, timer *time.Ticker, msg_channel chan Response) {
	for range timer.C {
		header := http.Header{}
		for i := 0; i < len(header_container.Objects); i++ {
			header_border := header_container.Objects[i].(*fyne.Container)
			header_grid := header_border.Objects[0].(*fyne.Container)
			enabled := header_border.Objects[1].(*widget.Check)
			name := header_grid.Objects[0].(*widget.Entry).Text
			value := header_grid.Objects[1].(*widget.Entry).Text
			regexp, _ := regexp.Compile(`^[A-Za-z\d[\]{}()<>\/@?=:";,-]*$`)
			if enabled.Checked && name != "" && regexp.MatchString(name) && value != "" {
				header.Add(name, value)
			}
		}
		ws, res, err := websocket.DefaultDialer.Dial(url, header)
		if err != nil {
			msg_channel <- Response{
				Headers: http.Header{},
				Status: "",
				Msg:     []byte{},
			}
		} else {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				msg_channel <- Response{
					Headers: res.Header,
					Status: res.Status,
					Msg:     []byte{},
				}
			} else {
				msg_channel <- Response{
					Headers: res.Header,
					Status: res.Status,
					Msg:     msg,
				}
			}
		}
	}
}

func main() {
	a := app.New()
	program := a.NewWindow("FetchTTP")
	program.Resize(fyne.NewSize(1280, 720))
	program.CenterOnScreen()
	req, _ := http.NewRequest("GET", "https://raw.githubusercontent.com/lorypelli/FetchTTP/main/icon.png", nil)
	c := &http.Client{}
	res, _ := c.Do(req)
	_, ok := a.(desktop.App)
	if ok {
		homeDirectory, _ := os.UserHomeDir()
		path := fmt.Sprintf("%s/FetchTTP/", homeDirectory)
		os.MkdirAll(path, os.ModePerm)
		filePath := filepath.Join(path, "icon.png")
		file, _ := os.Create(filePath)
		io.Copy(file, res.Body)
		file, _ = os.Open(filePath)
		stats, _ := os.Stat(filePath)
		size := stats.Size()
		fileBytes := make([]byte, size)
		fileSlice := fileBytes[:]
		file.Read(fileSlice)
		icon := fyne.NewStaticResource("icon.png", fileBytes)
		program.SetIcon(icon)
	}
	method := widget.NewSelect([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}, nil)
	method.Selected = "GET"
	url_http := widget.NewEntry()
	url_http.SetPlaceHolder("echo.zuplo.io")
	url_ws := widget.NewEntry()
	url_ws.SetPlaceHolder("echo.websocket.org")
	send := widget.NewButton("Send", nil)
	connect := widget.NewButton("Connect", nil)
	http_response := widget.NewLabel("")
	http_response.Wrapping = fyne.TextWrapWord
	ws_response := widget.NewList(nil, nil, nil)
	http_enabled := widget.NewCheck("", nil)
	ws_enabled := widget.NewCheck("", nil)
	http_header_name := widget.NewEntry()
	http_header_name.SetPlaceHolder("name")
	http_header_value := widget.NewEntry()
	http_header_value.SetPlaceHolder("value")
	ws_header_name := widget.NewEntry()
	ws_header_name.SetPlaceHolder("name")
	ws_header_value := widget.NewEntry()
	ws_header_value.SetPlaceHolder("value")
	http_header_name.OnChanged = func(s string) {
		if len(s) == 0 || len(http_header_value.Text) == 0 {
			http_enabled.SetChecked(false)
		} else {
			http_enabled.SetChecked(true)
		}
	}
	http_header_value.OnChanged = func(s string) {
		if len(s) == 0 || len(http_header_name.Text) == 0 {
			http_enabled.SetChecked(false)
		} else {
			http_enabled.SetChecked(true)
		}
	}
	ws_header_name.OnChanged = func(s string) {
		if len(s) == 0 || len(http_header_value.Text) == 0 {
			ws_enabled.SetChecked(false)
		} else {
			ws_enabled.SetChecked(true)
		}
	}
	ws_header_value.OnChanged = func(s string) {
		if len(s) == 0 || len(http_header_name.Text) == 0 {
			ws_enabled.SetChecked(false)
		} else {
			ws_enabled.SetChecked(true)
		}
	}
	http_header_grid := container.NewGridWithColumns(2, http_header_name, http_header_value)
	http_header_border := container.NewBorder(nil, nil, http_enabled, nil, http_header_grid)
	http_header_box := container.NewVBox(http_header_border)
	ws_header_grid := container.NewGridWithColumns(2, ws_header_name, ws_header_value)
	ws_header_border := container.NewBorder(nil, nil, ws_enabled, nil, ws_header_grid)
	ws_header_box := container.NewVBox(ws_header_border)
	tap_http := 0
	tap_ws := 0
	plus_http := widget.NewButton("+", func() {
		newEnabled := widget.NewCheck("", nil)
		newHeader_name := widget.NewEntry()
		newHeader_name.SetPlaceHolder("name")
		newHeader_value := widget.NewEntry()
		newHeader_value.SetPlaceHolder("value")
		newHeader_name.OnChanged = func(s string) {
			if len(s) == 0 || len(newHeader_value.Text) == 0 {
				newEnabled.SetChecked(false)
			} else {
				newEnabled.SetChecked(true)
			}
		}
		newHeader_value.OnChanged = func(s string) {
			if len(s) == 0 || len(newHeader_value.Text) == 0 {
				newEnabled.SetChecked(false)
			} else {
				newEnabled.SetChecked(true)
			}
		}
		http_header_grid = container.NewGridWithColumns(2, newHeader_name, newHeader_value)
		http_header_border = container.NewBorder(nil, nil, newEnabled, nil, http_header_grid)
		http_header_box.Add(http_header_border)
		tap_http -= 1
	})
	minus_http := widget.NewButton("-", func() {
		if len(http_header_box.Objects)-1 > tap_http {
			tap_http += 1
			http_header_box.Remove(http_header_box.Objects[-tap_http+1])
		}
	})
	plus_ws := widget.NewButton("+", func() {
		newEnabled := widget.NewCheck("", nil)
		newHeader_name := widget.NewEntry()
		newHeader_name.SetPlaceHolder("name")
		newHeader_value := widget.NewEntry()
		newHeader_value.SetPlaceHolder("value")
		newHeader_name.OnChanged = func(s string) {
			if len(s) == 0 || len(newHeader_value.Text) == 0 {
				newEnabled.SetChecked(false)
			} else {
				newEnabled.SetChecked(true)
			}
		}
		newHeader_value.OnChanged = func(s string) {
			if len(s) == 0 || len(newHeader_value.Text) == 0 {
				newEnabled.SetChecked(false)
			} else {
				newEnabled.SetChecked(true)
			}
		}
		ws_header_grid = container.NewGridWithColumns(2, newHeader_name, newHeader_value)
		ws_header_border = container.NewBorder(nil, nil, newEnabled, nil, http_header_grid)
		ws_header_box.Add(ws_header_border)
		tap_ws -= 1
	})
	minus_ws := widget.NewButton("-", func() {
		if len(ws_header_box.Objects)-1 > tap_ws {
			tap_ws += 1
			ws_header_box.Remove(ws_header_box.Objects[-tap_ws+1])
		}
	})
	reqbody := widget.NewMultiLineEntry()
	msg := widget.NewMultiLineEntry()
	http_options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(container.NewVBox(http_header_box, plus_http, minus_http))), container.NewTabItem("Body", reqbody))
	ws_options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(container.NewVBox(ws_header_box, plus_ws, minus_ws))), container.NewTabItem("Message", msg))
	http_response_headers := container.NewVBox()
	scroll_http_response := container.NewScroll(http_response)
	http_response_options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(http_response_headers)), container.NewTabItem("Response", scroll_http_response))
	http_response_status := widget.NewLabel("")
	ws_response_headers := container.NewVBox()
	scroll_ws_response := container.NewScroll(ws_response)
	ws_response_options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(ws_response_headers)), container.NewTabItem("Response", scroll_ws_response))
	ws_response_status := widget.NewLabel("")
	send.OnTapped = func() {
		method.Disable()
		url_http.Disable()
		send.Disable()
		http_response_headers.RemoveAll()
		http_response_options.SelectIndex(1)
		if len(url_http.Text) == 0 {
			urlWithHTTPS := fmt.Sprintf("https://%s", url_http.PlaceHolder)
			_, err := u.ParseRequestURI(urlWithHTTPS)
			if err == nil {
				status, headers, body := makeRequest(method.Selected, urlWithHTTPS, http_header_box, reqbody.Text)
				http_response_status.SetText(status)
				http_response_status.Refresh()
				for k, v := range headers {
					str, _ := j.Marshal(v)
					response_header := widget.NewLabel(fmt.Sprintf("%s: %s", k, str))
					response_header.Wrapping = fyne.TextWrapWord
					http_response_headers.Add(response_header)
				}
				http_response.SetText(string(body))
			}
		} else {
			urlWithHTTPS := url_http.Text
			if !strings.HasPrefix(strings.ToLower(urlWithHTTPS), "http") || !strings.HasPrefix(strings.ToLower(urlWithHTTPS), "https") {
				urlWithHTTPS = fmt.Sprintf("https://%s", url_http.Text)
			}
			_, err := u.ParseRequestURI(urlWithHTTPS)
			if err == nil {
				status, headers, body := makeRequest(method.Selected, urlWithHTTPS, http_header_box, reqbody.Text)
				http_response_status.SetText(status)
				http_response_status.Refresh()
				isImage := false
				for k, v := range headers {
					str, _ := j.Marshal(v)
					if strings.ToLower(k) == "content-type" && strings.Contains(strings.ToLower(string(str)), "image/png") || strings.Contains(strings.ToLower(string(str)), "image/jpeg") {
						isImage = true
					}
					response_header := widget.NewLabel(fmt.Sprintf("%s: %s", k, str))
					response_header.Wrapping = fyne.TextWrapWord
					http_response_headers.Add(response_header)
				}
				if isImage {
					http_response.Hide()
					img, _ := fyne.LoadResourceFromURLString(urlWithHTTPS)
					img_box := canvas.NewImageFromResource(img)
					img_box.FillMode = canvas.ImageFillContain
					scroll_http_response.Content = img_box
					scroll_http_response.Refresh()
				} else {
					http_response.SetText(string(body))
				}
			}
		}
		http_response.Refresh()
		method.Enable()
		url_http.Enable()
		send.Enable()
	}
	connect.OnTapped = func() {
		url_ws.Disable()
		ws_response_headers.RemoveAll()
		ws_response_options.SelectIndex(1)
		timer := time.NewTicker(500)
		ws_channel := make(chan Response)
		message := Response{}
		if connect.Text == "Disconnect" {
			url_ws.Enable()
			ws_response_status.SetText("")
			ws_response.Length = func() int {
				return 0
			}
			ws_response.Refresh()
			connect.SetText("Connect")
			return
		} else {
			connect.SetText("Disconnect")
		}
		if len(url_ws.Text) == 0 {
			urlWithWSS := url_ws.PlaceHolder
			if !strings.HasPrefix(strings.ToLower(urlWithWSS), "wss") {
				urlWithWSS = fmt.Sprintf("wss://%s", url_ws.PlaceHolder)
			}
			_, err := u.ParseRequestURI(urlWithWSS)
			if err == nil {
				var messages []string
				go ConnectWS(urlWithWSS, ws_header_box, msg.Text, timer, ws_channel)
				go func() {
					msg_number := 0
					oldMessage := ""
					message = <-ws_channel
					if message.Status == "" {
						url_ws.Enable()
						connect.SetText("Connect")
						return
					}
					ws_response_status.SetText(message.Status)
					ws_response_status.Refresh()
					for k, v := range message.Headers {
						str, _ := j.Marshal(v)
						response_header := widget.NewLabel(fmt.Sprintf("%s: %s", k, str))
						response_header.Wrapping = fyne.TextWrapWord
						ws_response_headers.Add(response_header)
					}
					for range timer.C {
						message = <-ws_channel
						newMessage := string(message.Msg)
						messages = append(messages, newMessage)
						if oldMessage == newMessage {
							return
						}
						if (len(ws_response_headers.Objects) != 0) {
							msg_number += 1
							ws_response.Length = func() int {
								return msg_number
							}
							ws_response.CreateItem = func() fyne.CanvasObject {
								ws_msg := widget.NewLabel("")
								ws_msg.Wrapping = fyne.TextWrapWord
								return ws_msg
							}
							ws_response.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
								item.(*widget.Label).SetText(messages[id])
							}
							ws_response.Refresh()
							oldMessage = string(message.Msg)
						} else {
							timer.Stop()
						}
					}
				}()
			}
		} else {
			urlWithWSS := url_ws.Text
			if !strings.HasPrefix(strings.ToLower(urlWithWSS), "wss") {
				urlWithWSS = fmt.Sprintf("wss://%s", url_ws.Text)
			}
			_, err := u.ParseRequestURI(urlWithWSS)
			if err == nil {
				var messages []string
				go ConnectWS(urlWithWSS, ws_header_box, msg.Text, timer, ws_channel)
				go func() {
					msg_number := 0
					oldMessage := ""
					message = <-ws_channel
					if message.Status == "" {
						url_ws.Enable()
						connect.SetText("Connect")
						return
					}
					ws_response_status.SetText(message.Status)
					ws_response_status.Refresh()
					for k, v := range message.Headers {
						str, _ := j.Marshal(v)
						response_header := widget.NewLabel(fmt.Sprintf("%s: %s", k, str))
						response_header.Wrapping = fyne.TextWrapWord
						ws_response_headers.Add(response_header)
					}
					for range timer.C {
						message = <-ws_channel
						newMessage := string(message.Msg)
						messages = append(messages, newMessage)
						if oldMessage == newMessage {
							return
						}
						if (len(ws_response_headers.Objects) != 0) {
							msg_number += 1
							ws_response.Length = func() int {
								return msg_number
							}
							ws_response.CreateItem = func() fyne.CanvasObject {
								ws_msg := widget.NewLabel("")
								ws_msg.Wrapping = fyne.TextWrapWord
								return ws_msg
							}
							ws_response.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
								item.(*widget.Label).SetText(messages[id])
							}
							ws_response.Refresh()
							oldMessage = string(message.Msg)
						} else {
							timer.Stop()
						}
					}
				}()
			}
		}
	}
	url_http.OnSubmitted = func(_ string) {
		method.Disable()
		url_http.Disable()
		send.Disable()
		http_response_headers.RemoveAll()
		http_response_options.SelectIndex(1)
		if len(url_http.Text) == 0 {
			urlWithHTTPS := fmt.Sprintf("https://%s", url_http.PlaceHolder)
			_, err := u.ParseRequestURI(urlWithHTTPS)
			if err == nil {
				status, headers, body := makeRequest(method.Selected, urlWithHTTPS, http_header_box, reqbody.Text)
				http_response_status.SetText(status)
				http_response_status.Refresh()
				for k, v := range headers {
					str, _ := j.Marshal(v)
					response_header := widget.NewLabel(fmt.Sprintf("%s: %s", k, str))
					response_header.Wrapping = fyne.TextWrapWord
					http_response_headers.Add(response_header)
				}
				http_response.SetText(string(body))
			}
		} else {
			urlWithHTTPS := url_http.Text
			if !strings.HasPrefix(strings.ToLower(urlWithHTTPS), "http") || !strings.HasPrefix(strings.ToLower(urlWithHTTPS), "https") {
				urlWithHTTPS = fmt.Sprintf("https://%s", url_http.Text)
			}
			_, err := u.ParseRequestURI(urlWithHTTPS)
			if err == nil {
				status, headers, body := makeRequest(method.Selected, urlWithHTTPS, http_header_box, reqbody.Text)
				http_response_status.SetText(status)
				http_response_status.Refresh()
				isImage := false
				for k, v := range headers {
					str, _ := j.Marshal(v)
					if strings.ToLower(k) == "content-type" && strings.Contains(strings.ToLower(string(str)), "image/png") || strings.Contains(strings.ToLower(string(str)), "image/jpeg") {
						isImage = true
					}
					response_header := widget.NewLabel(fmt.Sprintf("%s: %s", k, str))
					response_header.Wrapping = fyne.TextWrapWord
					http_response_headers.Add(response_header)
				}
				if isImage {
					http_response.Hide()
					img, _ := fyne.LoadResourceFromURLString(urlWithHTTPS)
					img_box := canvas.NewImageFromResource(img)
					img_box.FillMode = canvas.ImageFillContain
					scroll_http_response.Content = img_box
					scroll_http_response.Refresh()
				} else {
					http_response.SetText(string(body))
				}
			}
		}
		http_response.Refresh()
		method.Enable()
		url_http.Enable()
		send.Enable()
	}
	url_ws.OnSubmitted = func(_ string) {
		url_ws.Disable()
		ws_response_headers.RemoveAll()
		ws_response_options.SelectIndex(1)
		timer := time.NewTicker(time.Second)
		ws_channel := make(chan Response)
		message := Response{}
		connect.SetText("Disconnect")
		if len(url_ws.Text) == 0 {
			urlWithWSS := url_ws.PlaceHolder
			if !strings.HasPrefix(strings.ToLower(urlWithWSS), "wss") {
				urlWithWSS = fmt.Sprintf("wss://%s", url_ws.PlaceHolder)
			}
			_, err := u.ParseRequestURI(urlWithWSS)
			if err == nil {
				var messages []string
				go ConnectWS(urlWithWSS, ws_header_box, msg.Text, timer, ws_channel)
				go func() {
					msg_number := 0
					oldMessage := ""
					message = <-ws_channel
					if message.Status == "" {
						url_ws.Enable()
						connect.SetText("Connect")
						return
					}
					ws_response_status.SetText(message.Status)
					ws_response_status.Refresh()
					for k, v := range message.Headers {
						str, _ := j.Marshal(v)
						response_header := widget.NewLabel(fmt.Sprintf("%s: %s", k, str))
						response_header.Wrapping = fyne.TextWrapWord
						ws_response_headers.Add(response_header)
					}
					for range timer.C {
						message = <-ws_channel
						newMessage := string(message.Msg)
						messages = append(messages, newMessage)
						if oldMessage == newMessage {
							return
						}
						if (len(ws_response_headers.Objects) != 0) {
							msg_number += 1
							ws_response.Length = func() int {
								return msg_number
							}
							ws_response.CreateItem = func() fyne.CanvasObject {
								ws_msg := widget.NewLabel("")
								ws_msg.Wrapping = fyne.TextWrapWord
								return ws_msg
							}
							ws_response.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
								item.(*widget.Label).SetText(messages[id])
							}
							ws_response.Refresh()
							oldMessage = string(message.Msg)
						} else {
							timer.Stop()
						}
					}
				}()
			}
		} else {
			urlWithWSS := url_ws.Text
			if !strings.HasPrefix(strings.ToLower(urlWithWSS), "wss") {
				urlWithWSS = fmt.Sprintf("wss://%s", url_ws.Text)
			}
			_, err := u.ParseRequestURI(urlWithWSS)
			if err == nil {
				var messages []string
				go ConnectWS(urlWithWSS, ws_header_box, msg.Text, timer, ws_channel)
				go func() {
					msg_number := 0
					oldMessage := ""
					message = <-ws_channel
					if message.Status == "" {
						url_ws.Enable()
						connect.SetText("Connect")
						return
					}
					ws_response_status.SetText(message.Status)
					ws_response_status.Refresh()
					for k, v := range message.Headers {
						str, _ := j.Marshal(v)
						response_header := widget.NewLabel(fmt.Sprintf("%s: %s", k, str))
						response_header.Wrapping = fyne.TextWrapWord
						ws_response_headers.Add(response_header)
					}
					for range timer.C {
						message = <-ws_channel
						newMessage := string(message.Msg)
						messages = append(messages, newMessage)
						if oldMessage == newMessage {
							return
						}
						if (len(ws_response_headers.Objects) != 0) {
							msg_number += 1
							ws_response.CreateItem = func() fyne.CanvasObject {
								ws_msg := widget.NewLabel("")
								ws_msg.Wrapping = fyne.TextWrapWord
								return ws_msg
							}
							ws_response.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
								item.(*widget.Label).SetText(messages[id])
							}
							ws_response.Refresh()
							oldMessage = string(message.Msg)
						} else {
							timer.Stop()
						}
					}
				}()
			}
		}
	}
	http := container.NewBorder(container.NewBorder(nil, nil, method, send, container.NewBorder(nil, nil, nil, nil, url_http)), nil, nil, nil, container.NewHSplit(http_options, container.NewBorder(nil, nil, nil, http_response_status, http_response_options)))
	ws := container.NewBorder(container.NewBorder(nil, nil, nil, connect, container.NewBorder(nil, nil, nil, nil, url_ws)), nil, nil, nil, container.NewHSplit(ws_options, container.NewBorder(nil, nil, nil, ws_response_status, ws_response_options)))
	tabs := container.NewAppTabs(container.NewTabItem("HTTP", http), container.NewTabItem("WS", ws))
	tabs.SetTabLocation(container.TabLocationLeading)
	program.SetContent(tabs)
	program.ShowAndRun()
}
