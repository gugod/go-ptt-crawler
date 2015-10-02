package main

import (
	"os"
	"log"
	"regexp"
	"strconv"
	// "net/http"
	// "io/ioutil"
	"github.com/PuerkitoBio/goquery"
)

const PTT_URL = "https://www.ptt.cc"

type BoardIndexPage struct {
	page_number int
	url string
}

func harvest_board_indices(board_url string, board_name string)  []BoardIndexPage {
	var ret []BoardIndexPage

	doc, err := goquery.NewDocument(board_url)
	if err != nil { log.Fatal(err) }

	re := regexp.MustCompile("/bbs/"+board_name+"/index([0-9]+)\\.html")
	doc.Find("a[href^='/bbs/" + board_name + "/index']").Each(func (_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists { return }
		matched := re.FindStringSubmatch(href)
		if len(matched) == 0 { return }
		pn, err := strconv.Atoi(matched[1])
		if err != nil { log.Fatal(err) }
		ret = append( ret, BoardIndexPage{ pn, href } )
	})

	if (ret[0].page_number > ret[1].page_number) {
		ret[0],ret[1] = ret[1],ret[0]
	}

	for i := ret[0].page_number + 1; i < ret[1].page_number; i++ {
		ret = append(ret, BoardIndexPage{i, "/bbs/" + board_name + "/index" + strconv.Itoa(i) + ".html" } )
	}

	return ret;
}

func main() {
	board_name := os.Args[1]
	output_dir := os.Args[2]

	board_url := PTT_URL + "/bbs/" + board_name + "/index.html";
	output_board_dir := output_dir + "/" + board_name;

	os.MkdirAll(output_board_dir, os.ModeDir | os.ModePerm);

	board_indices := harvest_board_indices( board_url, board_name );
	log.Print(board_indices);
}
