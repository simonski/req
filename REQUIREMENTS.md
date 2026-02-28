EPIC: Platform Foundation and API Contract
ID: E1
DESCRIPTION: Establish the Go server, shared OpenAPI contract, SQLite persistence boundary, embedded website delivery, and build/test tooling for a single-binary application.
AC:
- A Go application named `req` exists as the primary deliverable.
- The server is the only component that accesses SQLite directly.
- An OpenAPI specification defines the REST API used by both CLI and website clients.
- The website is embedded into the Go binary with `go:embed`.
- A Makefile provides `make build`, `make test`, `make test-go`, and `make test-playwright`.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 1
DEPENDS-ON:

    STORY: Create server application skeleton
    ID: E1-S1
    DESCRIPTION: Implement the server entrypoint, configuration loading, routing skeleton, middleware hooks, and lifecycle management for the backend service.
    AC:
    - The server starts with a valid configuration.
    - HTTP routing is structured for versioned API endpoints.
    - Graceful startup and shutdown behavior is implemented.
    - Basic health and readiness endpoints exist.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON:

    STORY: Define the OpenAPI contract
    ID: E1-S2
    DESCRIPTION: Create and maintain an OpenAPI definition for authentication, project, input, requirement, conversation, decision, traceability, and specification operations.
    AC:
    - The OpenAPI document covers all REST operations required by the CLI and website.
    - Request and response schemas exist for all documented resources.
    - The CLI and website can be implemented without undocumented endpoints.
    - The OpenAPI document is versioned with the codebase.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E1-S1

    STORY: Establish SQLite persistence boundary
    ID: E1-S3
    DESCRIPTION: Implement database schema management and persistence layers so all durable state is isolated to the server.
    AC:
    - SQLite schema migration support exists.
    - The CLI does not read or write the database directly.
    - The website does not read or write the database directly.
    - Database access code is isolated behind server-side modules.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E1-S1

    STORY: Embed and serve the website
    ID: E1-S4
    DESCRIPTION: Package SPA assets into the Go binary and serve them from the same application as the API.
    AC:
    - Static site assets are embedded in the binary.
    - The server can serve the SPA shell and static assets.
    - SPA deep links resolve correctly through the server.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E1-S1

    STORY: Provide build and test automation
    ID: E1-S5
    DESCRIPTION: Add required Makefile targets and wire them to application build, Go tests, and browser tests.
    AC:
    - `make build` builds the application successfully.
    - `make test` runs all test suites.
    - `make test-go` runs backend and CLI automated tests.
    - `make test-playwright` runs browser-based tests for the website.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E1-S1

EPIC: Authentication and User Sessions
ID: E2
DESCRIPTION: Support user registration, login, logout, and authenticated access to project data across CLI and web interfaces.
AC:
- Users can register accounts.
- Users can log in and obtain an authenticated session.
- Users can log out and invalidate the active session.
- Protected endpoints require authentication.
- CLI and web session handling both use the server API.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 2
DEPENDS-ON: E1

    STORY: Implement user registration
    ID: E2-S1
    DESCRIPTION: Add account creation with validation, secure credential handling, and durable user records.
    AC:
    - A user can register with required credentials.
    - Duplicate or invalid registrations are rejected with clear errors.
    - User records are stored durably in SQLite.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E1-S2, E1-S3

    STORY: Implement login and logout flows
    ID: E2-S2
    DESCRIPTION: Add session creation and destruction for authenticated access from CLI and web clients.
    AC:
    - A valid user can log in successfully.
    - Invalid credentials are rejected.
    - A logged-in user can log out.
    - Authenticated requests succeed only with a valid session.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E2-S1

    STORY: Add CLI authentication commands
    ID: E2-S3
    DESCRIPTION: Provide terminal commands for registration, login, logout, and session persistence.
    AC:
    - Users can register from the CLI.
    - Users can log in from the CLI.
    - Users can log out from the CLI.
    - CLI commands persist and reuse session state securely.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E2-S1, E2-S2

    STORY: Add website authentication screens
    ID: E2-S4
    DESCRIPTION: Provide web UI flows for registration, login, logout, and session expiration handling.
    AC:
    - A user can register from the website.
    - A user can log in from the website.
    - A user can log out from the website.
    - Expired or invalid sessions redirect the user to authenticate again.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E1-S4, E2-S1, E2-S2

