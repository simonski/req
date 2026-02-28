README-WIGGUM.md

`wiggum` - a go binary which looks at beads and performs the work via a large language model, recording the output for the work carried out.

The instance of this wiggum will have a unique name "fred, pete, ralph, jane" which is what is used when assignment occurs.  Each instance is considered to be an agent.

1. take a piece of work from beads
    wiggum will have a name - use that as the assignee
    choose a bead you have been assigned to - or assign yourself to an appropriate bead
    an appropriate bead is
        (in development assigned to this wiggum, OR
        open, ready for development) AND
        appropriate priority AND
        not blocked by another bead AND
        not assigned to a different wiggum AND
        either unassigned OR 
        assigned to this wiggum 

2. take the work packet (The bead itself) and feed it into an LLM coding agent session
3. close teh bead
4. repeat

Usage

./wiggum loop -name fred -max 1
    performs the work in a loop -max times (0 = forever)

./wiggum loop -name fred -max 1 -dryrun
    simulates the work where it does not feed the packet to an LLM but prints to STDOUT only, sleeps for 1s then marks as complete in beads.

./wiggum 
    prints help

what sort of ralph wiggum loop woudl make beads keep working until it finished - can you write me a program that does that? go run wiggum.go

no, I want wiggum to look at beads directly using bd commands, find hte next best ticket and work on it

extend wiggum.go and parser.go so that they print a useful help usage if they are invoked with no command