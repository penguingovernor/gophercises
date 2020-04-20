package adventure

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type Story map[string]Storyblock

type Storyblock struct {
	Title      string              `json:"title"`
	Paragraphs []string            `json:"story"`
	Options    []StoryBlockOptions `json:"options"`
}

type StoryBlockOptions struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

func LoadStory(rd io.Reader) (*Story, error) {
	var st Story
	if err := json.NewDecoder(rd).Decode(&st); err != nil {
		return nil, err
	}
	return &st, nil
}

func (s Story) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Parse(storyHTML))

	resource := strings.TrimFunc(r.URL.Path, func(letter rune) bool {
		return letter == '/'
	})

	path, found := s[resource]
	if !found {
		http.Error(w, fmt.Sprintf("We couldn't find %s, sorry bud :/", r.URL.Path), http.StatusNotFound)
		return
	}
	tmpl.Execute(w, path)
}

const storyHTML = `<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <title>Choose your adventure!</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.8.2/css/bulma.min.css">
  <script defer src="https://use.fontawesome.com/releases/v5.3.1/js/all.js"></script>
</head>

<body>

  <section class="section">
    <div id="CurrentStory" class="container">
      <h1 id="StoryTitle" class="title">
        {{.Title}}
      </h1>

      {{range .Paragraphs}}
      <p class="StoryText">
        {{.}}
      </p>
      {{end}}
    </div>
  </section>

  <section class="section">
    <div class="container">
      <div id="StoryOptions" class="card">
        <div class="card-title">
          <p>Choose your adventure:</p>
        </div>
        <div class="card-content">
          <select id="storyOptions">
            {{range .Options}}
            <option value="{{.Arc}}">{{.Arc}}</option>
            {{end}}
          </select>
          <button onclick="goToPath()" class="button is-link">Go!</button>
        </div>
      </div>
    </div>
  </section>


  <footer class="page-footer blue">
    <div class="content has-text-centered">
      <span>Pssst... Want to start over? <a href="/intro">click here</a></span>
    </div>
  </footer>
  <script>
    const goToPath = () => {
      const selectElement = document.getElementById('storyOptions')
      const path = selectElement.options[selectElement.selectedIndex].value
      window.location.replace(path)
    }
  </script>
</body>

</html>
`
