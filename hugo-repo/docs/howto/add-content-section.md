---
title: "Add a Content Section"
description: "How to mount a new directory as a content section in an existing Hugo site using the /hugo-add-section command"
weight: 2
---

# Add a Content Section

This guide covers adding a new documentation directory to an existing Hugo site as a content section. Use the `/hugo-add-section` command to automate the process.

## When to use this

Use this when you have an existing Hugo site with module mounts and need to add a new directory. Common scenarios include:

- A new plugin or service was added to the repository and it has a `docs/` folder
- An existing directory's documentation should be included in the site
- You want to reorganize the content hierarchy by mounting a directory to a different section path

## Steps

### 1. Ensure the source directory exists

The directory you want to mount must exist and contain at least one markdown file. If the directory is new, create it with an `_index.md` file:

```markdown
---
title: "New Component"
description: "Documentation for the new component"
weight: 10
---

Overview content for the new component.
```

### 2. Run the `/hugo-add-section` command

Provide the source directory path and optionally the target section path:

```
/hugo-add-section new-component/docs components/new-component
```

If you omit the section argument, the command infers it from the directory path. For example, `new-component/docs` maps to a section based on the directory structure.

```
/hugo-add-section new-component/docs
```

### 3. Review the changes

The command makes the following changes:

- **hugo.toml**: Adds a new `[[module.mounts]]` entry:

  ```toml
  [[module.mounts]]
  source = 'new-component/docs'
  target = 'content/components/new-component'
  ```

- **Section index**: Creates `_index.md` in the source directory if one does not exist

- **Parent section**: Creates `_index.md` for the parent section if it does not exist (for example, `content/components/_index.md`)

- **GitHub Actions workflow**: Adds the new path to the deployment workflow's `paths` filter:

  ```yaml
  paths:
    - 'new-component/docs/**'
  ```

### 4. Verify the new section

Run Hugo's list command to confirm the new content appears:

```bash
hugo list all
```

Then start the development server to visually confirm:

```bash
hugo server --buildDrafts
```

Navigate to the new section's URL (for example, `http://localhost:1313/components/new-component/`) and verify it renders correctly.

## Examples

Mount a service's documentation:

```
/hugo-add-section services/auth/docs services/auth
```

Mount a tool's documentation with auto-detected section:

```
/hugo-add-section tools/cli/docs
```

Mount documentation to a custom location in the content tree:

```
/hugo-add-section lib/internal-sdk/docs reference/internal-sdk
```

## Troubleshooting

- **Section does not appear in navigation**: Check that the source directory has an `_index.md` file and that the parent section also has one.
- **Content renders but URL is wrong**: Verify the target path in the mount entry matches your expected URL structure.
- **Deployment does not trigger on changes**: Confirm the path filter in the GitHub Actions workflow includes the new directory.

## Related

- {{< ref "howto/set-up-hugo-site" >}} -- Initial Hugo site setup with module mounts
- {{< ref "reference/commands" >}} -- Full `/hugo-add-section` command specification
- {{< ref "explanation/component-interaction" >}} -- How commands use templates
