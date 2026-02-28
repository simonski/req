package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type issueType string

const (
	epicType  issueType = "epic"
	storyType issueType = "task"
)

var (
	issueRefPattern = regexp.MustCompile(`\bE\d+(?:-S\d+)?\b`)
	epicIDPattern   = regexp.MustCompile(`^E\d+$`)
	storyIDPattern  = regexp.MustCompile(`^E\d+-S\d+$`)
)

type issue struct {
	kind           issueType
	title          string
	titleLine      int
	id             string
	idLine         int
	description    string
	descriptionLine int
	acceptance     []fieldLine
	priority       string
	priorityLine   int
	dependsOn      []string
	dependsLine    int
	parentID       string
	parentLine     int
	lines          []fieldLine
}

type fieldLine struct {
	text string
	line int
}

func main() {
	var filePath string
	flag.StringVar(&filePath, "f", "", "requirements markdown file")
	flag.Parse()

	if filePath == "" {
		exitErr(errors.New("missing required -f flag"))
	}

	issues, err := parseRequirements(filePath)
	if err != nil {
		exitErr(err)
	}

	if err := validateRequirements(issues); err != nil {
		exitErr(err)
	}

	for i, item := range issues {
		if i > 0 {
			fmt.Print("\n\n")
		}
		fmt.Print(buildCommand(item))
	}
}

func parseRequirements(filePath string) ([]issue, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var (
		issues        []issue
		current       *issue
		currentEpicID string
		inAC          bool
		lineNo        int
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNo++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			continue
		}

		switch {
		case strings.HasPrefix(trimmed, "EPIC:"):
			if current != nil {
				if err := validateRequiredFields(*current); err != nil {
					return nil, err
				}
				issues = append(issues, *current)
				if current.kind == epicType {
					currentEpicID = current.id
				}
			}
			current = &issue{
				kind:      epicType,
				title:     cleanValue(trimmed, "EPIC:"),
				titleLine: lineNo,
			}
			current.lines = append(current.lines, fieldLine{text: current.title, line: lineNo})
			inAC = false

		case strings.HasPrefix(trimmed, "STORY:"):
			if current != nil {
				if err := validateRequiredFields(*current); err != nil {
					return nil, err
				}
				issues = append(issues, *current)
				if current.kind == epicType {
					currentEpicID = current.id
				}
			}
			if currentEpicID == "" {
				return nil, fmt.Errorf("line %d: story found before any epic", lineNo)
			}
			current = &issue{
				kind:       storyType,
				title:      cleanValue(trimmed, "STORY:"),
				titleLine:  lineNo,
				parentID:   currentEpicID,
				parentLine: lineNo,
			}
			current.lines = append(current.lines, fieldLine{text: current.title, line: lineNo})
			inAC = false

		case current == nil:
			continue

		case hasAnyPrefix(trimmed, "ID:"):
			current.id = cleanValue(trimmed, "ID:")
			current.idLine = lineNo
			current.lines = append(current.lines, fieldLine{text: current.id, line: lineNo})
			inAC = false

		case hasAnyPrefix(trimmed, "DESCRIPTION:", "DESCRIPION:", "DESCRITION:"):
			current.description = cleanValueAny(trimmed, "DESCRIPTION:", "DESCRIPION:", "DESCRITION:")
			current.descriptionLine = lineNo
			current.lines = append(current.lines, fieldLine{text: current.description, line: lineNo})
			inAC = false

		case hasAnyPrefix(trimmed, "AC:"):
			value := cleanValue(trimmed, "AC:")
			if value != "" {
				current.acceptance = append(current.acceptance, fieldLine{text: value, line: lineNo})
				current.lines = append(current.lines, fieldLine{text: value, line: lineNo})
			}
			inAC = true

		case hasAnyPrefix(trimmed, "PRIORITY:"):
			current.priority = strings.TrimSpace(cleanValue(trimmed, "PRIORITY:"))
			current.priorityLine = lineNo
			current.lines = append(current.lines, fieldLine{text: current.priority, line: lineNo})
			inAC = false

		case hasAnyPrefix(trimmed, "DEPENDS-ON:"):
			current.dependsOn = parseDependsOn(cleanValue(trimmed, "DEPENDS-ON:"))
			current.dependsLine = lineNo
			for _, dep := range current.dependsOn {
				current.lines = append(current.lines, fieldLine{text: dep, line: lineNo})
			}
			inAC = false

		case inAC && strings.HasPrefix(trimmed, "-"):
			value := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
			current.acceptance = append(current.acceptance, fieldLine{text: value, line: lineNo})
			current.lines = append(current.lines, fieldLine{text: value, line: lineNo})

		case inAC && isContinuationLine(raw):
			if len(current.acceptance) == 0 {
				current.acceptance = append(current.acceptance, fieldLine{text: trimmed, line: lineNo})
				current.lines = append(current.lines, fieldLine{text: trimmed, line: lineNo})
			} else {
				last := len(current.acceptance) - 1
				current.acceptance[last].text = current.acceptance[last].text + " " + trimmed
				current.lines = append(current.lines, fieldLine{text: trimmed, line: lineNo})
			}

		default:
			inAC = false
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if current != nil {
		if err := validateRequiredFields(*current); err != nil {
			return nil, err
		}
		issues = append(issues, *current)
	}

	return issues, nil
}

func validateRequirements(issues []issue) error {
	if len(issues) == 0 {
		return errors.New("no epics or stories found")
	}

	byID := make(map[string]issue, len(issues))
	for _, item := range issues {
		if existing, ok := byID[item.id]; ok {
			return fmt.Errorf("line %d: duplicate ID %q already defined at line %d", item.idLine, item.id, existing.idLine)
		}
		byID[item.id] = item
	}

	for _, item := range issues {
		switch item.kind {
		case epicType:
			if !epicIDPattern.MatchString(item.id) {
				return fmt.Errorf("line %d: epic ID %q must match E<number>", item.idLine, item.id)
			}
		case storyType:
			if !storyIDPattern.MatchString(item.id) {
				return fmt.Errorf("line %d: story ID %q must match E<number>-S<number>", item.idLine, item.id)
			}
			if item.parentID == "" {
				return fmt.Errorf("line %d: story %q is missing a parent epic", item.idLine, item.id)
			}
			parent, ok := byID[item.parentID]
			if !ok {
				return fmt.Errorf("line %d: story %q refers to missing parent epic %q", item.parentLine, item.id, item.parentID)
			}
			if parent.kind != epicType {
				return fmt.Errorf("line %d: story %q parent %q is not an epic", item.parentLine, item.id, item.parentID)
			}
			expectedPrefix := item.parentID + "-S"
			if !strings.HasPrefix(item.id, expectedPrefix) {
				return fmt.Errorf("line %d: story ID %q does not belong under epic %q", item.idLine, item.id, item.parentID)
			}
		}

		if err := validatePriority(item); err != nil {
			return err
		}

		if err := validateDependencies(item, byID); err != nil {
			return err
		}

		for _, line := range item.lines {
			refs := issueRefPattern.FindAllString(line.text, -1)
			for _, ref := range refs {
				if ref == item.id || ref == item.parentID {
					continue
				}
				if _, ok := byID[ref]; !ok {
					return fmt.Errorf("line %d: reference to undefined issue %q", line.line, ref)
				}
			}
		}
	}

	return nil
}

func validateRequiredFields(item issue) error {
	if item.title == "" {
		return fmt.Errorf("line %d: missing title", item.titleLine)
	}
	if item.id == "" {
		return fmt.Errorf("line %d: missing ID for %q", item.titleLine, item.title)
	}
	if item.description == "" {
		return fmt.Errorf("line %d: missing description for %q", item.idLine, item.id)
	}
	if item.priority == "" {
		return fmt.Errorf("line %d: missing priority for %q", item.idLine, item.id)
	}
	return nil
}

func validatePriority(item issue) error {
	if item.priorityLine == 0 {
		return fmt.Errorf("line %d: missing priority for %q", item.idLine, item.id)
	}

	n, err := strconv.Atoi(item.priority)
	if err != nil {
		return fmt.Errorf("line %d: priority %q must be an integer", item.priorityLine, item.priority)
	}
	if n < 1 {
		return fmt.Errorf("line %d: priority %q must be >= 1", item.priorityLine, item.priority)
	}
	return nil
}

func validateDependencies(item issue, byID map[string]issue) error {
	for _, dep := range item.dependsOn {
		if dep == item.id {
			return fmt.Errorf("line %d: issue %q cannot depend on itself", item.dependsLine, item.id)
		}
		if _, ok := byID[dep]; !ok {
			return fmt.Errorf("line %d: dependency on undefined issue %q", item.dependsLine, dep)
		}
	}
	return nil
}

func buildCommand(item issue) string {
	args := []string{
		"bd create",
		"--type", shellQuote(string(item.kind)),
		"--title", shellQuote(item.title),
		"--description", shellQuote(buildDescription(item)),
		"--priority", shellQuote(mapPriority(item.priority)),
		"--external-ref", shellQuote(item.id),
	}

	if len(item.acceptance) > 0 {
		args = append(args, "--acceptance", shellQuote(strings.Join(normalizeAcceptance(item.acceptance), " | ")))
	}

	if item.parentID != "" {
		args = append(args, "--parent", shellVar(item.parentID))
	}

	if len(item.dependsOn) > 0 {
		args = append(args, "--deps", shellDeps(item.dependsOn))
	}

	args = append(args, "--silent")

	return fmt.Sprintf("%s=$(%s)", shellName(item.id), strings.Join(args, " "))
}

func normalizeAcceptance(values []fieldLine) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		out = append(out, normalizeWhitespace(value.text))
	}
	return out
}

