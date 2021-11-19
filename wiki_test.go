package wiki

import (
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/testutil"
)

func TestWiki(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			Wiki,
		),
	)
	count := 0

	count++
	testutil.DoTestCase(markdown, testutil.MarkdownTestCase{
		No:          count,
		Description: "default",
		Markdown: strings.TrimSpace(`
		[[hello|world]]
		`),
		Expected: strings.TrimSpace(`
		<p><a href="hello">world</a></p>
		`),
	}, t)
}
