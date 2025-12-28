package pagination

import (
	"fmt"
	"hosting-kit/page"
	"hosting-resources-service/cmd/server/rest/gen"
)

func ToMetaData(pg page.Page, total int) gen.PageMetadata {
	doc := page.NewDocument(pg, total)

	return gen.PageMetadata{
		Number:          doc.Page,
		Size:            doc.PageSize,
		TotalElements:   doc.TotalCount,
		TotalPages:      doc.TotalPages,
		HasNextPage:     doc.HasNext,
		HasPreviousPage: doc.HasPrev,
	}
}

func ToLinks(baseURL string, pg page.Page, total int) gen.Links {
	doc := page.NewDocument(pg, total)
	links := make(gen.Links)

	makeHref := func(pageNum int) string {
		return fmt.Sprintf("%s?page=%d&pageSize=%d", baseURL, pageNum, doc.PageSize)
	}

	links["self"] = gen.Link{Href: makeHref(doc.Page)}
	links["first"] = gen.Link{Href: makeHref(1)}
	links["last"] = gen.Link{Href: makeHref(doc.TotalPages)}

	if doc.HasNext {
		links["next"] = gen.Link{Href: makeHref(doc.Page + 1)}
	}
	if doc.HasPrev {
		links["prev"] = gen.Link{Href: makeHref(doc.Page - 1)}
	}

	return links
}
