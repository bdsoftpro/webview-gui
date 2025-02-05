package main

import (
	"fmt"

	"github.com/bdsoftpro/webview-gui"
	"github.com/ncruces/zenity"
)

func main() {

	w := webview.New(true, true)
	defer w.Destroy()

	w.SetWindowEventsHandler("main", EventHandler)

	w.SetTitle("Basic Example")
	err := w.SetIcon("../asset/icon.png")
	if err != nil {
		fmt.Println(err.Error())
	}
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml("Thanks for using Golang Webview GUI!")
	w.Run()
}

func EventHandler(state webview.WindowState) {
	if state == webview.WindowClose {
		zenity.Info("Window Closed", zenity.NoIcon)
	}
}
