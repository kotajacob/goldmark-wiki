// Package ast defines a wiki link AST node to represent the wiki extension's
// link element.
package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

type Wiki struct {
	gast.BaseInline

	// Destination is a destination(URL) of this link.
	Destination []byte
}

// Dump implements Node.Dump.
func (n *Wiki) Dump(source []byte, level int) {
	m := map[string]string{}
	m["Destination"] = string(n.Destination)
	gast.DumpHelper(n, source, level, m, nil)
}

// KindWiki is a NodeKind of the Wiki node.
var KindWiki = gast.NewNodeKind("Wiki")

// Kind implements Node.Kind.
func (n *Wiki) Kind() gast.NodeKind {
	return KindWiki
}

// NewWiki returns a new Wiki node.
func NewWiki(dest []byte) *Wiki {
	c := &Wiki{
		BaseInline:  gast.BaseInline{},
		Destination: dest,
	}
	return c
}
