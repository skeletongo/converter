package main

import (
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var FM = template.FuncMap{
	"incr": func(n int) int {
		return n + 1
	},
	"add": func(n int, v int) int {
		return n + v
	},
	"ptstr": func(k int) string {
		return GetKindToProto(k)
	},
	"gostr": func(k int) string {
		return GetKindToGo(k)
	},
	"now": func() string {
		return time.Now().Format(time.RFC3339)
	},
	"keq": func(k int, s string) bool {
		return GetKindToGo(k) == s
	},
	"separator": func() int32 {
		return filepath.Separator
	},
	"filename": func(s string) string {
		return strings.TrimRight(filepath.Base(s), filepath.Ext(s))
	},
}
