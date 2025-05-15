package adventure

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var defaultHandlerTmpl = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.Title}}</title>
  </head>
  <body>
  <section class="page">
    <h1>{{.Title}}</h1>
    <div>
      {{range .Story}}
      <p>{{.}}</p>
      {{end}}
    </div>

    {{if .Options}}
    <h3>What do you want to do next?</h3>
    <ul>
      {{range .Options}}
      <li><a href="/{{.Arc}}">{{.Text}}</a></li>
      {{end}}
    </ul>

    {{end}}
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
        background: #FFFCF6;
        border: 1px solid #eee;
        box-shadow: 0 10px 6px -6px #777;
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
        text-decoration: none;
        color: #6295b5;
      }
      a:active,
      a:hover {
        color: #7792a2;
      }
      p {
        text-indent: 1em;
      }
    </style>
  </body>
</html>`

type HandlerOption func(h *handler)

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

type Story map[string]StoryArc

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

func NewHandler(s Story, opts ...HandlerOption) http.Handler {

	h := handler{s, tpl, defaultPathFn}

	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s      Story
	t      *template.Template
	pathFn func(h *http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)

	if story, ok := h.s[path]; ok {
		err := h.t.Execute(w, story)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)

		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)

}

func JsonStory(r io.Reader) (Story, error) {
	var story Story
	err := json.NewDecoder(r).Decode(&story)
	if err != nil {
		return nil, err
	}
	return story, nil
}
