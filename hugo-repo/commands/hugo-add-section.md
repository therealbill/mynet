---
name: hugo-add-section
description: Add a new content section to the Hugo site from a directory
arguments:
  - name: path
    description: "Path to the directory to mount (e.g., new-plugin/docs)"
    required: true
  - name: section
    description: "Target section path (e.g., plugins/new-plugin)"
    required: false
---

# Add Content Section

Mount a new directory into the Hugo site as a content section.

## Process

1. **Validate the source directory** exists and contains markdown files.

2. **Determine the target section:**
   - If `section` argument provided, use it as the mount target under `content/`
   - If not provided, infer from the directory path (e.g., `new-plugin/docs` -> `plugins/new-plugin`)

3. **Add mount entry** to `hugo.toml`:
   ```toml
   [[module.mounts]]
   source = 'new-plugin/docs'
   target = 'content/plugins/new-plugin'
   ```

4. **Create section index** if the source directory lacks `_index.md`:
   ```markdown
   ---
   title: "New Plugin"
   description: "New Plugin documentation"
   weight: 99
   bookCollapseSection: true
   ---
   ```

5. **Ensure parent section exists** — if mounting to `content/plugins/new-plugin`, verify `content/plugins/_index.md` exists.

6. **Update GitHub Actions path filter** — add the new path to the deployment workflow's `paths` list.

7. **Verify** by running `hugo list all` and checking the new section appears.

## Usage

```
/hugo-add-section <path> [section]
```

## Examples

```
/hugo-add-section new-plugin/docs                    # Auto-detect section
/hugo-add-section new-plugin/docs plugins/new-plugin # Explicit section
/hugo-add-section services/auth/docs services/auth   # Service docs
```
