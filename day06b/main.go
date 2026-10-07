package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}
type ChatResponse struct {
	Message Message `json:"message"`
}

const memoryFile = "day06b/memory.json" // 📔 the diary, on disk

// 📔 Load long-term memory (returns [] if the diary doesn't exist yet).
func loadMemory() []string {
	data, err := os.ReadFile(memoryFile)
	if err != nil {
		return []string{} // first run ever — empty diary
	}
	var facts []string
	json.Unmarshal(data, &facts)
	return facts
}

// 📔 Save a new fact to the diary (append + write back).
func saveMemory(facts []string) {
	data, _ := json.MarshalIndent(facts, "", "  ")
	os.WriteFile(memoryFile, data, 0644)
}

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
	// 1. 📔 Read the diary from disk (survives restarts).
	facts := loadMemory()
	fmt.Printf("📔 loaded %d fact(s) from long-term memory: %v\n\n", len(facts), facts)

	var history []Message

	// 2. Inject what we remember as a system message (put the page on the desk).
	if len(facts) > 0 {
		known := "Known facts about the user:\n"
		for _, f := range facts {
			known += "- " + f + "\n"
		}
		history = append(history, Message{Role: "system", Content: known})
	}

	// 3. Ask something that needs a fact from a PREVIOUS run.
	question := "Based on what you know about me, what do I love?"
	history = append(history, Message{Role: "user", Content: question})
	fmt.Printf("🙋 user: %s\n", question)

	reply, err := chat(history)
	if err != nil {
		fmt.Println("❌ couldn't reach the brain:", err)
		return
	}
	fmt.Printf("🧠 brain: %s\n\n", reply.Content)

	// 4. 📔 Save a new fact to the diary for FUTURE runs (only once).
	newFact := "The user's name is Ravi and they love Go."
	if len(facts) == 0 {
		facts = append(facts, newFact)
		saveMemory(facts)
		fmt.Println("💾 saved a fact to long-term memory. Run me again!")
	}
}
