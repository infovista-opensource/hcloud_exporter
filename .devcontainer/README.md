# Development Container for hcloud_exporter

This directory contains configuration for a Visual Studio Code development container that provides a consistent development environment for the hcloud_exporter project.

## Features

- Go 1.23 development environment
- Pre-installed Go tools:
  - gopls (Go language server)
  - staticcheck
  - delve (debugger)
  - revive (linter)
  - golangci-lint
- Git and GitHub CLI
- SSH agent forwarding
- Common development utilities

## Getting Started

1. Install [Visual Studio Code](https://code.visualstudio.com/)
2. Install the [Remote - Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension
3. Clone the repository: `git clone https://github.com/promhippie/hcloud_exporter.git`
4. Open the project in VS Code
5. When prompted, click "Reopen in Container" or run the "Remote-Containers: Reopen in Container" command from the command palette

## Environment Configuration

You can modify the environment variables in the `devcontainer.env` file to customize your development environment. For example, to test with a real Hetzner Cloud account, add your API token there.

## Building and Running

Inside the container, you can use standard Go commands or the project's Makefile:

```bash
# Generate files
make generate

# Build the project
make build

# Run the exporter
./bin/hcloud_exporter
```

## Customizing the Container

If you need to customize the development container:

- Modify `devcontainer.json` to change VS Code settings, extensions, or container options
- Modify `Dockerfile` to add dependencies or change the container setup
