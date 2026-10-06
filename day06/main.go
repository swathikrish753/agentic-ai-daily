package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// One line in the notepad: who said it + what they said.
type Message struct {
	Role    string `json:"role"` // "user" | "assistant" | "system"
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"` // 👈 the WHOLE notepad, every time
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Message Message `json:"message"` // the brain's reply
}

// Send the ENTIRE history and get the next reply.
func chat(history []Message) (Message, error) {
	reqBody := ChatRequest{Model: "llama3.2:3b", Messages: history, Stream: false}
	jsonBytes, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:11434/api/chat",
		"application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var parsed ChatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return Message{}, err
	}
	return parsed.Message, nil
}

func main() {
	// 📝 THE NOTEPAD. This slice IS the memory.
	var history []Message

	// We'll say two things; the 2nd depends on remembering the 1st.
	userTurns := []string{
		"Hi! My name is Ravi and I love Go.",
		"What's my name and what do I love?",
	}

	for _, text := range userTurns {
		// 1. Add the user's line to the notepad.
		history = append(history, Message{Role: "user", Content: text})
		fmt.Printf("🙋 user: %s\n", text)

		// 2. Send the WHOLE notepad, get a reply.
		reply, err := chat(history)
		if err != nil {
			fmt.Println("❌ couldn't reach the brain:", err)
			fmt.Println("   Is Ollama running? Try: ollama serve")
			return
		}

		// 3. Add the brain's reply to the notepad too (so it remembers next turn).
		history = append(history, reply)
		fmt.Printf("🧠 brain: %s\n\n", reply.Content)
	}

	fmt.Printf("📝 notepad now holds %d messages\n", len(history))
}
