package rootgrp

import (
	"fmt"
	"hosting-resources-service/cmd/server/rest/gen"
)

func toRoot(prefix string) gen.RootResource {
	links := make(gen.Links)

	makeHref := func(url string) string {
		return fmt.Sprintf("%s/%s", prefix, url)
	}

	links["self"] = gen.Link{Href: makeHref("")}
	links["pools"] = gen.Link{Href: makeHref("pools")}
	links["swagger"] = gen.Link{Href: makeHref("swagger/index.html")}

	return gen.RootResource{
		UnderscoreLinks: links,
	}
}
