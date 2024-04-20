package markdown

import "strings"

var replacer = strings.NewReplacer(
	"-",
	"\\-",
	"[",
	"\\[",
	"]",
	"\\]",
	"(",
	"\\(",
	")",
	"\\)",
	"~",
	"\\~",
	">",
	"\\>",
	"#",
	"\\#",
	"+",
	"\\+",
	"=",
	"\\=",
	"|",
	"\\|",
	"{",
	"\\{",
	"}",
	"\\}",
	".",
	"\\.",
	"!",
	"\\!",
)

func EscapeForMarkdown(src string) string {
	return replacer.Replace(src)
}
