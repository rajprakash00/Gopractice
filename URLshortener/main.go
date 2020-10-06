package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func main() {
	filePath := flag.String("f", "paths.yaml", "path/to/file")
	flag.Parse()

	mux := defaultMux()

	pathToURL := map[string]string{
		"/urlshort-godoc": "https://github.com/rajprakash00/Go-practice",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mapHandler := MapHandler(pathToURL, mux)

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("error opening file\n", err)

	}
	r, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("error reading file\n", err)
	}

	//YAML/JSON builder using mapHandler as fallback
	var pathHandler http.HandlerFunc
	switch {
	case path.Ext(*filePath) == ".json":
		pathHandler, err = JSONHandler(r, mapHandler)
		if err != nil {
			panic(err)
		}

	default:
		pathHandler, err = YAMLHandler(r, mapHandler)
		if err != nil {
			panic(err)
		}

	}

	fmt.Println("starting a server on :8000")
	http.ListenAndServe(":8000", pathHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerHello)
	return mux
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}
