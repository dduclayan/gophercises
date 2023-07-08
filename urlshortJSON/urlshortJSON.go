package urlshortJSON

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonData struct {
	// Notice that the var names have to be exported. TIL only EXPORTED fields will be unmarshalled. Same goes for JSON.
	Path string `json:"path"`
	Url  string `json:"url"`
}

// MapHandler will return a http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the Path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		fmt.Println(path)
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusPermanentRedirect)
		}
		// else
		fallback.ServeHTTP(w, r)
	}
}

// JSONHandler will parse the provided JSON and then return
// a http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the Path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
func JSONHandler(js []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathMap := make(map[string]string)
	if err := buildMap(js, pathMap); err != nil {
		return nil, err
	}
	fmt.Println(pathMap)
	return MapHandler(pathMap, fallback), nil
}

func buildMap(js []byte, pathMap map[string]string) error {
	var pathURLs []jsonData
	if err := json.Unmarshal(js, &pathURLs); err != nil {
		return err
	}
	for _, val := range pathURLs {
		fmt.Println(val)
		pathMap[val.Path] = val.Url
	}
	return nil
}
