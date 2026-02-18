---
title: "Commands Reference"
description: "Technical specifications for all 5 hugo-repo slash commands including arguments, behavior, and output"
weight: 3
---

# Commands Reference

The hugo-repo plugin includes 5 slash commands. Commands provide quick, targeted actions that users invoke directly. They follow a defined process and produce predictable output.

---

## /hugo-init

| Field | Value |
|-------|-------|
| Purpose | Scaffold a Hugo site for the current repository |

### Arguments

| Name | Required | Description |
|------|----------|-------------|
| theme | No | Theme to install (e.g., `ananke`, `book`, `docsy`) |

### Syntax

```
/hugo-init [theme]
```

### Examples

```
/hugo-init              # Analyze repo and scaffold with default theme
/hugo-init book         # Scaffold with Hugo Book theme
/hugo-init docsy        # Scaffold with Docsy theme
```

### Process

1. Scan the repository for directories containing markdown documentation (looks for `docs/` directories at any depth, `*.md` files in common locations)
2. Identify the repository type (monorepo, plugin marketplace, tool collection, single project)
3. Report findings to the user: discovered directories, recommended mounts, proposed hierarchy
4. Run `hugo new site . --force` and `hugo mod init github.com/<detected-from-git-remote>`
5. Generate `hugo.toml` with standard module mounts and custom mounts for each docs directory
6. Create section index pages: `content/_index.md`, parent section indexes; recursively scan each mounted directory and create `_index.md` for any subdirectory containing `.md` files but lacking an index; all generated indexes include `bookCollapseSection: true` for collapsible navigation
7. Create render hook at `layouts/_default/_markup/render-link.html` for portable links (handles `.md` extensions and resolves paths)
8. Install theme if specified or user accepts recommendation
9. Update `.gitignore` with `public/`, `resources/`, `.hugo_build.lock`
10. Verify by running `hugo server --buildDrafts` and reporting the result

### Output

- `hugo.toml` configuration file
- `go.mod` and `go.sum` module files
- `_index.md` section index pages for all sections (with collapsible navigation enabled)
- `layouts/_default/_markup/render-link.html` render hook
- Updated `.gitignore`
- Theme module dependency (if theme specified)

---

## /hugo-serve

| Field | Value |
|-------|-------|
| Purpose | Start the Hugo development server |

### Arguments

| Name | Required | Description |
|------|----------|-------------|
| flags | No | Additional hugo server flags |

### Syntax

```
/hugo-serve [flags]
```

### Examples

```
/hugo-serve                         # Start with default flags
/hugo-serve --disableFastRender     # Full rebuild on every change
/hugo-serve --port 1314             # Use a different port
```

### Process

1. Verify Hugo is installed (run `hugo version`). If not installed, suggest installation method for the current platform.
2. Resolve Hugo modules if `go.mod` exists (run `hugo mod get`)
3. Start the server: `hugo server --buildDrafts --navigateToChanged`
4. Report the local URL (typically `http://localhost:1313/`) and any warnings

### Default flags

| Flag | Purpose |
|------|---------|
| `--buildDrafts` | Include draft content (enabled by default) |
| `--navigateToChanged` | Browser navigates to changed page (enabled by default) |

### Optional flags

| Flag | Purpose |
|------|---------|
| `--disableFastRender` | Full rebuild on every change (use if seeing stale content) |
| `--port N` | Use a different port |

### Output

- Running Hugo development server process
- Local URL for browser access
- Any build warnings or errors

---

## /hugo-deploy

| Field | Value |
|-------|-------|
| Purpose | Set up or trigger Hugo site deployment |

### Arguments

| Name | Required | Description |
|------|----------|-------------|
| target | No | Deployment target: `pages` or `s3` (default: `pages`) |

### Syntax

```
/hugo-deploy [target]
```

### Examples

```
/hugo-deploy          # Check existing or set up GitHub Pages
/hugo-deploy pages    # Set up GitHub Pages deployment
/hugo-deploy s3       # Set up S3 deployment with OIDC
```

### Process (no existing workflow)

