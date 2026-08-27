package internal

import "net/url"

var wkhtmlQueryAllow = map[string]struct{}{
	"page-size":     {},
	"orientation":   {},
	"margin-top":    {},
	"margin-bottom": {},
	"margin-left":   {},
	"margin-right":  {},
	"dpi":           {},
}

// WkhtmlParamsFromQuery copies allowlisted keys only when orientation is set.
func WkhtmlParamsFromQuery(q url.Values) map[string]string {
	params := make(map[string]string)
	if q.Get("orientation") == "" {
		return params
	}
	for key := range wkhtmlQueryAllow {
		if !q.Has(key) {
			continue
		}
		params[key] = q.Get(key)
	}
	return params
}
