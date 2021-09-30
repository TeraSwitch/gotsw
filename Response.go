package gotsw

import "net/http"

type Response struct {
	*http.Response

	Metadata *Metadata
}

// Meta describes generic information about a response.
type Metadata struct {
	Total int `json:"total"`
}