func cleanValue(line, prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(line, prefix))
}

func cleanValueAny(line string, prefixes ...string) string {
	for _, prefix := range prefixes {
		if strings.HasPrefix(line, prefix) {
			return cleanValue(line, prefix)
		}
	}
	return strings.TrimSpace(line)
}

func hasAnyPrefix(line string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}

func isContinuationLine(raw string) bool {
	return strings.HasPrefix(raw, " ") || strings.HasPrefix(raw, "\t")
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func buildDescription(item issue) string {
	return fmt.Sprintf("[Logical ID: %s] %s", item.id, normalizeWhitespace(item.description))
}

func parseDependsOn(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		dep := strings.TrimSpace(part)
		if dep != "" {
			out = append(out, dep)
		}
	}
	return out
}

func mapPriority(input string) string {
	switch strings.TrimSpace(input) {
	case "1":
		return "0"
	case "2":
		return "1"
	case "3":
		return "2"
	case "4":
		return "3"
	case "5":
		return "4"
	default:
		return input
	}
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}

func shellName(value string) string {
	replacer := strings.NewReplacer("-", "_", " ", "_")
	return replacer.Replace(value)
}

func shellVar(logicalID string) string {
	return fmt.Sprintf("\"${%s}\"", shellName(logicalID))
}

func shellDeps(deps []string) string {
	parts := make([]string, 0, len(deps))
	for _, dep := range deps {
		parts = append(parts, fmt.Sprintf("${%s}", shellName(dep)))
	}
	return fmt.Sprintf("\"%s\"", strings.Join(parts, ","))
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
