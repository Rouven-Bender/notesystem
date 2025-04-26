package main

import (
	"os"
	"errors"
	"io"
	"reflect"
	"bytes"
)

var (
	ErrInvalidHeader error = errors.New("invalid header")
	ErrNoteWithoutBodyContent error = errors.New("note doesn't have body content")
)

func ParseFile(filename string) (*Note, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return Parse(f)
}

func Parse(reader io.Reader) (*Note, error) {
	n := &Note{}
	noteBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	
	if !reflect.DeepEqual(noteBytes[:4], []byte{'-','-','-','\n'}) {
		return nil, ErrInvalidHeader
	}

	header, content, found := bytes.Cut(noteBytes[4:], []byte{'-','-','-','\n'})
	if !found {
		_, _, found := bytes.Cut(noteBytes[4:], []byte{'-','-','-'})
		if found {
			return nil, ErrNoteWithoutBodyContent
		} else {
			return nil, ErrInvalidHeader
		}
	}
	n.NoteText = string(content)
	props, err := parseHeader(header)
	if err != nil {
		return nil, err
	}
	for _, p := range props {
		n.Header.parseProperty(&p)
	}

	return n, nil // CHANGE LATER
}

func parseHeader(header []byte) ([]property,error) {
	var (
		idxKeyStart int = -1
		idxKeyEnd int = -1
		idxValueStart int = -1
		idxValueEnd int = -1
		value bool
		cursor int = 0
		props []property = []property{}
	)
permareader:
	for {
		if !value {
			if idxKeyStart == -1 {
				idxKeyStart = cursor
			}
			switch header[cursor] {
			case ':':
				idxKeyEnd = cursor
				value = true
			case '\n':
				return nil, ErrInvalidHeader 
			default:
			}
			// advance cursor
			if cursor+1 != len(header) {
				cursor++
			} else {
				return nil, ErrInvalidHeader
			}
		} else {
			if idxValueStart == -1 {
				idxValueStart = cursor
			}
			switch header[cursor] {
			case '\n':
				if cursor+1 == len(header) || !reflect.DeepEqual(header[cursor+1:cursor+1+3], []byte{'\t', '-', ' '}) {
					idxValueEnd = cursor
					p := property{
						name: string(header[idxKeyStart:idxKeyEnd]),
						value: string(header[idxValueStart:idxValueEnd]),
					}
					props = append(props, p)
					idxKeyStart = -1
					idxKeyEnd = -1
					idxValueStart = -1
					idxValueEnd = -1
					value = false
				}
			}
			// advance cursor
			if cursor+1 != len(header) {
				cursor++
			} else {
				break permareader
			}
		}
	}
	return props, nil
}
