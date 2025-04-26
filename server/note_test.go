package main

import (
	"testing"
	"log"
	"bytes"
)

func TestParse(t *testing.T) {
	testdata := "---\nTitle: Test note\nTags: #test\nSource:\n\t- https://example.org\n---\nnote content"
	r := bytes.NewBufferString(testdata)

	x, e := Parse(r)
	log.Printf("%v ; %v", x, e)
}
func TestParseNoteWithoutBody(t *testing.T) {
	testdata := "---\nTitle: Test note\nTags: #test\nSource:\n\t- https://example.org\n---"
	r := bytes.NewBufferString(testdata)

	x, e := Parse(r)
	log.Printf("%v ; %v", x, e)
}
func TestParseHeader(t *testing.T) {
	testdata := []byte("Title: Test note\nTags: #test\nSource:\n\t- https://example.org\n")

	prop, err := parseHeader(testdata)
	if err != nil {
		t.Fatal(err)
	}

	for _, p := range prop {
		log.Printf("key:%s\nvalue:%s\n", p.name, p.value)
	}
}
