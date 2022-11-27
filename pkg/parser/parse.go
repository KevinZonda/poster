package parser

import (
	"github.com/KevinZonda/GoX/pkg/iox"
	"github.com/KevinZonda/GoX/pkg/stringx"
	"net/http"
	"strings"
)

type block struct {
	text []string
}

type RequestModel struct {
	Method    string
	Url       string
	Attribute []string
	Payload   string
}

const CommentPrefix = "--"
const TaskSeparator = "----"

func Parse(path string) {
	var bs []block
	lines, _ := iox.ReadAllLines(path)
	isNewBlock := true
	b := block{}
	for _, line := range lines {
		if strings.HasPrefix(line, TaskSeparator) {
			if isNewBlock {
				bs = append(bs, b)
				b = block{}
				isNewBlock = false
			}
			continue
		}
		text := strings.TrimSpace(line)
		if text == "" || strings.HasPrefix(text, CommentPrefix) {
			continue
		}
		b.text = append(b.text, text)
		isNewBlock = true
	}
	if isNewBlock {
		bs = append(bs, b)
	}
	for _, b := range bs {
		parseBlock(b)
	}

	return
}

var methods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
}

func methodAndUrl(line string) (string, string) {
	// GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD
	sep := strings.Index(line, " ")
	if sep == -1 {
		return "", ""
	}
	m := line[:sep]
	u := line[sep+1:]
	_, c := stringx.Contains(methods, m)
	if !c {
		return "", ""
	}
	return m, u
}

func parseBlock(b block) RequestModel {
	var rm RequestModel
	t := b.text

	// Parse Attribute
	// @Test
	for len(t) > 0 && strings.HasPrefix(t[0], "@") {
		rm.Attribute = append(rm.Attribute, t[0])
		t = t[1:]
	}
	if len(t) == 0 {
		return rm
	}

	// Parse Method & Url
	// GET https://api.google.com
	method, url := methodAndUrl(t[0])
	if method == "" {
		return rm
	}
	rm.Method = method
	rm.Url = url

	return rm
}
