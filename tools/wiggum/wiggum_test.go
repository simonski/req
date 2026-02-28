package main

import "testing"

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
