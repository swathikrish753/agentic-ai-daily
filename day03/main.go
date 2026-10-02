package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// The shape of the note we SLIDE UNDER THE DOOR.
type Request struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// The shape of the note we GET BACK (we only care about one field).
type Response struct {
	Response string `json:"response"`
}

// askBrain = slide a note under the door, return what comes back.
func askBrain(prompt string) (string, error) {
	// 1. Build the note as a struct, then turn it into JSON bytes.
	reqBody := Request{Model: "llama3.2:3b", Prompt: prompt, Stream: false}
	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 2. Slide it under the door (POST to Ollama's local address).
	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 3. Pick up the note that came back (read the raw bytes).
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 4. Unpack it and pull out just the text.
	var parsed Response
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", err
	}
	return parsed.Response, nil
}

func main() {
	prompt := "In one short sentence, what is a banana?"
	fmt.Printf("📝 note to brain: %s\n\n", prompt)

	answer, err := askBrain(prompt)
	if err != nil {
		fmt.Println("❌ couldn't reach the brain:", err)
		fmt.Println("   Is Ollama running? Try: ollama serve")
		return
	}
	fmt.Printf("🧠 brain says: %s\n", answer)
}
