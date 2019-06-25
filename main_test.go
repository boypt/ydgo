package main

import (
	"reflect"
	"testing"

	"github.com/antonholmquist/jason"
)

func Test_httpGet(t *testing.T) {
	ret := []byte(`{
    "errorCode": "0",
	"query": "good",
	"translation": [
		"好"
	],
	"basic": {
		"phonetic": "gʊd",
		"uk-phonetic": "gʊd",
		"us-phonetic": "ɡʊd",
		"uk-speech": "XXXX",
		"us-speech": "XXXX",
		"explains": [
			"好处",
			"好的",
			"好"
		]
	},
	"web": [
		{
			"key": "good",
			"value": [
				"良好",
				"善",
				"美好"
			]
		}
	],
	"dict": {
		"url": "yddict://m.youdao.com/dict?le=eng&q=good"
	},
	"webdict": {
		"url": "http://m.youdao.com/dict?le=eng&q=good"
	},
	"l": "EN2zh-CHS",
	"tSpeakUrl": "XXX",
	"speakUrl": "XXX"
}`)
	obj, _ := jason.NewObjectFromBytes(ret)
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want *jason.Object
	}{
		// TODO: Add test cases.
		{"mock", args{"http://www.mocky.io/v2/5d11ba66310000913e08cdcd"}, obj},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := httpGet(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printExplain(t *testing.T) {
	ret := []byte(`{
    "errorCode": "0",
	"query": "good",
	"translation": [
		"好"
	],
	"basic": {
		"phonetic": "gʊd",
		"uk-phonetic": "gʊd",
		"us-phonetic": "ɡʊd",
		"uk-speech": "XXXX",
		"us-speech": "XXXX",
		"explains": [
			"好处",
			"好的",
			"好"
		]
	},
	"web": [
		{
			"key": "good",
			"value": [
				"良好",
				"善",
				"美好"
			]
		}
	],
	"dict": {
		"url": "yddict://m.youdao.com/dict?le=eng&q=good"
	},
	"webdict": {
		"url": "http://m.youdao.com/dict?le=eng&q=good"
	},
	"l": "EN2zh-CHS",
	"tSpeakUrl": "XXX",
	"speakUrl": "XXX"
}`)
	obj, _ := jason.NewObjectFromBytes(ret)
	type args struct {
		q string
		v *jason.Object
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"good", args{"good", obj}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printExplain(tt.args.q, tt.args.v)
		})
	}
}
