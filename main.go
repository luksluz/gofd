package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"sync"

	sd "./ScanDir"
)

var (
	debug bool
	wg    sync.WaitGroup
)

func eventHandler(event sd.Events) {
	if event.EvType == sd.EVFILEADDED {
		// fmt.Println("added")
		if w != nil {
			w.Dispatch(func() {
				w.Eval(fmt.Sprintf(`window.app.files.push('%s')`, event.FileEv.Name))
			})
		}
	} else {
		if w != nil {
			sendFiles(event.Path)
		}

	}
	if debug {
		fmt.Println("debug")
	}

}

func showAllFiles() {
	if debug {
		wts := sd.ShowAllFiles()
		bts, err := json.Marshal(wts)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("wts.json", bts, 0666)
	}

}

func main() {
	debug = false
	argLength := len(os.Args)
	if argLength > 1 {
		args := strings.Join(os.Args, " ")
		if strings.Contains(args, "--DEBUG") || strings.Contains(args, "-d") {
			debug = true
		}
	}
	// wg.Add(1)

	User, err := user.Current()
	if err != nil {
		panic(err)
	}

	defaultDir := User.HomeDir + "\\downloads"
	sd.New(eventHandler, defaultDir)
	sd.Scan()

	startUi()
	sendDirs()
	sd.Wait()
}
