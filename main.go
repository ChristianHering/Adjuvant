package main

import (
	"syscall/js"
)

var currentPage string

func main() {
	//Expose "exported" functions to the JS runtime
	js.Global().Set("changeContentByID", js.FuncOf(changeContentByID))
	js.Global().Set("updateContentByID", js.FuncOf(updateContentByID))
	js.Global().Set("changeContentByClass", js.FuncOf(changeContentByClass))
	js.Global().Set("updateContentByClass", js.FuncOf(updateContentByClass))

	//WindowEventHandlers.onpopstate allows us to get a callback
	//whenever the current page changes (for backbutton detection)
	js.Global().Get("window").Set("onpopstate", js.FuncOf(onPopState))

	err := retryFunc(3, func() error { return loadWebPage(getCurrentPage()) })
	if err != nil {
		panic(err)
	}

	//TODO: Refresh page data every X seconds depending on the page

	//TODO: Predownload/parse templates for the current page

	<-make(chan int)
}

//retryFunc calls f 'tryCount' times until it returns a nil output.
func retryFunc(tryCount int, f func() error) (err error) {
	if tryCount == 0 {
		for {
			err = f()
			if err == nil {
				return nil
			}
		}
	} else {
		for i := 0; i < tryCount; i++ {
			err = f()
			if err == nil {
				return nil
			}
		}

		return err
	}
}
