package main

import (
	"strings"
	"testing"
)

func TestBuildCommandUsesNewlineSeparatedAcceptance(t *testing.T) {
	item := issue{
		kind:        storyType,
		title:       "Implement parser output",
		id:          "E2-S1",
		description: "Emit acceptance criteria as separate lines",
		priority:    "2",
		acceptance: []fieldLine{
			{text: "first criterion"},
			{text: "second criterion"},
		},
	}

	got := buildCommand(item)
	if !strings.Contains(got, "--acceptance 'first criterion\nsecond criterion\n"+additionalContextAC+"'") {
		t.Fatalf("buildCommand() missing newline separated acceptance: %q", got)
	}
	if strings.Contains(got, " | ") {
		t.Fatalf("buildCommand() unexpectedly contains pipe-separated acceptance: %q", got)
	}
}

func TestNormalizeAcceptanceAppendsAdditionalContextOnce(t *testing.T) {
	values := []fieldLine{
		{text: "first criterion"},
		{text: additionalContextAC},
	}

	got := normalizeAcceptance(values)
	if len(got) != 2 {
		t.Fatalf("normalizeAcceptance() len = %d, want 2 (%v)", len(got), got)
	}
}
