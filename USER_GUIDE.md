# User Guide

Use `req` in this order:

1. Pick a project.
2. Add raw input.
3. Review what exists.
4. Curate into requirements.
5. Accept or reject proposals.
6. Check traceability.
7. Generate a specification.

## Projects

```bash
req project create "Customer Portal"
req project list
req project use customer-portal
```

## Add Input

Requirement:

```bash
req add "Customers can reset their password."
```

Note:

```bash
req note "Need audit history for password changes."
```

Question:

```bash
req question "Should invited but inactive users be able to reset passwords?"
```

## Review Input

```bash
req list
req list --type raw
req show 42
req search "password reset"
```

## Curate

Curate one item:

```bash
req curate 42
```

Curate several items:

```bash
req curate 42 43 44
```

Curate by search:

```bash
req curate --from-search "password reset"
```

## Review Proposals

```bash
req review
req review --status proposed
req show requirement 17
req diff requirement 17
```

## Accept, Reject, Revise

Accept:

```bash
req accept requirement 17
```

Reject:

```bash
req reject requirement 18
```

Revise:

```bash
req revise requirement 17
```

Record a decision:

```bash
req decision add "Reset links expire after 15 minutes."
```

## Traceability

Trace a requirement:

```bash
req trace requirement 17
```

See history:

```bash
req history 17
```

Show a conversation:

```bash
req conversation show 9
```

List decisions:

```bash
req decision list
```

## Specification

Generate:

```bash
req spec generate
```

Show:

```bash
req spec show
req spec show --section goals
```

Export:

```bash
req spec export markdown
```

Refine from the spec:

```bash
req spec review
req spec trace 3
req curate --from-spec 3
```

## Short Reference

```bash
req project create "..."
req project list
req project use ...

req add "..."
req note "..."
req question "..."

req list
req show <id>
req search "..."

req curate <id>
req review
req accept requirement <id>
req reject requirement <id>
req revise requirement <id>

req trace requirement <id>
req history <id>
req conversation show <id>
req decision list

req spec generate
req spec show
req spec export markdown
```
