# gofish

Tool to automatically run various checks, for use when an LLM agent such as
[Claude Code](https://docs.anthropic.com/en/docs/claude-code) thinks it's
finished working.

## Available Checkers

### Go projects

- gofmt
- go vet
- go test
- staticcheck (only if there's a custom config file)

### JavaScript/TypeScript projects

- svelte-check (if listed as a dependency)
- prettier (if listed as a dependency)

## Installation

### Install from source

```bash
go install github.com/csmith/gofish/cmd/gofish@latest
```

### Build locally

```bash
git clone https://github.com/csmith/gofish.git
cd gofish
go build ./cmd/gofish
```

## Usage

### Command Line

Run `gofish` in any directory to scan its subdirectories. It can handle 'nested'
projects (e.g. having a `project.json` in a `frontend` directory).


### Claude Code Integration

gofish is designed to be run as a `Stop` hook in Claude Code.

#### Setup Hook

1. **Install gofish globally**:
   ```bash
   go install github.com/yourusername/gofish/cmd/gofish@latest
   ```

2. **Configure Claude Code settings**: Add the following to your Claude Code `settings.json`:
   ```json
   {
     "hooks": {
       "Stop": [
         {
           "matcher": "",
           "hooks": [
             {
               "type": "command",
               "command": "gofish"
             }
           ]
         }
       ]
     }
   }
   ```
   Alternatively, use the `/hooks` slash command to manually set it up.

The hook will automatically run gofish whenever Claude Code stops working on
your code, ensuring it doesn't leave anything in a mess.

For more information on Claude Code hooks, see the [official documentation](https://docs.anthropic.com/en/docs/claude-code/hooks).

## Exit Codes

- `0`: Success, no issues found
- `1`: Error running checkers
- `2`: Code quality issues detected