1. Check for existing deployment workflow at `.github/workflows/hugo-*.yml`
2. Ask user: GitHub Pages or S3
3. Generate the appropriate GitHub Actions workflow from template (`github-pages.yml.tmpl` or `s3-deploy.yml.tmpl`)
4. For S3: also generate the `[deployment]` section in `hugo.toml`
5. Commit the workflow file
6. Explain next steps (GitHub Pages source setting or AWS OIDC setup)

### Process (existing workflow)

1. Validate the workflow configuration
2. Check path filters include all mounted docs directories
3. Verify Hugo version is current
4. Report status and explain how to trigger deployment

### Output

- `.github/workflows/hugo-deploy.yml` (GitHub Pages target)
- `.github/workflows/hugo-s3-deploy.yml` (S3 target)
- Updated `hugo.toml` with deployment configuration (S3 target only)
- Post-setup instructions

### Post-setup requirements

**GitHub Pages**:

1. Settings > Pages > Source > set to "GitHub Actions"
2. Push a content change to trigger the first deployment

**S3**:

1. Set up AWS OIDC provider and IAM role
2. Add `CLOUDFRONT_DISTRIBUTION_ID` as a repository variable
3. Replace `ACCOUNT_ID` in the workflow with the AWS account ID
4. Push a commit to trigger the first deployment

---

## /hugo-add-section

| Field | Value |
|-------|-------|
| Purpose | Mount a new directory as a content section in the Hugo site |

### Arguments

| Name | Required | Description |
|------|----------|-------------|
| path | Yes | Path to the directory to mount (e.g., `new-plugin/docs`) |
| section | No | Target section path (e.g., `plugins/new-plugin`) |

### Syntax

```
/hugo-add-section <path> [section]
```

### Examples

```
/hugo-add-section new-plugin/docs                    # Auto-detect section
/hugo-add-section new-plugin/docs plugins/new-plugin # Explicit section
/hugo-add-section services/auth/docs services/auth   # Service docs
```

### Process

1. Validate the source directory exists and contains markdown files
2. Determine the target section: use the `section` argument if provided, otherwise infer from the directory path
3. Add mount entry to `hugo.toml`:
   ```toml
   [[module.mounts]]
   source = 'new-plugin/docs'
   target = 'content/plugins/new-plugin'
   ```
4. Create `_index.md` in the source directory if it does not exist
5. Ensure parent section has `_index.md` (e.g., if mounting to `content/plugins/new-plugin`, verify `content/plugins/_index.md` exists)
6. Update GitHub Actions path filter with the new path
7. Verify by running `hugo list all` and checking the new section appears

### Output

- Updated `hugo.toml` with new mount entry
- `_index.md` in the source directory (if created, includes `bookCollapseSection: true`)
- Parent section `_index.md` (if created)
- Updated deployment workflow path filter

---

## /hugo-fix-indexes

| Field | Value |
|-------|-------|
| Purpose | Scan content directories and create missing `_index.md` files |

### Arguments

| Name | Required | Description |
|------|----------|-------------|
| path | No | Optional path to scan (defaults to all mounted content) |

### Syntax

```
/hugo-fix-indexes [path]
```

### Examples

```
/hugo-fix-indexes            # Scan all content
/hugo-fix-indexes docs/      # Scan specific directory
```

### Process

1. Read `hugo.toml` to identify all module mounts targeting `content/`
2. Scan each mounted directory recursively for subdirectories that contain at least one `.md` file (excluding `_index.md`) and do NOT have an `_index.md` file
3. Create `_index.md` for each missing index:
   - Derive title from directory name using kebab-case to Title Case conversion
   - Include `bookCollapseSection: true` for collapsible sidebar navigation
   - Set weight based on alphabetical order within parent (starting at 10, incrementing by 10)
4. Report the list of created files

### Title derivation

| Directory name | Generated title |
|---------------|-----------------|
| `getting-started` | Getting Started |
| `api-reference` | API Reference |
| `howto` | Howto |

### Output

- `_index.md` files for all directories that were missing them
- Summary report of files created
