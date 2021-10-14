package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"syscall/js"
	"text/template"
)

//loadWebPage gets the data and template for a
//given path, and paints the page to the screen
func loadWebPage(path string) error {
	templateBody, err := getRequestBody(path + "/template")
	if err != nil {
		return err
	}

	jsonBody, err := getRequestBody(path + "/rest")
	if err != nil {
		return err
	}

	var m map[string]string

	err = json.Unmarshal(jsonBody, &m)
	if err != nil {
		m = nil
	}

	compiledTemplate, err := compileTemplate(templateBody, m)
	if err != nil {
		return err
	}

	//You must set the current page before navigating otherwise
	//onPopState() will catch the page change before the path's
	//set and initiate a new page change to the previous page
	currentPage = path

	setCurrentPage(path)

	displayHTML(compiledTemplate, js.Global().Get("document").Get("body"))

	return nil
}

//compileTemplate takes in a byte array, and returns
//a http template for rendering pages and subpages
func compileTemplate(requestBody []byte, templateData interface{}) (compiledTemplate string, err error) {
	t, err := template.New("").Parse(string(requestBody))
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, templateData)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

//getRequestBody returns a byte array from
//the body of a http.Get request for 'path'
func getRequestBody(path string) (requestBody []byte, err error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
