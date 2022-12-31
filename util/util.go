package util

import (
	"strings"

	"github.com/alecthomas/repr"
	"gopkg.in/yaml.v2"
)

func CloneMap[K comparable, V any](m1 map[K]V) map[K]V {
	m2 := map[K]V{}
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}

func Dedent(s string) string {
	lines := strings.Split(strings.TrimLeft(s, "\n"), "\n")
	repr.Println(lines)

	k := len(lines[0]) - len(strings.TrimLeft(lines[0], "\t "))

	for i, line := range lines {
		if k <= len(line) {
			lines[i] = line[k:]
		} else {
			lines[i] = ""
		}
	}

	return strings.Join(lines, "\n")
}

func ParseYAML(multiline string) map[string]any {
	var m map[string]any
	if err := yaml.Unmarshal([]byte(Dedent(multiline)), &m); err != nil {
		panic(err)
	}

	return m
}
