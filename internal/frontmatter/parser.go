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
		return parseYAML(src)
	}
	if strings.HasPrefix(src, "+++\n") || src == "+++" {
		return parseTOML(src)
	}
	return Result{FrontMatter: FrontMatter{}, Content: src}, nil
}

func parseYAML(src string) (Result, error) {
	rest := src[4:] // strip leading "---\n"
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return Result{FrontMatter: FrontMatter{}, Content: src}, nil
	}
	rawFM := rest[:idx]
	content := ""
	after := rest[idx+4:] // skip "\n---"
	if len(after) > 0 && after[0] == '\n' {
		content = after[1:]
	} else {
		content = after
	}

	var fm FrontMatter
	if err := yaml.Unmarshal([]byte(rawFM), &fm); err != nil {
		return Result{}, err
	}
	if fm == nil {
		fm = FrontMatter{}
	}
	return Result{FrontMatter: fm, Content: content}, nil
}

func parseTOML(src string) (Result, error) {
	rest := src[4:] // strip leading "+++\n"
	idx := strings.Index(rest, "\n+++")
	if idx == -1 {
		return Result{FrontMatter: FrontMatter{}, Content: src}, nil
	}
	rawFM := rest[:idx]
	content := ""
	after := rest[idx+4:] // skip "\n+++"
	if len(after) > 0 && after[0] == '\n' {
		content = after[1:]
	} else {
		content = after
	}

	var fm FrontMatter
	if err := toml.Unmarshal([]byte(rawFM), &fm); err != nil {
		return Result{}, err
	}
	if fm == nil {
		fm = FrontMatter{}
	}
	return Result{FrontMatter: fm, Content: content}, nil
}
