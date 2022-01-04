package main

import (
	"github.com/russross/blackfriday/v2"
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
	blackfriday.Run(f)
	d, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
		return
	}
	fs, err := os.OpenFile(d+"/assets/index.html", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer fs.Close()
	content := template.HTML(blackfriday.Run(f))
	mk := MK{Content: content}

	t, _ := template.ParseFiles(d + "/doc/index.html")
	err = t.Execute(fs, mk)
	if err != nil {
		return
	}
	fs.Sync()
}
