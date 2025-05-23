package main

import (
	"fmt"
	"strings"

	"github.com/Georgeygigz/go-quizes/parser"
)

func main() {
	s := `<p>Links:</p><ul><li><a href="foo">Foo</a><li><a href="/bar/baz">BarBaz</a></ul>`
	// s = `<a href="/dog">
	// 	<span>Something in a span</span>
	// 	Text not in a span
	// 	<b>Bold text!</b>
	// 	</a>`
	// s = `<a href="#">
	// 	Something here <a href="/dog">nested dog link</a>
	// 	</a>`
	// 	s = `<html>
	// <body>
	//   <h1>Hello!</h1>
	//   <a href="/other-page">A link to another page</a>
	// </body>
	// </html>`

	r := strings.NewReader(s)
	links, err := parser.Parse(r)
	if err != nil {
		fmt.Errorf("Error when parsing HTML file %s", err)
	}
	fmt.Printf("%+v\n", links)
}
