package main

import "C"
import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	webview "github.com/bdsoftpro/webview-gui"
	"github.com/bdsoftpro/webview-gui/pkg/dialog"
	"github.com/bdsoftpro/webview-gui/pkg/systray"
	"github.com/kardianos/osext"
	"github.com/ncruces/zenity"
)

//go:embed icon.png
var EF embed.FS

type App_ struct {
	Path, File, FileName string
}

var (
	iconData     []byte
	App          = &App_{}
	html         = `Thanks for using Golang Webview GUI!<br><br><input type="text" width='75%'>`
	ChildWindows []*exec.Cmd
)

func init() {
	var err error
	if App.Path, err = osext.ExecutableFolder(); err != nil {
		log.Fatal(err)
	}
	if App.File, err = osext.Executable(); err != nil {
		log.Fatal(err)
	} else {
		App.FileName = filepath.Base(App.File)
	}
}

func main() {
	var err error

	iconData, err = EF.ReadFile("icon.png")
	if err != nil {
		panic(err)
	}

	if !zenity.IsAvailable() {
		dialog.Message("Zenity is not installed.").Info()
	}

	// See NewInstance function
	new_window := flag.Bool("new_window", false, "")
	title := flag.String("title", "New Window", "")
	width := flag.Int("width", 200, "")
	height := flag.Int("height", 200, "")
	hint := flag.String("hint", "none", "")
	url := flag.String("url", "", "")
	flag.Parse()
	if *new_window {
		NewInstance(title, width, height, hint, url)
	}

	w := webview.New(true, false)
	defer w.Destroy()

	w.SetWindowEventsHandler("main", func(state webview.WindowState) {
		//fmt.Println(state)
		switch state {
		case webview.WindowClose:
			fmt.Println("Window state: WindowClose")
			w.Hide() // Click "Show" after
		case webview.WindowResize:
			// Example: save window size for restore in next launch
			fmt.Println("Window state: WindowResize")
		case webview.WindowMove:
			// Example: save window position for restore in next launch
			fmt.Println("Window state: WindowMove")
		case webview.WindowFocus:
			fmt.Println("Window state: WindowFocus")
		case webview.WindowBlur:
			fmt.Println("Window state: WindowBlur")
		case webview.WindowFullScreen:
			fmt.Println("Window state: WindowFullScreen")
		case webview.WindowExitFullScreen:
			fmt.Println("Window state: WindowExitFullScreen")
		case webview.WindowMinimize:
			fmt.Println("Window state: WindowMinimize")
		case webview.WindowUnminimize:
			fmt.Println("Window state: WindowUnminimize")
		case webview.WindowMaximize:
			fmt.Println("Window state: WindowMaximize")
		case webview.WindowUnmaximize:
			fmt.Println("Window state: WindowUnmaximize")

		}
	})

	w.SetContentStateHandler("main", func(state string) {
		fmt.Printf("document content state: %s\n", state)
	})

	systray.Register(onReady(w))
	w.SetTitle("Systray Example")
	w.SetIconBites(iconData, len(iconData))
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml(html)
	//w.Navigate("https://google.com")
	w.Run()
}

// NewInstance  - WebView does not support creating a new instance from the current application.
// Therefore, this is a possible option for creating a new window.
func NewInstance(title *string, width *int, height *int, hint *string, url *string) {
	w := webview.New(false, true)
	defer w.Destroy()
	w.SetTitle(*title)
	var _hint webview.Hint
	switch *hint {
	case "none":
		_hint = webview.HintNone
	case "fixed":
		_hint = webview.HintFixed
	}
	w.SetIconBites(iconData, len(iconData))
	w.SetSize(*width, *height, _hint)
	w.Navigate(*url)
	w.Focus()
	w.Run()
}

