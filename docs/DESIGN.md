# Prompt Rewrite

Build a single Go-installable binary named `req` published as `github.com/simonski/req`.

`req` is a requirements gathering and refinement system. Its purpose is to help users capture raw ideas, discuss and refine them, preserve traceability, and produce a structured specification that can be handed to a software delivery team or software factory.

## Product Goal

The product must support a workflow where a user starts with incomplete, mostly freeform text and gradually turns it into clear, machine-readable, traceable requirements.

The system should:

- capture raw requirements, notes, constraints, and questions
- preserve the original input without destructive rewriting
- support curation and refinement over time
- record the lineage of how raw ideas became structured requirements
- produce a final specification that can be decomposed into goals, epics, stories, and tasks

The main value of the system is not simple CRUD. The main value is turning ambiguous thought into a structured, reviewable specification with traceability.

## Core Product Principles

1. Raw user input is valuable and must be preserved.
2. Curation is incremental. The system should help users refine ideas over time rather than forcing a rigid upfront form.
3. Every important transformation must be traceable.
4. The server is the only component that accesses persistence.
5. The CLI and website must use the same API semantics and operate on the same underlying data model.

## Primary Interfaces

The system has three parts:

1. Server
2. CLI client
3. Embedded website

### Server

The server is the system of record. It exposes an OpenAPI-defined REST API and is solely responsible for persistence, authentication, authorization, business rules, and orchestration.

The server must:

- expose all core operations through REST endpoints described by OpenAPI
- persist application state in SQLite
- provide user registration, login, and logout
- support collaboration across multiple connected users
- provide a mechanism for near-real-time updates to connected clients

REST should be the default integration model. WebSocket or server-sent events may be used for live update delivery, but live updates should be layered on top of the core resource model instead of replacing it.

### CLI Client

The CLI is a power-user interface. It is intended for users who want precise, efficient control from the terminal.

Examples:

```bash
req add "Users can export project requirements as markdown"
req list
req show <id>
req curate <id>
req spec generate
```

The CLI must use the same server API contract as the website. It should not bypass the server or talk directly to SQLite.

The CLI should feel surgical and efficient rather than conversational by default, though selected commands may open an interactive flow when useful.

### Website

The website is the primary interface for most users.

The website must:

- be embedded into the Go binary using `go:embed`
- be implemented with simple web technology and kept operationally lightweight
- behave as a single-page application
- support collaborative use by multiple users
- reflect current system state without requiring constant manual refresh

The website should present requirements, discussions, decisions, and generated specification artifacts in a way that is easy to understand and navigate.

## Core Domain Model

The implementation must define a clear domain model. At minimum, model the following concepts:

- User
- Project
- Raw Input
- Requirement
- Requirement Draft or Curated Requirement
- Conversation or Chat Session
- Message or Turn
- Decision
- Trace Link
- Specification
- Specification Item

The exact naming may vary, but the system must distinguish between:

- original freeform input
- derived or curated requirement artifacts
- decisions made during refinement
- final specification output
- links that explain how one artifact was derived from another

## Required Workflows

The first implementation should support the following workflows end to end:

1. A user creates an account and authenticates.
2. A user creates or opens a project.
3. A user adds raw notes, constraints, requirements, and questions in freeform text.
4. The system stores those inputs durably and displays them clearly.
5. A user invokes curation/refinement for selected material.
6. The system records the full interaction history used during curation.
7. The system creates structured requirement artifacts without deleting the original text.
8. A user reviews and approves curated outputs.
9. A user generates a specification view or export.
10. The final specification can be organized into work-oriented structures such as goals, epics, stories, and tasks.

## LLM-Assisted Curation

The product may use an LLM to help curate and refine requirements, but the design must treat this as a controlled transformation step rather than opaque automation.

The system must:

- record prompts, responses, and relevant conversational context used for curation
- preserve lineage between source input and generated artifacts
- allow user review before important derived artifacts are treated as accepted truth
- store enough metadata to explain how a requirement or specification item was produced

LLM output should be treated as proposed or derived material unless the user explicitly accepts it.

## Traceability

Traceability is a first-class requirement.

The system must make it possible to answer questions such as:

- Which raw notes led to this requirement?
- Which conversation produced this decision?
- Which requirements support this specification item?
- What changed between earlier and later curated versions?

The design should prefer explicit trace links and append-only history over destructive replacement.

## Architecture Constraints

- The implementation language is Go.
- The install target is a single binary named `req`.
- The embedded website is served by the same Go application.
- The backend database is SQLite.
- Only the server accesses the database directly.
- The CLI and website both depend on the same OpenAPI-described HTTP API.

## Delivery Priorities

Prioritize a coherent vertical slice over feature breadth.

Phase 1 should prioritize:

- core domain model
- SQLite-backed server
- authentication
- basic project and requirement management
- append-only history of raw inputs and curated outputs
- traceability primitives
- CLI support for core actions
- minimal but usable SPA for viewing and editing data
- specification generation for a first useful export

Phase 2 may expand into:

- richer collaboration features
- stronger real-time updates
- more advanced curation flows
- richer spec export formats
- deeper decomposition into goals, epics, stories, and tasks

## Quality Gates

The repository should use a `Makefile` with at least:

```bash
make build
make test
make test-go
make test-playwright
```

Every shipped feature should include appropriate automated test coverage. Tests should pass before a feature is considered complete.

## Non-Goals For The First Cut

Do not optimize early for complexity that is not yet proven necessary.

Avoid overbuilding:

- a highly complex front-end stack
- direct database access from clients
- hidden or irreversible AI-driven rewrites
- multiple conflicting source-of-truth models

## Success Criteria

The implementation is successful if a user can:

1. capture messy requirements input
2. refine that input over time
3. understand how the refined output was produced
4. collaborate through the web interface
5. use the CLI for fast precise edits
6. generate a structured specification suitable for downstream delivery planning
