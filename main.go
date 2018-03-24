// Parses a set of html files and a css tachyons file
// and returns the css for the classes used in the html files
// It takes a list of html files from the stdIn and a single css file as argument
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var file string
var htmlFiles []string
var used []string

func main() {
	// read stdin for the list of html files to be parsed
	c := 0
	for {
		_, err := fmt.Scan(&file)
		if err != nil {
			htmlFiles = htmlFiles[:c]
			break
		}
		htmlFiles = append(htmlFiles, file)
		c += 1
	}

	//extract class names from html files, put it in used !used contains duplicates
	for i := range htmlFiles {
		file, err := os.Open(htmlFiles[i])
		if err != nil {
			fmt.Println("error opening file", err)
			continue
		}
		doc, err := goquery.NewDocumentFromReader(file)
		docClasses := doc.Find("[class]")
		for n := range docClasses.Nodes {
			cl, _ := docClasses.Eq(n).Attr("class")
			sp := strings.Split(cl, " ")
			// nasty hack
			var pseudo []string
			for c := range sp {
				pseudo = append(pseudo, sp[c]+":active")
				pseudo = append(pseudo, sp[c]+":after")
				pseudo = append(pseudo, sp[c]+":hover")
				pseudo = append(pseudo, sp[c]+":link")
				pseudo = append(pseudo, sp[c]+":focus")
				pseudo = append(pseudo, sp[c]+":visited")
				pseudo = append(pseudo, sp[c]+":before")
				pseudo = append(pseudo, sp[c]+":nth-child")
			}
			sp = append(sp, pseudo...)
			used = append(used, sp...)
		}
	}

	// parse css
	if len(os.Args) < 2 {
		fmt.Println("I need a path to an unminified tachyons file as an argument")
		return
	}

	css, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Error opening css file", err)
		return
	}
	parsed := make(map[string]bool)
	var standard []string
	// the css within the media query should be put in the respective list
	nsList := []string{`@media screen and (min-width: 30em) {`}
	mList := []string{`@media screen and (min-width: 30em) and (max-width: 60em) {`}
	lList := []string{`@media screen and (min-width: 60em) {`}
	s := bufio.NewScanner(css)
	for s.Scan() {
		for _, v := range used {
			if _, ok := parsed[v]; ok {
				continue
			}
			if bytes.Contains(s.Bytes(), []byte("."+v)) {
				if ok, _ := regexp.Match(`\.`+v+` |\.`+v+`:`, s.Bytes()); ok {
					parsed[v] = true
					if bytes.HasPrefix(s.Bytes(), []byte(" ")) {
						switch {
						case strings.HasSuffix(v, "-l"):
							lList = append(lList, s.Text())
						case strings.HasSuffix(v, "-m"):
							mList = append(mList, s.Text())
						case strings.HasSuffix(v, "-ns"):
							nsList = append(nsList, s.Text())
						}
					} else {
						standard = append(standard, s.Text())
					}
				}
				continue
			}
		}
	}
	nsList = append(nsList, "}")
	mList = append(mList, "}")
	lList = append(lList, "}")

	// print the used classes, and nested
	for _, v := range standard {
		fmt.Println(v)
	}
	fmt.Println(nested)
	// print responsive if they contain something
	hasClass := func(s []string) bool {
		if len(s) > 2 {
			return true
		}
		return false
	}
	if hasClass(nsList) {
		fmt.Println(strings.Join(nsList, "\n"))
	}
	if hasClass(mList) {
		fmt.Println(strings.Join(mList, "\n"))
	}
	if hasClass(lList) {
		fmt.Println(strings.Join(lList, "\n"))
	}
}

// nested is always included
const nested = `.nested-copy-line-height p, .nested-copy-line-height ul,
.nested-copy-line-height ol { line-height: 1.5; }
.nested-headline-line-height h1, .nested-headline-line-height h2,
.nested-headline-line-height h3, .nested-headline-line-height h4,
.nested-headline-line-height h5, .nested-headline-line-height h6 { line-height: 1.25; }
.nested-list-reset ul, .nested-list-reset ol { padding-left: 0; margin-left: 0; list-style-type: none; }
.nested-copy-indent p+p { text-indent: 1em; margin-top: 0; margin-bottom: 0; }
.nested-copy-separator p+p { margin-top: 1.5em; }
.nested-img img { width: 100%; max-width: 100%; display: block; }
.nested-links a { color: #357edd; transition: color .15s ease-in; }
.nested-links a:hover { color: #96ccff; transition: color .15s ease-in; }
.nested-links a:focus { color: #96ccff; transition: color .15s ease-in; }`
