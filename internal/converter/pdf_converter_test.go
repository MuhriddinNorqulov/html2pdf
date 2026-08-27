package converter

import (
	"slices"
	"testing"
)

func TestWkhtmlArgsDisableSmartShrinkingBoolean(t *testing.T) {
	args := wkhtmlArgs(map[string]string{
		"orientation":             "Landscape",
		"disable-smart-shrinking": "true",
		"dpi":                     "72",
	})

	if !slices.Contains(args, "--disable-smart-shrinking") {
		t.Fatalf("args %#v missing --disable-smart-shrinking", args)
	}
	if slices.Contains(args, "true") {
		t.Fatalf("args %#v must not pass true as a separate wkhtmltopdf argument", args)
	}
}

func TestWkhtmlArgsDisableSmartShrinkingFalseOmitted(t *testing.T) {
	args := wkhtmlArgs(map[string]string{
		"orientation":             "Landscape",
		"disable-smart-shrinking": "false",
	})

	if slices.Contains(args, "--disable-smart-shrinking") {
		t.Fatalf("args %#v must omit --disable-smart-shrinking when false", args)
	}
}

func TestWkhtmlArgsValueFlagsUnchanged(t *testing.T) {
	args := wkhtmlArgs(map[string]string{
		"page-size":   "A4",
		"orientation": "Landscape",
		"margin-top":  "0",
		"dpi":         "72",
	})

	want := []string{"--page-size", "A4", "--orientation", "Landscape", "--margin-top", "0", "--dpi", "72"}
	for _, item := range want {
		if !slices.Contains(args, item) {
			t.Fatalf("args %#v missing %q", args, item)
		}
	}
}
