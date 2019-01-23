package main

import (
	"bufio"
	"crypto/md5"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/fatih/color"
)

var VERSION string = "0.0-src" //set with ldflags
const (
	YDAPPKEY = "1d9b4cc7c9694745"
	YDSECKEY = "U9IEK5Qc4CMuWGvbsrwBXaeO6KO7xZwJ"
)

var debug bool = false

func httpGet(url string) *jason.Object {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		log.Fatalf("Fail to request api: %#v", err)
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		obj, err := jason.NewObjectFromReader(resp.Body)
		if err != nil {
			log.Fatalf("Fail to parse json: %#v", err)
			return nil
		}
		return obj
	}
	log.Fatalf("HTTP Non 200: %#v", resp)
	return nil
}

func printExplain(q string, v *jason.Object) {

	query, err := v.GetString("query")
	if err != nil {
		query = q
	}
	fmt.Fprintf(color.Output, color.HiWhiteString("%s    ", query))

	if basic, err := v.GetObject("basic"); err == nil {

		var phonetis []string
		if ph, err := basic.GetString("phonetic"); err == nil {
			phonetis = append(phonetis, fmt.Sprintf("[%s]", ph))
		}
		if ph, err := basic.GetString("us-phonetic"); err == nil {
			phonetis = append(phonetis, fmt.Sprintf("US:[%s]", ph))
		}
		if ph, err := basic.GetString("uk-phonetic"); err == nil {
			phonetis = append(phonetis, fmt.Sprintf("UK:[%s]", ph))
		}
		if len(phonetis) > 0 {
			color.Yellow(strings.Join(phonetis, " "))
			fmt.Println()
		}

		if expl, err := basic.GetStringArray("explains"); err == nil {
			color.Cyan("  Word Explanation:\n")
			for _, ex := range expl {
				color.HiWhite("      * %s", ex)
			}
			fmt.Println()
		}
	} else {
		fmt.Println()
	}

	if web, err := v.GetObjectArray("web"); err == nil {
		color.Cyan("  Web Reference:\n")
		for _, wres := range web {
			key, _ := wres.GetString("key")
			val, _ := wres.GetStringArray("value")
			color.HiCyan("    %s\n", key)
			for _, sval := range val {
				color.HiWhite("      * %s\n", sval)
			}
			fmt.Println()
		}

	}

	if transl, err := v.GetStringArray("translation"); err == nil {
		color.Cyan("  Translation:\n")
		for _, ex := range transl {
			color.Cyan("      * %s", ex)
		}
	}
}

func ydAPI(query string, from string) string {
	salt := rand.Int31()

	// assume query is in utf-8
	signstr := fmt.Sprintf("%s%s%d%s", YDAPPKEY, query, salt, YDSECKEY)
	sign := md5.Sum([]byte(signstr))
	uri := fmt.Sprintf(
		"https://openapi.youdao.com/api?appKey=%s&q=%s&from=%s&to=zh-CHS&salt=%d&sign=%x",
		YDAPPKEY, url.QueryEscape(query), from, salt, sign)

	if debug {
		log.Println(uri)
	}
	return uri
}

func interativeMode(from string) {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		query := scanner.Text()
		query = strings.TrimSpace(query)

		if len(query) == 0 {
			fmt.Print("\n> ")
			continue
		}

		if query == ":q" || query == "\\q" {
			break
		}

		printExplain(query, httpGet(ydAPI(query, from)))
		fmt.Print("> ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}

func showHelp() {
	color.HiWhite("ydgo version: %s\n", VERSION)
	color.HiRed("Usage:\n   %s [-i] word\n\n", os.Args[0])
	flag.Usage()
}

func main() {
	var interative, help bool
	var from string
	flag.BoolVar(&interative, "i", false, "interative mode, :q  \\q  EOF or Ctrl+C to exit.")
	flag.BoolVar(&help, "h", false, "show this help")
	flag.BoolVar(&debug, "d", false, "log api request")
	flag.StringVar(&from, "f", "EN", "translate-from language, default: EN")
	flag.Parse()

	if help {
		showHelp()
		return
	}

	if interative {
		interativeMode(from)
		return
	}

	if len(os.Args) < 2 {
		showHelp()
		return
	}

	query := strings.Join(os.Args[1:], " ")
	printExplain(query, httpGet(ydAPI(query, from)))
}
