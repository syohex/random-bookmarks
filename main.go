package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

var USERS = []string{
	"gfx",
	"mattn",
	"moznion",
	"mizchi",
}

type RSS struct {
	//Channel Channel `xml:"channel"`
	Items []Item `xml:"item"`
}

type TmplArg struct {
	User  string
	Items []Item
}

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

func rssURL(user string) string {
	return fmt.Sprintf("http://b.hatena.ne.jp/%s/rss", user)
}

var indexTmpl = `
<!doctype html>
<html>
<head>
<meta charset="utf-8">
<title>Random Bookmarks</title>
</head>
<body>
<h1>{{.User}} Bookmarks</h1>
<table>
  <thead>
    <th>ID</th><th>Title</th>
  </thead>
  <tbody>
{{range $i, $item := .Items}}
  <tr>
    <td>{{$i}}</td><td><a href="{{$item.Link}}">{{$item.Title}}</a></td>
  </tr>
{{end}}
  </tbody>
</body>
</html>
`

func main() {
	rd := rand.New(rand.NewSource(time.Now().Unix()))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		n := rd.Intn(len(USERS))

		user := USERS[n]
		url := rssURL(user)
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		decoder := xml.NewDecoder(resp.Body)

		rss := new(RSS)
		if err := decoder.Decode(rss); err != nil {
			log.Fatalln(err)
		}

		t := template.New("")
		tt, err := t.Parse(indexTmpl)
		if err != nil {
			log.Fatalln(err)
		}

		w.Header().Set("Content-Type", "text/html")
		tt.Execute(w, &TmplArg{
			User:  user,
			Items: rss.Items,
		})
	})
	http.ListenAndServe(":5000", nil)
}
