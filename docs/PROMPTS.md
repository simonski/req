0. <freeform converted to DESIGN.md>

1. Write a USER_GUIDE.md at the top based on a hypothetical implementation of this using hte docs/DESIGN.md.    Do not include how to run it, only from the perspective of a user in the terminal.

2. Using the DESIGN and USER_GUIDE, USER_GUIDE_LONGFORM write an example breakdown of implementation requirements as REQUIREMENTS.md in the format:

EPIC: title
ID: E1, E2, E3 etc
DESCRIPTION: description
AC: list of acceptance criteria
PRIORITY: 1-N (1 highest, do this first)
DEPENDS-ON: E2, E4

<indent for stories "in" the epic (the story ID should increment and be EPIC-STORY)>
    STORY: title
    ID: E1-S1, E1-S2, E1-S3 etc.
    DESCRIPTION: description
    AC: list of acceptance criteria
    PRIORITY: 1-N (1 highest, do this first)
    DEPENDS-ON: E1-S2

The intent is to take this output and model it in beads, or jira, or an issue tracker.  The scope is:
- ALL examples in the user guides
- ALL of the backend and frontend functionality as per the design

Note the DEPENDS-ON is a method of describing blocking features.
Ensure the acceptance critera contains
    - all tests pass
    - use make to verify tests
    - work in a branch that contains the EPIC and FEATURE name


3. Write/rewrite a parser go program that translates a requirements.md into beads commands (but do not call beads). It should just be a single go file runnable as "go run parser.go -f REQUIREMENTS.md" which writes to stdout all the beads commands with double- newlines between beads.   It should read yhe whole requirements, validate they are correct and have referntial integrity where they refer to other EPICS or STORIES, call out the error-line if there is one, exit 1 if there is a problem, or just print the commands and exit 0.



