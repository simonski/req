# User Guide

`req` is a terminal tool for capturing rough requirements, refining them over time, and turning them into a traceable specification.

This guide describes how to use `req` from the command line as a working user. It does not cover installation or server startup.

## Mental Model

Use `req` in four stages:

1. Capture raw thoughts quickly.
2. Organize and curate those thoughts into clearer requirements.
3. Review traceability so you can see how conclusions were reached.
4. Generate and inspect a structured specification.

The important rule is that `req` preserves original input. Curated requirements and specification items are derived artifacts, not destructive rewrites of what you first entered.

## Core Concepts

### Project

A project is the workspace for a body of requirements. It contains raw notes, curated requirements, discussion history, decisions, and generated specifications.

### Raw Input

Raw input is unstructured text you add as you think. It can be a requirement, question, note, rule, concern, or example.

Examples:

- "Users should be able to export reports as markdown."
- "Need audit history for changes."
- "Unsure whether guests can edit shared documents."

### Requirement

A requirement is a curated, clearer statement derived from one or more pieces of raw input. Requirements are intended to be reviewable and traceable.

### Conversation

A conversation records the interaction used to refine or clarify a requirement. This includes user prompts, assistant responses, and any accepted conclusions.

### Decision

A decision is a chosen outcome during refinement, such as selecting one interpretation over another or confirming a business rule.

### Specification

A specification is the structured output built from accepted requirements. It may be grouped into goals, epics, stories, and tasks.

## Typical Workflow

### 1. Create or Open a Project

Start by creating a project or selecting an existing one.

Examples:

```bash
req project create "Customer Portal"
req project list
req project use customer-portal
```

Once a project is active, most commands operate within that project.

### 2. Capture Raw Requirements

Add ideas as they occur. Do not try to make them perfect on first entry.

Examples:

```bash
req add "Customers can reset their password without contacting support."
req add "Password reset should expire after a short time."
req add "We may need MFA later."
req note "Need to check compliance rules for password handling."
req question "Should password reset be available to invited but inactive users?"
```

Use short entries when you want clean traceability. Use longer entries when you want to preserve context in one place.

### 3. Review What You Have Captured

Inspect the current state before refining anything.

Examples:

```bash
req list
req list --type raw
req show 42
req search "password reset"
```

At this stage you are usually checking for duplication, ambiguity, and open questions.

### 4. Curate Raw Input Into Requirements

Use curation when raw input is too vague, overlapping, or incomplete.

Examples:

```bash
req curate 42
req curate 42 43 44
req curate --from-search "password reset"
```

During curation, `req` may:

- identify distinct requirements within one note
- extract assumptions and constraints
- ask follow-up questions
- propose clearer requirement statements
- record decisions and unresolved issues

The result of curation is usually one or more proposed requirements linked back to the source material.

### 5. Review Proposed Changes

Do not assume curated output is automatically accepted. Review it.

Examples:

```bash
req review
req review --status proposed
req show requirement 17
req diff requirement 17
```

When reviewing, check:

- whether the requirement says what you intended
- whether important nuance was lost
- whether the source links are correct
- whether unresolved questions still need answers

### 6. Accept, Reject, or Revise

Once you review proposed output, decide what should become part of the working requirements set.

Examples:

```bash
req accept requirement 17
req reject requirement 18
req revise requirement 17
req decision add "Reset links expire after 15 minutes."
```

Accepted requirements become candidates for inclusion in the generated specification.

### 7. Inspect Traceability

One of the main reasons to use `req` is to understand where a requirement came from.

Examples:

```bash
req trace requirement 17
req history 17
req conversation show 9
req decision list
```

Use trace and history commands to answer questions like:

- Which raw notes produced this requirement?
- Which conversation clarified it?
- What changed between versions?
- Which decisions are still unresolved?

### 8. Generate a Specification

When enough requirements have been accepted, generate a specification view.

Examples:

```bash
req spec generate
req spec show
req spec show --section goals
req spec export markdown
```

The generated specification is expected to organize accepted requirements into a more delivery-oriented structure.

Depending on the project, this may include:

- goals
- epics
- stories
- tasks
- constraints
- open questions

### 9. Refine the Specification

Specification generation is not the end of the process. Continue iterating.

Examples:

```bash
req spec review
req spec trace 3
req curate --from-spec 3
```

If a spec item is unclear or unsupported, follow its trace links back to the source requirements and raw notes, then refine again.

## Working Style Recommendations

### Capture First, Perfect Later

Do not wait until an idea is polished. Enter it while it is still fresh, then curate it later.

### Keep Entries Focused

Small entries are easier to trace and curate than large mixed paragraphs. If one note contains multiple ideas, expect to split it during curation.

### Use Questions Explicitly

If something is uncertain, record it as a question rather than hiding the uncertainty inside a requirement.

### Review Trace Links

If a requirement cannot be traced back to clear source material, treat it as suspect until reviewed.

### Accept Deliberately

Accepted requirements should represent decisions you are willing to carry forward into the specification.

## Common Command Patterns

Capture:

```bash
req add "..."
req note "..."
req question "..."
```

Inspect:

```bash
req list
req show <id>
req search "<text>"
```

Curate:

```bash
req curate <id>
req review
req accept requirement <id>
req reject requirement <id>
```

Trace:

```bash
req trace requirement <id>
req history <id>
req conversation show <id>
```

Specification:

```bash
req spec generate
req spec show
req spec export markdown
```

## What Good Usage Looks Like

A healthy `req` project usually has:

- lots of preserved raw input
- curated requirements that are clearer than the source material
- explicit decisions where ambiguity was resolved
- trace links from specification items back to their origins
- unresolved questions called out instead of buried

If you are using `req` well, you should be able to move from a final spec item back to the exact notes and conversations that produced it.
