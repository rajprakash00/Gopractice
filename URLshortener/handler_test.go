package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	fallbackResponse = "fallback"
	paths            = "/github"
	url              = "https://github.com"
)

func assertResponse(t *testing.T, res *http.Response, body string) {
	t.Helper()
	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		t.Fatal("could not find response body", err)
	}

	if body != string(data) {
		t.Errorf("Expected response body to be %s & got %s", body, string(data))
	}
}

func assertStatusCode(t *testing.T, res *http.Response, code int) {
	t.Helper()
	if res.StatusCode != code {
		t.Errorf("Expected value to be %d & got %d", code, res.StatusCode)
	}
}

func assertURL(t *testing.T, res *http.Response, link string) {
	t.Helper()
	url, err := res.Location()

	if err != nil {
		t.Fatal("could not find route", err)
	}

	if url.String() != link {
		t.Errorf("Expected URL to be %s & got %s", url, link)
	}
}

func fallback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, fallbackResponse)
}

func createMapHandler(pathToURL map[string]string) http.HandlerFunc {

	return MapHandler(pathToURL, http.HandlerFunc(fallback))
}

func mapHandlerResult(pathToURL map[string]string, path string) *http.Response {
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	mapHandler := createMapHandler(pathToURL)
	mapHandler(response, req)

	return response.Result()
}

func TestMapHandler(t *testing.T) {
	pathToURL := map[string]string{paths: url}

	t.Run("it redirects to desired url", func(t *testing.T) {
		res := mapHandlerResult(pathToURL, paths)

		assertStatusCode(t, res, http.StatusFound)
		assertURL(t, res, url)
	})

	t.Run("it uses fallback for unknown routes", func(t *testing.T) {
		res := mapHandlerResult(pathToURL, "/error")

		assertResponse(t, res, fallbackResponse)
	})
}

/*--- Test methods for YAML/JSON handling---*/

type httpHandler func([]byte, http.Handler) (http.HandlerFunc, error)

func fileHandlerResult(handler httpHandler, data []byte, path string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		fmt.Println(err, "oops reading err")
	}
	response := httptest.NewRecorder()

	fileHandler := createFileHandler(handler, data)
	fileHandler(response, req)

	return response.Result()
}

func createFileHandler(handler httpHandler, data []byte) http.HandlerFunc {
	fileHandler, err := handler(data, http.HandlerFunc(fallback))

	if err != nil {
		fmt.Println(err)
		panic("could not create a file read handler for yaml/json")
	}
	return fileHandler
}

func TestYAMLHandler(t *testing.T) { //write yaml indentation correctly
	yml := `
- paths: /github
  url: https://github.com
`
	yaml := []byte(yml)
	t.Run("it redirects to desired url", func(t *testing.T) {
		res := fileHandlerResult(YAMLHandler, yaml, paths)

		assertStatusCode(t, res, http.StatusFound)
		assertURL(t, res, url)
	})

	t.Run("it uses fallback for unknown routes", func(t *testing.T) {
		res := fileHandlerResult(YAMLHandler, yaml, "/error")

		assertResponse(t, res, fallbackResponse)
	})
}

func TestJSONHandler(t *testing.T) {
	jsn := `[
		{
			"paths": "/github",
			"url": "https://github.com"
		}
		]`
	json := []byte(jsn)
	t.Run("it redirects to desired url", func(t *testing.T) {
		res := fileHandlerResult(JSONHandler, json, paths)

		assertStatusCode(t, res, http.StatusFound)
		assertURL(t, res, url)
	})

	t.Run("it uses fallback for unknown routes", func(t *testing.T) {
		res := fileHandlerResult(JSONHandler, json, "/error")

		assertResponse(t, res, fallbackResponse)
	})

}
