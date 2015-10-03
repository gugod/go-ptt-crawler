package main

import (
	"os"
	"io"
	"log"
	"regexp"
	"strconv"
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

const PTT_URL = "https://www.ptt.cc"

type BoardIndexPage struct {
	page_number int
	url string
}

type ArticlePage struct {
	id  string
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

func harvest_articles(url string, board_name string) []ArticlePage {
	var ret []ArticlePage

	doc, err := goquery.NewDocument(url)
	if err != nil { log.Fatal(err) }

	re := regexp.MustCompile("/(M\\.[0-9]+\\.A\\.[A-Z0-9]{3})\\.html")
	doc.Find("a[href*='/bbs/" + board_name + "/']").Each(func (_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists { return }
		matched := re.FindStringSubmatch(href)
		if len(matched) == 0 { return }
		ret = append(ret, ArticlePage{ matched[1], href })
	})

	return ret
}

func download_articles(articles []ArticlePage, output_board_dir string)  {
	for _, article := range articles {
		output_file := output_board_dir + "/" + article.id + ".html"

		output, err := os.Create(output_file)
		if err != nil { log.Fatal("Error while creating", output_file, "-", err)  }
		defer output.Close()

		res, err := http.Get( PTT_URL + article.url)
		if err != nil { log.Fatal(err) }
		defer res.Body.Close()

		_, err = io.Copy(output, res.Body)
		if err != nil {
			log.Fatal("Error while downloading", article.url, "-", err)
		}
		log.Println(output_file)
	}

}

func main() {
	board_name := os.Args[1]
	output_dir := os.Args[2]

	board_url := PTT_URL + "/bbs/" + board_name + "/index.html";
	output_board_dir := output_dir + "/" + board_name;

	os.MkdirAll(output_board_dir, os.ModeDir | os.ModePerm)

	board_indices := harvest_board_indices( board_url, board_name )
	for _,board := range board_indices {
		articles := harvest_articles( PTT_URL + board.url, board_name )
		download_articles( articles, output_board_dir )
	}
}
