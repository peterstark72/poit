package poit

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

//TrimInnerText returns inner text trimmed of space
func TrimInnerText(el *html.Node) string {
	return strings.TrimSpace(htmlquery.InnerText(el))
}
