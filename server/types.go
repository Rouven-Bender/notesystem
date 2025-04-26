package main

import (
	"strings"
	"regexp"
)

var notelink *regexp.Regexp = regexp.MustCompile(`^\[\[(?P<filename>.*)\]\]$`)

type Note struct {
	Header Metadata
	NoteText string
}

type Metadata struct {
	Title string
	Tags []string
	URLSources []string
	LinkSourceNote []string
}

func (m *Metadata) parseProperty(p *property) {
	switch p.name {
	case "Title":
		m.Title = p.StringValue()
	case "Tags":
		m.Tags = p.TagsValue()
	case "Source":
		m.URLSources, m.LinkSourceNote = p.SourceValue()
	default:
	}
}

type property struct {
	name string
	value string
}

func (p *property) StringValue() string {
	v, _ := strings.CutPrefix(p.value, " ")
	return v
}

// sourceNote is a list of filenames beeing linked to in the source property
func (p *property) SourceValue() (urls []string, sourceNote []string) {
	splitted := strings.Split(p.value, "\n")
	for _, s := range splitted {
		if s == "" {
			continue
		}
		s, found := strings.CutPrefix(s, "\t-")
		if !found {
			continue
		}
		s, _ = strings.CutPrefix(s, " ")

		if notelink.MatchString(s) {
			matches := notelink.FindStringSubmatch(s)
			filenameIdx := notelink.SubexpIndex("filename")
			sourceNote = append(sourceNote, matches[filenameIdx])
			continue
		}
		urls = append(urls, s)
	}
	return
}

func (p *property) TagsValue() []string {
	initial := p.StringValue()
	splitted := strings.Split(initial, " ")
	tags := []string{}
	for _, s := range splitted {
		if s == "" {
			continue
		}
		if !strings.HasPrefix(s, "#") {
			continue
		}
		tags = append(tags, s)
	}
	return tags
}
