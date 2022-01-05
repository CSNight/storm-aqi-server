package main

import (
	"bytes"
	toc "github.com/abhinav/goldmark-toc"
	ht "github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

type MK struct {
	//Content string
	Content template.HTML
	Toc     template.HTML
}

func main() {
	f, err := ioutil.ReadFile("./docs/doc.md")
	if err != nil {
		log.Println(err.Error())
		return
	}
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, highlighting.NewHighlighting(
			highlighting.WithStyle("dracula"),
			highlighting.WithFormatOptions(
				ht.WithLineNumbers(true),
			),
		)),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
			parser.WithHeadingAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	if err = md.Convert(f, &buf); err != nil {
		return
	}
	doc := md.Parser().Parse(text.NewReader(f))
	tree, err := toc.Inspect(doc, f)
	if err != nil {
		// handle the error
	}
	var bufToc bytes.Buffer
	list := toc.RenderList(tree)
	err = md.Renderer().Render(&bufToc, f, list)
	if err != nil {
		return
	}
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
		return
	}
	fs, err := os.OpenFile(pwd+"/assets/index.html", os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer fs.Close()
	mk := MK{Content: template.HTML(buf.String()), Toc: template.HTML(bufToc.String())}

	t, _ := template.ParseFiles(pwd + "/docs/index.html")
	err = t.Execute(fs, mk)
	if err != nil {
		return
	}
	err = fs.Sync()
	if err != nil {
		return
	}
}
