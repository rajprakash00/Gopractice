package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

//MapURL ...
type MapURL struct {
	Path string `yaml:"paths" json:"paths"`
	URL  string `yaml:"url" json:"url"`
}

//MapHandler method to implement map handling of paths to URLs
//if path not provided then fallback http handler will be called instead
//MapHandler...
func MapHandler(pathToURL map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathToURL[r.URL.Path]; ok {
			http.Redirect(w, r, path, http.StatusFound)
		}
		fallback.ServeHTTP(w, r)
	}
}

//MapBuilder ...
func MapBuilder(mapURL []MapURL) map[string]string {
	yamlMap := make(map[string]string)
	for _, mapPath := range mapURL {
		yamlMap[mapPath.Path] = mapPath.URL
	}
	return yamlMap
}

//YAMLHandler ...
func YAMLHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var mapURL []MapURL
	err := yaml.Unmarshal(data, &mapURL)

	if err != nil {
		return nil, err
	}

	yamlMap := MapBuilder(mapURL)

	return MapHandler(yamlMap, fallback), nil

	//an alternative to not using MapHandler method

	/*return func(w http.ResponseWriter, r *http.Request) {
		for _, mapPath := range mapURL {
			if mapPath.Path == r.URL.Path {
				http.Redirect(w, r, mapPath.URL, http.StatusFound)
				return
			}
		}
		fallback.ServeHTTP(w, r)
	}*/
}

//JSONHandler ...
func JSONHandler(data []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var mapURL []MapURL
	err := json.Unmarshal(data, &mapURL)
	if err != nil {
		return nil, err
	}

	jsonMap := MapBuilder(mapURL)

	return MapHandler(jsonMap, fallback), nil
}
