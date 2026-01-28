package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Parser interface {
	Parse(path string) (*Document, error)
	Extensions() []string
}

type Document struct {
	Title      string
	Content    string
	Attributes map[string]interface{}
	RawDate    string
}

var parserRegistry []Parser
var parserByExtension = make(map[string]Parser)

func registerParser(p Parser) {
	parserRegistry = append(parserRegistry, p)
	for _, ext := range p.Extensions() {
		parserByExtension[ext] = p
	}
}

func NewParser(path string) (Parser, error) {
	ext := strings.ToLower(filepath.Ext(path))

	if parser, ok := parserByExtension[ext]; ok {
		return parser, nil
	}

	return nil, fmt.Errorf("unsupported file extension: %s", ext)
}

func IsContentFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := parserByExtension[ext]
	return ok
}

func ParseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"January 2, 2006",
		"Jan 2, 2006",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func ExtractDate(doc *Document, sourcePath string) time.Time {
	if doc.RawDate != "" {
		if date, err := ParseDate(doc.RawDate); err == nil {
			return date
		}
	}

	if info, err := os.Stat(sourcePath); err == nil {
		return info.ModTime()
	}

	return time.Now()
}
