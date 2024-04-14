package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	port = flag.String("port", "8080", "port for server to listen on")

	tmpl *template.Template
)

func init() {
	tmpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.Title}}</h1>
      {{range .Paragraphs}}
        <p>{{.}}</p>
      {{end}}
      <ul>
      {{range .Options}}
        <li><a href="/{{.Chapter}}">{{.Text}}</a></li>
      {{end}}
      </ul>
    </section>
    <style>
      body {
        font-family: helvetica, arial;
      }
      h1 {
        text-align:center;
        position:relative;
      }
      .page {
        width: 80%;
        max-width: 500px;
        margin: auto;
        margin-top: 40px;
        margin-bottom: 40px;
        padding: 80px;
        background: #FCF6FC;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #797;
      }
      ul {
        border-top: 1px dotted #ccc;
        padding: 10px 0 0 0;
        -webkit-padding-start: 0;
      }
      li {
        padding-top: 10px;
      }
      a,
      a:visited {
        text-decoration: underline;
        color: #555;
      }
      a:active,
      a:hover {
        color: #222;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </body>
</html>
`

type handler struct {
	s Story
}

func newHandler(s Story) http.Handler {
	return handler{s}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	// "/intro" -> "intro"
	path = path[1:]
	log.Info().Msgf("path = %v\n", path)

	if chapter, ok := h.s[path]; ok {
		err := tmpl.Execute(w, chapter)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

var (
	storyPath = flag.String("path", "story.json", "path to story json file")
)

func main() {
	flag.Parse()

	storyData, err := loadStory(*storyPath)
	if err != nil {
		log.Fatal().Msgf("failed to load story: %v\n", err)
	}

	h := newHandler(storyData)

	fmt.Printf("starting the server on port :%s\n", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", *port), h); err != nil {
		log.Fatal().Msgf("failed to start server: %v\n", err)
	}
}
