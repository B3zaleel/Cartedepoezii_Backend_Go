package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Represents an item of a page.
type Item struct {
	id string;
}

// Represents a page of items
type PageSpec struct {
	Span int
	After, Before string
	PopTop bool
}

// Retrieves a section of a list of items with a given span.
func RangeSlicer(span string, items []Item) ([]Item, error) {
	valid_span, err := regexp.MatchString("\\d+(,\\d+)?", span)
	if !valid_span {
		return nil, err
	}
	parts := strings.Split(span, ",")
	size, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	idx := 0
	if len(parts) > 1 {
		idx, err = strconv.Atoi(parts[1])
	}
	start := size * idx
	end := start + size
	return items[start:end], err
}

// Extracts a section of a data list based on anchors.
func ExtractPage(items []Item, pageSpec PageSpec) ([]Item, error) {
	start := 0
	end := 0
	if len(pageSpec.After) > 0 && len(pageSpec.Before) > 0 {
		return nil, errors.New("Only one page anchor needed.")
	} else if len(pageSpec.After) > 0 {
		for i := 0; i < len(items); i++ {
			if items[i].id == pageSpec.After {
				start = i + 1
				break
			}
		}
		end = start + pageSpec.Span
		if start + pageSpec.Span >= len(items) {
			end = len(items)
		}
	} else if len(pageSpec.Before) > 0 {
		for i := 0; i < len(items); i++ {
			if items[i].id == pageSpec.Before {
				end = i + 1
				break
			}
		}
		start = end - pageSpec.Span
		if start < 0 {
			start = 0
		}
	} else if pageSpec.PopTop {
		start = 0
		end = start + pageSpec.Span
		if start + pageSpec.Span >= len(items) {
			end = len(items)
		}
	} else {
		end = len(items)
		start = end - pageSpec.Span
		if start < 0 {
			start = 0
		}
	}
	return items[start:end], nil
}

// Retrieves a page spec from a gin Context.
func GetPageSpec(c *gin.Context, popTop bool) (page *PageSpec, err error) {
	after := c.DefaultQuery("after", "")
	before := c.DefaultQuery("before", "")
	spanStr := c.DefaultQuery("span", "12")
	span, err := strconv.ParseUint(spanStr, 10, 32)
	if err != nil {
		return nil, err
	}
	page = &PageSpec{
		After: after,
		Before: before,
		Span: int(span),
		PopTop: popTop,
	}
	return page, nil
}
