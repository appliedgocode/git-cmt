# git-cmt

An AI-powered Git commit message generator that analyzes your staged changes and creates conventional commit messages using Claude AI.

## Overview

git-cmt automatically generates meaningful commit messages based on your staged changes using Anthropic's Claude AI. It follows the [Conventional Commits](https://www.conventionalcommits.org/) specification and provides an interactive commit experience.

This project was created for the [2025-08-17 issue of The Applied Go Weekly Newsletter](https://newsletter.appliedgo.net/archive/2025-08-17-commit-messages-that-write-themselves/)

Feel free to play with the code to change providers and models, try different prompts, and more. It's just around 110 lines of code. 

BTW, parts of this README are LLM-generated, so I apologize in advance for the verbosity, but I didn't want to cut down the amount of information too much, to increase the chance that you'll be finding what you're looking for. 

## Features

- 🤖 **AI-powered**: Uses Claude 3.5 Haiku to analyze code changes (but you can switch models and providers using LangChainGo's selection of provider-specific sub-packages)
- 📝 **Conventional Commits**: Generates messages following the `type(scope): description` format
- 🎯 **Smart Analysis**: Understands code changes and creates contextually appropriate messages
- ⚡ **Interactive**: Opens your editor for final review before committing
- 🔍 **Diff-aware**: Only analyzes currently staged changes
- 📏 **Length-aware**: Keeps commit messages concise (50 chars max for description)

## Installation

### Prerequisites

- Go 1.25.0 or later (surely works with earlier versions, too, just change the Go version in `go.mod` and test it out)
- Git
- An Anthropic API key (or, if you change the code to use another provider, their API key. Or, if you use a local AI... no API key)

### Build from source

#### Without cloning the project

```bash
go install github.com/appliedgocode/git-cmt
```

#### With cloning the project

```bash
git clone https://github.com/appliedgocode/git-cmt
cd git-cmt
go build
```

### Install globally

```bash
go install
```

## Usage

### Setup

1. Set your Anthropic API key:
   ```bash
   export ANTHROPIC_API_KEY="your-api-key-here"
   ```

### Basic Usage

1. Stage your changes:
   ```bash
   git add <whatever you have to add>
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

If you're new to Conventional Commits, [start here](https://www.conventionalcommits.org/en/v1.0.0/).

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

### Customization: Play with the code!

- **Change the provider and/or model**
	- Change `anthropic.WithModel()` to use different Claude models
	- Use a different LangChainGo subpackage (`openai`, `ollama`,...) to switch the provider
	- Use the `openai` package and `openai.WithURL()` to use an OpenAI-compatible API of a third-party provider
- **Update the prompt template in `generateMessage()` for custom behavior**
- **Let `git-cmt` present an empty commit message if the call to the LLM fails,** instead of erroring out
- **Make provider and model configurable** in a config file
- **Use a secrets manager to securely obtain the API key** instead of an environment variable

## Error Handling

- **No staged changes**: Exits with helpful message
- **API key missing**: Prompts to set `ANTHROPIC_API_KEY`
- **API failures**: Provides detailed error messages
- **Invalid JSON**: Shows raw response for debugging

## Development

### Dependencies

- `github.com/tmc/langchaingo` - LLM integration
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
go build
```

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
