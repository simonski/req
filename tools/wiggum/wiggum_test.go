package main

import (
	"strings"
	"testing"
)

func TestBuildAgentCommandRequiresArgs(t *testing.T) {
	t.Parallel()

	if _, err := buildAgentCommand(nil); err == nil {
		t.Fatal("expected an error for missing args")
	}
}

func TestBuildAgentCommandUsesShellForSingleArgument(t *testing.T) {
	t.Parallel()

	cmd, err := buildAgentCommand([]string{"codex --help"})
	if err != nil {
		t.Fatalf("buildAgentCommand returned error: %v", err)
	}

	gotArgs := cmd.Args
	wantArgs := []string{"sh", "-c", "codex --help"}
	if len(gotArgs) != len(wantArgs) {
		t.Fatalf("Args len = %d, want %d (%v)", len(gotArgs), len(wantArgs), gotArgs)
	}
	for i := range wantArgs {
		if gotArgs[i] != wantArgs[i] {
			t.Fatalf("Args[%d] = %q, want %q (all args: %v)", i, gotArgs[i], wantArgs[i], gotArgs)
		}
	}
}

func TestBuildAgentCommandUsesDirectExecForMultipleArguments(t *testing.T) {
	t.Parallel()

	cmd, err := buildAgentCommand([]string{"codex", "--help"})
	if err != nil {
		t.Fatalf("buildAgentCommand returned error: %v", err)
	}

	gotArgs := cmd.Args
	wantArgs := []string{"codex", "--help"}
	if len(gotArgs) != len(wantArgs) {
		t.Fatalf("Args len = %d, want %d (%v)", len(gotArgs), len(wantArgs), gotArgs)
	}
	for i := range wantArgs {
		if gotArgs[i] != wantArgs[i] {
			t.Fatalf("Args[%d] = %q, want %q (all args: %v)", i, gotArgs[i], wantArgs[i], gotArgs)
		}
	}
}

func TestRenderPromptUsesDefaultTemplate(t *testing.T) {
	t.Setenv("WIGGUM_PROMPT_TEMPLATE", "")

	got := renderPrompt("BEAD")
	want := "Perform the following:\nBEAD"
	if got != want {
		t.Fatalf("renderPrompt() = %q, want %q", got, want)
	}
}

func TestRenderPromptUsesTemplatePlaceholder(t *testing.T) {
	t.Setenv("WIGGUM_PROMPT_TEMPLATE", "Please do this:\n<BEAD>\nThanks")

	got := renderPrompt("BEAD")
	want := "Please do this:\nBEAD\nThanks"
	if got != want {
		t.Fatalf("renderPrompt() = %q, want %q", got, want)
	}
}

func TestRenderWorkPacketIsPromptWrapped(t *testing.T) {
	t.Setenv("WIGGUM_PROMPT_TEMPLATE", "")

	item := issue{
		ID:          "bd-123",
		Title:       "Tidy parser help",
		Description: "Update the help output",
		Acceptance:  "help prints usage|tests pass",
		Status:      "open",
		Priority:    2,
		IssueType:   "task",
	}

	got := renderPrompt(renderWorkPacket(item, "ralph", "feature/ralph/tidy-parser-help", false, false, false))
	if !strings.HasPrefix(got, "Perform the following:\nWiggum: ralph\n") {
		t.Fatalf("wrapped packet missing default prompt prefix: %q", got)
	}
	if !strings.Contains(got, "ID: bd-123\n") {
		t.Fatalf("wrapped packet missing bead contents: %q", got)
	}
}