EPIC: Project Management
ID: E3
DESCRIPTION: Allow users to create, list, select, and open projects as the container for requirements work.
AC:
- A user can create a project.
- A user can list accessible projects.
- A user can open or select a project in both CLI and web.
- Project-scoped operations are isolated to the active project.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 1
DEPENDS-ON: E1, E2

    STORY: Create project domain model and storage
    ID: E3-S1
    DESCRIPTION: Define the project resource and persistence model used to scope all requirements activity.
    AC:
    - Projects have a durable identifier and display name.
    - Projects can be created and fetched from SQLite.
    - Projects own related raw inputs, requirements, decisions, conversations, and specifications.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E1-S3, E2-S2

    STORY: Implement project REST endpoints
    ID: E3-S2
    DESCRIPTION: Add API endpoints for creating, listing, and retrieving projects.
    AC:
    - The API supports project creation.
    - The API supports project listing for the authenticated user.
    - The API supports retrieval of a single project.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E3-S1, E1-S2

    STORY: Implement CLI project commands
    ID: E3-S3
    DESCRIPTION: Support `req project create`, `req project list`, and `req project use` flows in the terminal.
    AC:
    - `req project create "Customer Portal"` creates a project.
    - `req project list` lists available projects.
    - `req project use customer-portal` sets the active project context.
    - Commands fail clearly when the target project does not exist.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E3-S2

    STORY: Implement web project selection UI
    ID: E3-S4
    DESCRIPTION: Provide website screens for creating a project, listing projects, and switching the active project.
    AC:
    - Users can create projects from the website.
    - Users can browse available projects from the website.
    - Users can switch the active project in the website.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E1-S4, E3-S2

EPIC: Raw Input Capture
ID: E4
DESCRIPTION: Support durable capture of freeform requirements, notes, questions, constraints, and related input without destructive rewriting.
AC:
- Users can add raw input items to a project.
- The system distinguishes between requirement-like input, notes, and questions.
- Original raw input is preserved even after curation.
- Raw input is visible in both CLI and web interfaces.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 1
DEPENDS-ON: E1, E2, E3

    STORY: Define raw input model
    ID: E4-S1
    DESCRIPTION: Model raw input as a first-class resource with type, content, author, timestamps, and project linkage.
    AC:
    - Raw input records store content, type, project, author, and creation time.
    - Supported types include requirement, note, and question.
    - Raw input records are immutable or append-only for lineage purposes.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E3-S1

    STORY: Implement raw input creation API
    ID: E4-S2
    DESCRIPTION: Add REST endpoints to create and retrieve raw input items for a project.
    AC:
    - The API supports creating requirement-like raw input.
    - The API supports creating note-type raw input.
    - The API supports creating question-type raw input.
    - The API returns stored items with stable identifiers.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E4-S1, E3-S2

    STORY: Implement CLI input capture commands
    ID: E4-S3
    DESCRIPTION: Support `req add`, `req note`, and `req question` in the terminal.
    AC:
    - `req add "Customers can reset their password."` stores a raw requirement item.
    - `req note "Need audit history for password changes."` stores a raw note item.
    - `req question "Should invited but inactive users be able to reset passwords?"` stores a raw question item.
    - Commands require an active project or explicit project selection.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E4-S2, E3-S3

    STORY: Implement web input capture UI
    ID: E4-S4
    DESCRIPTION: Provide web forms or controls for adding raw requirement, note, and question items.
    AC:
    - A user can add a raw requirement from the website.
    - A user can add a note from the website.
    - A user can add a question from the website.
    - Newly created items appear in the project view without manual refresh.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E1-S4, E4-S2, E3-S4

