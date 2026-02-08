---
title: "Getting Started with hugo-repo"
description: "End-to-end walkthrough from scaffolding a Hugo site to deploying it, using the plugin's commands, skills, and agents together"
weight: 1
---

# Getting Started with hugo-repo

This tutorial walks you through the full lifecycle of a Hugo documentation site using the hugo-repo plugin. You will scaffold a site, add a content section, customize the theme, preview locally, and deploy to GitHub Pages. By the end, you will have a working Hugo site built from your repository's documentation directories.

## What you will learn

- How to scaffold a Hugo site with the `/hugo-init` command
- How to add content sections from subdirectories with `/hugo-add-section`
- How to customize a Hugo theme using the hugo-themes skill
- How to preview your site locally with `/hugo-serve`
- How to deploy to GitHub Pages with `/hugo-deploy`

## Prerequisites

- A Git repository with at least one `docs/` directory containing markdown files
- Hugo installed locally (version 0.142.0 or later recommended)
- Go installed (required for Hugo modules)
- A GitHub repository with GitHub Pages available

## Step 1: Scaffold the Hugo site

Run the `/hugo-init` command in your repository. This command analyzes the repository structure, identifies documentation directories, and generates the Hugo configuration.

```
/hugo-init
```

The command performs the following actions:

- Scans the repository for directories containing markdown documentation
- Reports findings and proposes a site hierarchy
- Runs `hugo new site . --force` and `hugo mod init`
- Generates `hugo.toml` with module mounts for each discovered docs directory
- Creates `_index.md` section index pages for every section in the content tree
- Updates `.gitignore` to exclude `public/`, `resources/`, and `.hugo_build.lock`

If you want to install a specific theme during scaffolding, pass the theme name:

```
/hugo-init book
```

### Checkpoint

After this step, you should have:

- A `hugo.toml` file in the repository root with module mounts
- A `go.mod` and `go.sum` file
- An `_index.md` file in each content section
- A `.gitignore` with Hugo-specific entries

Verify by checking the generated configuration:

```bash
cat hugo.toml
```

You should see the standard mounts (content, static, layouts, data, assets, archetypes) plus custom mounts for each docs directory found in the repository. The hugo-site-architect agent performs the repository analysis behind this command, using its Read, Glob, and Grep tools to discover documentation directories.

## Step 2: Add a content section

Suppose you add a new component to your repository and it has its own `docs/` directory. Use the `/hugo-add-section` command to mount it into the Hugo site:

```
/hugo-add-section new-component/docs components/new-component
```

This command:

- Validates that `new-component/docs` exists and contains markdown files
- Adds a `[[module.mounts]]` entry to `hugo.toml` mapping `new-component/docs` to `content/components/new-component`
- Creates an `_index.md` in the source directory if one does not exist
- Ensures the parent section (`content/components/`) has its own `_index.md`
- Updates the GitHub Actions path filter to include the new directory

If you omit the second argument, the command infers the section path from the directory path:

```
/hugo-add-section new-component/docs
```

### Checkpoint

After this step, verify the new mount appears in `hugo.toml`:

```bash
grep -A2 'new-component' hugo.toml
```

You should see:

```toml
[[module.mounts]]
source = 'new-component/docs'
target = 'content/components/new-component'
```

Also confirm that `new-component/docs/_index.md` exists with proper front matter containing title, description, and weight fields.

## Step 3: Customize the theme

The hugo-themes skill provides knowledge about Hugo theme customization. Ask Claude to help you customize your theme:

```
I want to override the header partial in my Hugo theme to add a custom logo
```

Claude will use the hugo-themes skill to guide you through the template lookup order. The key principle is that files in your project's `layouts/` directory override the same paths in the theme. For example, to override the header:

1. Identify the theme's header partial path (typically `layouts/partials/header.html`)
2. Create the same file in your project's `layouts/partials/` directory
3. Modify the copy to include your changes

To add custom CSS, create an `assets/scss/custom.scss` file and include it via Hugo Pipes in your head partial override:

```html
{{ $style := resources.Get "scss/custom.scss" | toCSS | minify | fingerprint }}
<link rel="stylesheet" href="{{ $style.RelPermalink }}">
```

### Checkpoint

After this step, you should have:

- One or more files in `layouts/` that override theme templates
- Optionally, custom styles in `assets/` or `static/css/`
- No modifications to the theme's own files (all changes are overrides)

## Step 4: Preview locally

Start the Hugo development server with the `/hugo-serve` command:

```
/hugo-serve
```

This command:

- Verifies Hugo is installed
- Resolves Hugo modules if `go.mod` exists
- Starts `hugo server --buildDrafts --navigateToChanged`
- Reports the local URL (typically `http://localhost:1313/`)

Open the reported URL in your browser. The development server watches for file changes and live-reloads the browser. If you edit a content file, the browser automatically navigates to the changed page.

If you see stale content after making changes, try adding the `--disableFastRender` flag:

```
/hugo-serve --disableFastRender
```

### Checkpoint

After this step, you should see:

- Your site rendered at `http://localhost:1313/`
- All mounted content sections appearing in the navigation
- Your theme customizations applied
- Draft content visible (since `--buildDrafts` is enabled)

If content is missing, ask Claude for help. The hugo-build-doctor agent specializes in diagnosing missing content, which is commonly caused by a missing `_index.md`, an incorrect mount configuration, or `draft: true` in front matter.

## Step 5: Deploy to GitHub Pages

Set up deployment with the `/hugo-deploy` command:

```
/hugo-deploy pages
```

This command:

- Checks for an existing deployment workflow in `.github/workflows/`
- Generates `.github/workflows/hugo-deploy.yml` using the `github-pages.yml.tmpl` template
- Configures content-aware path filtering so that only documentation changes trigger builds
- Adds path entries for all mounted docs directories
- Explains the remaining manual steps

After the workflow is created, complete the GitHub Pages setup:

1. Go to your repository's Settings page on GitHub
2. Navigate to Pages in the left sidebar
3. Set Source to "GitHub Actions"
4. Push a commit that changes content to trigger the first deployment

The generated workflow uses Hugo Extended edition (for SCSS support), pins the Hugo version, caches Hugo modules, builds with `--minify`, and deploys using the official GitHub Pages actions.

### Checkpoint

After this step, you should have:

- `.github/workflows/hugo-deploy.yml` in your repository
- GitHub Pages source set to "GitHub Actions" in repository settings
- A successful deployment after pushing a content change
- Your site accessible at `https://USERNAME.github.io/REPO/`

## Next steps

Now that your site is deployed, explore these topics:

- {{< ref "howto/add-content-section" >}} to mount additional documentation directories
- {{< ref "howto/deploy-to-s3" >}} if you need AWS S3 hosting instead of GitHub Pages
- {{< ref "howto/customize-theme" >}} for deeper theme customization techniques
- {{< ref "reference/skills" >}} for the full list of skills available in this plugin
- {{< ref "explanation/architecture" >}} to understand how all 17 components work together

## Troubleshooting

If you encounter problems at any step, the hugo-build-doctor agent can help diagnose the issue. Describe the error or symptom to Claude, and the agent will systematically investigate using its Read, Bash, Glob, and Grep tools.

Common issues during initial setup:

- **Hugo not found**: Install Hugo from [gohugo.io/installation](https://gohugo.io/installation/)
- **Module errors**: Run `hugo mod init github.com/your/repo` if `go.mod` is missing
- **Missing sections**: Ensure every mounted directory and its parent sections have `_index.md` files
- **Theme not found**: Run `hugo mod get` to download module dependencies
- **Build fails in CI but works locally**: Check that the Hugo version in the workflow matches your local version
