package main

import (
	"os"
	"log"
	"regexp"
	// "net/http"
	// "io/ioutil"
	"github.com/PuerkitoBio/goquery"
)

const PTT_URL = "https://www.ptt.cc"

func harvest_board_indices(board_url string, board_name string)  {
	log.Printf("url: %s", board_url)
	log.Printf("name: %s", board_name)

	doc, err := goquery.NewDocument(board_url)
	if err != nil { log.Fatal(err) }

	re := regexp.MustCompile("/bbs/"+board_name+"/index([0-9]+)\\.html")
	doc.Find("a[href^='/bbs/" + board_name + "/index']").Each(func (_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists { return }
		matched := re.FindStringSubmatch(href)
		if len(matched) == 0 { return }
		log.Printf("url: %s => %s", href, matched[1])
	})

	return;
}

func main() {
	board_name := os.Args[1]
	output_dir := os.Args[2]

	board_url := PTT_URL + "/bbs/" + board_name + "/index.html";
	output_board_dir := output_dir + "/" + board_name;

	os.MkdirAll(output_board_dir, os.ModeDir | os.ModePerm);

	harvest_board_indices( board_url, board_name );

}