EPIC: Review, Listing, Search, and Item Inspection
ID: E5
DESCRIPTION: Allow users to list project content, inspect individual items, filter by type, and search by text across captured requirements material.
AC:
- Users can list available content in a project.
- Users can filter raw inputs by type.
- Users can show a single item by ID.
- Users can search content text within a project.
- CLI and web both expose these workflows.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 1
DEPENDS-ON: E4

    STORY: Implement listing and filtering APIs
    ID: E5-S1
    DESCRIPTION: Add API support for listing project items with type filters, status filters, and pagination or equivalent.
    AC:
    - The API supports listing all relevant items in a project.
    - The API supports filtering raw items by type.
    - The API supports filtering reviewed or proposed requirement items by status.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E4-S2

    STORY: Implement item retrieval APIs
    ID: E5-S2
    DESCRIPTION: Add API support to retrieve a specific raw input, requirement, conversation, decision, or specification item by ID.
    AC:
    - A single item can be retrieved by stable ID.
    - Missing IDs return a clear not-found response.
    - The response includes related metadata needed for display.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E4-S2

    STORY: Implement search APIs
    ID: E5-S3
    DESCRIPTION: Add text search across project content to support targeted review and curation.
    AC:
    - A user can search project content by freeform text.
    - Search results include enough context to identify matching items.
    - Search can be scoped to the active project.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E4-S2

    STORY: Implement CLI review and inspection commands
    ID: E5-S4
    DESCRIPTION: Support `req list`, `req list --type raw`, `req show 42`, and `req search "password reset"` flows in the terminal.
    AC:
    - `req list` displays project items.
    - `req list --type raw` displays raw input items only.
    - `req show 42` displays an individual item with details.
    - `req search "password reset"` returns matching items.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E5-S1, E5-S2, E5-S3, E3-S3

    STORY: Implement web review and search UI
    ID: E5-S5
    DESCRIPTION: Provide screens to browse items, filter content, inspect item details, and search within a project.
    AC:
    - The website shows project items in a browsable list.
    - Users can filter by raw input type and supported statuses.
    - Users can inspect a single item detail view.
    - Users can search content from the website.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E1-S4, E5-S1, E5-S2, E5-S3, E3-S4

EPIC: Curation and Requirement Derivation
ID: E6
DESCRIPTION: Transform raw input into proposed requirements through controlled curation flows that preserve source lineage and conversation history.
AC:
- Users can curate one or more items into proposed requirements.
- Curation can be started from explicit IDs and from search results.
- The system records conversation turns used during curation.
- Derived requirements remain linked to their source raw inputs.
- Curated output is proposed until accepted by a user.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 1
DEPENDS-ON: E4, E5

    STORY: Define curated requirement and conversation models
    ID: E6-S1
    DESCRIPTION: Model proposed requirements, conversation sessions, conversation turns, and links to source material.
    AC:
    - Proposed requirements have status and stable identity.
    - Conversations and conversation turns are stored durably.
    - Source raw inputs are linked to derived requirements.
    - Multiple raw inputs can feed one curation flow.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E4-S1

    STORY: Implement curation orchestration service
    ID: E6-S2
    DESCRIPTION: Build the backend workflow that takes selected inputs, invokes curation logic or LLM-assisted refinement, and stores proposed requirements plus conversation history.
    AC:
    - A curation request can target one item or multiple items.
    - The workflow stores prompts, responses, or equivalent curation interactions.
    - The workflow stores resulting proposed requirements.
    - Original source items are not mutated or deleted.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E6-S1, E5-S2

    STORY: Implement curation REST endpoints
    ID: E6-S3
    DESCRIPTION: Add API endpoints to launch curation by item IDs and search-driven selections and to retrieve created proposals.
    AC:
    - The API supports curation by explicit item IDs.
    - The API supports curation by search query or server-side resolved selection.
    - The API returns created proposed requirements and related conversation IDs.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E6-S2, E1-S2

    STORY: Implement CLI curation commands
    ID: E6-S4
    DESCRIPTION: Support `req curate 42`, `req curate 42 43 44`, and `req curate --from-search "password reset"` in the terminal.
    AC:
    - `req curate 42` curates a single item.
    - `req curate 42 43 44` curates multiple items together.
    - `req curate --from-search "password reset"` curates matching items.
    - The CLI reports created proposed requirements or follow-up actions.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E6-S3, E5-S3

    STORY: Implement web curation UI
    ID: E6-S5
    DESCRIPTION: Provide web workflows to curate selected items and display generated proposals.
    AC:
    - Users can launch curation from selected raw items in the website.
    - Users can curate one or multiple items from the website.
    - Proposed requirements appear in the website when curation completes.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E1-S4, E6-S3, E5-S5

