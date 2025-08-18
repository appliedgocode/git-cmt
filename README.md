# git-cmt

An AI-powered Git commit message generator that analyzes your staged changes and creates conventional commit messages using Claude AI.

## Overview

git-cmt automatically generates meaningful commit messages based on your staged changes using Anthropic's Claude AI. It follows the [Conventional Commits](https://www.conventionalcommits.org/) specification and provides an interactive commit experience.

## Features

- 🤖 **AI-powered**: Uses Claude 3.5 Haiku to analyze code changes
- 📝 **Conventional Commits**: Generates messages following the `type(scope): description` format
- 🎯 **Smart Analysis**: Understands code changes and creates contextually appropriate messages
- ⚡ **Interactive**: Opens your editor for final review before committing
- 🔍 **Diff-aware**: Only analyzes currently staged changes
- 📏 **Length-aware**: Keeps commit messages concise (50 chars max for description)

## Installation

### Prerequisites

- Go 1.25.0 or later
- Git
- Anthropic API key

### Build from source

```bash
git clone https://github.com/appliedgocode/commit-ai.git
cd commit-ai
go build -o git-cmt
```

### Install globally

```bash
# Move to a directory in your PATH
sudo mv git-cmt /usr/local/bin/
```

## Usage

### Setup

1. Set your Anthropic API key:
   ```bash
   export ANTHROPIC_API_KEY="your-api-key-here"
   ```

2. Add to your shell profile (`.bashrc`, `.zshrc`, etc.) for persistence:
   ```bash
   echo 'export ANTHROPIC_API_KEY="your-api-key-here"' >> ~/.zshrc
   source ~/.zshrc
   ```

### Basic Usage

1. Stage your changes:
   ```bash
   git add .
   ```

2. Generate and commit:
   ```bash
   ./git-cmt
   ```

3. Review and edit the generated message in your default editor
4. Save and close to complete the commit

## How it works

1. **Diff Analysis**: Reads staged changes using `git diff --cached -b`
2. **AI Processing**: Sends the diff to Claude AI with structured prompts
3. **Message Generation**: Creates a commit object with type, scope, and description
4. **Interactive Commit**: Opens your editor with the generated message for review
5. **Final Commit**: Executes `git commit` with your approved message

## Commit Message Format

The AI generates messages following this structure:

```
type(scope): description
```

**Types**: feat, fix, docs, style, refactor, test, chore  
**Scope**: Optional component/module name  
**Description**: Clear, concise summary (max 50 chars)

## Examples

### Feature Addition
```bash
$ ./git-cmt
# Generated: feat(auth): add OAuth2 login integration
```

### Bug Fix
```bash
$ ./git-cmt
# Generated: fix(api): resolve null pointer in user validation
```

### Documentation
```bash
$ ./git-cmt
# Generated: docs(readme): update installation instructions
```

## Configuration

### Environment Variables

- `ANTHROPIC_API_KEY`: Required for Claude API access
- `EDITOR`: Controls which editor opens for message review (defaults to system default)

### Customization

The AI model and prompts can be modified in `main.go`:
- Change `anthropic.WithModel()` to use different Claude models
- Update the prompt template in `generateMessage()` for custom behavior

## Error Handling

- **No staged changes**: Exits with helpful message
- **API key missing**: Prompts to set `ANTHROPIC_API_KEY`
- **API failures**: Provides detailed error messages
- **Invalid JSON**: Shows raw response for debugging

## Development

### Dependencies

- `github.com/tmc/langchaingo v0.1.13` - LLM integration
- Claude AI via Anthropic API

### Project Structure

```
├── main.go          # Core application logic
├── go.mod          # Go module definition
├── go.sum          # Dependency checksums
└── README.md       # This file
```

### Building

```bash
go build -o git-cmt
```

### Testing

```bash
go run main.go  # Run directly
go build        # Build binary
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test with real commits
5. Submit a pull request

## License

This project is open source. Check the repository for specific license details.

## Troubleshooting

### Common Issues

**"No staged changes found"**
- Ensure you've run `git add` to stage files

**"ANTHROPIC_API_KEY not set"**
- Set the environment variable with your API key

**"LLM request failed"**
- Check your internet connection
- Verify your API key is valid
- Ensure you have credits on your Anthropic account

**Editor not opening**
- Set your preferred editor: `export EDITOR="code --wait"`

### Debug Mode

Run with verbose logging to see detailed output:
```bash
./git-cmt 2>&1 | tee debug.log
```