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

func getStagedChanges() (string, error) {
	// Use git diff --cached -b to get staged changes ignoring whitespace
	cmd := exec.Command("git", "diff", "--cached", "-b")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}

	diff := string(output)
	if diff == "" {
		return "", fmt.Errorf("no staged changes found")
	}

	// Limit diff size for LLM (keep first 3KB)
	if len(diff) > 3072 {
		diff = diff[:3072] + "\n... (truncated)"
	}

	return diff, nil
}

type Commit struct {
	Type    string `json:"type"`    // feat, fix, docs, etc.
	Scope   string `json:"scope"`   // optional component
	Message string `json:"message"` // the actual description
}

func generateMessage(changes string) (Commit, error) {
	// Create LLM client - easily swap providers here
	llm, err := openai.New(
		openai.WithModel("gpt-5-mini-2025-08-07"),
		openai.WithToken(os.Getenv("OPENAI_API_KEY")),
	)
	if err != nil {
		return Commit{}, fmt.Errorf("failed to create LLM client: %w", err)
	}

	prompt := fmt.Sprintf(`You are a git commit message generator.
    Analyze changes and output JSON with:
    - type: feat|fix|docs|style|refactor|test|chore
    - scope: affected component (optional)
    - message: clear description (50 chars max)
    
    Changes:
    %s
    
    Return ONLY valid JSON, no other text.`, changes)

	resp, err := llms.GenerateFromSinglePrompt(
		context.Background(),
		llm,
		prompt,
		llms.WithTemperature(1),
	)
	if err != nil {
		return Commit{}, fmt.Errorf("LLM request failed: %w", err)
	}

	var commit Commit
	if err := json.Unmarshal([]byte(resp), &commit); err != nil {
		return Commit{}, fmt.Errorf("failed to parse JSON response: %w (raw response: %q)", err, resp)
	}

	return commit, nil
}

func main() {
	changes, err := getStagedChanges()
	if err != nil {
		log.Fatalf("Failed to get staged changes: %v", err)
	}

	log.Printf("Staged diff found: %q", changes)
	log.Printf("Generating message for changes")

	commit, err := generateMessage(changes)
	if err != nil {
		log.Fatalf("Failed to generate commit message: %v", err)
	}

	log.Printf("Parsed commit: %+v", commit)

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