EPIC: Requirement Review, Decisions, and Status Management
ID: E7
DESCRIPTION: Let users review proposed requirements, compare changes, accept or reject proposals, revise requirements, and record explicit decisions.
AC:
- Users can list proposed requirements for review.
- Users can inspect a proposed requirement and compare changes.
- Users can accept, reject, or revise proposed requirements.
- Users can create and list decisions tied to the project and related items.
- Accepted requirements are marked as ready for specification generation.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 1
DEPENDS-ON: E6

    STORY: Implement requirement statuses and review model
    ID: E7-S1
    DESCRIPTION: Add statuses and transitions for proposed, accepted, rejected, and revised requirement artifacts.
    AC:
    - Requirement status values are modeled explicitly.
    - Allowed transitions are enforced by the server.
    - Revision creates a new tracked version or lineage-preserving artifact.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E6-S1

    STORY: Implement review and diff APIs
    ID: E7-S2
    DESCRIPTION: Add endpoints to list proposed requirements, fetch individual requirement details, and compare changes.
    AC:
    - The API supports listing requirements by status.
    - The API supports retrieving a specific requirement.
    - The API supports diff or change comparison for reviewed requirements.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E7-S1, E5-S2

    STORY: Implement accept, reject, and revise APIs
    ID: E7-S3
    DESCRIPTION: Add backend actions to accept, reject, and revise requirements while preserving history.
    AC:
    - A user can accept a proposed requirement.
    - A user can reject a proposed requirement.
    - A user can revise a requirement and preserve previous state.
    - Accepted requirements are available for specification generation.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E7-S1, E7-S2

    STORY: Implement decision tracking APIs
    ID: E7-S4
    DESCRIPTION: Support creation and retrieval of decisions made during curation and review.
    AC:
    - A user can add a decision statement.
    - Decisions are stored with project context and authorship.
    - Decisions can be listed and associated with relevant requirements or conversations.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E6-S1

    STORY: Implement CLI review and decision commands
    ID: E7-S5
    DESCRIPTION: Support `req review`, `req review --status proposed`, `req show requirement 17`, `req diff requirement 17`, `req accept requirement 17`, `req reject requirement 18`, `req revise requirement 17`, and `req decision add/list`.
    AC:
    - `req review` lists reviewable proposed requirements.
    - `req review --status proposed` filters to proposed requirements.
    - `req show requirement 17` shows requirement details.
    - `req diff requirement 17` shows changes or derivation differences.
    - `req accept requirement 17` marks a requirement accepted.
    - `req reject requirement 18` marks a requirement rejected.
    - `req revise requirement 17` creates a revision flow.
    - `req decision add "Reset links expire after 15 minutes."` stores a decision.
    - `req decision list` displays project decisions.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E7-S2, E7-S3, E7-S4

    STORY: Implement web review and decision UI
    ID: E7-S6
    DESCRIPTION: Provide website views and actions for reviewing proposed requirements, inspecting diffs, and managing decisions.
    AC:
    - Users can browse proposed requirements in the website.
    - Users can inspect requirement details and changes.
    - Users can accept, reject, and revise from the website.
    - Users can add and browse decisions from the website.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E1-S4, E7-S2, E7-S3, E7-S4

