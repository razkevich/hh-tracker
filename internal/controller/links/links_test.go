package links

import (
	"github.com/go-test/deep"
	"net/url"
	"testing"
)

func TestSelfLink(t *testing.T) {
	for name, tc := range map[string]struct {
		host   string
		path   string
		query  url.Values
		expect url.Values
	}{
		"emptyQuery": {
			host:   "foo",
			path:   "bar",
			query:  map[string][]string{},
			expect: map[string][]string{},
		},
		"queryWithOffsetAndLimit": {
			host: "foo",
			path: "bar",
			query: map[string][]string{
				"page[offset]": {"100"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"100"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
		},
		"queryWithOffsetNoLimit": {
			host: "foo",
			path: "bar",
			query: map[string][]string{
				"page[offset]": {"100"},
				"foo":          {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"100"},
				"foo":          {"bar"},
			},
		},
		"queryWithLimitNoOffset": {
			host: "foo",
			path: "bar",
			query: map[string][]string{
				"page[limit]": {"30"},
				"foo":         {"bar"},
			},
			expect: map[string][]string{
				"page[limit]": {"30"},
				"foo":         {"bar"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, got := createSelfLink(tc.query)

			if diff := deep.Equal(got, tc.expect); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestFirstPageLink(t *testing.T) {
	for name, tc := range map[string]struct {
		host   string
		path   string
		query  url.Values
		limit  int64
		expect url.Values
	}{
		"emptyQuery": {
			host:  "foo",
			path:  "bar",
			limit: 0,
			query: map[string][]string{},
			expect: map[string][]string{
				"page[offset]": {"0"},
			},
		},
		"queryWithOffsetAndLimit": {
			host:  "foo",
			path:  "bar",
			limit: 30,
			query: map[string][]string{
				"page[offset]": {"100"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"0"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
		},
		"queryWithOffsetNoLimit": {
			host:  "foo",
			path:  "bar",
			limit: 0,
			query: map[string][]string{
				"page[offset]": {"100"},
				"foo":          {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"0"},
				"foo":          {"bar"},
			},
		},
		"queryWithLimitNoOffset": {
			host:  "foo",
			path:  "bar",
			limit: 30,
			query: map[string][]string{
				"page[limit]": {"30"},
				"foo":         {"bar"},
			},
			expect: map[string][]string{
				"page[limit]":  {"30"},
				"page[offset]": {"0"},
				"foo":          {"bar"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, got := createFirstPageLink(tc.query, tc.limit)

			if diff := deep.Equal(got, tc.expect); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestPreviousPageLink(t *testing.T) {
	for name, tc := range map[string]struct {
		host         string
		path         string
		query        url.Values
		offset       int64
		count        int64
		expectNoLink bool
		limit        int64
		expect       url.Values
	}{
		"emptyQuery": {
			host:         "foo",
			path:         "bar",
			query:        map[string][]string{},
			expectNoLink: true,
		},
		"pagesLeft": {
			host:   "foo",
			path:   "bar",
			offset: 10,
			limit:  6,
			count:  100,
			query: map[string][]string{
				"page[offset]": {"10"},
				"page[limit]":  {"6"},
				"foo":          {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"4"},
				"page[limit]":  {"6"},
				"foo":          {"bar"},
			},
		},
		"offsetButNoLimit": {
			host:   "foo",
			path:   "bar",
			offset: 10,
			count:  100,
			query: map[string][]string{
				"page[offset]": {"10"},
				"foo":          {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"0"},
				"foo":          {"bar"},
			},
		},
		"noPreviousPages": {
			host:         "foo",
			path:         "bar",
			count:        100,
			expectNoLink: true,
			offset:       0,
			limit:        5,
			query: map[string][]string{
				"page[offset]": {"3"},
				"page[limit]":  {"5"},
				"foo":          {"bar"},
			},
		},
		"noOffset": {
			host:         "foo",
			path:         "bar",
			count:        100,
			expectNoLink: true,
			limit:        30,
			query: map[string][]string{
				"page[limit]": {"30"},
				"foo":         {"bar"},
			},
		},
		"noLimit": {
			host:   "foo",
			path:   "bar",
			count:  100,
			offset: 20,
			query: map[string][]string{
				"page[offset]": {"20"},
				"foo":          {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"0"},
				"foo":          {"bar"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			include, got := createPrevPageLink(tc.query, tc.offset, tc.limit, tc.count)

			if tc.expectNoLink && include {
				t.Errorf("No link expected, but one was generated: %v", got)
			} else if !tc.expectNoLink && !include {
				t.Errorf("Link was expected, but none was generated: %v", got)
			} else if diff := deep.Equal(got, tc.expect); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestNextPageLink(t *testing.T) {
	for name, tc := range map[string]struct {
		host         string
		path         string
		query        url.Values
		offset       int64
		limit        int64
		expectNoLink bool
		totalCount   int64
		expect       url.Values
	}{
		"emptyQuery": {
			host:         "foo",
			path:         "bar",
			expectNoLink: true,
			query:        map[string][]string{},
		},
		"pagesLeft": {
			host:       "foo",
			path:       "bar",
			totalCount: 100,
			offset:     10,
			limit:      30,
			query: map[string][]string{
				"page[offset]": {"10"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"40"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
		},
		"noPagesLeft": {
			host:       "foo",
			path:       "bar",
			totalCount: 1000,
			offset:     972,
			limit:      30,
			query: map[string][]string{
				"page[offset]": {"972"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
			expectNoLink: true,
		},
		"noOffset": {
			host:       "foo",
			path:       "bar",
			totalCount: 100,
			limit:      30,
			query: map[string][]string{
				"page[limit]": {"30"},
				"foo":         {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"30"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
		},
		"noLimit": {
			host:       "foo",
			path:       "bar",
			totalCount: 100,
			offset:     20,
			query: map[string][]string{
				"page[offset]": {"20"},
				"foo":          {"bar"},
			},
			expectNoLink: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			include, got := createNextPageLink(tc.query, tc.offset, tc.limit, tc.totalCount)

			if tc.expectNoLink && include {
				t.Errorf("No link expected, but one was generated: %v", got)
			} else if !tc.expectNoLink && !include {
				t.Errorf("Link was expected, but none was generated: %v", got)
			} else if diff := deep.Equal(got, tc.expect); diff != nil {
				t.Error(diff)
			}
		})
	}
}

func TestLastPageLink(t *testing.T) {
	for name, tc := range map[string]struct {
		host       string
		path       string
		query      url.Values
		limit      int64
		totalCount int64
		expect     url.Values
	}{
		"emptyQuery": {
			host:  "foo",
			path:  "bar",
			query: map[string][]string{},
			expect: map[string][]string{
				"page[offset]": {"0"},
			},
		},
		"limitLowerThanTotal": {
			host:       "foo",
			path:       "bar",
			totalCount: 1000,
			limit:      30,
			query: map[string][]string{
				"page[offset]": {"972"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"990"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
		},
		"NoRemainderInDivision": {
			host:       "foo",
			path:       "bar",
			totalCount: 1000,
			limit:      20,
			query: map[string][]string{
				"page[offset]": {"643"},
				"page[limit]":  {"20"},
				"foo":          {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"980"},
				"page[limit]":  {"20"},
				"foo":          {"bar"},
			},
		},
		"limitHigherThanTotal": {
			host:       "foo",
			path:       "bar",
			totalCount: 20,
			limit:      30,
			query: map[string][]string{
				"page[limit]": {"30"},
				"foo":         {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"0"},
				"page[limit]":  {"30"},
				"foo":          {"bar"},
			},
		},
		"noLimit": {
			host: "foo",
			path: "bar",
			query: map[string][]string{
				"foo": {"bar"},
			},
			expect: map[string][]string{
				"page[offset]": {"0"},
				"foo":          {"bar"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, got := createlastPageLink(tc.query, tc.limit, tc.totalCount)

			if diff := deep.Equal(got, tc.expect); diff != nil {
				t.Error(diff)
			}
		})
	}
}
