package main

import "fmt"

func main() {
	goal := "find a banana"
	fmt.Printf("🎯 Goal: %s\n\n", goal)

	// The places our agent can look — these are its "hands" (tools).
	places := []string{"basket", "fridge"}

	// 🔁 THE LOOP — the thing that makes this an agent.
	for _, place := range places {
		fmt.Printf("🧠 think: let me check the %s\n", place)

		found := (place == "fridge") // pretend the banana is in the fridge

		fmt.Printf("✋ do: looking in the %s...\n", place)

		if found { // 👀 saw a banana — decide to STOP
			fmt.Println("👀 see: banana is here! 🍌 Done!")
			return
		}
		fmt.Println("👀 see: no banana here. Try the next place.")
	}

	fmt.Println("Checked everywhere, no banana. 😞")
}
