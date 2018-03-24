package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/anaskhan96/soup"
)

func checkfile(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fetchlink(url string) ([]string, string) {
	var alllink, txt []string
	resq, err := soup.Get(url)
	if err != nil {
		fmt.Println("download failed " + url)
		fmt.Println(err)
		return nil, "nil"
	}
	fmt.Println("download "+url+" succeed")
	html := soup.HTMLParse(resq)
	bookName := html.Find("div", "id", "info").Find("h1").Text()
	links := html.Find("div", "id", "list").FindAll("a")
	for _, link := range links {
		txt = strings.Split(link.Attrs()["href"], "/")
		if len(txt) > 0 {
			alllink = append(alllink, strings.Trim(txt[len(txt)-1], ".html"))
		}
	}
	fmt.Println(bookName)
	return alllink, bookName
}

func duplicate(dustr []string) []string {
	var duped []string
	seen := make(map[string]bool)
	for _, str := range dustr {
		if !seen[str] {
			seen[str] = true
		}
	}
	for str := range seen {
		duped = append(duped, str)
	}
	return duped
}

func fetchnovel(baseurl, dir string, url string, threads chan bool, num chan bool) {
	var f *os.File
	var err error
	var fname, link string
	fname = dir + "/" + url + ".txt"
	link = baseurl + url + ".html"
	defer func() {
		f.Close()
		<-threads
		num <- true
	}()
	threads <- true
	if checkfile(fname) {
		f, err = os.OpenFile(fname, os.O_WRONLY, 0666)
		fmt.Println("文件存在")
		check(err)
		return
	}
	f, err = os.Create(fname)
	fmt.Println("创建文件")
	check(err)
	resq, err := soup.Get(link)
	resq = strings.Replace(resq, "<br />", "\n", -1)
	check(err)
	html := soup.HTMLParse(resq)
	bookName := html.Find("div", "class", "bookname").Find("h1")
	_, err = fmt.Fprintf(f, "%s \n", bookName.Text())
	check(err)
	doc := html.FindAll("div", "id", "content")
	for _, txt := range doc {
		_, err = fmt.Fprintf(f, "%s \n", txt.Text())
		check(err)
	}
	return
}

func createDir(dir string) {
	err := os.MkdirAll(dir, 0666)
	if err != nil {
		fmt.Println("create directory wrong")
		panic(err)
	} else {
		fmt.Println("directory create sucessed")
	}

}

func main() {
	thread := make(chan bool, 20)
	num := make(chan bool)

	for _, url := range os.Args[1:] {
		alllink, name := fetchlink(url)
		if alllink != nil {
			duplicatealllink := duplicate(alllink)
			fmt.Println(duplicatealllink)
			i := 1
			dir := "./output/" + name
			createDir(dir)
			for _, duurl := range duplicatealllink {
				go fetchnovel(url, dir, duurl, thread, num)
				i++
			}
			for j := 1; j < i; j++ {
				<-num
				fmt.Println(i, j)
			}
			merge(dir, dir+".txt")

		}

	}
}

func merge(path string, file string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fmt.Println(path, len(path))
		fmt.Println(f.Size())
		writeTo(path, file)
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}
func writeTo(readFile, writeFile string) {
	var f *os.File
	var err1 error
	buf, err := ioutil.ReadFile(readFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file error: %s\n", err)
	}
	if checkfile(writeFile) {
		f, err1 = os.OpenFile(writeFile, os.O_APPEND|os.O_WRONLY, 0666)

	} else {
		f, err1 = os.Create(writeFile)
	}
	if err1 != nil {
		fmt.Println("write failed")
		panic(err1)
	}
	f.Write(buf)
	defer f.Close()
	fmt.Println("shuaizhao")
}