EPIC: Traceability and History
ID: E8
DESCRIPTION: Provide full lineage from raw input through curation, decisions, revisions, and final specification items.
AC:
- Users can trace a requirement back to raw inputs and conversations.
- Users can inspect history for a requirement or related item.
- Users can show a conversation by ID.
- The system preserves append-only or equivalent auditable history.
- Specification items can be traced back to accepted requirements.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 1
DEPENDS-ON: E6, E7

    STORY: Define trace link model
    ID: E8-S1
    DESCRIPTION: Implement explicit trace links between raw input, conversations, requirements, decisions, and specification items.
    AC:
    - Trace links can connect source and derived artifacts.
    - Trace links identify relationship type.
    - Trace links are queryable by source and target item.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E6-S1, E7-S4

    STORY: Implement history storage and retrieval
    ID: E8-S2
    DESCRIPTION: Preserve historical versions and event records for curation, review, revision, and specification generation.
    AC:
    - Requirement history is stored durably.
    - Users can retrieve a chronological history for an item.
    - Revision and acceptance events are included in history.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E7-S1, E7-S3

    STORY: Implement traceability REST endpoints
    ID: E8-S3
    DESCRIPTION: Add endpoints to retrieve trace graphs, history views, and conversation details for related artifacts.
    AC:
    - The API supports tracing a requirement.
    - The API supports retrieving item history.
    - The API supports retrieving a conversation by ID.
    - The API supports tracing from specification items to accepted requirements.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E8-S1, E8-S2

    STORY: Implement CLI traceability commands
    ID: E8-S4
    DESCRIPTION: Support `req trace requirement 17`, `req history 17`, and `req conversation show 9` in the terminal.
    AC:
    - `req trace requirement 17` shows source raw inputs and related artifacts.
    - `req history 17` shows a chronological history of the item.
    - `req conversation show 9` shows conversation turns for the curation session.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E8-S3, E7-S5

    STORY: Implement web traceability views
    ID: E8-S5
    DESCRIPTION: Provide website displays for lineage, history, and conversation detail.
    AC:
    - Users can view trace links for a requirement in the website.
    - Users can view change history for a requirement or related item.
    - Users can inspect a curation conversation in the website.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E1-S4, E8-S3, E7-S6

