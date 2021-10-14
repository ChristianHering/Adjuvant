package main

import (
	"encoding/json"
	"log"
	"syscall/js"
)

//Replaces the innerHTML of the html attribute '0' with
//that located at the path '1'. Optionally, you can
//pass a rest url for use when compiling the template
//
//Values:
//   0- (string) HTML ID of the attribute to change
//   1- (string) Path to the template
//   2- (bool) Whether or not to fetch rest data
//   3- (string) Path to the rest endpoint
//   4- (int) Refresh rate in seconds (rest bool must be true)
func changeContentByID(this js.Value, inputs []js.Value) interface{} {
	if len(inputs) < 4 {
		log.Println("Argument/paramater count mismatch")

		return nil
	}

	body, err := getRequestBody(inputs[1].String())
	if err != nil {
		log.Println("Failed to fetch template:", err)

		return nil
	}

	var template string

	if inputs[2].Bool() {
		body, err := getRequestBody(inputs[3].String())
		if err != nil {
			log.Println("Failed to get rest data for changeByID call:", err)

			return nil
		}

		var m map[string]string

		err = json.Unmarshal(body, &m)
		if err != nil {
			log.Println("Parsing rest data failed:", err)

			return nil
		}

		template, err = compileTemplate(body, m)
	} else {
		template, err = compileTemplate(body, nil)
	}
	if err != nil { //for compileTemplate
		log.Println("Template compilation failed:", err)
	}

	displayHTML(template, js.Global().Get("document").Call("getElementById", inputs[0].String()))

	return nil
}

//refreshes data from rest server
func updateContentByID(this js.Value, inputs []js.Value) interface{} {
	//document.getElementById
	js.Global().Get("document").Call(inputs[0].String(), js.FuncOf(onPopState))

	return nil
}

func changeContentByClass(this js.Value, inputs []js.Value) interface{} {
	//document.getElementsByClassName

	return nil
}

//refreshes data from rest server
func updateContentByClass(this js.Value, inputs []js.Value) interface{} {
	//document.getElementsByClassName

	return nil
}

//onPopState recieves a callback from javascript
//when the page's path changes. This essentially
//imitates a backbutton callback function
func onPopState(this js.Value, inputs []js.Value) interface{} {
	//If this function (or any 'onpopstate' handler) blocks,
	//it will immediately induce a goroutine deadlock and crash
	//our wasm binary. Run everything in a seperate goroutine.
	go func() {
		page := getCurrentPage()

		if currentPage != page {
			err := loadWebPage(page)
			if err != nil {
				log.Println("Internet Sucks:", err)
			}
		}
	}()

	return nil
}
