---
name: hugo-fix-indexes
description: Scan content directories and create missing _index.md files
arguments:
  - name: path
    description: "Optional path to scan (defaults to all mounted content)"
    required: false
---

# Fix Missing Index Files

Scan content directories and create `_index.md` files where missing to ensure proper section navigation and collapsible sidebars.

## Process

1. **Read hugo.toml** to identify all module mounts targeting `content/`

2. **Scan each mounted directory recursively** for subdirectories that:
   - Contain at least one `.md` file (excluding `_index.md`)
   - Do NOT have an `_index.md` file

3. **Create `_index.md`** for each missing index:
   - Derive title from directory name using kebab-case to Title Case conversion
   - Include `bookCollapseSection: true` for collapsible sidebar navigation
   - Set weight based on alphabetical order within parent (starting at 10, incrementing by 10)

4. **Report** the list of created files

## Title Derivation

Directory names are converted to titles:
- `getting-started` → "Getting Started"
- `api-reference` → "API Reference"
- `howto` → "Howto" (single words preserved)

## Usage

```
/hugo-fix-indexes            # Scan all content
/hugo-fix-indexes docs/      # Scan specific directory
```

## Examples

```
/hugo-fix-indexes
# Scanning hugo.toml for content mounts...
# Found mounts: docs/, plugins/auth/docs/, plugins/core/docs/
#
# Missing index files:
#   docs/advanced/        -> created _index.md (title: "Advanced")
#   docs/advanced/config/ -> created _index.md (title: "Config")
#   plugins/auth/docs/api/-> created _index.md (title: "Api")
#
# Created 3 index files.
```
