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
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/websocket"
)

func SendRequest(method string, url string, header_container *fyne.Container, query_container *fyne.Container, body string) (string, http.Header, []byte) {
	data := []byte(body)
	for i := 0; i < len(query_container.Objects); i++ {
		query_border := query_container.Objects[i].(*fyne.Container)
		query_grid := query_border.Objects[0].(*fyne.Container)
		enabled := query_border.Objects[1].(*widget.Check)
		name := query_grid.Objects[0].(*widget.Entry).Text
		value := query_grid.Objects[1].(*widget.Entry).Text
		if enabled.Checked && name != "" && value != "" {
			if (strings.Contains(url, "?")) {
				url += fmt.Sprintf("&%s=%s", name, value)
			} else {
				url += fmt.Sprintf("?%s=%s", name, value)
			}
		}
	}
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
	Status  string
	Msg     []byte
}

func ConnectWS(url string, header_container *fyne.Container, query_container *fyne.Container, msg string, timer *time.Ticker, msg_channel chan Response) {
	for range timer.C {
		for i := 0; i < len(query_container.Objects); i++ {
			query_border := query_container.Objects[i].(*fyne.Container)
			query_grid := query_border.Objects[0].(*fyne.Container)
			enabled := query_border.Objects[1].(*widget.Check)
			name := query_grid.Objects[0].(*widget.Entry).Text
			value := query_grid.Objects[1].(*widget.Entry).Text
			if enabled.Checked && name != "" && value != "" {
				if (strings.Contains(url, "?")) {
					url += fmt.Sprintf("&%s=%s", name, value)
				} else {
					url += fmt.Sprintf("?%s=%s", name, value)
				}
			}
		}
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
				Status:  "",
				Msg:     []byte{},
			}
		} else {
			ws.WriteMessage(websocket.TextMessage, []byte(msg))
			_, msg, err := ws.ReadMessage()
			if err != nil {
				msg_channel <- Response{
					Headers: res.Header,
					Status:  res.Status,
					Msg:     []byte{},
				}
			} else {
				msg_channel <- Response{
					Headers: res.Header,
					Status:  res.Status,
					Msg:     msg,
				}
			}
		}
	}
}

