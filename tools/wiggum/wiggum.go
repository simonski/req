package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	defaultReadyLimit     = 100
	defaultDryRunSleep    = 10 * time.Second
	defaultPromptTemplate = "Perform the following:\n<BEAD>"
)

type issue struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Acceptance   string    `json:"acceptance_criteria"`
	Status       string    `json:"status"`
	Priority     int       `json:"priority"`
	IssueType    string    `json:"issue_type"`
	ExternalRef  string    `json:"external_ref"`
	Assignee     string    `json:"assignee"`
	Owner        string    `json:"owner"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Parent       string    `json:"parent"`
	Dependencies []dep     `json:"dependencies"`
}

type dep struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Status         string `json:"status"`
	DependencyType string `json:"dependency_type"`
	ExternalRef    string `json:"external_ref"`
}

var logicalIDPattern = regexp.MustCompile(`(?i)\[Logical ID:\s*([^\]]+)\]`)

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "loop":
		if err := runLoop(os.Args[2:]); err != nil {
			exitErr(err)
		}
	case "check":
		if err := runCheck(os.Args[2:]); err != nil {
			exitErr(err)
		}
	case "agent":
		if err := runAgent(os.Args[2:]); err != nil {
			exitErr(err)
		}
	case "-h", "--help", "help":
		printUsage()
	default:
		printUsage()
		exitErr(fmt.Errorf("unknown command %q", os.Args[1]))
	}
}

func runLoop(args []string) error {
	fs := flag.NewFlagSet("loop", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	var (
		name       string
		max        int
		dryRun     bool
		readyLimit int
		sleepDur   time.Duration
		noBranch   bool
	)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: wiggum loop -name fred -max 1 [-dryrun]\n\n")
		fmt.Fprintf(fs.Output(), "Chooses the next best ready bead for the named wiggum, assigns it,\n")
		fmt.Fprintf(fs.Output(), "creates or switches to a branch for the work, performs the work,\n")
		fmt.Fprintf(fs.Output(), "closes the bead, and repeats.\n\n")
		fmt.Fprintf(fs.Output(), "Flags:\n")
		fs.PrintDefaults()
		fmt.Fprintf(fs.Output(), "\nEnvironment:\n")
		fmt.Fprintf(fs.Output(), "  WIGGUM_WORK_CMD   Optional command used to process a work packet in non-dry-run mode.\n")
		fmt.Fprintf(fs.Output(), "                    The command receives the work packet on stdin.\n")
	}

	fs.StringVar(&name, "name", "", "unique wiggum name used for assignment")
	fs.IntVar(&max, "max", 1, "maximum number of issues to process (0 = forever)")
	fs.BoolVar(&dryRun, "dryrun", false, "simulate work, sleep briefly, and close the bead")
	fs.IntVar(&readyLimit, "limit", defaultReadyLimit, "maximum ready issues to consider per iteration")
	fs.DurationVar(&sleepDur, "sleep", defaultDryRunSleep, "sleep duration used during dry-run work simulation")
	fs.BoolVar(&noBranch, "no-branch", false, "do not create or switch git branches")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if name == "" {
		fs.Usage()
		return errors.New("missing required -name flag")
	}

	processed := 0
	for max == 0 || processed < max {
		ready, err := loadReady(readyLimit)
		if err != nil {
			return err
		}

		chosen, ok := chooseIssue(ready, name)
		if !ok {
			if processed == 0 {
				return errors.New("no suitable ready issue found")
			}
			return nil
		}

		if err := assignIssue(chosen.ID, name); err != nil {
			return err
		}

		full, err := showIssue(chosen.ID)
		if err != nil {
			return err
		}

		branchName := branchNameFor(full)
		branched := false
		if !noBranch && !dryRun {
			if err := switchBranch(branchName); err != nil {
				return err
			}
			branched = true
		}

		printWorkPacket(full, name, branchName, true, branched, dryRun)

		if dryRun {
			time.Sleep(sleepDur)
		} else {
			if err := performWork(full, name, branchName); err != nil {
				return err
			}
		}

		if err := closeIssue(full.ID, name, dryRun); err != nil {
			return err
		}

		processed++
	}

	return nil
}

func runCheck(args []string) error {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)

	var (
		name       string
		readyLimit int
	)

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: wiggum check -name fred\n\n")
		fmt.Fprintf(fs.Output(), "Shows the next bead wiggum would work on, including the derived branch\n")
		fmt.Fprintf(fs.Output(), "name and work packet, without mutating beads state or git state.\n\n")
		fmt.Fprintf(fs.Output(), "Flags:\n")
		fs.PrintDefaults()
	}

	fs.StringVar(&name, "name", "", "unique wiggum name used for assignment filtering")
	fs.IntVar(&readyLimit, "limit", defaultReadyLimit, "maximum ready issues to consider")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if name == "" {
		fs.Usage()
		return errors.New("missing required -name flag")
	}

	ready, err := loadReady(readyLimit)
	if err != nil {
		return err
	}

	chosen, ok := chooseIssue(ready, name)
	if !ok {
		fmt.Printf("Wiggum: %s\n", name)
		fmt.Println("Next: none")
		fmt.Println("Reason: no suitable ready issue found")
		return nil
	}

	full, err := showIssue(chosen.ID)
	if err != nil {
		return err
	}

	branchName := branchNameFor(full)
	fmt.Print(renderWorkPacket(full, name, branchName, false, false, false))
	return nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  wiggum agent \"entire command\"")
	fmt.Println("  wiggum agent command [arg...]")
	fmt.Println("  wiggum check -name fred")
	fmt.Println("  wiggum loop -name fred -max 1")
	fmt.Println("  wiggum loop -name fred -max 1 -dryrun")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  agent   Run an interactive coding agent command with stdio passed through.")
	fmt.Println("  check   Show what wiggum would do next without changing beads or git.")
	fmt.Println("  loop    Claim the next best ready bead, work it, close it, and repeat.")
	fmt.Println()
	fmt.Println("Run `wiggum agent -h`, `wiggum check -h`, or `wiggum loop -h` for command flags.")
}

func runAgent(args []string) error {
	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: wiggum agent \"entire command\"\n")
		fmt.Fprintf(fs.Output(), "       wiggum agent command [arg...]\n\n")
		fmt.Fprintf(fs.Output(), "Starts a child process and passes stdin, stdout, and stderr through directly,\n")
		fmt.Fprintf(fs.Output(), "so interactive coding agents behave as if they were launched without wiggum.\n")
	}

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	cmd, err := buildAgentCommand(fs.Args())
	if err != nil {
		fs.Usage()
		return err
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("agent command failed: %w", err)
	}
	return nil
}

func buildAgentCommand(args []string) (*exec.Cmd, error) {
	switch len(args) {
	case 0:
		return nil, errors.New("missing agent command")
	case 1:
		return exec.Command("sh", "-c", args[0]), nil
	default:
		return exec.Command(args[0], args[1:]...), nil
	}
}

func loadReady(limit int) ([]issue, error) {
	args := []string{"ready", "--json", "--limit", fmt.Sprintf("%d", limit)}
	output, err := runBD(args...)
	if err != nil {
		return nil, err
	}

	var issues []issue
	if err := json.Unmarshal(output, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

func chooseIssue(issues []issue, name string) (issue, bool) {
	candidates := make([]issue, 0, len(issues))
	for _, item := range issues {
		if item.Assignee != "" && !strings.EqualFold(item.Assignee, name) {
			continue
		}
		candidates = append(candidates, item)
	}

	if len(candidates) == 0 {
		return issue{}, false
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		leftAssigned := strings.EqualFold(candidates[i].Assignee, name)
		rightAssigned := strings.EqualFold(candidates[j].Assignee, name)
		if leftAssigned != rightAssigned {
			return leftAssigned
		}

		leftLeaf := isLeafWork(candidates[i].IssueType)
		rightLeaf := isLeafWork(candidates[j].IssueType)
		if leftLeaf != rightLeaf {
			return leftLeaf
		}

		if candidates[i].Priority != candidates[j].Priority {
			return candidates[i].Priority < candidates[j].Priority
		}
		if !candidates[i].CreatedAt.Equal(candidates[j].CreatedAt) {
			return candidates[i].CreatedAt.Before(candidates[j].CreatedAt)
		}
		return candidates[i].ID < candidates[j].ID
	})

	return candidates[0], true
}

func isLeafWork(issueType string) bool {
	switch issueType {
	case "task", "bug", "feature", "chore", "merge-request":
		return true
	default:
		return false
	}
}

func assignIssue(id, name string) error {
	_, err := runBD("update", id, "--assignee", name, "--status", "in_progress")
	if err != nil {
		return fmt.Errorf("failed to assign %s to %s: %w", id, name, err)
	}
	return nil
}

func showIssue(id string) (issue, error) {
	output, err := runBD("show", "--json", id)
	if err != nil {
		return issue{}, err
	}

	var issues []issue
	if err := json.Unmarshal(output, &issues); err != nil {
		return issue{}, err
	}
	if len(issues) == 0 {
		return issue{}, fmt.Errorf("bd show returned no issue for %s", id)
	}
	return issues[0], nil
}

func performWork(item issue, name, branch string) error {
	workCmd := strings.TrimSpace(os.Getenv("WIGGUM_WORK_CMD"))
	if workCmd == "" {
		return errors.New("WIGGUM_WORK_CMD is not set for non-dry-run execution")
	}

	startedAt := time.Now()
	prompt := renderPrompt(renderWorkPacket(item, name, branch, true, branch != "", false))
	transcript := &lockedBuffer{}
	cmd := exec.Command("sh", "-c", workCmd)
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Stdout = io.MultiWriter(os.Stdout, transcript)
	cmd.Stderr = io.MultiWriter(os.Stderr, transcript)
	cmd.Env = append(os.Environ(),
		"WIGGUM_NAME="+name,
		"WIGGUM_ISSUE_ID="+item.ID,
		"WIGGUM_LOGICAL_ID="+item.ExternalRef,
		"WIGGUM_BRANCH="+branch,
	)

	err := cmd.Run()
	completedAt := time.Now()
	exitCode := exitCodeFor(err)
	if logErr := writeAgentLog(item.ID, name, branch, prompt, transcript.String(), startedAt, completedAt, exitCode); logErr != nil {
		if err != nil {
			return fmt.Errorf("worker command failed for %s: %w (also failed to write log: %v)", item.ID, err, logErr)
		}
		return fmt.Errorf("failed to write agent log for %s: %w", item.ID, logErr)
	}
	if err != nil {
		return fmt.Errorf("worker command failed for %s: %w", item.ID, err)
	}
	return nil
}

func closeIssue(id, name string, dryRun bool) error {
	reason := "completed by wiggum " + name
	if dryRun {
		reason = "dryrun completed by wiggum " + name
	}
	_, err := runBD("close", id, "--reason", reason)
	if err != nil {
		return fmt.Errorf("failed to close %s: %w", id, err)
	}
	return nil
}

func branchNameFor(item issue) string {
	logicalID := item.ExternalRef
	if logicalID == "" {
		if match := logicalIDPattern.FindStringSubmatch(item.Description); len(match) == 2 {
			logicalID = match[1]
		}
	}

	prefix := slugify(logicalID)
	title := slugify(item.Title)
	switch {
	case prefix != "" && title != "":
		return prefix + "-" + title
	case title != "":
		return title
	case prefix != "":
		return prefix
	default:
		return slugify(item.ID)
	}
}

func slugify(input string) string {
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" {
		return ""
	}

	var b strings.Builder
	lastDash := false
	for _, r := range input {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			lastDash = false
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}

	return strings.Trim(b.String(), "-")
}

func switchBranch(branch string) error {
	if branch == "" {
		return errors.New("empty branch name")
	}

	if err := exec.Command("git", "rev-parse", "--verify", branch).Run(); err == nil {
		cmd := exec.Command("git", "switch", branch)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to switch to branch %s: %s", branch, strings.TrimSpace(string(output)))
		}
		return nil
	}

	cmd := exec.Command("git", "switch", "-c", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create branch %s: %s", branch, strings.TrimSpace(string(output)))
	}
	return nil
}

func printWorkPacket(item issue, name, branch string, assigned, branched, dryRun bool) {
	fmt.Print(renderPrompt(renderWorkPacket(item, name, branch, assigned, branched, dryRun)))
}

func renderWorkPacket(item issue, name, branch string, assigned, branched, dryRun bool) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Wiggum: %s\n", name)
	fmt.Fprintf(&b, "ID: %s\n", item.ID)
	if item.ExternalRef != "" {
		fmt.Fprintf(&b, "Logical ID: %s\n", item.ExternalRef)
	}
	fmt.Fprintf(&b, "Type: %s\n", item.IssueType)
	fmt.Fprintf(&b, "Priority: %d\n", item.Priority)
	fmt.Fprintf(&b, "Status: %s\n", item.Status)
	fmt.Fprintf(&b, "Assignee: %s\n", item.Assignee)
	fmt.Fprintf(&b, "Owner: %s\n", item.Owner)
	fmt.Fprintf(&b, "Title: %s\n", item.Title)
	fmt.Fprintf(&b, "Branch: %s\n", branch)
	fmt.Fprintf(&b, "Assigned: %t\n", assigned)
	fmt.Fprintf(&b, "Branch Switched: %t\n", branched)
	fmt.Fprintf(&b, "Dry Run: %t\n\n", dryRun)

	fmt.Fprintln(&b, "Description:")
	fmt.Fprintln(&b, strings.TrimSpace(item.Description))
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "Acceptance Criteria:")
	for _, part := range splitAcceptance(item.Acceptance) {
		fmt.Fprintf(&b, "- %s\n", part)
	}
	fmt.Fprintln(&b)

	if len(item.Dependencies) > 0 {
		fmt.Fprintln(&b, "Dependencies:")
		for _, dep := range item.Dependencies {
			ref := dep.ID
			if dep.ExternalRef != "" {
				ref = dep.ExternalRef + " (" + dep.ID + ")"
			}
			fmt.Fprintf(&b, "- %s: %s [%s]\n", ref, dep.Title, dep.Status)
		}
		fmt.Fprintln(&b)
	}

	fmt.Fprintln(&b, "Next:")
	fmt.Fprintf(&b, "- Start work on %s\n", item.ID)
	fmt.Fprintln(&b, "- Run tests with make before closing the issue")
	return b.String()
}

func renderPrompt(packet string) string {
	template := strings.TrimSpace(os.Getenv("WIGGUM_PROMPT_TEMPLATE"))
	if template == "" {
		template = defaultPromptTemplate
	}
	if strings.Contains(template, "<BEAD>") {
		return strings.ReplaceAll(template, "<BEAD>", packet)
	}
	if strings.HasSuffix(template, "\n") {
		return template + packet
	}
	return template + "\n" + packet
}

type lockedBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (l *lockedBuffer) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.Write(p)
}

func (l *lockedBuffer) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.b.String()
}

func writeAgentLog(beadID, name, branch, prompt, transcript string, startedAt, completedAt time.Time, exitCode int) error {
	logPath := filepath.Join("logs", sanitizeLogComponent(name), logFileName(beadID, branch))
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return err
	}

	var b strings.Builder
	fmt.Fprintln(&b, "FULL PROMPT")
	fmt.Fprintln(&b, "-----------")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, prompt)
	if !strings.HasSuffix(prompt, "\n") {
		fmt.Fprintln(&b)
	}
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "FULL STDIN/STDOUT")
	fmt.Fprintln(&b, "-----------")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "STDIN:")
	fmt.Fprintln(&b, prompt)
	if !strings.HasSuffix(prompt, "\n") {
		fmt.Fprintln(&b)
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "STDOUT:")
	fmt.Fprintln(&b, transcript)
	if !strings.HasSuffix(transcript, "\n") {
		fmt.Fprintln(&b)
	}
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "TIME STARTED:")
	fmt.Fprintln(&b, "-----------")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, startedAt.Format(time.RFC3339Nano))
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "TIME COMPLETED:")
	fmt.Fprintln(&b, "-----------")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, completedAt.Format(time.RFC3339Nano))
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "EXIT CODE:")
	fmt.Fprintln(&b, "-----------")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, exitCode)
	fmt.Fprintln(&b)

	fmt.Fprintln(&b, "BEAD ID:")
	fmt.Fprintln(&b, "-----------")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, beadID)

	return os.WriteFile(logPath, []byte(b.String()), 0o644)
}

func logFileName(beadID, branch string) string {
	return sanitizeLogComponent(beadID) + "-" + sanitizeLogComponent(branch) + ".log"
}

func sanitizeLogComponent(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown"
	}

	var b strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			lastDash = false
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
			lastDash = false
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case r == '.' || r == '_' || r == '-':
			if !(r == '-' && lastDash) {
				b.WriteRune(r)
				lastDash = r == '-'
			}
		default:
			if !lastDash && b.Len() > 0 {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}

	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "unknown"
	}
	return out
}

func exitCodeFor(err error) int {
	if err == nil {
		return 0
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}
	return -1
}

func splitAcceptance(input string) []string {
	parts := strings.Split(input, "|")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func runBD(args ...string) ([]byte, error) {
	cmd := exec.Command("bd", args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return nil, errors.New(msg)
	}
	return stdout.Bytes(), nil
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