func onReady(w webview.WebView) func() {
	return func() {
		systray.SetTitle("Tray")
		systray.SetIcon(iconData)
		go func() {

			mShowTitle := systray.AddMenuItem("GetTitle", "")
			mGetSizePosition := systray.AddMenuItem("Get Size and Position", "")
			mHide := systray.AddMenuItem("Hide", "")
			mShow := systray.AddMenuItem("Show", "")
			mFocus := systray.AddMenuItem("Focus", "")
			mMove := systray.AddMenuItem("Move", "")
			mMaximize := systray.AddMenuItem("Maximize", "")
			mUnmaximize := systray.AddMenuItem("Unmaximize", "")
			mMinimize := systray.AddMenuItem("Minimize", "")
			mUnminimize := systray.AddMenuItem("Unminimize", "")
			mSetFullScreen := systray.AddMenuItem("SetFullScreen", "")
			mExitFullScreen := systray.AddMenuItem("ExitFullScreen", "")
			mSetAlwaysOnTopTrue := systray.AddMenuItem("SetAlwaysOnTop True", "")
			mSetAlwaysOnTopFalse := systray.AddMenuItem("SetAlwaysOnTop False", "")
			mSetBorderless := systray.AddMenuItem("SetBorderless", "")
			mSetBordered := systray.AddMenuItem("SetBordered", "")
			mSetUserAgent := systray.AddMenuItem("SetUserAgent", "")
			mSetDraggable := systray.AddMenuItem("SetDraggable", "")
			mUnSetDraggable := systray.AddMenuItem("UnSetDraggable", "")
			mGetHtml := systray.AddMenuItem("GetHtml", "")
			mGetUrl := systray.AddMenuItem("GetUrl", "")
			mGetPageTitle := systray.AddMenuItem("GetPageTitle", "")
			mOpenBrowser := systray.AddMenuItem("Open google.com", "")
			systray.AddSeparator()
			mNewWindow := systray.AddMenuItem("New Window", "")
			systray.AddSeparator()
			mQuit := systray.AddMenuItem("Quit                 ", "")

			for {
				select {
				case <-mQuit.ClickedCh:
					// signal not working on Windows
					//syscall.Kill(syscall.Getpid(), syscall.SIGINT)

					// kill child processes
					KillChildWindows()
					w.Terminate()
					os.Exit(1)
				case <-mShowTitle.ClickedCh:
					zenity.Info(w.GetTitle(), zenity.Title("Info"), zenity.NoIcon, GetAttachId())
				case <-mGetSizePosition.ClickedCh:
					width, height, hint := w.GetSize()
					x, y := w.GetPosition()
					zenity.Info(fmt.Sprintf("Size:%dx%d %v. \nPosition: X:%d Y:%d", width, height, hint, x, y), zenity.Title("Info"),
						zenity.NoIcon, GetAttachId())
				case <-mHide.ClickedCh:
					w.Dispatch(func() {
						w.Hide()
					})
				case <-mShow.ClickedCh:
					w.Dispatch(func() {
						w.Show()
					})
				case <-mFocus.ClickedCh:
					w.Dispatch(func() {
						w.Focus()
					})
				case <-mMove.ClickedCh:
					w.Dispatch(func() {
						w.Move(100, 100)
					})
				case <-mMaximize.ClickedCh:
					w.Dispatch(func() {
						w.Maximize()
					})
				case <-mUnmaximize.ClickedCh:
					w.Dispatch(func() {
						w.Unmaximize()
					})
				case <-mMinimize.ClickedCh:
					w.Dispatch(func() {
						w.Minimize()
					})
				case <-mUnminimize.ClickedCh:
					w.Dispatch(func() {
						w.Unminimize()
					})
				case <-mSetBorderless.ClickedCh:
					w.Dispatch(func() {
						w.SetBorderless()
					})
				case <-mSetBordered.ClickedCh:
					w.Dispatch(func() {
						w.SetBordered()
					})
				case <-mSetAlwaysOnTopTrue.ClickedCh:
					w.Dispatch(func() {
						w.SetAlwaysOnTop(true)
					})
				case <-mSetAlwaysOnTopFalse.ClickedCh:
					w.Dispatch(func() {
						w.SetAlwaysOnTop(false)
					})
				case <-mSetFullScreen.ClickedCh:
					w.Dispatch(func() {
						w.SetFullScreen()
					})
				case <-mExitFullScreen.ClickedCh:
					w.Dispatch(func() {
						w.ExitFullScreen()
					})
				case <-mSetUserAgent.ClickedCh:
					w.Dispatch(func() {
						w.SetUserAgent("My User Agent")
						w.Navigate("https://httpbin.org/user-agent")
					})
				case <-mSetDraggable.ClickedCh:
					w.Dispatch(func() {
						w.SetDraggable("drg")
						w.SetHtml(`<div style="width:97%; height: 20px; border: 1px solid blue;padding: 5px; text-align: center;" id="drg">Drag me</div>`)
						w.SetBorderless()
					})
				case <-mUnSetDraggable.ClickedCh:
					w.Dispatch(func() {
						w.UnSetDraggable("drg")
						w.SetHtml(html)
						w.SetBordered()
					})
				case <-mGetHtml.ClickedCh:
					zenity.Info(w.GetHtml(), zenity.Title("Info"), zenity.NoIcon, GetAttachId())

				case <-mGetUrl.ClickedCh:
					zenity.Info(w.GetUrl(), zenity.Title("Info"), zenity.NoIcon, GetAttachId())

				case <-mGetPageTitle.ClickedCh:
					w.Dispatch(func() {
						w.Navigate("https://golang.org")
						w.SetContentStateHandler("page-title", func(state string) {
							if state == "complete" {
								zenity.Info(w.GetPageTitle(), zenity.Title("Info"), zenity.NoIcon, GetAttachId())
								w.UnSetContentStateHandler("page-title")
							}
						})
					})

				case <-mOpenBrowser.ClickedCh:
					w.OpenUrlInBrowser("https://google.com")

				case <-mNewWindow.ClickedCh:
					p := []string{
						"-new_window",
						"-title", "About",
						"-width", "300",
						"-height", "200",
						"-hint", "fixed",
						"-url", "file://" + App.Path + "/../asset/about.html",
					}
					cmd := exec.Command(App.File, p...)
					//fmt.Println(cmd.Args)

					// signal not working on Windows
					/*c := make(chan os.Signal, 2)
					signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM) //
					go func() {
						<-c
						cmd.Process.Kill()
					}()*/

					if err := cmd.Start(); err != nil {
						zenity.Error(err.Error(), zenity.Title("Error"), zenity.ErrorIcon, GetAttachId())
					} else {
						ChildWindows = append(ChildWindows, cmd)
					}

				}
			}
		}()
	}
}

func KillChildWindows() {
	for _, w := range ChildWindows {
		w.Process.Kill()
	}
}

// GetAttachId This function is required only on Linux operating systems.
// On Linux zenity does not return focus to parent window
// For other systems, this is just a fake implementation
//
//	var wid interface{}
//	switch runtime.GOOS {
//	case "linux":
//		wid = os.Getpid()
//	case "darwin":
//		wid = os.Getpid() // or example "My.app".
//	case "windows":
//		wid = uintptr(w.Window())
//	}
func GetAttachId() zenity.Option {
	if runtime.GOOS == "linux" {
		return zenity.Attach(os.Getpid())
	}
	return zenity.Modal()
}