func main() {
	a := app.New()
	program := a.NewWindow("FetchTTP")
	cookies := a.NewWindow("Cookies")
	program.Resize(fyne.NewSize(1280, 720))
	cookies.Resize(fyne.NewSize(500, 205))
	cookies.SetFixedSize(true)
	program.CenterOnScreen()
	cookies.CenterOnScreen()
	cookies.SetCloseIntercept(func() {
		cookies.Hide()
	})
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
		cookies.SetIcon(icon)
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
	http_header_enabled := widget.NewCheck("", nil)
	ws_header_enabled := widget.NewCheck("", nil)
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
			http_header_enabled.SetChecked(false)
		} else {
			http_header_enabled.SetChecked(true)
		}
	}
	http_header_value.OnChanged = func(s string) {
		if len(s) == 0 || len(http_header_name.Text) == 0 {
			http_header_enabled.SetChecked(false)
		} else {
			http_header_enabled.SetChecked(true)
		}
	}
	ws_header_name.OnChanged = func(s string) {
		if len(s) == 0 || len(ws_header_value.Text) == 0 {
			ws_header_enabled.SetChecked(false)
		} else {
			ws_header_enabled.SetChecked(true)
		}
	}
	ws_header_value.OnChanged = func(s string) {
		if len(s) == 0 || len(ws_header_name.Text) == 0 {
			ws_header_enabled.SetChecked(false)
		} else {
			ws_header_enabled.SetChecked(true)
		}
	}
	http_header_grid := container.NewGridWithColumns(2, http_header_name, http_header_value)
	http_header_border := container.NewBorder(nil, nil, http_header_enabled, nil, http_header_grid)
	http_header_box := container.NewVBox(http_header_border)
	ws_header_grid := container.NewGridWithColumns(2, ws_header_name, ws_header_value)
	ws_header_border := container.NewBorder(nil, nil, ws_header_enabled, nil, ws_header_grid)
	ws_header_box := container.NewVBox(ws_header_border)
	tap_header_http := 0
	tap_header_ws := 0
	plus_header_http := widget.NewButton("+", func() {
		http_newEnabled := widget.NewCheck("", nil)
		http_newHeader_name := widget.NewEntry()
		http_newHeader_name.SetPlaceHolder("name")
		http_newHeader_value := widget.NewEntry()
		http_newHeader_value.SetPlaceHolder("value")
		http_newHeader_name.OnChanged = func(s string) {
			if len(s) == 0 || len(http_newHeader_value.Text) == 0 {
				http_newEnabled.SetChecked(false)
			} else {
				http_newEnabled.SetChecked(true)
			}
		}
		http_newHeader_value.OnChanged = func(s string) {
			if len(s) == 0 || len(http_newHeader_value.Text) == 0 {
				http_newEnabled.SetChecked(false)
			} else {
				http_newEnabled.SetChecked(true)
			}
		}
		http_header_grid = container.NewGridWithColumns(2, http_newHeader_name, http_newHeader_value)
		http_header_border = container.NewBorder(nil, nil, http_newEnabled, nil, http_header_grid)
		http_header_box.Add(http_header_border)
		tap_header_http -= 1
	})
	minus_header_http := widget.NewButton("-", func() {
		if len(http_header_box.Objects)-1 > tap_header_http {
			tap_header_http += 1
			http_header_box.Remove(http_header_box.Objects[-tap_header_http+1])
		}
	})
	plus_header_ws := widget.NewButton("+", func() {
		ws_newEnabled := widget.NewCheck("", nil)
		ws_newHeader_name := widget.NewEntry()
		ws_newHeader_name.SetPlaceHolder("name")
		ws_newHeader_value := widget.NewEntry()
		ws_newHeader_value.SetPlaceHolder("value")
		ws_newHeader_name.OnChanged = func(s string) {
			if len(s) == 0 || len(ws_newHeader_value.Text) == 0 {
				ws_newEnabled.SetChecked(false)
			} else {
				ws_newEnabled.SetChecked(true)
			}
		}
		ws_newHeader_value.OnChanged = func(s string) {
			if len(s) == 0 || len(ws_newHeader_value.Text) == 0 {
				ws_newEnabled.SetChecked(false)
			} else {
				ws_newEnabled.SetChecked(true)
			}
		}
		ws_header_grid = container.NewGridWithColumns(2, ws_newHeader_name, ws_newHeader_value)
		ws_header_border = container.NewBorder(nil, nil, ws_newEnabled, nil, ws_header_grid)
		ws_header_box.Add(ws_header_border)
		tap_header_ws -= 1
	})
	minus_header_ws := widget.NewButton("-", func() {
		if len(ws_header_box.Objects)-1 > tap_header_ws {
			tap_header_ws += 1
			ws_header_box.Remove(ws_header_box.Objects[-tap_header_ws+1])
		}
	})
	http_query_enabled := widget.NewCheck("", nil)
	ws_query_enabled := widget.NewCheck("", nil)
	http_query_name := widget.NewEntry()
	http_query_name.SetPlaceHolder("name")
	http_query_value := widget.NewEntry()
	http_query_value.SetPlaceHolder("value")
	ws_query_name := widget.NewEntry()
	ws_query_name.SetPlaceHolder("name")
	ws_query_value := widget.NewEntry()
	ws_query_value.SetPlaceHolder("value")
	http_query_name.OnChanged = func(s string) {
		if len(s) == 0 || len(http_query_value.Text) == 0 {
			http_query_enabled.SetChecked(false)
		} else {
			http_query_enabled.SetChecked(true)
		}
	}
	http_query_value.OnChanged = func(s string) {
		if len(s) == 0 || len(http_query_name.Text) == 0 {
			http_query_enabled.SetChecked(false)
		} else {
			http_query_enabled.SetChecked(true)
		}
	}
	ws_query_name.OnChanged = func(s string) {
		if len(s) == 0 || len(ws_query_value.Text) == 0 {
			ws_query_enabled.SetChecked(false)
		} else {
			ws_query_enabled.SetChecked(true)
		}
	}
	ws_query_value.OnChanged = func(s string) {
		if len(s) == 0 || len(ws_query_name.Text) == 0 {
			ws_query_enabled.SetChecked(false)
		} else {
			ws_query_enabled.SetChecked(true)
		}
	}
	http_query_grid := container.NewGridWithColumns(2, http_query_name, http_query_value)
	http_query_border := container.NewBorder(nil, nil, http_query_enabled, nil, http_query_grid)
	http_query_box := container.NewVBox(http_query_border)
	ws_query_grid := container.NewGridWithColumns(2, ws_query_name, ws_query_value)
	ws_query_border := container.NewBorder(nil, nil, ws_query_enabled, nil, ws_query_grid)
	ws_query_box := container.NewVBox(ws_query_border)
	tap_query_http := 0
	tap_query_ws := 0
	plus_query_http := widget.NewButton("+", func() {
		http_newEnabled := widget.NewCheck("", nil)
		http_newQuery_name := widget.NewEntry()
		http_newQuery_name.SetPlaceHolder("name")
		http_newQuery_value := widget.NewEntry()
		http_newQuery_value.SetPlaceHolder("value")
		http_newQuery_name.OnChanged = func(s string) {
			if len(s) == 0 || len(http_newQuery_value.Text) == 0 {
				http_newEnabled.SetChecked(false)
			} else {
				http_newEnabled.SetChecked(true)
			}
		}
		http_newQuery_value.OnChanged = func(s string) {
			if len(s) == 0 || len(http_newQuery_value.Text) == 0 {
				http_newEnabled.SetChecked(false)
			} else {
				http_newEnabled.SetChecked(true)
			}
		}
		http_query_grid = container.NewGridWithColumns(2, http_newQuery_name, http_newQuery_value)
		http_query_border = container.NewBorder(nil, nil, http_newEnabled, nil, http_query_grid)
		http_query_box.Add(http_query_border)
		tap_query_http -= 1
	})
	minus_query_http := widget.NewButton("-", func() {
		if len(http_query_box.Objects)-1 > tap_query_http {
			tap_query_http += 1
			http_query_box.Remove(http_query_box.Objects[-tap_query_http+1])
		}
	})
	plus_query_ws := widget.NewButton("+", func() {
		ws_newEnabled := widget.NewCheck("", nil)
		ws_newQuery_name := widget.NewEntry()
		ws_newQuery_name.SetPlaceHolder("name")
		ws_newQuery_value := widget.NewEntry()
		ws_newQuery_value.SetPlaceHolder("value")
		ws_newQuery_name.OnChanged = func(s string) {
			if len(s) == 0 || len(ws_newQuery_value.Text) == 0 {
				ws_newEnabled.SetChecked(false)
			} else {
				ws_newEnabled.SetChecked(true)
			}
		}
		ws_newQuery_value.OnChanged = func(s string) {
			if len(s) == 0 || len(ws_newQuery_value.Text) == 0 {
				ws_newEnabled.SetChecked(false)
			} else {
				ws_newEnabled.SetChecked(true)
			}
		}
		ws_query_grid = container.NewGridWithColumns(2, ws_newQuery_name, ws_newQuery_value)
		ws_query_border = container.NewBorder(nil, nil, ws_newEnabled, nil, ws_query_grid)
		ws_query_box.Add(ws_query_border)
		tap_query_ws -= 1
	})
	minus_query_ws := widget.NewButton("-", func() {
		if len(ws_query_box.Objects)-1 > tap_query_ws {
			tap_query_ws += 1
			ws_query_box.Remove(ws_query_box.Objects[-tap_query_ws+1])
		}
	})
	reqbody := widget.NewMultiLineEntry()
	reqbody.SetPlaceHolder("Request Body")
	msg := widget.NewMultiLineEntry()
	msg.SetPlaceHolder("Message Content")
	ws_send := widget.NewButton("Send", nil)
	http_cookie_box := container.NewVBox(widget.NewButton("Add Cookie", func() {
		key := widget.NewEntry()
		key.SetPlaceHolder("Key")
		value := widget.NewEntry()
		value.SetPlaceHolder("Value")
		domain := widget.NewEntry()
		domain.SetPlaceHolder("Domain")
		path := widget.NewEntry()
		path.SetPlaceHolder("Path")
		calendar := widget.NewEntry()
		calendar.SetPlaceHolder("Expires")
		secure := widget.NewCheck("Secure", nil)
		http_only := widget.NewCheck("HTTPOnly", nil)
		cookies.SetContent(container.NewVBox(container.NewBorder(nil, nil, nil, nil, container.NewGridWithRows(5, container.NewGridWithColumns(2, key, value), container.NewGridWithColumns(2, domain, path), calendar, container.NewGridWithColumns(2, secure, http_only), widget.NewButton("Save", func() {
			c := &http.Cookie{
				Name:    key.Text,
				Value:   value.Text,
				Domain:  domain.Text,
				Path:    path.Text,
				Secure: secure.Checked,
				HttpOnly: http_only.Checked,
			}
			err := c.Valid()
			if err != nil {
				dialog.ShowError(err, cookies)
				return
			}
			c_str := c.String()
			c_enabled := widget.NewCheck("", nil)
			c_enabled.SetChecked(true)
			c_name := widget.NewEntry()
			c_name.SetPlaceHolder("Name")
			c_name.SetText("Cookie")
			c_name.OnChanged = func(s string) {
				if len(s) == 0 || len(c_name.Text) == 0 {
					c_enabled.SetChecked(false)
				} else {
					c_enabled.SetChecked(true)
				}
			}
			c_value := widget.NewEntry()
			c_value.SetPlaceHolder("Value")
			c_value.SetText(c_str)
			c_value.OnChanged = func(s string) {
				if len(s) == 0 || len(c_value.Text) == 0 {
					c_enabled.SetChecked(false)
				} else {
					c_enabled.SetChecked(true)
				}
			}
			http_header_box.Objects = append(http_header_box.Objects, container.NewBorder(nil, nil, c_enabled, nil, container.NewGridWithColumns(2, c_name, c_value)))
			http_header_box.Refresh()
			cookies.Hide()
		})))))
		cookies.Show()
	}))
	ws_cookie_box := container.NewVBox(widget.NewButton("Add Cookie", func() {
		key := widget.NewEntry()
		key.SetPlaceHolder("Key")
		value := widget.NewEntry()
		value.SetPlaceHolder("Value")
		domain := widget.NewEntry()
		domain.SetPlaceHolder("Domain")
		path := widget.NewEntry()
		path.SetPlaceHolder("Path")
		calendar := widget.NewEntry()
		calendar.SetPlaceHolder("Expires")
		secure := widget.NewCheck("Secure", nil)
		http_only := widget.NewCheck("HTTPOnly", nil)
		cookies.SetContent(container.NewVBox(container.NewBorder(nil, nil, nil, nil, container.NewGridWithRows(5, container.NewGridWithColumns(2, key, value), container.NewGridWithColumns(2, domain, path), calendar, container.NewGridWithColumns(2, secure, http_only), widget.NewButton("Save", func() {
			c := &http.Cookie{
				Name:    key.Text,
				Value:   value.Text,
				Domain:  domain.Text,
				Path:    path.Text,
				Secure: secure.Checked,
				HttpOnly: http_only.Checked,
			}
			err := c.Valid()
			if err != nil {
				dialog.ShowError(err, cookies)
				return
			}
			c_str := c.String()
			c_enabled := widget.NewCheck("", nil)
			c_enabled.SetChecked(true)
			c_name := widget.NewEntry()
			c_name.SetPlaceHolder("Name")
			c_name.SetText("Cookie")
			c_name.OnChanged = func(s string) {
				if len(s) == 0 || len(c_name.Text) == 0 {
					c_enabled.SetChecked(false)
				} else {
					c_enabled.SetChecked(true)
				}
			}
			c_value := widget.NewEntry()
			c_value.SetPlaceHolder("Value")
			c_value.SetText(c_str)
			c_value.OnChanged = func(s string) {
				if len(s) == 0 || len(c_value.Text) == 0 {
					c_enabled.SetChecked(false)
				} else {
					c_enabled.SetChecked(true)
				}
			}
			ws_header_box.Objects = append(ws_header_box.Objects, container.NewBorder(nil, nil, c_enabled, nil, container.NewGridWithColumns(2, c_name, c_value)))
			ws_header_box.Refresh()
			cookies.Hide()
		})))))
		cookies.Show()
	}))
	http_options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(container.NewVBox(http_header_box, plus_header_http, minus_header_http))), container.NewTabItem("Query", container.NewScroll(container.NewVBox(http_query_box, plus_query_http, minus_query_http))), container.NewTabItem("Cookies", http_cookie_box), container.NewTabItem("Body", reqbody))
	ws_options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(container.NewVBox(ws_header_box, plus_header_ws, minus_header_ws))), container.NewTabItem("Query", container.NewScroll(container.NewVBox(ws_query_box, plus_query_ws, minus_query_ws))), container.NewTabItem("Cookies", ws_cookie_box), container.NewTabItem("Message", container.NewBorder(nil, ws_send, nil, nil, msg)))
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
				status, headers, body := SendRequest(method.Selected, urlWithHTTPS, http_header_box, http_query_box, reqbody.Text)
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
				status, headers, body := SendRequest(method.Selected, urlWithHTTPS, http_header_box, http_query_box, reqbody.Text)
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
				go ConnectWS(urlWithWSS, ws_header_box, ws_query_box, msg.Text, timer, ws_channel)
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
						if len(ws_response_headers.Objects) != 0 {
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
				go ConnectWS(urlWithWSS, ws_header_box, ws_query_box, msg.Text, timer, ws_channel)
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
						if len(ws_response_headers.Objects) != 0 {
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
				status, headers, body := SendRequest(method.Selected, urlWithHTTPS, http_header_box, http_query_box, reqbody.Text)
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
				status, headers, body := SendRequest(method.Selected, urlWithHTTPS, http_header_box, http_query_box, reqbody.Text)
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
				go ConnectWS(urlWithWSS, ws_header_box, ws_query_box, msg.Text, timer, ws_channel)
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
						if len(ws_response_headers.Objects) != 0 {
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
				go ConnectWS(urlWithWSS, ws_header_box, ws_query_box, msg.Text, timer, ws_channel)
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
						if len(ws_response_headers.Objects) != 0 {
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
	program.Show()
	a.Run()
}
