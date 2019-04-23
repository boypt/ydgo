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

	"github.com/BurntSushi/toml"
	"github.com/antonholmquist/jason"
	"github.com/fatih/color"
)

var VERSION string = "0.0-src" //set with ldflags

type Config struct {
	YDAppId  string
	YDAppSec string
}

var debug bool = false
var config Config

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

	errorCode, _ := v.GetString("errorCode")
	if errorCode != "0" {
		fmt.Fprintf(color.Output, color.HiRedString("ErrorCode: %s\nRefer to: http://ai.youdao.com/docs/doc-trans-api.s#p08\n", errorCode))
		return
	}

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
	signstr := fmt.Sprintf("%s%s%d%s", config.YDAppId, query, salt, config.YDAppSec)
	sign := md5.Sum([]byte(signstr))
	uri := fmt.Sprintf(
		"https://openapi.youdao.com/api?appKey=%s&q=%s&from=%s&to=zh-CHS&salt=%d&sign=%x",
		config.YDAppId, url.QueryEscape(query), from, salt, sign)

	if debug {
		log.Printf("Req: %s", uri)
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
	var interative bool
	var configfile, from string
	flag.BoolVar(&interative, "i", false, "interative mode, :q  \\q  EOF or Ctrl+C to exit.")
	flag.BoolVar(&debug, "d", false, "log api request")
	flag.StringVar(&from, "f", "EN", "translate-from language, default: EN")
	flag.StringVar(&configfile, "c", "~/.ydgo", "translate-from language, default: EN")
	flag.Parse()

	configfile = NormalizePath(configfile)
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing: ", configfile)
	}

	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}

	if interative || flag.NArg() == 0 {
		interativeMode(from)
		return
	}

	query := strings.Join(flag.Args(), " ")
	printExplain(query, httpGet(ydAPI(query, from)))
}
