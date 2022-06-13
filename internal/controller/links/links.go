package links

import (
	"fmt"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/controller/helper"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CreateSelfLink Create Self Link
func CreateSelfLink(r *http.Request, baseURL string) helper.Links {
	links := helper.Links{}

	if include, self := createSelfLink(r.URL.Query()); include {
		links.Self = createURLFromQuery(r, baseURL, self)
	}

	return links
}

// CreatePaginationLinks create pagination links
func CreatePaginationLinks(r *http.Request, baseURL string, offset int64, limit int64, totalCount int64) helper.Links {
	links := helper.Links{}

	if offset != 0 || limit != 0 {
		if include, self := createCurrentLink(r.URL.Query(), limit, offset); include {
			links.Current = createURLFromQuery(r, baseURL, self)
		}

		if include, first := createFirstPageLink(r.URL.Query(), limit); include {
			links.First = createURLFromQuery(r, baseURL, first)
		}

		if include, next := createNextPageLink(r.URL.Query(), offset, limit, totalCount); include {
			links.Next = createURLFromQuery(r, baseURL, next)
		} else {
			links.Next = "null"
		}

		if include, prev := createPrevPageLink(r.URL.Query(), offset, limit, totalCount); include {
			links.Prev = createURLFromQuery(r, baseURL, prev)
		} else {
			links.Prev = "null"
		}

		if include, last := createlastPageLink(r.URL.Query(), limit, totalCount); include {
			links.Last = createURLFromQuery(r, baseURL, last)
		}
	} else if include, self := createSelfLink(r.URL.Query()); include {
		links.Self = createURLFromQuery(r, baseURL, self)
	}

	return links
}
func createURLFromQuery(r *http.Request, baseURL string, query url.Values) string {
	if len(query) > 0 {
		queryString := valuesToString(query)
		return fmt.Sprintf("%s%s?%s", baseURL, r.URL.Path, queryString)
	}

	return fmt.Sprintf("%s%s", baseURL, r.URL.Path)
}

func createCurrentLink(query url.Values, limit int64, offset int64) (bool, url.Values) {
	queryParams := url.Values{}
	for k, v := range query {
		if k != "page[limit]" && k != "page[offset]" {
			queryParams[k] = v
		}
	}

	queryParams["page[offset]"] = []string{strconv.Itoa(int(offset))}
	queryParams["page[limit]"] = []string{strconv.Itoa(int(limit))}

	return true, queryParams
}

func createlastPageLink(query url.Values, limitPtr int64, count int64) (bool, url.Values) {
	if limitPtr == 0 || count-limitPtr <= 0 {
		// no limit or limit greater than total count, so last page is the same as the first page.
		return createFirstPageLink(query, limitPtr)
	}

	limit := limitPtr

	queryParams := url.Values{}
	for k, v := range query {
		if k != "page[limit]" && k != "page[offset]" {
			queryParams[k] = v
		}
	}

	lastPageOffset := (count / limit) * limit
	if count%limit == 0 {
		lastPageOffset -= limit
	}

	queryParams["page[offset]"] = []string{strconv.FormatInt(lastPageOffset, 10)}
	queryParams["page[limit]"] = []string{strconv.Itoa(int(limit))}

	return true, queryParams
}

func createPrevPageLink(query url.Values, offset int64, limit int64, count int64) (bool, url.Values) {
	if offset == 0 {
		// no offset, so no previous page .
		return false, nil
	} else if limit == 0 {
		// no actual page size, so previous page will be the first page
		return createFirstPageLink(query, limit)
	}

	// there's a previous page, so decrease the offset accordingly
	newOffset := offset - limit
	if newOffset < 0 {
		newOffset = 0
	}
	if newOffset > count {
		// previous page would actually be higher than the number of results available, so
		// point the user to the last actual page instead.
		return createlastPageLink(query, limit, count)
	}

	queryParams := url.Values{}
	for k, v := range query {
		if k != "page[limit]" && k != "page[offset]" {
			queryParams[k] = v
		}
	}

	queryParams["page[offset]"] = []string{strconv.Itoa(int(newOffset))}
	queryParams["page[limit]"] = []string{strconv.Itoa(int(limit))}

	return true, queryParams
}

func createNextPageLink(query url.Values, offsetPtr int64, limit int64, count int64) (bool, url.Values) {
	if limit == 0 {
		// no page size, so no next page.
		return false, nil
	}

	queryParams := url.Values{}
	for k, v := range query {
		if k != "page[limit]" && k != "page[offset]" {
			queryParams[k] = v
		}
	}

	var offset int64
	if offsetPtr != 0 {
		offset = offsetPtr
	} else {
		offset = 0
	}
	if offset+limit >= count {
		// we're already on the last page, there is no next page
		return false, nil
	}

	// there's a next page, so increase the offset accordingly
	offset += limit
	queryParams["page[offset]"] = []string{strconv.Itoa(int(offset))}
	queryParams["page[limit]"] = []string{strconv.Itoa(int(limit))}

	return true, queryParams
}

func createFirstPageLink(query url.Values, limit int64) (bool, url.Values) {
	queryParams := url.Values{}
	for k, v := range query {
		if k != "page[limit]" && k != "page[offset]" {
			queryParams[k] = v
		}
	}

	queryParams["page[offset]"] = []string{"0"}

	if limit != 0 {
		queryParams["page[limit]"] = []string{strconv.Itoa(int(limit))}
	}

	return true, queryParams
}

func valuesToString(m url.Values) string {
	var sb strings.Builder

	for k, vs := range m {
		for _, v := range vs {
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
			sb.WriteString("&")
		}
	}
	// trim the last &
	s := sb.String()
	sLen := len(s)
	if sLen > 0 {
		s = s[:sLen-1]
	}

	return s
}

func createSelfLink(query url.Values) (bool, url.Values) {
	return true, query
}
