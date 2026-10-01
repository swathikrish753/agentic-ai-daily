package main

import (
	"fmt"
	"strings"
)

// What the agent decides to do next.
type Action struct {
	Go   string // place to check next; empty means "stop"
	Done bool
	Msg  string
}

// 🧠 THE BRAIN: looks at the LAST thing it saw, and reacts.
// This is what makes it an agent instead of a to-do list.
func decide(lastSeen string) Action {
	switch {
	case lastSeen == "": // nothing seen yet → the very first move
		return Action{Go: "kitchen"}
	case strings.Contains(lastSeen, "shiny banana"):
		return Action{Done: true, Msg: "Found it! 🍌"}
	case strings.Contains(lastSeen, "go to the fridge"):
		return Action{Go: "fridge"}
	default:
		return Action{Done: true, Msg: "No clue where it is. Giving up. 😞"}
	}
}

// ✋ THE HANDS: go look in a place and report what's there.
func look(place string) string {
	world := map[string]string{
		"kitchen": "no banana, but there's a note: go to the fridge",
		"fridge":  "a shiny banana is sitting here",
	}
	if whatIsThere, ok := world[place]; ok {
		return whatIsThere
	}
	return "an empty room"
}

func main() {
	lastSeen := "" // we've observed nothing yet

	// 🔁 THE LOOP — now the SEE feeds back into the next THINK.
	for step := 1; step <= 10; step++ { // step cap = safety belt
		action := decide(lastSeen) // 🧠 think (reacting to what we saw)

		if action.Done { // brain decided to stop
			fmt.Printf("✅ %s\n", action.Msg)
			return
		}

		fmt.Printf("Step %d\n", step)
		fmt.Printf("  🧠 think: I'll check the %s\n", action.Go)
		lastSeen = look(action.Go)              // ✋ do
		fmt.Printf("  👀 see: %s\n\n", lastSeen) // 👀 observe → feeds next think
	}
	fmt.Println("Ran out of steps.")
}
