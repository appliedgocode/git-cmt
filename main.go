package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/go-git/go-git/v5"
    "github.com/tmc/langchaingo/llms"
    "github.com/tmc/langchaingo/llms/openai"
)

func getStagedChanges() string {
    repo, err := git.PlainOpen(".")
    if err != nil {
        return ""
    }
    
    worktree, _ := repo.Worktree()
    status, _ := worktree.Status()
    
    var diffSummary strings.Builder
    for filepath, fileStatus := range status {
        if fileStatus.Staging != git.Unmodified {
            diffSummary.WriteString(fmt.Sprintf("- %s: %v\n", 
                filepath, fileStatus.Staging))
        }
    }
    
    return diffSummary.String()
}

type Commit struct {
    Type    string `json:"type"`    // feat, fix, docs, etc.
    Scope   string `json:"scope"`   // optional component
    Message string `json:"message"` // the actual description
}

func generateMessage(changes string) Commit {
    // Create LLM client - easily swap providers here
    llm, err := openai.New(
        openai.WithModel("gpt-oss-120b"),
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
    
    resp, _ := llms.GenerateFromSinglePrompt(
        context.Background(),
        llm,
        prompt,
        llms.WithTemperature(0.3), // Lower = more consistent
    )
    
    var commit Commit
    json.Unmarshal([]byte(resp), &commit)
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