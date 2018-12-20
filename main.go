package main

import (
	"crypto/md5"
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

var VERSION = "0.0-src" //set with ldflags

const (
	YDAPPKEY = "1d9b4cc7c9694745"
	YDSECKEY = "U9IEK5Qc4CMuWGvbsrwBXaeO6KO7xZwJ"
)

func httpGet(url string) *jason.Object {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		obj, _ := jason.NewObjectFromReader(resp.Body)
		return obj
	} else {
		log.Fatalf("HTTP Non 200: %v", resp)
	}
	return nil
}

func PrintExplain(v *jason.Object) {

	query, _ := v.GetString("query")
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

func ydApi(query string) string {
	salt := rand.Int31()

	// assume query is in utf-8
	signstr := fmt.Sprintf("%s%s%d%s", YDAPPKEY, query, salt, YDSECKEY)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(signstr)))
	yd_api := fmt.Sprintf("https://openapi.youdao.com/api?appKey=%s&q=%s&from=auto&to=zh-CHS&salt=%d&sign=%s",
		YDAPPKEY, url.QueryEscape(query), salt, sign)

	return yd_api
}

func main() {
	if len(os.Args) < 2 {
		color.HiWhite("%s version: %s", os.Args[0], VERSION)
		color.HiRed("Usage:\n   %s word", os.Args[0])
		return
	}
	query := strings.Join(os.Args[1:], " ")
	PrintExplain(httpGet(ydApi(query)))
}
