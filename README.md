# markd

A self-hosted bookmark manager with CLI tool, REST API, HTMX web UI, background workers, 
observable infrastructure you could deploy in a homelab.

## What it does

**CLI** (`markcli`):

- Add bookmarks with URL, title, description, and tags
- List all bookmarks
- Get a bookmark by ID
- WIP: Delete a bookmark
- WIP: Add and remove tags on existing bookmarks
- WIP: Search bookmarks by title, URL, or tag
- WIP: Concurrent URL health checking with configurable worker count,
  request timeout, retry with backoff, and failure reporting
- WIP: Import bookmarks from Netscape HTML exports, CSV, or JSON
- WIP: Export bookmarks to JSON, CSV, or Netscape HTML
- WIP: Authenticate against the server and store the API token locally

**REST API** (`markd` server — WIP):

**Web UI** (WIP):

**Background workers** (WIP):

**Deployment** (WIP):


## Installation

You need Go installed.

```bash
git clone https://github.com/def4alt/markd.git
cd markd
go build -o markcli ./cmd/markcli
go build -o markd ./cmd/markd
```

## Usage

### CLI

```bash
# Add a bookmark
markcli add https://go.dev --title "Go" --tag programming

# List bookmarks
markcli list

# Get a bookmark by ID
markcli get cfe951bb-21a6-405d-ac9e-bd5055c875a6

# Add tags
markcli tag cfe951bb-21a6-405d-ac9e-bd5055c875a6 go backend

# Search
markcli search go

# Check all bookmarks
markcli check --workers 20 --timeout 2s

# Check bookmarks with a specific tag
markcli check --tag go

# Import and export
markcli import bookmarks.html
markcli export --format json
markcli export --format csv

# Authenticate against the server
markcli login http://markd.local
```

