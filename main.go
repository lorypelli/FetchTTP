package main

import (
	"bytes"
	"encoding/json"
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

func makeRequest(method string, url string, header_container *fyne.Container, body string) (string, http.Header, string) {
	data := []byte(body)
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return "", http.Header{}, ""
	}
	req.Header.Set("User-Agent", "FetchTTP")
	for i := 0; i < len(header_container.Objects); i++ {
		header_border := header_container.Objects[i].(*fyne.Container)
		header_grid := header_border.Objects[0].(*fyne.Container)
		enabled := header_border.Objects[1].(*widget.Check)
		name := header_grid.Objects[0].(*widget.Entry).Text
		value := header_grid.Objects[1].(*widget.Entry).Text
		regexp, _ := regexp.Compile(`^[A-Za-z[\]{}()<>\/@?=:";,-]*$`)
		if enabled.Checked && name != "" && regexp.MatchString(name) && value != "" {
			req.Header.Add(name, value)
		}
	}
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return "", http.Header{}, ""
	}
	bytes, _ := io.ReadAll(res.Body)
	resbody := string(bytes)
	return res.Status, res.Header, resbody
}

func main() {
	a := app.New()
	program := a.NewWindow("FetchTTP")
	program.Resize(fyne.NewSize(1280, 720))
	program.CenterOnScreen()
	method := widget.NewSelect([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}, nil)
	method.Selected = "GET"
	url := widget.NewEntry()
	url.SetPlaceHolder("echo.zuplo.io/")
	send := widget.NewButton("Send", nil)
	response := widget.NewLabel("")
	response.Wrapping = fyne.TextWrapWord
	enabled := widget.NewCheck("", nil)
	header_name := widget.NewEntry()
	header_name.SetPlaceHolder("name")
	header_value := widget.NewEntry()
	header_value.SetPlaceHolder("value")
	header_name.OnChanged = func(s string) {
		if len(s) == 0 || len(header_value.Text) == 0 {
			enabled.SetChecked(false)
		} else {
			enabled.SetChecked(true)
		}
	}
	header_value.OnChanged = func(s string) {
		if len(s) == 0 || len(header_name.Text) == 0 {
			enabled.SetChecked(false)
		} else {
			enabled.SetChecked(true)
		}
	}
	header_grid := container.NewGridWithColumns(2, header_name, header_value)
	header_border := container.NewBorder(nil, nil, enabled, nil, header_grid)
	header_box := container.NewVBox(header_border)
	tap := 0
	plus := widget.NewButton("+", func() {
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
		header_grid = container.NewGridWithColumns(2, newHeader_name, newHeader_value)
		header_border = container.NewBorder(nil, nil, newEnabled, nil, header_grid)
		header_box.Add(header_border)
		tap -= 1
	})
	minus := widget.NewButton("-", func() {
		if len(header_box.Objects)-1 > tap {
			tap += 1
			header_box.Remove(header_box.Objects[-tap+1])
		}
	})
	reqbody := widget.NewMultiLineEntry()
	options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(container.NewVBox(header_box, plus, minus))), container.NewTabItem("Body", reqbody))
	response_headers := container.NewVBox()
	scroll_response := container.NewScroll(response)
	response_options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(response_headers)), container.NewTabItem("Response", scroll_response))
	response_status := widget.NewLabel("")
	send.OnTapped = func() {
		response_headers.RemoveAll()
		response_options.SelectIndex(1)
		if len(url.Text) == 0 {
			urlWithHTTPS := fmt.Sprintf("https://%s", url.PlaceHolder)
			_, err := u.ParseRequestURI(urlWithHTTPS)
			if err == nil {
				status, headers, body := makeRequest(method.Selected, urlWithHTTPS, header_box, reqbody.Text)
				response_status.SetText(fmt.Sprint(status))
				response_status.Refresh()
				for k, v := range headers {
					str, _ := json.Marshal(v)
					response_header := widget.NewLabel(fmt.Sprintf("%s: %s", k, str))
					response_header.Wrapping = fyne.TextWrapWord
					response_headers.Add(response_header)
				}
				response.SetText(body)
			}
		} else {
			urlWithHTTPS := url.Text
			if !strings.HasPrefix(strings.ToLower(urlWithHTTPS), "http") || !strings.HasPrefix(strings.ToLower(urlWithHTTPS), "https") {
				urlWithHTTPS = fmt.Sprintf("https://%s", url.Text)
			}
			_, err := u.ParseRequestURI(urlWithHTTPS)
			if err == nil {
				status, headers, body := makeRequest(method.Selected, urlWithHTTPS, header_box, reqbody.Text)
				response_status.SetText(fmt.Sprint(status))
				response_status.Refresh()
				isImage := false
				for k, v := range headers {
					str, _ := json.Marshal(v)
					if strings.ToLower(k) == "content-type" && strings.Contains(strings.ToLower(string(str)), "image/png") || strings.Contains(strings.ToLower(string(str)), "image/jpeg") {
						isImage = true
					}
					response_header := widget.NewLabel(fmt.Sprintf("%s: %s", k, str))
					response_header.Wrapping = fyne.TextWrapWord
					response_headers.Add(response_header)
				}
				if isImage {
					response.Hide()
					img, _ := fyne.LoadResourceFromURLString(urlWithHTTPS)
					img_box := canvas.NewImageFromResource(img)
					img_box.FillMode = canvas.ImageFillContain
					scroll_response.Content = img_box
					scroll_response.Refresh()
				} else {
					response.SetText(body)
				}
			}
		}
		response.Refresh()
	}
	program.SetContent(container.NewBorder(container.NewBorder(nil, nil, method, send, container.NewBorder(nil, nil, nil, nil, url)), nil, nil, nil, container.NewHSplit(options, container.NewBorder(nil, nil, nil, response_status, response_options))))
	program.ShowAndRun()
}
