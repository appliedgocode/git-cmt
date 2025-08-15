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
	// git diff --cached -b gets staged changes, ignoring whitespace
	cmd := exec.Command("git", "diff", "--cached", "-b")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}

	diff := string(output)
	if diff == "" {
		return "", fmt.Errorf("no staged changes found")
	}

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

func APIToken(path string) string {
	gopass := exec.Command("gopass", path)
	token, _ := gopass.Output() // leave error handling to caller
	return string(token)
}

func generateMessage(changes string) (Commit, error) {
	// Easily swap providers here by using another subpackage
	llm, err := openai.New(
		openai.WithBaseURL("https://api.studio.nebius.com/v1/"),
		openai.WithModel("mistralai/Devstral-Small-2505"),
		openai.WithToken(APIToken("api/nebius/cli")),
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
		llms.WithTemperature(0),
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

	log.Printf("Staged diff found; generating message for changes")

	commit, err := generateMessage(changes)
	message := ""
	if err != nil {
		log.Printf("Failed to generate commit message: %v", err)
	} else {
		message = commit.Type
		if commit.Scope != "" {
			message += "(" + commit.Scope + ")"
		}
		message += ": " + commit.Message
		log.Printf("Parsed commit: %+v", commit)
	}

	var cmd *exec.Cmd
	if len(message) == 0 {
		cmd = exec.Command("git", "commit", "-e")
	} else {
		cmd = exec.Command("git", "commit", "-e", "-m", message)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed committing the changes: %s", err)
	}
}