EPIC: Specification Generation and Export
ID: E9
DESCRIPTION: Generate structured specifications from accepted requirements, support sectioned viewing, traceability from spec items, review workflows, and export to Markdown.
AC:
- Users can generate a specification from accepted requirements.
- Users can view the full specification.
- Users can view sections such as goals.
- Users can export the specification as Markdown.
- Users can trace specification items back to source requirements.
- Users can continue refining requirements from specification review.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 1
DEPENDS-ON: E7, E8

    STORY: Define specification and specification item models
    ID: E9-S1
    DESCRIPTION: Model generated specifications and their component items, including structural grouping and trace links to accepted requirements.
    AC:
    - Specifications have stable identity and generation metadata.
    - Specification items are individually addressable.
    - Specification items can map to accepted requirements.
    - Sections such as goals can be represented.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E8-S1, E7-S3

    STORY: Implement specification generation service
    ID: E9-S2
    DESCRIPTION: Build backend logic that produces a structured specification from accepted requirements.
    AC:
    - The service uses only accepted requirements as generation inputs unless configured otherwise.
    - Generated outputs are stored durably.
    - A new specification can be generated repeatedly without losing prior history.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E9-S1, E7-S3

    STORY: Implement specification retrieval and export APIs
    ID: E9-S3
    DESCRIPTION: Add endpoints to generate, show, section-filter, trace, review, and export specifications.
    AC:
    - The API supports generating a specification.
    - The API supports retrieving the current specification.
    - The API supports retrieving a section such as goals.
    - The API supports exporting Markdown.
    - The API supports tracing a specification item.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E9-S2, E8-S3

    STORY: Implement CLI specification commands
    ID: E9-S4
    DESCRIPTION: Support `req spec generate`, `req spec show`, `req spec show --section goals`, `req spec export markdown`, `req spec review`, `req spec trace 3`, and `req curate --from-spec 3`.
    AC:
    - `req spec generate` creates or refreshes a specification.
    - `req spec show` displays the specification.
    - `req spec show --section goals` displays the goals section only.
    - `req spec export markdown` exports the spec in Markdown.
    - `req spec review` identifies issues or unresolved items for further work.
    - `req spec trace 3` traces a specification item to accepted requirements.
    - `req curate --from-spec 3` starts a refinement flow from a specification item.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E9-S3, E8-S4, E6-S4

    STORY: Implement web specification UI
    ID: E9-S5
    DESCRIPTION: Provide website views for generating, browsing, section-filtering, tracing, reviewing, and exporting specifications.
    AC:
    - Users can generate a specification from the website.
    - Users can browse the full specification from the website.
    - Users can filter by specification section from the website.
    - Users can export the specification as Markdown from the website.
    - Users can trace a specification item from the website.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E1-S4, E9-S3, E8-S5

EPIC: Collaborative Web Experience and Live Updates
ID: E10
DESCRIPTION: Deliver the website as the primary collaborative interface with SPA navigation and near-real-time synchronization across multiple users.
AC:
- The website behaves as a single-page application.
- Multiple users can work on the same project.
- Changes made by one user become visible to others without constant manual refresh.
- The website covers projects, raw inputs, review, curation results, traceability, decisions, and specifications.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 3
DEPENDS-ON: E3, E4, E5, E6, E7, E8, E9

    STORY: Build SPA shell and navigation
    ID: E10-S1
    DESCRIPTION: Implement the client-side shell, routing, and state transitions for the embedded website.
    AC:
    - Users can navigate between project, item, review, traceability, and specification screens without full page reloads.
    - Direct navigation to supported SPA routes works through server fallback.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E1-S4, E3-S4, E5-S5

    STORY: Implement live update transport
    ID: E10-S2
    DESCRIPTION: Add WebSocket or server-sent event support so the website can reflect current shared state.
    AC:
    - Connected web clients receive updates when relevant project data changes.
    - Live updates are scoped to the active project or relevant subscriptions.
    - Failure of the live transport degrades gracefully to manual refresh or polling behavior.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E10-S1, E6-S5, E7-S6, E8-S5, E9-S5

    STORY: Implement multi-user project synchronization
    ID: E10-S3
    DESCRIPTION: Ensure concurrent user activity produces coherent project views and durable updates.
    AC:
    - Two users can add or review items in the same project without data loss.
    - The website surfaces the latest committed state from the server.
    - Concurrent changes are handled consistently by the backend.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E10-S2, E3-S2, E4-S2, E7-S3

