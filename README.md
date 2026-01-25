# atlassian

A fast, lightweight command-line interface for Atlassian products (Jira & Confluence) built in Go.

## Features

### Jira
- **Issue Management**: Get, create, update, and search issues
- **Assignments**: Assign/unassign users to issues
- **Story Points**: Set story points on issues
- **Sprints**: List boards, sprints, and move issues between sprints
- **Users**: Search for users to get accountId
- **Fields**: Discover custom field IDs (Story Points, Sprint, etc.)
- **Comments**: List and add comments to issues
- **Transitions**: View available transitions and change issue status
- **Sprint & My Issues**: Quick access to sprint issues and personal assignments

### Confluence
- **Space Management**: List spaces and get space details
- **Page Operations**: List, get, create, and update pages
- **Search**: Search content using CQL (Confluence Query Language)

### General
- **Multiple Output Formats**: Text (default) and JSON
- **Single Binary**: No runtime dependencies required

## Installation

### Prerequisites

- Go 1.21 or higher

### Build from Source

```bash
cd atlassian

make build

make install
```

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `JIRA_TOKEN` | For Jira | Jira API Bearer token |
| `JIRA_BASE_URL` | For Jira | Jira instance URL (e.g., `https://jira.company.com`) |
| `CONFLUENCE_TOKEN` | For Confluence | Confluence API Bearer token |
| `CONFLUENCE_BASE_URL` | For Confluence | Confluence instance URL (e.g., `https://confluence.company.com`) |

## Usage

### Jira Commands

All Jira commands are under the `jira` subcommand:

#### Get Issue Details

```bash
atlassian jira get PROJECT-123
atlassian jira get PROJECT-123 -o json
```

#### Create Issue

```bash
atlassian jira create --project MYPROJ --type Story --summary "New feature"
atlassian jira create -p MYPROJ -t Bug -s "Fix login" -d "Description here"

echo "Long description..." | atlassian jira create -p MYPROJ -t Story -s "Title" --stdin
```

#### Update Issue

```bash
atlassian jira update PROJECT-123 --summary "Updated title"
atlassian jira update PROJECT-123 --description "New description"
atlassian jira update PROJECT-123 --assignee user@email.com
atlassian jira update PROJECT-123 --points 5
atlassian jira update PROJECT-123 --sprint 123

# Combine multiple updates
atlassian jira update PROJECT-123 -a user@email.com -p 5 --sprint 123

cat description.txt | atlassian jira update PROJECT-123 --stdin
```

#### Assign Issue

```bash
atlassian jira assign PROJECT-123 user@email.com
atlassian jira assign PROJECT-123 5b10ac8d82e05b22cc7d4ef5
atlassian jira assign PROJECT-123 --unassign
```

#### Search Users

```bash
atlassian jira users --query "john"
atlassian jira users -q "john@company.com" -o json
```

#### List Boards

```bash
atlassian jira boards
atlassian jira boards --project MYPROJ
atlassian jira boards -o json
```

#### List Sprints

```bash
atlassian jira sprints --board 123
atlassian jira sprints -b 123 --state active
atlassian jira sprints -b 123 --state future -o json
```

#### List Fields

```bash
atlassian jira fields
atlassian jira fields --name "Story Points"
atlassian jira fields --custom -o json
```

#### Search Issues (JQL)

```bash
atlassian jira search "project = MYPROJ AND status = 'In Progress'"
atlassian jira search "assignee = currentUser()" --max 100
```

#### My Issues

```bash
atlassian jira my-issues
atlassian jira my-issues -o json
```

#### Sprint Issues

```bash
atlassian jira sprint --project MYPROJ
atlassian jira sprint -p MYPROJ
```

#### Comments

```bash
atlassian jira comment list PROJECT-123
atlassian jira comment add PROJECT-123 "This is my comment"
```

#### Transitions

```bash
atlassian jira transition list PROJECT-123
atlassian jira transition do PROJECT-123 "In Progress"
```

### Confluence Commands

All Confluence commands are under the `confluence` subcommand (alias: `conf`):

#### List Spaces

```bash
atlassian confluence spaces
atlassian conf spaces --limit 50
atlassian conf spaces -o json
```

#### Get Space Details

```bash
atlassian conf spaces MYSPACE
atlassian conf spaces MYSPACE -o json
```

#### List Pages in a Space

```bash
atlassian conf pages --space MYSPACE
atlassian conf pages -s MYSPACE --limit 50
```

#### Get Page by ID

```bash
atlassian conf get 123456
atlassian conf get 123456 -o json
atlassian conf get 123456 --body-format storage
atlassian conf get 123456 --body-format view
```

#### Search Content (CQL)

```bash
atlassian conf search "space=MYSPACE"
atlassian conf search "type=page AND title~'Documentation'"
atlassian conf search "text~'API'" --limit 50
```

#### Create Page

```bash
atlassian conf create --space MYSPACE --title "New Page"
atlassian conf create -s MYSPACE -t "Child Page" --parent 123456

echo "<p>Page content</p>" | atlassian conf create -s MYSPACE -t "Page" --stdin
```

#### Update Page

```bash
atlassian conf update 123456 --title "New Title"
atlassian conf update 123456 --message "Updated via CLI"

echo "<p>New content</p>" | atlassian conf update 123456 --stdin
```

## Output Formats

### Text (Default)

Human-readable markdown-style tables:

```
| Field        | Value                |
| ------------ | -------------------- |
| ID           | 123456               |
| Title        | My Page              |
| Status       | current              |
```

### JSON

Machine-readable JSON output:

```bash
atlassian jira get PROJECT-123 -o json
atlassian conf spaces MYSPACE -o json
```

## Project Structure

```
atlassian/
├── main.go
├── go.mod
├── Makefile
├── cmd/
│   ├── root.go
│   ├── jira/
│   │   ├── jira.go
│   │   ├── get.go
│   │   ├── create.go
│   │   ├── update.go
│   │   ├── search.go
│   │   ├── myissues.go
│   │   ├── sprint.go
│   │   ├── sprints.go
│   │   ├── boards.go
│   │   ├── assign.go
│   │   ├── users.go
│   │   ├── fields.go
│   │   ├── comment.go
│   │   └── transition.go
│   └── confluence/
│       ├── confluence.go
│       ├── spaces.go
│       ├── pages.go
│       ├── get.go
│       ├── search.go
│       ├── create.go
│       └── update.go
├── internal/
│   ├── jira/
│   │   ├── client.go
│   │   ├── issues.go
│   │   ├── comments.go
│   │   ├── transitions.go
│   │   ├── users.go
│   │   ├── fields.go
│   │   └── agile.go
│   └── confluence/
│       ├── client.go
│       ├── spaces.go
│       ├── pages.go
│       └── search.go
└── bin/
    └── atlassian
```

## Development

### Building

```bash
make build          # Build binary to bin/atlassian
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

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [viper](https://github.com/spf13/viper) - Configuration management

## License

MIT
