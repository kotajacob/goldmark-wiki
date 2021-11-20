# goldmark-wiki [![godocs.io](https://godocs.io/git.sr.ht/~kota/goldmark-wiki?status.svg)](https://godocs.io/git.sr.ht/~kota/goldmark-wiki) [![builds.sr.ht status](https://builds.sr.ht/~kota/goldmark-wiki.svg)](https://builds.sr.ht/~kota/goldmark-wiki)

Adds wiki style link support to goldmark. Includes a parsers and HTML renderer.
I use this extension's parser in my [gemtext renderer
library](https://git.sr.ht/~kota/goldmark-gemtext).

A wiki style link uses the form `[[click here]]` where "click here" will be both
the link's destination and label. A pipe character `|` can provide a seperate
destination and label: `[[destination page|click here]]`

```go
import (
    "bytes"
    "github.com/yuin/goldmark"
    "github.com/yuin/goldmark/extension"
)

md := goldmark.New(
  goldmark.WithExtensions(wiki.Wiki),
)
var buf bytes.Buffer
if err := md.Convert(source, &buf); err != nil {
    panic(err)
}
```
