package main

import (
	"log"
	"strings"
	"syscall/js"
)

//displayHTML inputs the passed 'compiledTemplate'
//string into document.body.innerHTML for rendering
func displayHTML(compiledTemplate string, htmlElement js.Value) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("JS Sucks:", err)
		}
	}()

	htmlElement.Set("innerHTML", compiledTemplate)

	return
}

//setCurrentPage pushes the given url
//onto the browser's history stack
func setCurrentPage(path string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("JS Sucks:", err)
		}
	}()

	js.Global().Get("history").Call("pushState", "", "", path)

	return
}

//getCurrentPage returns the address of the
//page the browser is currently looking at
func getCurrentPage() (path string) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("JS Sucks:", err)
		}
	}()

	path = js.Global().Get("window").Get("location").Get("href").String()

	return strings.TrimSuffix(path, "/")
}
