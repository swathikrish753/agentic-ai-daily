package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"`
}
type Response struct {
	Response string `json:"response"`
}

// 🗺️ THE PLAN: we force the brain to return steps as a JSON array.
type Plan struct {
	Steps []string `json:"steps"`
}

func askBrainJSON(prompt string) (string, error) {
	reqBody := Request{Model: "llama3.2:3b", Prompt: prompt, Stream: false, Format: "json"}
	jsonBytes, _ := json.Marshal(reqBody)
	resp, err := http.Post("http://localhost:11434/api/generate",
		"application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var parsed Response
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", err
	}
	return parsed.Response, nil
}

func main() {
	goal := "Make a cup of tea"

	// ── PHASE 1: PLAN ───────────────────────────────────────────
	planPrompt := fmt.Sprintf(`You are a planning agent. Break the goal into clear, ordered steps.

Goal: "%s"

Reply ONLY as JSON in this exact shape:
{"steps": ["step 1", "step 2", "step 3"]}`, goal)

	fmt.Printf("🎯 goal: %s\n\n", goal)

	rawPlan, err := askBrainJSON(planPrompt)
	if err != nil {
		fmt.Println("❌ couldn't reach the brain:", err)
		return
	}

	var plan Plan
	if err := json.Unmarshal([]byte(rawPlan), &plan); err != nil {
		fmt.Println("⚠️ plan wasn't valid JSON:", err)
		return
	}

	fmt.Printf("🗺️  the brain made a %d-step plan:\n", len(plan.Steps))
	for i, step := range plan.Steps {
		fmt.Printf("   %d. %s\n", i+1, step)
	}

	// ── PHASE 2: EXECUTE ────────────────────────────────────────
	fmt.Println("\n🔁 executing the plan:")
	for i, step := range plan.Steps {
		// In a real agent, each step would call a tool (Day 4).
		// Here we just "do" it to show the walk-the-list structure.
		fmt.Printf("   ▶ step %d: %s ... done ✅\n", i+1, step)
	}

	fmt.Println("\n🎉 goal complete.")
}
