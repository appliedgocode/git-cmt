package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func getStagedChanges() string {
	// Use git diff --cached -b to get staged changes ignoring whitespace
	cmd := exec.Command("git", "diff", "--cached", "-b")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to get git diff: %v", err)
		return ""
	}

	diff := string(output)
	if diff == "" {
		log.Printf("No staged changes found")
		return ""
	}

	// Limit diff size for LLM (keep first 3KB)
	if len(diff) > 3072 {
		diff = diff[:3072] + "\n... (truncated)"
	}

	log.Printf("Staged diff found: %q", diff)
	return diff
}

type Commit struct {
	Type    string `json:"type"`    // feat, fix, docs, etc.
	Scope   string `json:"scope"`   // optional component
	Message string `json:"message"` // the actual description
}

func generateMessage(changes string) Commit {
	log.Printf("Generating message for changes: %q", changes)

	// Create LLM client - easily swap providers here
	llm, err := openai.New(
		openai.WithModel("gpt-5-mini-2025-08-07"),
		openai.WithToken(os.Getenv("OPENAI_API_KEY")),
	)
	if err != nil {
		log.Fatal("Failed to create LLM client:", err)
	}

	prompt := fmt.Sprintf(`You are a git commit message generator.
    Analyze changes and output JSON with:
    - type: feat|fix|docs|style|refactor|test|chore
    - scope: affected component (optional)
    - message: clear description (50 chars max)
    
    Changes:
    %s
    
    Return ONLY valid JSON, no other text.`, changes)

	log.Printf("Sending prompt to LLM: %q", prompt)

	resp, err := llms.GenerateFromSinglePrompt(
		context.Background(),
		llm,
		prompt,
		llms.WithTemperature(1),
	)
	if err != nil {
		log.Printf("LLM request failed: %v", err)
		return Commit{}
	}

	log.Printf("LLM response: %q", resp)

	var commit Commit
	if err := json.Unmarshal([]byte(resp), &commit); err != nil {
		log.Printf("Failed to parse JSON response: %v", err)
		log.Printf("Raw response was: %q", resp)
		return Commit{}
	}

	log.Printf("Parsed commit: %+v", commit)
	return commit
}

func main() {
	changes := getStagedChanges()
	if changes == "" {
		log.Fatal("No staged changes")
	}

	commit := generateMessage(changes)

	// Build conventional format
	output := commit.Type
	if commit.Scope != "" {
		output += "(" + commit.Scope + ")"
	}
	output += ": " + commit.Message

	fmt.Println("\n✨ Generated commit message:")
	fmt.Println(output)

	// Optional: copy to clipboard
	fmt.Println("\nRun: git commit -m \"" + output + "\"")
}
