// Package wiki is an extension for goldmark.
// https://github.com/yuin/goldmark
package wiki

import (
	"bytes"

	"git.sr.ht/~kota/goldmark-wiki/ast"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type wiki struct{}

// Wiki is a goldmark.Extender implementation.
var Wiki = &wiki{}

// New returns a new extension. Useless, but included for compatibility.
func New() goldmark.Extender {
	return &wiki{}
}

// Extend implements goldmark.Extender.
func (w *wiki) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewParser(), 199),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewHTMLRenderer(), 199),
	))
}

type wikiParser struct{}

// NewParser returns a new parser.InlineParser that can parse wiki link syntax.
func NewParser() parser.InlineParser {
	p := &wikiParser{}
	return p
}

// Trigger returns characters that trigger this parser.
func (p *wikiParser) Trigger() []byte {
	return []byte{'['}
}

var (
	parseOpen  = []byte("[[")
	parsePipe  = []byte{'|'}
	parseClose = []byte("]]")
)

// Parse a wiki style link using the form: [[click here]]
// "click here" will be both the link destination and label.
// A pipe character "|" may be used after the link destination to specify a
// different label: [[destination page|click here]]
func (p *wikiParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, seg := block.PeekLine()
	if !bytes.HasPrefix(line, parseOpen) {
		return nil
	}

	stop := bytes.Index(line, parseClose)
	if stop < 0 {
		return nil // Link must close on the same line.
	}
	seg = text.NewSegment(seg.Start+2, seg.Start+stop)

	n := ast.NewWiki(block.Value(seg))
	if idx := bytes.Index(n.Destination, parsePipe); idx >= 0 {
		// If there's text after a pipe symbol.
		n.Destination = n.Destination[:idx]      // Pre pipe.
		seg = seg.WithStart(seg.Start + idx + 1) // Post pipe.
	}

	if len(n.Destination) == 0 || seg.Len() == 0 {
		return nil // Ensure destination and label are not empty.
	}

	n.AppendChild(n, gast.NewTextSegment(seg))
	block.Advance(stop + 2)
	return n
}

type wikiHTMLRenderer struct{}

// NewHTMLRenderer returns a new HTMLRenderer.
func NewHTMLRenderer() renderer.NodeRenderer {
	r := &wikiHTMLRenderer{}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *wikiHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindWiki, r.renderWiki)
}

func (r *wikiHTMLRenderer) renderWiki(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*ast.Wiki)
	if entering {
		_, _ = w.WriteString("<a href=\"")
		// NOTE: URLs are not fully checked for safety!
		_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		_ = w.WriteByte('"')
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, html.LinkAttributeFilter)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</a>")
	}
	return gast.WalkContinue, nil
}
