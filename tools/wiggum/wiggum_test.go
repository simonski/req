package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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

func TestRunLoopRequiresAgentWhenNotDryRun(t *testing.T) {
	t.Parallel()

	err := runLoop([]string{"-name", "ralph", "-max", "1"})
	if err == nil || !strings.Contains(err.Error(), "missing required -agent flag") {
		t.Fatalf("runLoop() error = %v, want missing required -agent flag", err)
	}
}

func TestRunLoopDryRunFlagRequiresInteger(t *testing.T) {
	t.Parallel()

	err := runLoop([]string{"-name", "ralph", "-dryrun"})
	if err == nil {
		t.Fatal("runLoop() error = nil, want non-nil")
	}
}

func TestDryRunCommandUsesProvidedSeconds(t *testing.T) {
	t.Parallel()

	got := dryRunCommand(7)
	want := "echo 'dry-run; sleep 7'"
	if got != want {
		t.Fatalf("dryRunCommand() = %q, want %q", got, want)
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

func TestLogFileNameSanitizesPathLikeCharacters(t *testing.T) {
	t.Parallel()

	got := logFileName("bd/123", "feature/ralph/fix:thing")
	want := "bd-123-feature-ralph-fix-thing.log"
	if got != want {
		t.Fatalf("logFileName() = %q, want %q", got, want)
	}
}

func TestWriteAgentLogIncludesRequiredSections(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	startedAt := time.Date(2026, 2, 28, 11, 0, 0, 0, time.UTC)
	completedAt := startedAt.Add(2 * time.Minute)

	if err := writeAgentLog("bd/123", "ralph", "feature/ralph/fix:thing", "Perform the following:\nBEAD", "agent output\n", startedAt, completedAt, 7); err != nil {
		t.Fatalf("writeAgentLog() error = %v", err)
	}

	logPath := filepath.Join(tempDir, "logs", "ralph", "bd-123-feature-ralph-fix-thing.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", logPath, err)
	}

	got := string(data)
	for _, want := range []string{
		"FULL PROMPT\n-----------",
		"FULL STDIN/STDOUT\n-----------",
		"TIME STARTED:\n-----------",
		"TIME COMPLETED:\n-----------",
		"EXIT CODE:\n-----------",
		"BEAD ID:\n-----------",
		"STDIN:\nPerform the following:\nBEAD",
		"STDOUT:\nagent output\n",
		"2026-02-28T11:00:00Z",
		"2026-02-28T11:02:00Z",
		"\n7\n",
		"\nbd/123\n",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("log content missing %q in %q", want, got)
		}
	}
}

func TestPerformWorkWritesLogOnWorkerFailure(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)
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

	err := performWork(item, "ralph", "feature/ralph/tidy-parser-help", true, "printf 'agent says hi\\n'; cat >/dev/null; exit 7")
	if err == nil {
		t.Fatal("performWork() error = nil, want non-nil")
	}

	logPath := filepath.Join(tempDir, "logs", "ralph", "bd-123-feature-ralph-tidy-parser-help.log")
	data, readErr := os.ReadFile(logPath)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", logPath, readErr)
	}

	got := string(data)
	if !strings.Contains(got, "agent says hi\n") {
		t.Fatalf("log missing worker output: %q", got)
	}
	if !strings.Contains(got, "\n7\n") {
		t.Fatalf("log missing exit code 7: %q", got)
	}
	if !strings.Contains(got, "Perform the following:\nWiggum: ralph\n") {
		t.Fatalf("log missing prompt: %q", got)
	}
}

func TestPerformDryRunWorkWritesLog(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)
	t.Setenv("WIGGUM_PROMPT_TEMPLATE", "")

	item := issue{
		ID:          "bd-456",
		Title:       "Simulate parser help",
		Description: "Update the help output",
		Acceptance:  "help prints usage|tests pass",
		Status:      "open",
		Priority:    2,
		IssueType:   "task",
	}

	if err := performDryRunWork(item, "jane", "feature/jane/simulate-parser-help", false, 7); err != nil {
		t.Fatalf("performDryRunWork() error = %v", err)
	}

	logPath := filepath.Join(tempDir, "logs", "jane", "bd-456-feature-jane-simulate-parser-help.log")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", logPath, err)
	}

	got := string(data)
	if !strings.Contains(got, "dry-run; sleep 7\n") {
		t.Fatalf("log missing dry-run process output: %q", got)
	}
	if !strings.Contains(got, "\n0\n") {
		t.Fatalf("log missing exit code 0: %q", got)
	}
	if !strings.Contains(got, "Dry Run: true\n") {
		t.Fatalf("log missing dry-run prompt contents: %q", got)
	}
	if !strings.Contains(got, "Branch Switched: false\n") {
		t.Fatalf("log missing correct branch switched value: %q", got)
	}
}
