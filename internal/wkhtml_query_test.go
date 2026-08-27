package internal

import (
	"net/url"
	"testing"
)

func TestWkhtmlParamsFromQueryEmptyWithoutOrientation(t *testing.T) {
	q := url.Values{}
	q.Set("page-size", "A4")
	q.Set("margin-top", "15mm")
	got := WkhtmlParamsFromQuery(q)
	if len(got) != 0 {
		t.Fatalf("got %#v, want empty (ConvertA4 must stay unchanged)", got)
	}
}

func TestWkhtmlParamsFromQueryLandscapeAllowlist(t *testing.T) {
	q := url.Values{}
	q.Set("orientation", "Landscape")
	q.Set("page-size", "A4")
	q.Set("margin-top", "0")
	q.Set("dpi", "72")
	q.Set("zoom", "1.25")
	q.Set("disable-smart-shrinking", "true")
	got := WkhtmlParamsFromQuery(q)
	if got["orientation"] != "Landscape" || got["dpi"] != "72" || got["margin-top"] != "0" {
		t.Fatalf("got %#v", got)
	}
	if _, ok := got["zoom"]; ok {
		t.Fatal("zoom must be dropped")
	}
	if _, ok := got["disable-smart-shrinking"]; ok {
		t.Fatal("disable-smart-shrinking must be dropped")
	}
}

func TestWkhtmlParamsFromQueryMobilityLandscapeNoDpi(t *testing.T) {
	q := url.Values{}
	q.Set("page-size", "A4")
	q.Set("orientation", "Landscape")
	got := WkhtmlParamsFromQuery(q)
	if len(got) != 2 {
		t.Fatalf("got %#v, want only page-size+orientation", got)
	}
	if got["page-size"] != "A4" || got["orientation"] != "Landscape" {
		t.Fatalf("got %#v", got)
	}
}
