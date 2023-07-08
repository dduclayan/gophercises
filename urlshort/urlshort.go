package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

type yamlData struct {
	// Notice that the var names have to be exported. TIL only EXPORTED fields will be unmarshalled. Same goes for JSON.
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
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
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusPermanentRedirect)
		}
		// else
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// a http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the Path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - Path: /some-Path
//     Url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathMap := make(map[string]string)
	if err := buildMap(yml, pathMap); err != nil {
		return nil, err
	}
	return MapHandler(pathMap, fallback), nil
}

func buildMap(yml []byte, pathMap map[string]string) error {
	var pathURLs []yamlData
	if err := yaml.Unmarshal(yml, &pathURLs); err != nil {
		return err
	}
	for _, val := range pathURLs {
		pathMap[val.Path] = val.Url
	}
	return nil
}
