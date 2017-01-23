package up

import "regexp"

type Response string

type Handler func(*Request) Response

type Route struct {
	pattern string
	handler Handler
}

var func404 = func(r *Request) Response {
	return Response("Not Found")
}

func extractURIParameters(r *regexp.Regexp, s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]

	}
	return captures
}
