---
name: hugo-init
description: Scaffold a Hugo site for the current repository
arguments:
  - name: theme
    description: "Optional theme to install (e.g., ananke, book, docsy)"
    required: false
---

# Initialize Hugo Site

Analyze the repository structure and scaffold a Hugo site with appropriate module mounts.

## Process

1. **Scan the repository** for directories containing markdown documentation:
   - Look for `docs/` directories at any depth
   - Look for `*.md` files in common documentation locations
   - Identify the repository type (monorepo, plugin marketplace, tool collection, single project)

2. **Report findings** to the user:
   - List all discovered documentation directories
   - Recommend which should be mounted as content sections
   - Propose the site hierarchy

3. **Initialize Hugo:**
   ```bash
   hugo new site . --force  # --force since repo already exists
   hugo mod init github.com/<detected-from-git-remote>
   ```

4. **Generate `hugo.toml`** with:
   - Standard module mounts (content, static, layouts, data, assets)
   - Custom mounts for each discovered docs directory
   - Sensible defaults (title from repo name, baseURL placeholder)

5. **Create section index pages:**
   - `content/_index.md` (home page)
   - Parent section `_index.md` for each mount group (e.g., `content/plugins/_index.md`)
   - **Scan each mounted directory recursively** — create `_index.md` for any subdirectory containing `.md` files but lacking an index
   - All generated indexes include `bookCollapseSection: true` for collapsible sidebar navigation (Book theme)
   - Derive titles from directory names using kebab-case to Title Case conversion

6. **Create render hook for portable links:**
   - Create `layouts/_default/_markup/render-link.html` from template
   - This handles `.md` extensions in links and resolves paths correctly
   - Enables standard markdown links without Hugo-specific shortcodes

7. **Install theme** (if specified or if user accepts recommendation):
   ```bash
   hugo mod get github.com/theNewDynamic/gohugo-theme-ananke
   ```

8. **Update `.gitignore`:**
   ```
   public/
   resources/
   .hugo_build.lock
   ```

9. **Verify** by running `hugo server --buildDrafts` and reporting the result.

## Usage

```
/hugo-init [theme]
```

## Examples

```
/hugo-init              # Analyze repo and scaffold with default theme
/hugo-init book         # Scaffold with Hugo Book theme
/hugo-init docsy        # Scaffold with Docsy theme
```
