# jira-cli

A fast, lightweight command-line interface for Jira operations built in Go.

## Features

- **Issue Management**: Get, create, update, and search issues
- **Comments**: List and add comments to issues
- **Transitions**: View available transitions and change issue status
- **Sprint & My Issues**: Quick access to sprint issues and personal assignments
- **Multiple Output Formats**: Text (default) and JSON
- **Single Binary**: No runtime dependencies required

## Installation

### Prerequisites

- Go 1.21 or higher
- `JIRA_TOKEN` environment variable set with your Jira API token

### Build from Source

```bash
# Clone or navigate to the project
cd jira-cli

# Build the binary
make build

# Or build and install to /usr/local/bin
make install
```

### Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `JIRA_TOKEN` | Yes | - | Jira API Bearer token |
| `JIRA_BASE_URL` | Yes | - | Jira instance URL (e.g., `https://jira.company.com`) |

## Usage

### Get Issue Details

```bash
jira-cli get PROJECT-123
jira-cli get PROJECT-123 -o json    # JSON output
```

### Create Issue

```bash
jira-cli create --project MYPROJ --type Story --summary "New feature"
jira-cli create -p MYPROJ -t Bug -s "Fix login" -d "Description here"

# Read description from stdin (useful for long descriptions)
echo "Long description..." | jira-cli create -p MYPROJ -t Story -s "Title" --stdin
```

### Update Issue

```bash
jira-cli update PROJECT-123 --summary "Updated title"
jira-cli update PROJECT-123 --description "New description"

# Read description from stdin
cat description.txt | jira-cli update PROJECT-123 --stdin
```

### Search Issues (JQL)

```bash
jira-cli search "project = MYPROJ AND status = 'In Progress'"
jira-cli search "assignee = currentUser()" --max 100
```

### My Issues

```bash
jira-cli my-issues                 # All my open issues
jira-cli my-issues -o json         # JSON output
```

### Sprint Issues

```bash
jira-cli sprint --project MYPROJ   # Issues in current sprint
jira-cli sprint -p MYPROJ          # Short form
```

### Comments

```bash
jira-cli comment list PROJECT-123
jira-cli comment add PROJECT-123 "This is my comment"
```

### Transitions

```bash
jira-cli transition list PROJECT-123              # See available transitions
jira-cli transition do PROJECT-123 "In Progress"  # Change status
```

## Output Formats

### Text (Default)

Human-readable markdown-style tables:

```
## PROJECT-123

| Campo | Valor |
|-------|-------|
| **Summary** | Issue title |
| **Status** | In Progress |
```

### JSON

Machine-readable JSON output:

```bash
jira-cli get PROJECT-123 -o json
```

```json
{
  "key": "PROJECT-123",
  "fields": {
    "summary": "Issue title",
    "status": { "name": "In Progress" }
  }
}
```

## Project Structure

```
jira-cli/
├── main.go                 # Entry point
├── go.mod                  # Go module definition
├── Makefile                # Build automation
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go             # Root command and config
│   ├── get.go              # Get issue
│   ├── create.go           # Create issue
│   ├── update.go           # Update issue
│   ├── search.go           # Search with JQL
│   ├── myissues.go         # My assigned issues
│   ├── sprint.go           # Sprint issues
│   ├── comment.go          # Comment operations
│   └── transition.go       # Transition operations
├── internal/
│   └── jira/               # Jira API client
│       ├── client.go       # HTTP client
│       ├── issues.go       # Issue operations
│       ├── comments.go     # Comment operations
│       └── transitions.go  # Transition operations
└── bin/
    └── jira-cli            # Compiled binary
```

## Development

### Building

```bash
make build          # Build binary to bin/jira-cli
make install        # Build and install to /usr/local/bin
make clean          # Remove build artifacts
```

### Testing

```bash
make test           # Run tests
make test-coverage  # Run tests with coverage
```

### Cross-compilation

```bash
make build-all      # Build for macOS, Linux, Windows
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) v1.10.2 - CLI framework
- [viper](https://github.com/spf13/viper) v1.21.0 - Configuration management

## License

MIT