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
		return parseBlock(src, "---", yaml.Unmarshal)
	}
	if strings.HasPrefix(src, "+++\n") || src == "+++" {
		return parseBlock(src, "+++", toml.Unmarshal)
	}
	return Result{FrontMatter: FrontMatter{}, Content: src}, nil
}

func parseBlock(src, delim string, unmarshal func([]byte, any) error) (Result, error) {
	rest := src[len(delim)+1:] // 開始デリミタ行（"---\n" 等）を除去
	endMarker := "\n" + delim
	idx := strings.Index(rest, endMarker)
	if idx == -1 {
		return Result{FrontMatter: FrontMatter{}, Content: src}, nil
	}

	rawFM := rest[:idx]
	after := rest[idx+len(endMarker):]
	content := strings.TrimPrefix(after, "\n")

	var fm FrontMatter
	if err := unmarshal([]byte(rawFM), &fm); err != nil {
		return Result{}, err
	}
	if fm == nil {
		fm = FrontMatter{}
	}
	return Result{FrontMatter: fm, Content: content}, nil
}
