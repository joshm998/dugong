package parser

import (
	"bytes"
	"fmt"
	"path/filepath"

	"git.sr.ht/~shulhan/asciidoctor-go"
)

type AsciiDocParser struct{}

func init() {
	registerParser(&AsciiDocParser{})
}

func (p *AsciiDocParser) Parse(path string) (*Document, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	doc, err := asciidoctor.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AsciiDoc: %w", err)
	}

	var htmlBuf bytes.Buffer
	err = doc.ToHTMLEmbedded(&htmlBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to HTML: %w", err)
	}

	attributes := make(map[string]interface{})
	for key, value := range doc.Attributes.Entry {
		attributes[key] = value
	}

	dateStr := ""
	if val, ok := doc.Attributes.Entry["date"]; ok {
		dateStr = val
	}

	return &Document{
		Title:      doc.Title.String(),
		Content:    htmlBuf.String(),
		Attributes: attributes,
		RawDate:    dateStr,
	}, nil
}

func (p *AsciiDocParser) Extensions() []string {
	return []string{".adoc", ".asciidoc", ".asc"}
}
