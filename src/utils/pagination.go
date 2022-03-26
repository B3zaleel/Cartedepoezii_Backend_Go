package utils

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
)

// Represents an item of a page.
type Item struct {
	id string;
}

// Represents a page of items
type PageSpec struct {
	span int
	after, before string
	popTop bool
}

// Retrieves a section of a list of items with a given span.
func RangeSlicer(span string, items []Item) ([]Item, error) {
	//
	valid_span, err := regexp.MatchString("\\d+(,\\d+)?", span)
	if !valid_span {
		return nil, err
	}
	parts := bytes.Split([]byte(span), []byte(","))
	size, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return nil, err
	}
	idx := 0
	if len(parts) > 1 {
		idx, err = strconv.Atoi(string(parts[1]))
	}
	start := size * idx
	end := start + size
	return items[start:end], err
}

// Extracts a section of a data list based on anchors.
func ExtractPage(items []Item, pageSpec PageSpec) ([]Item, error) {
	start := 0
	end := 0
	if len(pageSpec.after) > 0 && len(pageSpec.before) > 0 {
		return nil, errors.New("Only one page anchor needed.")
	} else if len(pageSpec.after) > 0 {
		for i := 0; i < len(items); i++ {
			if items[i].id == pageSpec.after {
				start = i + 1
				break
			}
		}
		end = start + pageSpec.span
		if start + pageSpec.span >= len(items) {
			end = len(items)
		}
	} else if len(pageSpec.before) > 0 {
		for i := 0; i < len(items); i++ {
			if items[i].id == pageSpec.before {
				end = i + 1
				break
			}
		}
		start = end - pageSpec.span
		if start < 0 {
			start = 0
		}
	} else if pageSpec.popTop {
		start = 0
		end = start + pageSpec.span
		if start + pageSpec.span >= len(items) {
			end = len(items)
		}
	} else {
		end = len(items)
		start = end - pageSpec.span
		if start < 0 {
			start = 0
		}
	}
	return items[start:end], nil
}
