package main

import (
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

func makeRequest(method string, url string, header_container *fyne.Container) string {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return ""
	} else {
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
		res, _ := c.Do(req)
		bytes, _ := io.ReadAll(res.Body)
		body := string(bytes)
		return body
	}
}

func main() {
	a := app.New()
	program := a.NewWindow("FetchTTP")
	program.Resize(fyne.NewSize(1280, 720))
	program.CenterOnScreen()
	method := widget.NewSelect([]string{"GET", "POST", "PUT", "PATCH", "DELETE"}, nil)
	method.Selected = "GET"
	url := widget.NewEntry()
	url.SetPlaceHolder("https://echo.zuplo.io/")
	send := widget.NewButton("Send", nil)
	response := widget.NewLabel("")
	response.Wrapping = fyne.TextWrapWord
	enabled := widget.NewCheck("", nil)
	enabled.SetChecked(true)
	header_name := widget.NewEntry()
	header_name.SetPlaceHolder("name")
	header_value := widget.NewEntry()
	header_value.SetPlaceHolder("value")
	header_grid := container.NewGridWithColumns(2, header_name, header_value)
	header_border := container.NewBorder(nil, nil, enabled, nil, header_grid)
	header_box := container.NewVBox(header_border)
	tap := 0
	plus := widget.NewButton("+", func() {
		newHeader_name := widget.NewEntry()
		newHeader_name.SetPlaceHolder("name")
		newHeader_value := widget.NewEntry()
		newHeader_value.SetPlaceHolder("value")
		header_grid = container.NewGridWithColumns(2, newHeader_name, newHeader_value)
		header_border = container.NewBorder(nil, nil, enabled, nil, header_grid)
		header_box.Add(header_border)
		tap -= 1
	})
	minus := widget.NewButton("-", func() {
		if len(header_box.Objects) - 1 > tap {
			tap += 1
			header_box.Remove(header_box.Objects[-tap])
		}
	})
	options := container.NewAppTabs(container.NewTabItem("Headers", container.NewScroll(container.NewVBox(header_box, plus, minus))), container.NewTabItem("Body", widget.NewMultiLineEntry()))
	scroll_response := container.NewScroll(response)
	send.OnTapped = func() {
		if len(url.Text) == 0 {
			_, err := u.ParseRequestURI(url.PlaceHolder)
			if err == nil {
				response.SetText(makeRequest(method.Selected, url.PlaceHolder, header_box))
			}
		} else {
			_, err := u.ParseRequestURI(url.PlaceHolder)
			if err == nil {
				if strings.HasSuffix(strings.ToLower(url.Text), ".png") || strings.HasSuffix(strings.ToLower(url.Text), ".jpg") || strings.HasSuffix(strings.ToLower(url.Text), ".jpeg") {
					response.Hide()
					img, _ := fyne.LoadResourceFromURLString(url.Text)
					img_box := canvas.NewImageFromResource(img)
					img_box.FillMode = canvas.ImageFillContain
					scroll_response.Content = img_box
					scroll_response.Refresh()
				} else {
					response.SetText(makeRequest(method.Selected, url.Text, header_box))
				}
			}
		}
		response.Refresh()
	}
	program.SetContent(container.NewBorder(container.NewBorder(nil, nil, method, send, container.NewBorder(nil, nil, nil, nil, url)), nil, nil, nil, container.NewHSplit(options, scroll_response)))
	program.ShowAndRun()
}
