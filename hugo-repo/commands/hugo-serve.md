---
name: hugo-serve
description: Start the Hugo development server
arguments:
  - name: flags
    description: "Additional hugo server flags"
    required: false
---

# Start Hugo Development Server

Start `hugo server` with appropriate flags for the current project.

## Process

1. **Verify Hugo is installed:**
   ```bash
   hugo version
   ```
   If not installed, suggest installation method for the current platform.

2. **Resolve Hugo modules** if `go.mod` exists:
   ```bash
   hugo mod get
   ```

3. **Start the server:**
   ```bash
   hugo server --buildDrafts --navigateToChanged
   ```

4. **Report** the local URL (typically `http://localhost:1313/`) and any warnings.

## Usage

```
/hugo-serve [flags]
```

## Flags

| Flag | Purpose |
|------|---------|
| `--buildDrafts` | Include draft content (default: on) |
| `--navigateToChanged` | Browser navigates to changed page |
| `--disableFastRender` | Full rebuild on every change (use if seeing stale content) |
| `--port 1314` | Use a different port |

## Troubleshooting

If the server fails to start:

- Check `hugo.toml` exists and is valid
- If using modules, ensure `go.mod` exists (`hugo mod init`)
- If theme not found, run `hugo mod get`
- Check for port conflicts (`--port` to change)
