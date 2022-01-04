package main

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

type MK struct {
	//Content string
	Content template.HTML
}

func main() {
	f, err := ioutil.ReadFile("./docs/doc.md")
	if err != nil {
		log.Println(err.Error())
		return
	}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(f, &buf); err != nil {
		panic(err)
	}
	d, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
		return
	}
	fs, err := os.OpenFile(d+"/assets/index.html", os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer fs.Close()
	mk := MK{Content: template.HTML(buf.String())}

	t, _ := template.ParseFiles(d + "/docs/index.html")
	err = t.Execute(fs, mk)
	if err != nil {
		return
	}
	fs.Sync()
}
