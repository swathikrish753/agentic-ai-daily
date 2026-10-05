package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// What we SEND. New field: Format forces Ollama to reply in valid JSON.
type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format"` // "json" = reply MUST be valid JSON
}

// Ollama wraps the model's text answer in .Response (still just a string).
type Response struct {
	Response string `json:"response"`
}

// 🧾 THE FORM: the shape we force the brain's answer into.
type ToolChoice struct {
	Tool   string `json:"tool"`
	Reason string `json:"reason"`
}

func askBrainJSON(prompt string) (string, error) {
	reqBody := Request{
		Model:  "llama3.2:3b",
		Prompt: prompt,
		Stream: false,
		Format: "json", // 👈 the magic flag
	}
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
	return parsed.Response, nil // this string is now itself JSON
}

func main() {
	userAsk := "Ugh, I'm bored. Make me laugh."

	// We describe the form (the shape) right in the prompt.
	prompt := fmt.Sprintf(`You are an agent. Choose ONE tool for the user's request.

Available tools:
- get_weather: check today's weather
- tell_joke: tell a short funny joke

User request: "%s"

Reply ONLY as JSON in this exact shape:
{"tool": "<tool name>", "reason": "<short why>"}`, userAsk)

	fmt.Printf("🙋 user: %s\n\n", userAsk)

	rawAnswer, err := askBrainJSON(prompt)
	if err != nil {
		fmt.Println("❌ couldn't reach the brain:", err)
		return
	}
	fmt.Printf("🧠 brain's raw JSON: %s\n", rawAnswer)

	// 🎯 THE PAYOFF: parse into a struct. No string-guessing anywhere.
	var choice ToolChoice
	if err := json.Unmarshal([]byte(rawAnswer), &choice); err != nil {
		fmt.Println("⚠️ brain didn't return valid JSON:", err)
		return
	}

	// Now we have clean, typed data.
	fmt.Printf("✅ parsed → tool=%q, reason=%q\n", choice.Tool, choice.Reason)

	switch choice.Tool {
	case "tell_joke":
		fmt.Println("👀 result: Why do programmers prefer dark mode? Light attracts bugs.")
	case "get_weather":
		fmt.Println("👀 result: It's 28°C and sunny in Bangalore.")
	default:
		fmt.Printf("🤷 unknown tool: %q\n", choice.Tool)
	}
}
