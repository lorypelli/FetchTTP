package main

import (
	"bytes"
	j "encoding/json"
	"fmt"
	"io"
	"net/http"
	u "net/url"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
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

func main() {
	a := app.New()
	program := a.NewWindow("FetchTTP")
	program.Resize(fyne.NewSize(1280, 720))
	program.CenterOnScreen()
	method := widget.NewSelect([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}, nil)
	method.Selected = "GET"
	url_http := widget.NewEntry()
	url_http.SetPlaceHolder("echo.zuplo.io")
	url_ws := widget.NewEntry()
	url_ws.SetPlaceHolder("echo.zuplo.io")
	send := widget.NewButton("Send", nil)
	connect := widget.NewButton("Connect", nil)
	http_response := widget.NewLabel("")
	http_response.Wrapping = fyne.TextWrapWord
	ws_response := widget.NewLabel("")
	ws_response.Wrapping = fyne.TextWrapWord
	enabled := widget.NewCheck("", nil)
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
			enabled.SetChecked(false)
		} else {
			enabled.SetChecked(true)
		}
	}
	http_header_value.OnChanged = func(s string) {
		if len(s) == 0 || len(http_header_name.Text) == 0 {
			enabled.SetChecked(false)
		} else {
			enabled.SetChecked(true)
		}
	}
	http_header_grid := container.NewGridWithColumns(2, http_header_name, http_header_value)
	http_header_border := container.NewBorder(nil, nil, enabled, nil, http_header_grid)
	http_header_box := container.NewVBox(http_header_border)
	tap_http := 0
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
	reqbody := widget.NewMultiLineEntry()
	http_options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(container.NewVBox(http_header_box, plus_http, minus_http))), container.NewTabItem("Body", reqbody))
	http_response_headers := container.NewVBox()
	scroll_http_response := container.NewScroll(http_response)
	http_response_options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(http_response_headers)), container.NewTabItem("Response", scroll_http_response))
	http_response_status := widget.NewLabel("")
	send.OnTapped = func() {
		http_response_headers.RemoveAll()
		http_response_options.SelectIndex(1)
		if len(url_http.Text) == 0 {
			urlWithHTTPS := fmt.Sprintf("https://%s", url_http.PlaceHolder)
			_, err := u.ParseRequestURI(urlWithHTTPS)
			if err == nil {
				status, headers, body := makeRequest(method.Selected, urlWithHTTPS, http_header_box, reqbody.Text)
				http_response_status.SetText(fmt.Sprint(status))
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
				http_response_status.SetText(fmt.Sprint(status))
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
	}
	connect.OnTapped = func() {
		url_ws.Disable()
		http_response_headers.RemoveAll()
		http_response_options.SelectIndex(1)
		if len(url_http.Text) == 0 {
			urlWithHTTPS := fmt.Sprintf("https://%s", url_http.PlaceHolder)
			_, err := u.ParseRequestURI(urlWithHTTPS)
			if err == nil {
				status, headers, body := makeRequest(method.Selected, urlWithHTTPS, http_header_box, reqbody.Text)
				http_response_status.SetText(fmt.Sprint(status))
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
				http_response_status.SetText(fmt.Sprint(status))
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
	}
	url_http.OnSubmitted = func(_ string) {
		http_response_headers.RemoveAll()
		http_response_options.SelectIndex(1)
		if len(url_http.Text) == 0 {
			urlWithHTTPS := fmt.Sprintf("https://%s", url_http.PlaceHolder)
			_, err := u.ParseRequestURI(urlWithHTTPS)
			if err == nil {
				status, headers, body := makeRequest(method.Selected, urlWithHTTPS, http_header_box, reqbody.Text)
				http_response_status.SetText(fmt.Sprint(status))
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
				http_response_status.SetText(fmt.Sprint(status))
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
	}
	http := container.NewBorder(container.NewBorder(nil, nil, method, send, container.NewBorder(nil, nil, nil, nil, url_http)), nil, nil, nil, container.NewHSplit(http_options, container.NewBorder(nil, nil, nil, http_response_status, http_response_options)))
	ws := container.NewBorder(container.NewBorder(nil, nil, nil, connect, container.NewBorder(nil, nil, nil, nil, url_ws)), nil, nil, nil, container.NewHSplit(http_options, container.NewBorder(nil, nil, nil, http_response_status, http_response_options)))
	tabs := container.NewAppTabs(container.NewTabItem("HTTP", http), container.NewTabItem("WS", ws))
	tabs.SetTabLocation(container.TabLocationLeading)
	program.SetContent(tabs)
	program.ShowAndRun()
}
