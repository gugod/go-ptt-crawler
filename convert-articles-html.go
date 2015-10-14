package main

import (
	"os"
	"log"
	"strings"
	"encoding/json"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

type PushLine struct {
	Tag        string `json:"tag"`
	UserId     string `json:"userid"`
	Content    string `json:"content"`
	IpDateTime string `json:"ipdatetime"`
}

type MetaLine []string

type PttArticle struct {
	Body string     `json:"body"`
	Meta []MetaLine `json:"meta"`
	Push []PushLine `json:"push"`
}

func convert_article_html(path string)  PttArticle {
	var e error
	var fr *os.File
	var doc *goquery.Document
	var body string
	var pushes []PushLine
	var meta []MetaLine

	if fr, e = os.Open(path); e != nil {
		log.Fatal(e)
	}
	if doc, e = goquery.NewDocumentFromReader(fr); e != nil {
		log.Fatal(e)
	}

	main_content_dom := doc.Find("#main-content")
	main_content_dom.Find(".article-metaline").Each(
		func(i int, s *goquery.Selection) {
			var p = MetaLine{ s.Find(".article-meta-tag").Text(), s.Find(".article-meta-value").Text() }
			meta = append(meta, p)
		});

	main_content_dom.Find("div.push").Each(
		func(i int, s *goquery.Selection) {
			pushes = append(pushes, PushLine{
				Tag: s.Find(".push-tag").Text(),
				UserId: s.Find(".push-userid").Text(),
				Content: s.Find(".push-content").Text(),
				IpDateTime: s.Find(".push-ipdatetime").Text(),
			})
		});

	body = doc.Find("#main-content").Text()

	return PttArticle{Body: body, Meta: meta, Push: pushes}
}

func visit(path string, f os.FileInfo, err error) error {
	var output_path string
	var e error
	var fw *os.File

	if ! strings.HasSuffix(path, ".html") { return nil }

	ptt_article := convert_article_html(path)

	var b []byte
	if b, e = json.MarshalIndent(ptt_article, "", "   "); e != nil {
		log.Fatal(e)
		return nil
	}

	output_path = strings.TrimSuffix(path, ".html") + ".json"
	if fw, e = os.Create(output_path); e != nil {
		log.Fatal(e)
		return nil
	}
	log.Printf("%s", output_path)
	fw.Write(b)
	return nil
}

func main() {
	ptt_dir := os.Args[1]

	log.Printf("arg: %s", ptt_dir)
	err := filepath.Walk(ptt_dir, visit)
	log.Printf("filepath.Walk() returned %v\n", err)
}