EPIC: CLI Experience and Shared Client Semantics
ID: E11
DESCRIPTION: Ensure the terminal interface is a first-class power-user client with clear output, consistent command behavior, and full parity for the documented workflows.
AC:
- All commands shown in the short and long user guides exist or have implemented equivalents.
- CLI commands use the same API semantics as the web client.
- Command output is readable and action-oriented.
- Error handling guides the user toward the next command or fix.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 2
DEPENDS-ON: E2, E3, E4, E5, E6, E7, E8, E9

    STORY: Implement CLI command structure
    ID: E11-S1
    DESCRIPTION: Organize the CLI into project, input, review, traceability, decision, conversation, and specification command groups.
    AC:
    - Command groups map cleanly to documented workflows.
    - Help output documents command usage and arguments.
    - Invalid command usage returns actionable error messages.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E1-S2

    STORY: Implement shared API client for CLI
    ID: E11-S2
    DESCRIPTION: Build a reusable API client layer for terminal commands using the same contract as the website.
    AC:
    - All CLI commands use HTTP APIs rather than database access.
    - Authentication and active project context are handled consistently across commands.
    - API errors are rendered clearly in the terminal.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E11-S1, E1-S2, E2-S2

    STORY: Provide CLI output formats for review workflows
    ID: E11-S3
    DESCRIPTION: Format lists, details, diffs, trace output, and specification output for efficient terminal use.
    AC:
    - List output is readable for `req list`, `req review`, and `req decision list`.
    - Detail output is readable for `req show`, `req trace`, `req history`, and `req conversation show`.
    - Diff output is readable for `req diff requirement <id>`.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E11-S2, E5-S4, E7-S5, E8-S4, E9-S4

EPIC: Quality, Testing, and Non-Functional Requirements
ID: E12
DESCRIPTION: Validate the system with automated tests that cover backend logic, CLI workflows, and website behavior for the documented feature set.
AC:
- Automated tests cover core domain rules, API behavior, CLI commands, and website workflows.
- The documented example flows in the user guides are represented in automated test coverage.
- Required Makefile targets execute successfully.
- Feature delivery is gated on passing tests.
- All tests pass.
- Use `make` to verify tests.
- Work in a branch that contains the EPIC and FEATURE name.
PRIORITY: 2
DEPENDS-ON: E1, E2, E3, E4, E5, E6, E7, E8, E9, E10, E11

    STORY: Add backend unit and integration tests
    ID: E12-S1
    DESCRIPTION: Test persistence, authentication, project management, input capture, curation, review, traceability, and specification generation on the server side.
    AC:
    - Core domain logic is covered by unit tests.
    - HTTP API behavior is covered by integration tests.
    - SQLite-backed workflows are tested.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E2, E3, E4, E5, E6, E7, E8, E9

    STORY: Add CLI workflow tests
    ID: E12-S2
    DESCRIPTION: Test the terminal commands shown in the user guides against representative project workflows.
    AC:
    - Tests cover project creation and selection.
    - Tests cover `req add`, `req note`, and `req question`.
    - Tests cover `req list`, `req show`, and `req search`.
    - Tests cover `req curate`, `req review`, `req accept`, `req reject`, and `req revise`.
    - Tests cover `req trace`, `req history`, `req conversation show`, and `req decision list`.
    - Tests cover `req spec generate`, `req spec show`, `req spec export markdown`, `req spec review`, and `req spec trace`.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 1
    DEPENDS-ON: E11, E2-S3, E3-S3, E4-S3, E5-S4, E6-S4, E7-S5, E8-S4, E9-S4

    STORY: Add end-to-end website tests
    ID: E12-S3
    DESCRIPTION: Use Playwright or equivalent to validate the web application workflows across authentication, project work, review, traceability, and specification generation.
    AC:
    - End-to-end tests cover login and project access.
    - End-to-end tests cover raw input creation and review.
    - End-to-end tests cover curation and proposal review.
    - End-to-end tests cover traceability and specification viewing.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 2
    DEPENDS-ON: E10, E2-S4, E3-S4, E4-S4, E5-S5, E6-S5, E7-S6, E8-S5, E9-S5

    STORY: Add regression coverage for collaborative updates
    ID: E12-S4
    DESCRIPTION: Validate multi-user and live update behavior in automated tests where feasible.
    AC:
    - Tests cover at least one multi-user shared project scenario.
    - Tests validate that updates propagate to the website without a full reload.
    - All tests pass.
    - Use `make` to verify tests.
    - Work in a branch that contains the EPIC and FEATURE name.
    PRIORITY: 3
    DEPENDS-ON: E10-S2, E10-S3, E12-S3
