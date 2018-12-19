package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/fatih/color"
)

var VERSION = "0.0-src" //set with ldflags

const (
	YD_APINAME = "YouDaoCV"
	YD_APIKEY  = "659600698"
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
	color.HiWhite("%s\n\n", query)

	if basic, err := v.GetObject("basic"); err == nil {

		if ph, err := basic.GetString("phonetic"); err == nil {
			color.Yellow("  [%s]\n", ph)
		}
		if us, err := basic.GetString("us-phonetic"); err == nil {
			color.Yellow("  US:[%s]\n", us)
		}
		if uk, err := basic.GetString("uk-phonetic"); err == nil {
			color.Yellow("  UK:[%s]\n", uk)
		}
		fmt.Println()

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

func main() {
	if len(os.Args) < 2 {
		color.HiWhite("%s version: %s", os.Args[0], VERSION)
		color.HiRed("Usage:\n   %s word", os.Args[0])
		return
	}
	query := strings.Join(os.Args[1:], " ")
	url := fmt.Sprintf(
		"http://fanyi.youdao.com/openapi.do?keyfrom=%s&key=%s&type=data&doctype=json&version=1.2&q=%s",
		YD_APINAME,
		YD_APIKEY,
		url.QueryEscape(query))

	PrintExplain(httpGet(url))
}
