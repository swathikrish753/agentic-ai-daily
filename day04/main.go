package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}
type Response struct {
	Response string `json:"response"`
}

// Day 3's brain, unchanged: text in -> text out.
func askBrain(prompt string) (string, error) {
	reqBody := Request{Model: "llama3.2:3b", Prompt: prompt, Stream: false}
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

// A Tool = a name, a description the BRAIN reads, and the real code we run.
type Tool struct {
	Name        string
	Description string
	Run         func() string
}

func main() {
	// 🧰 The agent's toolbox.
	tools := []Tool{
		{
			Name:        "get_weather",
			Description: "check today's weather",
			Run:         func() string { return "It's 28°C and sunny in Bangalore. ☀️" },
		},
		{
			Name:        "tell_joke",
			Description: "tell a short funny joke",
			Run:         func() string { return "Why do programmers prefer dark mode? Because light attracts bugs. 🐛" },
		},
	}

	// 👇 Change this line and watch the brain pick a different tool.
	userAsk := "Ugh, I'm bored. Make me laugh."

	// Build the menu of tools for the brain to read.
	var menu strings.Builder
	for _, t := range tools {
		fmt.Fprintf(&menu, "- %s: %s\n", t.Name, t.Description)
	}

	prompt := fmt.Sprintf(`You are an agent that must choose ONE tool for the user's request.

Available tools:
%s
User request: "%s"

Reply with ONLY the tool name, nothing else.`, menu.String(), userAsk)

	fmt.Printf("🙋 user asks: %s\n\n", userAsk)

	// 🧠 THINK: the LLM picks a tool.
	choice, err := askBrain(prompt)
	if err != nil {
		fmt.Println("❌ couldn't reach the brain:", err)
		return
	}
	choice = strings.ToLower(strings.TrimSpace(choice))
	fmt.Printf("🧠 brain picked: %q\n", choice)

	// 🔎 Find the tool the brain named.
	var picked *Tool
	for i := range tools {
		if strings.Contains(choice, tools[i].Name) {
			picked = &tools[i]
			break
		}
	}
	if picked == nil {
		fmt.Println("🤷 brain didn't name a known tool. (We make this reliable on Day 5.)")
		return
	}

	// ✋ DO: run the REAL Go code behind the chosen tool.
	fmt.Printf("✋ running tool: %s\n", picked.Name)
	fmt.Printf("👀 result: %s\n", picked.Run())
}
