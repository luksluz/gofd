package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/user"

	sd "./ScanDir"
	"github.com/google/uuid"
	"github.com/zserge/webview"
)

var (
	w webview.WebView
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func startUi() {
	// fb, _ := ioutil.ReadFile("index.html")
	// vue, _ := ioutil.ReadFile("static\\vue.js")
	// app, _ := ioutil.ReadFile("static\\app.js")
	// css, _ := ioutil.ReadFile("static\\style.css")
	fb := MustAsset("index.html")
	vue := MustAsset("static/vue.js")
	app := MustAsset("static/app.js")
	css := MustAsset("static/style.css")

	w = webview.New(webview.Settings{
		URL:                    `data:text/html,` + url.PathEscape(string(fb)),
		Title:                  "WatcherDown v0.0.1",
		Width:                  800,
		Height:                 600,
		ExternalInvokeCallback: uiEvents,
	})

	w.Dispatch(func() {
		w.InjectCSS(string(css))
		w.Eval(string(vue))
		w.Eval(string(app))
		if _, err := os.Stat("data"); !os.IsNotExist(err) {
			w.Eval("app.page = 'watcher'")
		}
		// w.Eval(`app.page = 'welcome`)

		usr, _ := user.Current()
		// // w.Eval( fmt.Sprintf("window.app.user.name = '%s'", User.Name))
		j, _ := json.Marshal(&User{usr.Name, 27})
		w.Eval(fmt.Sprintf(`window.app.user = %s`, string(j)))
	})

	sendDirs()
	defer quit()

	w.Run()
}

type Response struct {
	Event string `json:"event"`
	Value string `json:"value"`
}

func sendDirs() {
	var dirs []string
	for k := range sd.Wts.Watchers {
		dirs = append(dirs, k)
	}
	w.Dispatch(func() {
		j, _ := json.Marshal(dirs)
		w.Eval(fmt.Sprintf("app.dirs = %s", string(j)))
	})
}

func sendFiles(dir string) {
	w.Dispatch(func() {
		tFiles := sd.Wts.Watchers[dir].Files
		var files []string
		for _, file := range tFiles {
			files = append(files, file.Name)
		}
		j, _ := json.Marshal(files)
		w.Eval(fmt.Sprintf(`window.app.files = %s`, string(j)))
	})
}

func uiEvents(wv webview.WebView, data string) {
	var resp Response
	json.Unmarshal([]byte(data), &resp)

	// log.Println("here")

	switch resp.Event {
	case "noWelcome":
		if _, err := os.Stat("data"); os.IsNotExist(err) {
			id := uuid.New()
			uid, _ := id.MarshalBinary()
			ioutil.WriteFile("data", uid, 0666)
		}

	case "getFiles":
		// log.Println(resp.Value)
		sendFiles(resp.Value)
	case "searchDir":
		w.Dialog(webview.DialogTypeOpen, 0, "Localizar diretorio", "")
	case "addDir":
		if _, err := os.Stat(resp.Value); os.IsNotExist(err) {
			//TODO: SEND SAME ERROR MSG TO GUI
		} else {
			log.Println(resp.Value)
			sd.Wts.Watchers[resp.Value] = sd.Watcher{Path: resp.Value}
			sd.Scan()
			sendDirs()
			//TODO: SEND SUCESS MSG BACK
			// w.Dispatch(func(){
			// 	w.Eval(`window.app.`)
			// })
		}
	}
}

func quit() {
	os.Exit(0)
}
