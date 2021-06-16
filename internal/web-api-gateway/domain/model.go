package domain

type UrlMapping struct {
	urlMapping map[string]string
}

func NewUrlMapping(mapping map[string]string) *UrlMapping {
	u := new(UrlMapping)
	u.urlMapping = mapping
	return u
}

func (u *UrlMapping) Get(mapping string) string {
	return u.urlMapping[mapping]
}
