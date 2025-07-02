# SurveyBot Development Guide

## Tool Management

This project uses [mise](https://mise.jdx.dev/) for managing development tools and runtime versions.

### Setup

Install mise if you haven't already:
```bash
curl https://mise.jdx.dev/install.sh | sh
```

Install the project's tools:
```bash
mise install
```

### Configuration

The project's tool versions are defined in `mise.toml`. This ensures all developers use consistent versions of:
- Go runtime
- Other development tools

### Usage

Mise automatically activates the correct tool versions when you enter the project directory. You can also manually run:
```bash
mise use
```

To see current tool versions:
```bash
mise current
```