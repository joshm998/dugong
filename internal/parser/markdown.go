package parser

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/frontmatter"
)

type MarkdownParser struct{}

func init() {
	registerParser(&MarkdownParser{})
}

func (p *MarkdownParser) Parse(path string) (*Document, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.Linkify,
			extension.TaskList,
			&frontmatter.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	ctx := parser.NewContext()
	var buf bytes.Buffer
	err = md.Convert(content, &buf, parser.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to convert markdown: %w", err)
	}

	fm := frontmatter.Get(ctx)
	attributes := make(map[string]interface{})
	title := ""
	dateStr := ""

	if fm != nil {
		err := fm.Decode(&attributes)
		if err == nil {
			if t, ok := attributes["title"].(string); ok {
				title = t
			}
			if d, ok := attributes["date"].(string); ok {
				dateStr = d
			}
			if tags, ok := attributes["tags"].([]interface{}); ok {
				strTags := make([]string, len(tags))
				for i, tag := range tags {
					strTags[i] = fmt.Sprint(tag)
				}
				attributes["tags"] = strTags
			}
		}
	}

	return &Document{
		Title:      title,
		Content:    buf.String(),
		Attributes: attributes,
		RawDate:    dateStr,
	}, nil
}

func (p *MarkdownParser) Extensions() []string {
	return []string{".md", ".markdown", ".mdoc"}
}
