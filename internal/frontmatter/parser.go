package frontmatter

import (
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type FrontMatter = map[string]any

type Result struct {
	FrontMatter FrontMatter
	Content     string
}

func Parse(src string) (Result, error) {
	if strings.HasPrefix(src, "---\n") || src == "---" {
		return parseBlock(src, "---", func(b []byte, v any) error { return yaml.Unmarshal(b, v) })
	}
	if strings.HasPrefix(src, "+++\n") || src == "+++" {
		return parseBlock(src, "+++", func(b []byte, v any) error { return toml.Unmarshal(b, v) })
	}
	return Result{FrontMatter: FrontMatter{}, Content: src}, nil
}

func parseBlock(src, delim string, unmarshal func([]byte, any) error) (Result, error) {
	rest := src[len(delim)+1:] // デリミタ行と改行を除去
	idx := strings.Index(rest, "\n"+delim)
	if idx == -1 {
		return Result{FrontMatter: FrontMatter{}, Content: src}, nil
	}

	rawFM := rest[:idx]
	after := rest[idx+len(delim)+1:] // "\n" + delim を除去
	content := ""
	if len(after) > 0 && after[0] == '\n' {
		content = after[1:]
	} else {
		content = after
	}

	var fm FrontMatter
	if err := unmarshal([]byte(rawFM), &fm); err != nil {
		return Result{}, err
	}
	if fm == nil {
		fm = FrontMatter{}
	}
	return Result{FrontMatter: fm, Content: content}, nil
}
