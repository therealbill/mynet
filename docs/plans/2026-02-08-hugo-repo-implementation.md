# Hugo Repo Plugin Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create the `hugo-repo` plugin with 7 skills, 2 agents, 4 commands, and 4 templates for Hugo site management in GitHub repositories.

**Architecture:** All components are markdown files following the established plugin conventions in this marketplace. Skills use YAML frontmatter with third-person trigger descriptions. Agents use the example-block description format. Commands use the argument-declaration format. Templates are annotated starter files.

**Tech Stack:** Markdown with YAML frontmatter. Hugo knowledge sourced from Hugo documentation. No code dependencies.

**Design document:** `docs/plans/2026-02-08-hugo-repo-design.md`

**Reference files for conventions:**

- Agent format: `research/agents/academic-researcher.md`
- Command format: `timelord/commands/tl-deploy.md`
- Skill format: `gnu-make/skills/makefile-fundamentals/SKILL.md`
- Plugin.json format: `timelord/.claude-plugin/plugin.json`

---

### Task 1: Plugin scaffolding and marketplace registration

**Files:**

- Create: `hugo-repo/.claude-plugin/plugin.json`
- Create: `hugo-repo/.gitignore`
- Modify: `.claude-plugin/marketplace.json`

**Step 1:** Create plugin.json.

```json
{
  "name": "hugo-repo",
  "version": "1.0.0",
  "description": "Hugo site management for GitHub repositories — multi-directory content aggregation, theme customization, data-driven content, and deployment to GitHub Pages or AWS S3",
  "author": {
    "name": "Bill",
    "url": "https://github.com/therealbill"
  },
  "repository": "https://github.com/therealbill/mynet",
  "license": "SEE LICENSE IN ../LICENSE",
  "keywords": [
    "hugo",
    "static-site",
    "github-pages",
    "s3",
    "documentation",
    "module-mounts",
    "shortcodes",
    "themes"
  ]
}
```

**Step 2:** Create .gitignore.

```
.DS_Store
**/.DS_Store
```

**Step 3:** Add marketplace entry to `.claude-plugin/marketplace.json`. Add to the `plugins` array:

```json
{
  "name": "hugo-repo",
  "source": "./hugo-repo",
  "description": "Hugo site management for GitHub repositories — multi-directory content aggregation, theme customization, data-driven content, and deployment to GitHub Pages or AWS S3",
  "version": "1.0.0",
  "keywords": ["hugo", "static-site", "github-pages", "s3", "documentation", "module-mounts"]
}
```

**Step 4:** Commit.

```bash
git add hugo-repo/.claude-plugin/plugin.json hugo-repo/.gitignore .claude-plugin/marketplace.json
git commit -m "add hugo-repo plugin scaffolding and marketplace registration"
```

---

### Task 2: hugo-fundamentals skill

**Files:**

- Create: `hugo-repo/skills/hugo-fundamentals/SKILL.md`

**Step 1:** Create the skill file with the following content.

````markdown
---
name: Hugo Fundamentals
description: >-
  This skill should be used when the user asks to create a Hugo site, set up
  Hugo, start a new Hugo project, understand Hugo directory structure, configure
  hugo.toml or hugo.yaml, work with Hugo front matter, create content types or
  archetypes, run hugo server for local development, or understand Hugo content
  organization including section pages and page bundles. Also applies when the
  user is new to Hugo and needs guidance on basic concepts, URL structure, or
  the relationship between content directories and rendered pages.
version: 1.0.0
---

## Overview

Hugo is a fast static site generator written in Go. A Hugo project has a standard directory structure, content written in markdown with front matter metadata, and configuration in `hugo.toml` (or `hugo.yaml`/`hugo.json`). Understanding this structure is essential before working with advanced features like module mounts or custom themes.

## Site Initialization

Create a new Hugo site:

```bash
hugo new site mysite
cd mysite
hugo mod init github.com/username/mysite
```

The `hugo mod init` step initializes Hugo modules, required for module mounts and theme-as-module installation. Always initialize modules even if not immediately using mounts.

## Directory Structure

```
mysite/
├── archetypes/       # Content templates for `hugo new`
├── assets/           # Files processed by Hugo Pipes (SCSS, JS)
├── content/          # Site content (markdown files)
├── data/             # Data files (YAML, JSON, TOML)
├── i18n/             # Internationalization strings
├── layouts/          # Templates (override theme layouts here)
├── static/           # Static files copied as-is (images, CSS, JS)
├── themes/           # Installed themes
└── hugo.toml         # Site configuration
```

Key points:

- `content/` is the only directory that holds page content
- `layouts/` in the project root overrides theme layouts (template lookup order)
- `static/` files are copied to the output root without processing
- `assets/` files are processed by Hugo Pipes (minification, fingerprinting, SCSS compilation)

## Configuration

Hugo supports `hugo.toml`, `hugo.yaml`, and `hugo.json`. TOML is the default and most common in Hugo documentation.

```toml
baseURL = 'https://example.com/'
languageCode = 'en-us'
title = 'My Site'
theme = 'my-theme'

[params]
  description = 'A Hugo site'
  author = 'Author Name'

[menu]
  [[menu.main]]
    name = 'Home'
    url = '/'
    weight = 1
  [[menu.main]]
    name = 'About'
    url = '/about/'
    weight = 2
```

Configuration can also be split into a `config/` directory with per-environment files (`config/_default/hugo.toml`, `config/production/hugo.toml`).

## Content and Front Matter

Content files are markdown with YAML or TOML front matter:

```markdown
---
title: "Getting Started"
date: 2024-01-15
draft: false
weight: 10
description: "How to get started with the project"
tags: ["setup", "quickstart"]
---

Content goes here in standard markdown.
```

Common front matter fields:

- `title` — Page title (required)
- `date` — Publication date
- `draft` — If true, excluded from production builds (visible with `--buildDrafts`)
- `weight` — Sort order within a section (lower = first)
- `description` — Used in meta tags and list pages
- `tags`, `categories` — Taxonomies for classification
- `layout` — Override the default template for this page
- `aliases` — Redirect URLs to this page

## Section Pages vs Leaf Bundles

Hugo has two types of content organization:

**Section pages** use `_index.md` — they represent a list page for a directory:

```
content/
├── _index.md              # Home page
├── blog/
│   ├── _index.md          # Blog list page
│   ├── first-post.md      # Blog post
│   └── second-post.md     # Blog post
└── docs/
    ├── _index.md          # Docs list page
    ├── getting-started.md
    └── configuration.md
```

**Leaf bundles** use `index.md` — they represent a single page with co-located resources:

```
content/
└── blog/
    └── my-post/
        ├── index.md       # The page content
        ├── hero.jpg       # Page resource (accessible via .Resources)
        └── data.csv       # Page resource
```

Key distinction: `_index.md` = section (has children), `index.md` = leaf (no children, has resources).

## URL Structure

Hugo maps content organization to URLs:

| File Path | URL |
|-----------|-----|
| `content/_index.md` | `/` |
| `content/about.md` | `/about/` |
| `content/blog/_index.md` | `/blog/` |
| `content/blog/first-post.md` | `/blog/first-post/` |
| `content/docs/guide/setup.md` | `/docs/guide/setup/` |

Override with `url` front matter or `[permalinks]` configuration.

## Archetypes

Archetypes are templates for `hugo new`:

```bash
hugo new content blog/my-post.md
```

Hugo looks for archetypes in order: `archetypes/blog.md`, then `archetypes/default.md`, then theme archetypes.

Example archetype (`archetypes/blog.md`):

```markdown
---
title: "{{ replace .File.ContentBaseName "-" " " | title }}"
date: {{ .Date }}
draft: true
tags: []
---
```

## Development Server

```bash
# Start dev server with drafts visible
hugo server --buildDrafts

# With live reload navigating to changed file
hugo server --buildDrafts --navigateToChanged

# Disable fast render if seeing stale content
hugo server --buildDrafts --disableFastRender
```

The dev server watches for changes and live-reloads the browser. Default URL: `http://localhost:1313/`.

## Building for Production

```bash
# Build the site
hugo

# Build with specific environment
hugo --environment production

# Build with minification
hugo --minify
```

Output goes to `public/` by default. Add `public/` and `resources/` to `.gitignore`.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Spaces in content filenames | Use hyphens: `my-page.md` not `my page.md` |
| Missing `_index.md` in sections | Every directory that should appear as a section needs `_index.md` |
| Forgetting `hugo mod init` | Required for module mounts and theme-as-module |
| Using `index.md` when `_index.md` needed | `index.md` = leaf bundle (no children), `_index.md` = section (has children) |
| Draft content not showing | Use `hugo server --buildDrafts` or set `draft: false` |
````

**Step 2:** Commit.

```bash
git add hugo-repo/skills/hugo-fundamentals/SKILL.md
git commit -m "add hugo-fundamentals skill"
```

---

### Task 3: hugo-module-mounts skill

**Files:**

- Create: `hugo-repo/skills/hugo-module-mounts/SKILL.md`

**Step 1:** Create the skill file with the following content.

````markdown
---
name: Hugo Module Mounts
description: >-
  This skill should be used when the user asks about Hugo module mounts,
  mounting multiple directories into Hugo's content tree, aggregating docs from
  scattered directories, using hugo mod init, configuring module.mounts in
  hugo.toml, building a unified documentation site from a monorepo, mapping
  subdirectory docs/ folders into a single Hugo site, or understanding Hugo's
  union file system. Also applies when the user has a repository with multiple
  docs/ directories that need to be presented as a single site. This is the
  core skill for multi-directory content aggregation.
version: 1.0.0
---

## Overview

Hugo's module mounts map arbitrary directories into Hugo's content tree using the `[[module.mounts]]` configuration. This creates a union file system — multiple source directories appear as a single content tree without symlinks or build-time copying. This is Hugo's official mechanism for multi-directory content aggregation.

Use module mounts when a repository has documentation scattered across multiple subdirectories (per-plugin, per-service, per-tool) and needs a unified Hugo site.

## Prerequisites

Hugo modules require Go and module initialization:

```bash
# Initialize Hugo modules (required once per project)
hugo mod init github.com/username/repo
```

This creates `go.mod` and `go.sum` files. Commit both to the repository.

## Basic Mount Configuration

In `hugo.toml`, mounts map a `source` directory to a `target` location in Hugo's virtual filesystem:

```toml
[[module.mounts]]
source = 'content'
target = 'content'

[[module.mounts]]
source = 'static'
target = 'static'

[[module.mounts]]
source = 'layouts'
target = 'layouts'

[[module.mounts]]
source = 'data'
target = 'data'

[[module.mounts]]
source = 'assets'
target = 'assets'

[[module.mounts]]
source = 'i18n'
target = 'i18n'

[[module.mounts]]
source = 'archetypes'
target = 'archetypes'
```

When you define any mount, Hugo's default mounts are replaced. Always include the standard mounts above when adding custom mounts, or the default directories will stop working.

## Multi-Directory Content Aggregation

The core pattern: mount each subdirectory's `docs/` folder into the content tree.

### Example: Plugin Marketplace

Repository structure:

```
repo/
├── hugo.toml
├── content/           # Root site content (home page, about, etc.)
├── plugin-a/
│   ├── docs/          # Plugin A documentation
│   │   ├── _index.md
│   │   └── usage.md
│   └── src/
├── plugin-b/
│   ├── docs/
│   │   ├── _index.md
│   │   └── api.md
│   └── src/
└── plugin-c/
    ├── docs/
    └── src/
```

Mount configuration:

```toml
# Standard mounts (required when defining any custom mounts)
[[module.mounts]]
source = 'content'
target = 'content'

[[module.mounts]]
source = 'static'
target = 'static'

[[module.mounts]]
source = 'layouts'
target = 'layouts'

[[module.mounts]]
source = 'data'
target = 'data'

[[module.mounts]]
source = 'assets'
target = 'assets'

# Plugin documentation mounts
[[module.mounts]]
source = 'plugin-a/docs'
target = 'content/plugins/plugin-a'

[[module.mounts]]
source = 'plugin-b/docs'
target = 'content/plugins/plugin-b'

[[module.mounts]]
source = 'plugin-c/docs'
target = 'content/plugins/plugin-c'
```

Result: Hugo sees a unified content tree:

```
content/
├── _index.md                    # From content/
├── plugins/
│   ├── plugin-a/
│   │   ├── _index.md           # From plugin-a/docs/
│   │   └── usage.md
│   ├── plugin-b/
│   │   ├── _index.md           # From plugin-b/docs/
│   │   └── api.md
│   └── plugin-c/
│       └── ...
```

### Example: Monorepo Services

```toml
# Service documentation mounts
[[module.mounts]]
source = 'services/api-gateway/docs'
target = 'content/services/api-gateway'

[[module.mounts]]
source = 'services/auth-service/docs'
target = 'content/services/auth-service'

[[module.mounts]]
source = 'services/worker/docs'
target = 'content/services/worker'
```

### Example: CLI Tool Collection

```toml
# Tool documentation mounts
[[module.mounts]]
source = 'tools/migrate/docs'
target = 'content/tools/migrate'

[[module.mounts]]
source = 'tools/validate/docs'
target = 'content/tools/validate'

[[module.mounts]]
source = 'tools/generate/docs'
target = 'content/tools/generate'
```

## Section Index Pages

Each mounted directory needs an `_index.md` to appear as a proper section. The parent directory also needs one.

Create `content/plugins/_index.md` (or `content/services/_index.md`):

```markdown
---
title: "Plugins"
description: "Documentation for all plugins"
weight: 10
---

Browse documentation for individual plugins below.
```

Each mounted directory needs its own `_index.md` (inside the source `docs/` folder):

```markdown
---
title: "Plugin A"
description: "Plugin A documentation"
weight: 1
---

Plugin A overview content.
```

## Mount Precedence

When files conflict (same path from multiple mounts), the first mount wins:

```toml
# This mount takes precedence for any conflicting paths
[[module.mounts]]
source = 'content'
target = 'content'

# This mount's files are used only when they don't conflict with above
[[module.mounts]]
source = 'shared-content'
target = 'content'
```

Order mounts from highest to lowest priority.

## Mounting Non-Content Directories

Mounts work for any Hugo component directory, not just content:

```toml
# Mount shared layouts
[[module.mounts]]
source = 'layouts'
target = 'layouts'

[[module.mounts]]
source = 'shared/layouts'
target = 'layouts'

# Mount data files from a subdirectory
[[module.mounts]]
source = 'data'
target = 'data'

[[module.mounts]]
source = 'plugin-a/data'
target = 'data/plugins/plugin-a'
```

Valid mount targets: `archetypes`, `assets`, `content`, `data`, `i18n`, `layouts`, `static`.

## Adding a New Mount

When adding documentation for a new subdirectory:

1. Create the `docs/` directory with `_index.md`
2. Add the mount entry to `hugo.toml`
3. Ensure the parent section has `_index.md`
4. Run `hugo server` to verify

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Defining custom mounts without standard mounts | Always include mounts for `content`, `static`, `layouts`, `data`, `assets` when adding any custom mount |
| Missing `_index.md` in mounted directories | Every mounted directory needs `_index.md` to be a proper Hugo section |
| Missing parent section `_index.md` | If mounting to `content/plugins/foo`, ensure `content/plugins/_index.md` exists |
| Forgetting `hugo mod init` | Module mounts require module initialization |
| Absolute paths in source | Use paths relative to the repository root |
| Not committing `go.mod` and `go.sum` | These are required for reproducible builds |

## Debugging Mounts

```bash
# Verify Hugo sees all mounted content
hugo list all

# Check which files Hugo resolves
hugo config mounts

# Verbose build output to see mount resolution
hugo --verbose
```
````

**Step 2:** Commit.

```bash
git add hugo-repo/skills/hugo-module-mounts/SKILL.md
git commit -m "add hugo-module-mounts skill"
```

---

### Task 4: hugo-themes skill

**Files:**

- Create: `hugo-repo/skills/hugo-themes/SKILL.md`

**Step 1:** Create the skill file with the following content.

````markdown
---
name: Hugo Themes
description: >-
  This skill should be used when the user asks about Hugo themes, installing a
  theme, customizing a theme, overriding theme layouts or partials, Hugo
  template lookup order, creating a Hugo theme from scratch, Hugo asset
  pipeline (SCSS, PostCSS), Hugo baseof.html, or theme configuration
  parameters. Also applies when evaluating which theme to use, switching
  themes, or troubleshooting theme-related rendering issues.
version: 1.0.0
---

## Overview

Hugo themes control the visual presentation and layout structure of a site. Most users install a community theme and customize it by overriding specific templates and styles. Full theme creation from scratch is less common but documented here at reference level.

## Installing a Theme as a Hugo Module

The recommended approach — no git submodules, version-managed:

```bash
# Initialize modules if not already done
hugo mod init github.com/username/mysite

# Add theme as module dependency
hugo mod get github.com/theNewDynamic/gohugo-theme-ananke
```

In `hugo.toml`:

```toml
[module]
  [[module.imports]]
    path = 'github.com/theNewDynamic/gohugo-theme-ananke'
```

Update the theme:

```bash
hugo mod get -u github.com/theNewDynamic/gohugo-theme-ananke
```

## Template Lookup Order

Hugo resolves templates in this order (first match wins):

1. `layouts/` in the project root (your overrides)
2. `layouts/` in the theme

This means: to override any theme template, create the same file path under your project's `layouts/` directory.

## Overriding Theme Templates

### Override a Partial

If the theme has `themes/mytheme/layouts/partials/header.html`, override it by creating:

```
layouts/partials/header.html
```

Your version replaces the theme's version completely. To extend rather than replace, copy the theme's partial and modify it.

### Override a Layout

Theme provides `themes/mytheme/layouts/_default/single.html`. Override with:

```
layouts/_default/single.html
```

### Override for Specific Sections

Override the list template only for the `blog` section:

```
layouts/blog/list.html
```

Hugo's template lookup searches section-specific templates before `_default/`.

## Customizing CSS

### Method 1: Custom CSS File

Add to `static/css/custom.css` and include in the head partial:

```html
<!-- layouts/partials/head.html -->
{{ partial "head/meta.html" . }}
{{ partial "head/css.html" . }}
<link rel="stylesheet" href="{{ "css/custom.css" | relURL }}">
```

### Method 2: Hugo Pipes (Recommended)

Process SCSS/CSS through Hugo's asset pipeline:

```html
<!-- layouts/partials/head.html -->
{{ $style := resources.Get "scss/main.scss" | toCSS | minify | fingerprint }}
<link rel="stylesheet" href="{{ $style.RelPermalink }}" integrity="{{ $style.Data.Integrity }}">
```

Place source files in `assets/scss/main.scss`.

### Method 3: Theme Parameters

Many themes support CSS customization via config:

```toml
[params]
  customCSS = ["css/custom.css"]
  colorScheme = "dark"
```

Check the theme's documentation for supported parameters.

## Theme Configuration Parameters

Themes define custom parameters in `[params]`:

```toml
[params]
  # Common theme parameters
  logo = "/images/logo.png"
  favicon = "/images/favicon.ico"
  description = "Site description"
  dateFormat = "January 2, 2006"

  # Social links (theme-specific)
  [params.social]
    github = "username"
    twitter = "username"

  # Navigation (theme-specific)
  [params.nav]
    showBreadcrumbs = true
    showTableOfContents = true
```

## Creating a Theme from Scratch

### Scaffold

```bash
hugo new theme mytheme
```

Creates:

```
themes/mytheme/
├── archetypes/
│   └── default.md
├── layouts/
│   ├── _default/
│   │   ├── baseof.html
│   │   ├── list.html
│   │   └── single.html
│   ├── partials/
│   │   ├── footer.html
│   │   ├── head.html
│   │   └── header.html
│   └── index.html
├── static/
│   ├── css/
│   └── js/
└── theme.toml
```

### baseof.html — The Master Template

All pages inherit from this:

```html
<!DOCTYPE html>
<html lang="{{ .Site.LanguageCode }}">
<head>
  {{ partial "head.html" . }}
</head>
<body>
  {{ partial "header.html" . }}
  <main>
    {{ block "main" . }}{{ end }}
  </main>
  {{ partial "footer.html" . }}
</body>
</html>
```

### single.html — Individual Pages

```html
{{ define "main" }}
<article>
  <h1>{{ .Title }}</h1>
  <time>{{ .Date.Format "January 2, 2006" }}</time>
  {{ .Content }}
</article>
{{ end }}
```

### list.html — Section List Pages

```html
{{ define "main" }}
<h1>{{ .Title }}</h1>
{{ .Content }}
<ul>
  {{ range .Pages }}
  <li>
    <a href="{{ .RelPermalink }}">{{ .Title }}</a>
    <p>{{ .Summary }}</p>
  </li>
  {{ end }}
</ul>
{{ end }}
```

For detailed theme creation patterns, see Hugo's template documentation. The hugo-content-authoring skill covers shortcodes and render hooks that complement theme layouts.
````

**Step 2:** Commit.

```bash
git add hugo-repo/skills/hugo-themes/SKILL.md
git commit -m "add hugo-themes skill"
```

---

### Task 5: hugo-content-authoring skill

**Files:**

- Create: `hugo-repo/skills/hugo-content-authoring/SKILL.md`

**Step 1:** Create the skill file with the following content.

````markdown
---
name: Hugo Content Authoring
description: >-
  This skill should be used when the user asks about Hugo shortcodes, creating
  custom shortcodes, Hugo render hooks, custom rendering for links or images
  or headings or code blocks, Hugo taxonomies (tags, categories, custom),
  page bundles and page resources, content formatting beyond basic markdown,
  creating callout boxes or tabs or admonitions in Hugo, or customizing how
  Hugo processes markdown elements. Also applies when the user wants to add
  reusable content components to their Hugo site.
version: 1.0.0
---

## Overview

Hugo provides content authoring features beyond standard markdown: shortcodes for reusable content components, render hooks to customize how markdown elements are processed, taxonomies for content classification, and page bundles for co-located assets. These features work independently of theme choice.

## Built-in Shortcodes

Hugo includes several shortcodes:

```markdown
<!-- Syntax-highlighted code with line numbers -->
{{</* highlight go "linenos=true,hl_lines=3" */>}}
func main() {
    fmt.Println("Hello")
    fmt.Println("Highlighted line")
}
{{</* /highlight */>}}

<!-- Figure with caption -->
{{</* figure src="/images/photo.jpg" title="Caption" alt="Description" */>}}

<!-- Cross-reference another page -->
[Link text]({{</* ref "other-page.md" */>}})

<!-- Relative reference -->
[Related]({{</* relref "sibling-page.md" */>}})

<!-- GitHub Gist -->
{{</* gist username gist-id */>}}
```

## Creating Custom Shortcodes

Shortcodes live in `layouts/shortcodes/`. The filename becomes the shortcode name.

### Simple Shortcode with Parameters

`layouts/shortcodes/callout.html`:

```html
<div class="callout callout-{{ .Get "type" | default "info" }}">
  <strong>{{ .Get "title" | default "" }}</strong>
  {{ .Inner | markdownify }}
</div>
```

Usage in content:

```markdown
{{</* callout type="warning" title="Important" */>}}
This is a warning message with **markdown** support.
{{</* /callout */>}}
```

### Positional Parameters

`layouts/shortcodes/badge.html`:

```html
<span class="badge badge-{{ .Get 0 }}">{{ .Get 1 }}</span>
```

Usage: `{{</* badge "success" "Stable" */>}}`

### Shortcode Accessing Page Variables

`layouts/shortcodes/last-modified.html`:

```html
<time datetime="{{ .Page.Lastmod.Format "2006-01-02" }}">
  Last updated: {{ .Page.Lastmod.Format "January 2, 2006" }}
</time>
```

### Common Shortcode Patterns

**Tabs component** (`layouts/shortcodes/tabs.html` and `layouts/shortcodes/tab.html`):

```html
<!-- layouts/shortcodes/tabs.html -->
<div class="tabs">
  <div class="tab-buttons">
    {{ range $i, $e := .Scratch.Get "tabs" }}
    <button class="tab-btn{{ if eq $i 0 }} active{{ end }}" data-tab="{{ $i }}">{{ .name }}</button>
    {{ end }}
  </div>
  <div class="tab-panels">
    {{ .Inner }}
  </div>
</div>
```

**Code block with filename** (`layouts/shortcodes/file.html`):

```html
<div class="code-file">
  <div class="code-filename">{{ .Get "name" }}</div>
  {{ .Inner | markdownify }}
</div>
```

Usage:

```markdown
{{</* file name="config.yaml" */>}}
```yaml
key: value
```
{{</* /file */>}}
```

## Render Hooks

Render hooks customize how Hugo processes standard markdown elements. Place them in `layouts/_markup/`.

### Custom Link Rendering

`layouts/_markup/render-link.html`:

```html
<a href="{{ .Destination | safeURL }}"
  {{- with .Title }} title="{{ . }}"{{ end }}
  {{- if strings.HasPrefix .Destination "http" }} target="_blank" rel="noopener"{{ end }}>
  {{- .Text | safeHTML -}}
</a>
```

This makes external links open in a new tab automatically.

### Custom Image Rendering

`layouts/_markup/render-image.html`:

```html
<figure>
  <img src="{{ .Destination | safeURL }}" alt="{{ .Text }}"
    {{- with .Title }} title="{{ . }}"{{ end }} loading="lazy">
  {{- with .Title }}
  <figcaption>{{ . }}</figcaption>
  {{- end }}
</figure>
```

### Custom Heading Rendering

`layouts/_markup/render-heading.html`:

```html
<h{{ .Level }} id="{{ .Anchor }}">
  {{ .Text | safeHTML }}
  <a class="heading-anchor" href="#{{ .Anchor }}">#</a>
</h{{ .Level }}>
```

### Custom Code Block Rendering

`layouts/_markup/render-codeblock.html`:

```html
{{ $lang := .Type }}
<div class="code-block" data-lang="{{ $lang }}">
  <div class="code-header">
    <span class="code-lang">{{ $lang }}</span>
    <button class="copy-btn">Copy</button>
  </div>
  {{ .Inner | highlight $lang "" }}
</div>
```

## Taxonomies

Hugo supports tags and categories by default. Add custom taxonomies in `hugo.toml`:

```toml
[taxonomies]
  tag = 'tags'
  category = 'categories'
  author = 'authors'
  technology = 'technologies'
```

Use in front matter:

```markdown
---
title: "My Post"
tags: ["go", "hugo"]
categories: ["tutorials"]
authors: ["bill"]
technologies: ["hugo", "github-actions"]
---
```

Hugo auto-generates list pages at `/tags/`, `/categories/`, `/authors/`, `/technologies/`.

## Page Bundles and Resources

Leaf bundles co-locate a page with its assets:

```
content/blog/my-post/
├── index.md          # Page content
├── hero.jpg          # Page resource
├── diagram.svg       # Page resource
└── data.csv          # Page resource
```

Access resources in templates:

```html
{{ with .Resources.GetMatch "hero.jpg" }}
  <img src="{{ .RelPermalink }}" alt="Hero">
{{ end }}

{{ range .Resources.Match "*.svg" }}
  <img src="{{ .RelPermalink }}" alt="{{ .Title }}">
{{ end }}
```

## Table of Contents

Hugo auto-generates a table of contents from headings:

```html
<!-- In a template -->
{{ .TableOfContents }}
```

Configure depth in `hugo.toml`:

```toml
[markup]
  [markup.tableOfContents]
    startLevel = 2
    endLevel = 4
```
````

**Step 2:** Commit.

```bash
git add hugo-repo/skills/hugo-content-authoring/SKILL.md
git commit -m "add hugo-content-authoring skill"
```

---

### Task 6: hugo-data-templates skill

**Files:**

- Create: `hugo-repo/skills/hugo-data-templates/SKILL.md`

**Step 1:** Create the skill file with the following content.

````markdown
---
name: Hugo Data Templates
description: >-
  This skill should be used when the user asks about Hugo's data directory,
  loading YAML or JSON or TOML data files in Hugo templates, using .Site.Data
  to access structured data, data-driven navigation or menus, rendering
  datasets or bibliographies or registries from data files, generating content
  from data, using range to iterate over data collections, Hugo
  resources.GetRemote for external data, or organizing data files for large
  projects. Also applies when the user wants to drive site content or navigation
  from structured data rather than hardcoding it in templates.
version: 1.0.0
---

## Overview

Hugo's `data/` directory holds structured data files (YAML, JSON, TOML) accessible in templates via `.Site.Data`. This enables data-driven patterns: navigation from data files, content generated from datasets, registries rendered as pages, and bibliographies organized as structured data.

## Basic Data Access

Place a data file in the `data/` directory:

```yaml
# data/team.yaml
- name: Alice
  role: Lead Developer
  github: alice
- name: Bob
  role: Designer
  github: bob
```

Access in templates:

```html
<ul>
{{ range .Site.Data.team }}
  <li>
    <strong>{{ .name }}</strong> — {{ .role }}
    <a href="https://github.com/{{ .github }}">GitHub</a>
  </li>
{{ end }}
</ul>
```

## Data File Formats

Hugo supports YAML, JSON, and TOML in the data directory:

```yaml
# data/config.yaml
site_name: My Project
features:
  - name: Fast
    icon: zap
  - name: Secure
    icon: shield
```

```json
// data/config.json
{
  "site_name": "My Project",
  "features": [
    {"name": "Fast", "icon": "zap"},
    {"name": "Secure", "icon": "shield"}
  ]
}
```

```toml
# data/config.toml
site_name = "My Project"

[[features]]
name = "Fast"
icon = "zap"

[[features]]
name = "Secure"
icon = "shield"
```

All three produce the same `.Site.Data.config` object.

## Nested Data Directories

Data files in subdirectories create nested access:

```
data/
├── plugins/
│   ├── plugin-a.yaml
│   └── plugin-b.yaml
├── navigation.yaml
└── team.yaml
```

Access nested data:

```html
{{ range $name, $data := .Site.Data.plugins }}
  <h3>{{ $data.title }}</h3>
  <p>{{ $data.description }}</p>
{{ end }}
```

## Data-Driven Patterns

### Navigation from Data

```yaml
# data/navigation.yaml
main:
  - title: Home
    url: /
    weight: 1
  - title: Documentation
    url: /docs/
    weight: 2
    children:
      - title: Getting Started
        url: /docs/getting-started/
      - title: Configuration
        url: /docs/configuration/
  - title: About
    url: /about/
    weight: 3
```

```html
<!-- layouts/partials/nav.html -->
<nav>
  {{ range sort .Site.Data.navigation.main "weight" }}
  <div class="nav-item">
    <a href="{{ .url }}">{{ .title }}</a>
    {{ with .children }}
    <ul class="subnav">
      {{ range . }}
      <li><a href="{{ .url }}">{{ .title }}</a></li>
      {{ end }}
    </ul>
    {{ end }}
  </div>
  {{ end }}
</nav>
```

### Plugin or Service Registry

```yaml
# data/registry.yaml
plugins:
  - name: hugo-repo
    description: Hugo site management for GitHub repositories
    version: 1.0.0
    keywords: [hugo, static-site, deployment]
    status: stable
  - name: gnu-make
    description: GNU Make best practices skills
    version: 1.0.0
    keywords: [make, build]
    status: stable
```

```html
<!-- layouts/partials/registry.html -->
<div class="registry">
  {{ range .Site.Data.registry.plugins }}
  <div class="registry-card">
    <h3>{{ .name }} <span class="version">v{{ .version }}</span></h3>
    <p>{{ .description }}</p>
    <div class="tags">
      {{ range .keywords }}<span class="tag">{{ . }}</span>{{ end }}
    </div>
    <span class="status status-{{ .status }}">{{ .status }}</span>
  </div>
  {{ end }}
</div>
```

### Research Bibliography

```yaml
# data/bibliography.yaml
sources:
  - id: smith2024
    authors: "Smith, J. and Lee, K."
    title: "Advances in Static Site Generation"
    year: 2024
    journal: "Web Engineering Review"
    doi: "10.1234/wer.2024.001"
    tags: [static-sites, performance]

  - id: johnson2023
    authors: "Johnson, M."
    title: "Content Management in Monorepos"
    year: 2023
    journal: "Software Architecture Journal"
    doi: "10.5678/saj.2023.042"
    tags: [monorepo, documentation]
```

```html
<!-- layouts/shortcodes/cite.html -->
{{ $id := .Get 0 }}
{{ range .Site.Data.bibliography.sources }}
  {{ if eq .id $id }}
    <a href="https://doi.org/{{ .doi }}" class="citation">({{ .authors }}, {{ .year }})</a>
  {{ end }}
{{ end }}
```

Usage in content: `{{</* cite "smith2024" */>}}`

### Feature Comparison Matrix

```yaml
# data/comparison.yaml
features:
  - name: Module Mounts
    hugo: true
    jekyll: false
    mkdocs: false
  - name: Built-in SCSS
    hugo: true
    jekyll: true
    mkdocs: false
  - name: Data Directory
    hugo: true
    jekyll: true
    mkdocs: false
  - name: Custom Shortcodes
    hugo: true
    jekyll: false
    mkdocs: true
```

```html
<table>
  <tr><th>Feature</th><th>Hugo</th><th>Jekyll</th><th>MkDocs</th></tr>
  {{ range .Site.Data.comparison.features }}
  <tr>
    <td>{{ .name }}</td>
    <td>{{ if .hugo }}✓{{ else }}✗{{ end }}</td>
    <td>{{ if .jekyll }}✓{{ else }}✗{{ end }}</td>
    <td>{{ if .mkdocs }}✓{{ else }}✗{{ end }}</td>
  </tr>
  {{ end }}
</table>
```

### Changelog from Data

```yaml
# data/changelog.yaml
releases:
  - version: 1.2.0
    date: 2024-03-15
    changes:
      - type: feature
        description: Added S3 deployment support
      - type: fix
        description: Fixed module mount ordering
  - version: 1.1.0
    date: 2024-02-01
    changes:
      - type: feature
        description: Added custom shortcodes
```

## Remote Data

Hugo can fetch data from URLs at build time:

```html
{{ $data := resources.GetRemote "https://api.github.com/repos/user/repo" }}
{{ with $data }}
  {{ $json := .Content | transform.Unmarshal }}
  <p>Stars: {{ $json.stargazers_count }}</p>
{{ end }}
```

Use sparingly — remote data adds build-time latency and external dependencies.

## Data File Organization

For large datasets, organize by domain:

```
data/
├── navigation/
│   ├── main.yaml
│   ├── footer.yaml
│   └── sidebar.yaml
├── plugins/
│   ├── code-quality.yaml
│   └── web-development.yaml
├── team.yaml
└── site.yaml
```

Access: `.Site.Data.navigation.main`, `.Site.Data.plugins.code_quality` (hyphens become underscores).
````

**Step 2:** Commit.

```bash
git add hugo-repo/skills/hugo-data-templates/SKILL.md
git commit -m "add hugo-data-templates skill"
```

---

### Task 7: hugo-github-actions skill

**Files:**

- Create: `hugo-repo/skills/hugo-github-actions/SKILL.md`

**Step 1:** Create the skill file with the following content.

````markdown
---
name: Hugo GitHub Actions
description: >-
  This skill should be used when the user asks about deploying Hugo to GitHub
  Pages, setting up a GitHub Actions workflow for Hugo, Hugo CI/CD pipeline,
  content-aware path filtering for Hugo builds, caching Hugo modules in CI,
  Hugo PR preview deployments, Hugo build workflow, or configuring GitHub Pages
  deployment for a Hugo site. Also applies when troubleshooting GitHub Actions
  workflow failures related to Hugo builds.
version: 1.0.0
---

## Overview

Deploy Hugo sites to GitHub Pages using GitHub Actions with content-aware path filtering. In repositories where Hugo documentation lives alongside code, path filtering ensures site rebuilds only trigger when content, configuration, or theme files change — not on code-only commits.

## Standard GitHub Pages Workflow

Create `.github/workflows/hugo-deploy.yml`:

```yaml
name: Deploy Hugo site to GitHub Pages

on:
  push:
    branches: [master, main]
    paths:
      # Hugo content and configuration
      - 'content/**'
      - 'layouts/**'
      - 'static/**'
      - 'assets/**'
      - 'data/**'
      - 'themes/**'
      - 'hugo.toml'
      - 'hugo.yaml'
      - 'hugo.json'
      - 'config/**'
      - 'go.mod'
      - 'go.sum'
      # Mounted docs directories (add one per mount)
      # - 'plugin-a/docs/**'
      # - 'plugin-b/docs/**'
      # The workflow file itself
      - '.github/workflows/hugo-deploy.yml'
  workflow_dispatch:  # Allow manual triggering

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

defaults:
  run:
    shell: bash

jobs:
  build:
    runs-on: ubuntu-latest
    env:
      HUGO_VERSION: '0.142.0'
      HUGO_ENVIRONMENT: production
    steps:
      - name: Install Hugo CLI
        run: |
          wget -O ${{ runner.temp }}/hugo.deb https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-amd64.deb
          sudo dpkg -i ${{ runner.temp }}/hugo.deb

      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0

      - name: Setup Pages
        id: pages
        uses: actions/configure-pages@v5

      - name: Install Node.js dependencies
        run: |
          [[ -f package-lock.json || -f npm-shrinkwrap.json ]] && npm ci || true

      - name: Cache Hugo modules
        uses: actions/cache@v4
        with:
          path: |
            ~/.cache/hugo_cache
            /tmp/hugo_cache
          key: ${{ runner.os }}-hugo-${{ hashFiles('go.sum') }}
          restore-keys: |
            ${{ runner.os }}-hugo-

      - name: Build with Hugo
        run: |
          hugo \
            --minify \
            --baseURL "${{ steps.pages.outputs.base_url }}/"

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## Content-Aware Path Filtering

The `paths` filter is critical for repos where Hugo content lives alongside code. Without it, every code commit triggers a site rebuild.

Add one path entry per mounted docs directory:

```yaml
on:
  push:
    paths:
      # Standard Hugo directories
      - 'content/**'
      - 'layouts/**'
      - 'static/**'
      - 'assets/**'
      - 'data/**'
      - 'themes/**'
      - 'hugo.toml'
      - 'go.mod'
      - 'go.sum'
      # Per-plugin/service docs (one per mount)
      - 'plugin-a/docs/**'
      - 'plugin-b/docs/**'
      - 'services/api/docs/**'
```

When adding a new mount, also add its path to the workflow filter.

## Hugo Extended Edition

The extended edition is required for SCSS/PostCSS processing. The workflow above uses `hugo_extended_` in the download URL. If the theme does not use SCSS, the standard edition works and is smaller:

```yaml
# Standard edition (no SCSS)
wget -O ${{ runner.temp }}/hugo.deb https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_linux-amd64.deb
```

## Hugo Version Pinning

Pin the Hugo version in the workflow (`HUGO_VERSION: '0.142.0'`). Update deliberately, not automatically, since Hugo occasionally has breaking changes in template handling.

Check for updates: `https://github.com/gohugoio/hugo/releases`

## PR Preview Deployments

For previewing documentation changes before merging, add a PR workflow that builds but does not deploy:

```yaml
name: Hugo PR Preview

on:
  pull_request:
    paths:
      - 'content/**'
      - 'layouts/**'
      - 'data/**'
      - 'hugo.toml'
      # Add mounted paths

jobs:
  build:
    runs-on: ubuntu-latest
    env:
      HUGO_VERSION: '0.142.0'
    steps:
      - name: Install Hugo CLI
        run: |
          wget -O ${{ runner.temp }}/hugo.deb https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-amd64.deb
          sudo dpkg -i ${{ runner.temp }}/hugo.deb

      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0

      - name: Build with Hugo
        run: hugo --minify --baseURL "https://example.com/"

      - name: Report
        run: |
          echo "## Hugo Build Summary" >> $GITHUB_STEP_SUMMARY
          echo "Pages: $(find public -name '*.html' | wc -l)" >> $GITHUB_STEP_SUMMARY
          echo "Total size: $(du -sh public | cut -f1)" >> $GITHUB_STEP_SUMMARY
```

## GitHub Pages Setup

Before the workflow deploys successfully:

1. Go to repository Settings > Pages
2. Set Source to "GitHub Actions"
3. The first workflow run will create the deployment

## Common Issues

| Issue | Fix |
|-------|-----|
| 404 on deployed site | Check `baseURL` matches your GitHub Pages URL |
| CSS/JS not loading | Ensure `baseURL` has trailing slash and uses HTTPS |
| Build succeeds but site empty | Check `public/` contains HTML files; verify content is not all `draft: true` |
| Module resolution fails | Ensure `go.mod` and `go.sum` are committed |
| SCSS compilation fails | Use `hugo_extended_` in the download URL |
| Old content after deploy | GitHub Pages CDN may cache; wait a few minutes or check cache headers |
````

**Step 2:** Commit.

```bash
git add hugo-repo/skills/hugo-github-actions/SKILL.md
git commit -m "add hugo-github-actions skill"
```

---

### Task 8: hugo-s3-deployment skill

**Files:**

- Create: `hugo-repo/skills/hugo-s3-deployment/SKILL.md`

**Step 1:** Create the skill file with the following content.

````markdown
---
name: Hugo S3 Deployment
description: >-
  This skill should be used when the user asks about deploying Hugo to AWS S3,
  S3 static site hosting, Hugo with CloudFront, the hugo deploy command,
  configuring Hugo's deployment section, GitHub Actions with AWS OIDC for Hugo,
  S3 bucket static website configuration, CloudFront cache invalidation after
  Hugo deploy, or choosing between GitHub Pages and S3 for Hugo hosting.
  Also applies when GitHub Pages limitations (auth, custom headers, size)
  require an alternative deployment target.
version: 1.0.0
---

## Overview

AWS S3 + CloudFront is an alternative to GitHub Pages for hosting Hugo sites. Choose S3 when you need: custom authentication, multiple sites from one repo, custom HTTP headers, sites larger than GitHub Pages limits (1GB), or deployment to a private CDN.

Hugo has a built-in `hugo deploy` command that syncs the built site to S3 with intelligent caching and content-type handling.

## When to Choose S3 over GitHub Pages

| Factor | GitHub Pages | S3 + CloudFront |
|--------|-------------|-----------------|
| Cost | Free | ~$1-5/month for small sites |
| Setup | Minimal | Moderate (S3, CloudFront, IAM) |
| Custom domain | Yes (with limits) | Yes (full control) |
| HTTPS | Automatic | Via CloudFront |
| Auth/access control | Public only | IAM, signed URLs, WAF |
| Custom headers | No | Yes |
| Size limit | 1GB | Unlimited |
| Multiple sites | One per repo | Unlimited |
| Build location | GitHub Actions only | Any CI |

## Hugo Deploy Configuration

Add to `hugo.toml`:

```toml
[deployment]
  [[deployment.targets]]
    name = "production"
    URL = "s3://my-site-bucket?region=us-east-1"

  [[deployment.matchers]]
    pattern = "^.+\\.(js|css|svg|ttf|woff|woff2)$"
    cacheControl = "max-age=31536000, immutable"
    gzip = true

  [[deployment.matchers]]
    pattern = "^.+\\.(png|jpg|jpeg|gif|webp)$"
    cacheControl = "max-age=31536000, immutable"
    gzip = false

  [[deployment.matchers]]
    pattern = "^.+\\.(html|xml|json)$"
    cacheControl = "max-age=300"
    gzip = true
```

Deploy locally:

```bash
hugo deploy --target production
```

## S3 Bucket Setup

Create an S3 bucket configured for static website hosting:

```bash
# Create bucket
aws s3 mb s3://my-site-bucket --region us-east-1

# Enable static website hosting
aws s3 website s3://my-site-bucket \
  --index-document index.html \
  --error-document 404.html

# Set bucket policy for public access (if not using CloudFront OAI)
aws s3api put-bucket-policy --bucket my-site-bucket --policy '{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": "*",
    "Action": "s3:GetObject",
    "Resource": "arn:aws:s3:::my-site-bucket/*"
  }]
}'
```

## GitHub Actions with OIDC (Recommended)

OIDC federation eliminates long-lived AWS credentials. GitHub Actions authenticates directly with AWS using short-lived tokens.

### IAM Setup

Create an OIDC provider and IAM role (one-time setup):

```bash
# Create OIDC provider for GitHub Actions
aws iam create-open-id-connect-provider \
  --url https://token.actions.githubusercontent.com \
  --client-id-list sts.amazonaws.com \
  --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1
```

Create IAM role trust policy (`trust-policy.json`):

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "Federated": "arn:aws:iam::ACCOUNT_ID:oidc-provider/token.actions.githubusercontent.com"
    },
    "Action": "sts:AssumeRoleWithWebIdentity",
    "Condition": {
      "StringEquals": {
        "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
      },
      "StringLike": {
        "token.actions.githubusercontent.com:sub": "repo:USERNAME/REPO:ref:refs/heads/master"
      }
    }
  }]
}
```

Create the role with S3 and CloudFront permissions:

```bash
aws iam create-role \
  --role-name hugo-deploy \
  --assume-role-policy-document file://trust-policy.json

aws iam put-role-policy \
  --role-name hugo-deploy \
  --policy-name hugo-deploy-policy \
  --policy-document '{
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": ["s3:PutObject", "s3:DeleteObject", "s3:ListBucket"],
        "Resource": [
          "arn:aws:s3:::my-site-bucket",
          "arn:aws:s3:::my-site-bucket/*"
        ]
      },
      {
        "Effect": "Allow",
        "Action": "cloudfront:CreateInvalidation",
        "Resource": "arn:aws:cloudfront::ACCOUNT_ID:distribution/DISTRIBUTION_ID"
      }
    ]
  }'
```

### GitHub Actions Workflow

Create `.github/workflows/hugo-s3-deploy.yml`:

```yaml
name: Deploy Hugo site to S3

on:
  push:
    branches: [master, main]
    paths:
      - 'content/**'
      - 'layouts/**'
      - 'static/**'
      - 'assets/**'
      - 'data/**'
      - 'themes/**'
      - 'hugo.toml'
      - 'go.mod'
      - 'go.sum'
      - '.github/workflows/hugo-s3-deploy.yml'
  workflow_dispatch:

permissions:
  id-token: write
  contents: read

env:
  HUGO_VERSION: '0.142.0'
  AWS_REGION: 'us-east-1'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Install Hugo CLI
        run: |
          wget -O ${{ runner.temp }}/hugo.deb https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-amd64.deb
          sudo dpkg -i ${{ runner.temp }}/hugo.deb

      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0

      - name: Configure AWS credentials (OIDC)
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: arn:aws:iam::ACCOUNT_ID:role/hugo-deploy
          aws-region: ${{ env.AWS_REGION }}

      - name: Cache Hugo modules
        uses: actions/cache@v4
        with:
          path: |
            ~/.cache/hugo_cache
            /tmp/hugo_cache
          key: ${{ runner.os }}-hugo-${{ hashFiles('go.sum') }}
          restore-keys: |
            ${{ runner.os }}-hugo-

      - name: Build with Hugo
        run: hugo --minify

      - name: Deploy to S3
        run: hugo deploy --target production --maxDeletes 100

      - name: Invalidate CloudFront cache
        run: |
          aws cloudfront create-invalidation \
            --distribution-id ${{ vars.CLOUDFRONT_DISTRIBUTION_ID }} \
            --paths "/*"
```

Store `CLOUDFRONT_DISTRIBUTION_ID` as a repository variable (Settings > Secrets and variables > Actions > Variables), not a secret — it is not sensitive.

## Environment-Based Deployment

Use separate buckets for staging and production:

```toml
[deployment]
  [[deployment.targets]]
    name = "staging"
    URL = "s3://my-site-staging?region=us-east-1"

  [[deployment.targets]]
    name = "production"
    URL = "s3://my-site-production?region=us-east-1"
```

```bash
hugo deploy --target staging
hugo deploy --target production
```

## Common Issues

| Issue | Fix |
|-------|-----|
| Access Denied on deploy | Check IAM role trust policy matches repo/branch |
| OIDC auth fails | Verify `id-token: write` permission in workflow |
| Old content after deploy | Run CloudFront invalidation |
| MIME types wrong | Hugo deploy handles this; check `deployment.matchers` config |
| 403 on site access | Check S3 bucket policy or CloudFront OAI configuration |
````

**Step 2:** Commit.

```bash
git add hugo-repo/skills/hugo-s3-deployment/SKILL.md
git commit -m "add hugo-s3-deployment skill"
```

---

### Task 9: Agents

**Files:**

- Create: `hugo-repo/agents/hugo-site-architect.md`
- Create: `hugo-repo/agents/hugo-build-doctor.md`

**Step 1:** Create hugo-site-architect agent.

```markdown
---
name: hugo-site-architect
description: >
  Hugo site architecture and scaffolding specialist for repositories.
  Use when the user asks to set up a Hugo site for a repository, plan
  documentation structure, scaffold Hugo configuration with module mounts,
  design site architecture for a monorepo or multi-project repo, migrate
  from another static site generator (Jekyll, MkDocs, Docusaurus), or
  integrate multiple Hugo features (mounts, themes, data, shortcodes)
  for a specific project.

  <example>
  Context: User wants to add a Hugo documentation site to their monorepo
  user: "Set up a Hugo site for this repo — we have docs in each service directory"
  assistant: "I'll use the hugo-site-architect agent to analyze your repo structure, identify all docs directories, and generate the Hugo configuration with module mounts."
  <commentary>
  Multi-directory Hugo setup requires analyzing repo structure and generating mount configuration.
  </commentary>
  </example>

  <example>
  Context: User wants to migrate from MkDocs to Hugo
  user: "We're using MkDocs but want to switch to Hugo for better performance"
  assistant: "I'll use the hugo-site-architect agent to plan the migration — mapping your MkDocs structure to Hugo's content model and preserving your URL scheme."
  <commentary>
  Static site generator migration requires understanding both source and target systems.
  </commentary>
  </example>

  <example>
  Context: User needs to add a new section to an existing Hugo site
  user: "We just added a new microservice and need its docs in the Hugo site"
  assistant: "I'll use the hugo-site-architect agent to add the module mount, create the section index, and update the navigation."
  <commentary>
  Adding a new mounted section requires coordinating config, content, and navigation changes.
  </commentary>
  </example>
model: sonnet
color: blue
tools: ["Read", "Write", "Edit", "Bash", "Glob", "Grep", "WebSearch", "WebFetch"]
---

You are a Hugo site architect who designs and scaffolds Hugo sites for repositories with complex content structures. You think like a documentation engineer — systematic about content organization, precise about configuration, and practical about what the site needs to do.

**Site Analysis:**

- Start by analyzing the repository structure: find all directories containing markdown or documentation
- Identify the content hierarchy: what should be top-level sections, what should be nested
- Check for existing Hugo configuration, themes, or static site generator configs
- Assess whether the repo needs module mounts (multiple docs directories) or a simple content tree

**Hugo Configuration:**

- Generate `hugo.toml` with appropriate module mounts mapping each docs directory to the content tree
- Include standard mounts (content, static, layouts, data, assets) alongside custom mounts — omitting standard mounts breaks Hugo
- Set sensible defaults: baseURL, title, language, pagination
- Configure menus and navigation based on the content structure
- Pin Hugo version and use Hugo modules (not git submodules) for themes

**Content Organization:**

- Create `_index.md` section pages for every section in the content tree, including parent directories of mounted content
- Set appropriate `weight` values for section ordering
- Design the URL structure to be logical and stable
- Ensure every mounted directory has an `_index.md` in its source docs/ folder

**Theme Selection:**

- Recommend themes that support documentation sites (Docsy, Book, Geekdoc, Ananke) based on the project's needs
- Configure theme parameters for the specific project
- Set up layout overrides only when the theme doesn't support a needed feature

**Deployment:**

- Generate GitHub Actions workflow with content-aware path filtering
- Include path entries for every mounted docs directory
- Configure for GitHub Pages or S3 based on the user's needs

**Process:**

1. Analyze repository structure — what directories have content?
2. Design the content hierarchy — how should it map to the site?
3. Generate Hugo configuration with module mounts
4. Scaffold directory structure and section index pages
5. Recommend and configure a theme
6. Set up deployment workflow

**Do Not:**

- Create Hugo configuration without first analyzing the repo structure
- Omit standard mounts when adding custom mounts (this is a common and breaking mistake)
- Hardcode content in templates when data files would be more maintainable
- Recommend themes without considering the project's actual needs
- Set up deployment without content-aware path filtering in repos with code alongside docs
```

**Step 2:** Create hugo-build-doctor agent.

```markdown
---
name: hugo-build-doctor
description: >
  Hugo build diagnostic and troubleshooting specialist.
  Use when Hugo build fails, hugo server shows errors or warnings, the
  deployed site renders incorrectly, content is missing from the built site,
  GitHub Actions deployment workflow fails, the user needs help interpreting
  cryptic Hugo error messages, builds are unexpectedly slow, or assets
  (CSS, images, JS) are not loading or processing correctly.

  <example>
  Context: Hugo build fails with a template error
  user: "Hugo build is failing with 'execute of template failed: template: index.html:12:14'"
  assistant: "I'll use the hugo-build-doctor agent to diagnose the template error — Hugo error messages reference internal template execution and need careful interpretation."
  <commentary>
  Hugo template errors are notoriously cryptic, requiring analysis of the template chain.
  </commentary>
  </example>

  <example>
  Context: Content missing from deployed site
  user: "I added docs for the new service but they don't show up on the site"
  assistant: "I'll use the hugo-build-doctor agent to check the mount configuration, section index pages, and build output to find why the content isn't appearing."
  <commentary>
  Missing content usually means a mount issue, missing _index.md, or draft content.
  </commentary>
  </example>

  <example>
  Context: GitHub Actions deployment fails
  user: "The Hugo deploy workflow is failing in CI but works locally"
  assistant: "I'll use the hugo-build-doctor agent to compare the CI environment with your local setup — common issues include Hugo version mismatches, missing modules, and baseURL configuration."
  <commentary>
  CI/local divergence in Hugo builds has specific common causes that need systematic checking.
  </commentary>
  </example>
model: sonnet
color: yellow
tools: ["Read", "Bash", "Glob", "Grep"]
---

You are a Hugo build diagnostics specialist who systematically identifies and fixes Hugo build problems. You think like a debugger — gather evidence first, form hypotheses, then verify before suggesting fixes.

**Diagnostic Approach:**

- Read the exact error message before suggesting anything — Hugo errors contain line numbers and template names that point directly to the problem
- Check `hugo.toml` configuration for common issues: missing standard mounts, incorrect baseURL, theme misconfiguration
- Verify content structure: `_index.md` vs `index.md`, front matter validity, proper section hierarchy
- Check theme compatibility: does the theme support the Hugo version being used?
- For deployment issues: compare CI Hugo version with local, check module resolution, verify path filtering

**Common Issue Patterns:**

Template errors:
- "execute of template failed" — check the referenced template file and line number
- "can't evaluate field X" — the template references a variable or method that doesn't exist in the context
- "partial not found" — the partial file is missing or the path is wrong

Content issues:
- Pages not appearing — check `draft: true`, missing `_index.md`, mount not configured
- Wrong layout applied — check template lookup order, section vs single vs list
- Broken links — check `ref` and `relref` shortcodes point to valid content paths

Mount issues:
- Standard mounts missing — when custom mounts are defined, all default mounts must be explicitly included
- Mount source path wrong — paths are relative to repo root
- Missing section index — every mounted directory needs `_index.md`

Build/deploy issues:
- SCSS fails in CI — need Hugo extended edition
- Modules not resolving — `go.mod` and `go.sum` must be committed
- baseURL wrong — causes broken CSS/JS paths in production

**Process:**

1. Read the exact error message or symptom description
2. Check the most likely cause first (don't shotgun suggestions)
3. Read the relevant files (config, templates, content) to verify
4. Provide the specific fix with file path and line reference
5. Explain why it broke so the user can prevent it next time

**Do Not:**

- Suggest multiple possible fixes without reading the files first
- Skip reading the actual error message
- Assume the issue without checking — "it's probably X" is not diagnosis
- Recommend reinstalling Hugo or clearing caches as a first step
- Ignore the difference between local and CI environments when debugging deployment issues
```

**Step 3:** Commit.

```bash
git add hugo-repo/agents/
git commit -m "add hugo-site-architect and hugo-build-doctor agents"
```

---

### Task 10: Commands

**Files:**

- Create: `hugo-repo/commands/hugo-init.md`
- Create: `hugo-repo/commands/hugo-serve.md`
- Create: `hugo-repo/commands/hugo-deploy.md`
- Create: `hugo-repo/commands/hugo-add-section.md`

**Step 1:** Create hugo-init command.

```markdown
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
   - Verify each mounted directory has its own `_index.md`

6. **Install theme** (if specified or if user accepts recommendation):
   ```bash
   hugo mod get github.com/theNewDynamic/gohugo-theme-ananke
   ```

7. **Update `.gitignore`:**
   ```
   public/
   resources/
   .hugo_build.lock
   ```

8. **Verify** by running `hugo server --buildDrafts` and reporting the result.

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
```

**Step 2:** Create hugo-serve command.

```markdown
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
```

**Step 3:** Create hugo-deploy command.

```markdown
---
name: hugo-deploy
description: Set up or trigger Hugo site deployment
arguments:
  - name: target
    description: "Deployment target: pages or s3 (default: pages)"
    required: false
---

# Deploy Hugo Site

Set up or trigger deployment to GitHub Pages or AWS S3.

## Process

1. **Check for existing deployment workflow:**
   ```bash
   ls .github/workflows/hugo-*.yml 2>/dev/null
   ```

2. **If no workflow exists:**
   - Ask user: GitHub Pages or S3?
   - Generate the appropriate GitHub Actions workflow
   - For S3: also generate the `[deployment]` section in `hugo.toml`
   - Commit the workflow file
   - Explain next steps (GitHub Pages source setting, or AWS OIDC setup)

3. **If workflow exists:**
   - Validate the workflow configuration
   - Check path filters include all mounted docs directories
   - Verify Hugo version is current
   - Report status and explain how to trigger deployment

## Usage

```
/hugo-deploy [target]
```

## Examples

```
/hugo-deploy          # Check existing or set up GitHub Pages
/hugo-deploy pages    # Set up GitHub Pages deployment
/hugo-deploy s3       # Set up S3 deployment with OIDC
```

## Post-Setup

### GitHub Pages

After creating the workflow:

1. Go to Settings > Pages
2. Set Source to "GitHub Actions"
3. Push a commit that changes content to trigger the first deploy

### S3

After creating the workflow:

1. Set up AWS OIDC provider and IAM role (see hugo-s3-deployment skill)
2. Add `CLOUDFRONT_DISTRIBUTION_ID` as a repository variable
3. Push a commit to trigger the first deploy
```

**Step 4:** Create hugo-add-section command.

```markdown
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
   - If not provided, infer from the directory path (e.g., `new-plugin/docs` → `plugins/new-plugin`)

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
```

**Step 5:** Commit.

```bash
git add hugo-repo/commands/
git commit -m "add hugo-init, hugo-serve, hugo-deploy, hugo-add-section commands"
```

---

### Task 11: Templates

**Files:**

- Create: `hugo-repo/templates/hugo.toml.tmpl`
- Create: `hugo-repo/templates/github-pages.yml.tmpl`
- Create: `hugo-repo/templates/s3-deploy.yml.tmpl`
- Create: `hugo-repo/templates/_index.md.tmpl`

**Step 1:** Create hugo.toml template.

```toml
# Hugo Site Configuration
# Generated by hugo-repo plugin

baseURL = '{{BASE_URL}}'
languageCode = 'en-us'
title = '{{SITE_TITLE}}'

# Theme (installed as Hugo module)
# [module]
#   [[module.imports]]
#     path = 'github.com/theNewDynamic/gohugo-theme-ananke'

# Module Mounts
# When defining ANY custom mounts, you MUST include all standard mounts
# or those directories will stop working.

# Standard mounts (always include these)
[[module.mounts]]
source = 'content'
target = 'content'

[[module.mounts]]
source = 'static'
target = 'static'

[[module.mounts]]
source = 'layouts'
target = 'layouts'

[[module.mounts]]
source = 'data'
target = 'data'

[[module.mounts]]
source = 'assets'
target = 'assets'

[[module.mounts]]
source = 'archetypes'
target = 'archetypes'

# Custom content mounts — add one per docs directory
# [[module.mounts]]
# source = 'plugin-a/docs'
# target = 'content/plugins/plugin-a'

# Site parameters
[params]
  description = '{{SITE_DESCRIPTION}}'

# Markup configuration
[markup]
  [markup.goldmark]
    [markup.goldmark.renderer]
      unsafe = false
  [markup.tableOfContents]
    startLevel = 2
    endLevel = 4

# Output formats
[outputs]
  home = ['HTML', 'RSS']
  section = ['HTML', 'RSS']
  page = ['HTML']

# Pagination
[pagination]
  pagerSize = 20
```

**Step 2:** Create GitHub Pages workflow template.

```yaml
# Hugo GitHub Pages Deployment
# Generated by hugo-repo plugin
#
# Prerequisites:
# 1. Go to Settings > Pages > Source > set to "GitHub Actions"
# 2. Add mounted docs directory paths to the 'paths' filter below

name: Deploy Hugo site to GitHub Pages

on:
  push:
    branches: [master, main]
    paths:
      # Hugo content and configuration
      - 'content/**'
      - 'layouts/**'
      - 'static/**'
      - 'assets/**'
      - 'data/**'
      - 'themes/**'
      - 'hugo.toml'
      - 'hugo.yaml'
      - 'hugo.json'
      - 'config/**'
      - 'go.mod'
      - 'go.sum'
      # Mounted docs directories (add one line per mount)
      # - 'plugin-a/docs/**'
      # - 'plugin-b/docs/**'
      # This workflow
      - '.github/workflows/hugo-deploy.yml'
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

defaults:
  run:
    shell: bash

jobs:
  build:
    runs-on: ubuntu-latest
    env:
      HUGO_VERSION: '0.142.0'
      HUGO_ENVIRONMENT: production
    steps:
      - name: Install Hugo CLI
        run: |
          wget -O ${{ runner.temp }}/hugo.deb https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-amd64.deb
          sudo dpkg -i ${{ runner.temp }}/hugo.deb

      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0

      - name: Setup Pages
        id: pages
        uses: actions/configure-pages@v5

      - name: Install Node.js dependencies
        run: |
          [[ -f package-lock.json || -f npm-shrinkwrap.json ]] && npm ci || true

      - name: Cache Hugo modules
        uses: actions/cache@v4
        with:
          path: |
            ~/.cache/hugo_cache
            /tmp/hugo_cache
          key: ${{ runner.os }}-hugo-${{ hashFiles('go.sum') }}
          restore-keys: |
            ${{ runner.os }}-hugo-

      - name: Build with Hugo
        run: |
          hugo \
            --minify \
            --baseURL "${{ steps.pages.outputs.base_url }}/"

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: ./public

  deploy:
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

**Step 3:** Create S3 deploy workflow template.

```yaml
# Hugo S3 Deployment
# Generated by hugo-repo plugin
#
# Prerequisites:
# 1. Set up AWS OIDC provider and IAM role (see hugo-s3-deployment skill)
# 2. Replace ACCOUNT_ID with your AWS account ID
# 3. Add CLOUDFRONT_DISTRIBUTION_ID as a repository variable
# 4. Add mounted docs directory paths to the 'paths' filter below

name: Deploy Hugo site to S3

on:
  push:
    branches: [master, main]
    paths:
      - 'content/**'
      - 'layouts/**'
      - 'static/**'
      - 'assets/**'
      - 'data/**'
      - 'themes/**'
      - 'hugo.toml'
      - 'go.mod'
      - 'go.sum'
      # Mounted docs directories (add one line per mount)
      # - 'plugin-a/docs/**'
      # - 'plugin-b/docs/**'
      - '.github/workflows/hugo-s3-deploy.yml'
  workflow_dispatch:

permissions:
  id-token: write
  contents: read

env:
  HUGO_VERSION: '0.142.0'
  AWS_REGION: 'us-east-1'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Install Hugo CLI
        run: |
          wget -O ${{ runner.temp }}/hugo.deb https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_extended_${HUGO_VERSION}_linux-amd64.deb
          sudo dpkg -i ${{ runner.temp }}/hugo.deb

      - name: Checkout
        uses: actions/checkout@v4
        with:
          submodules: recursive
          fetch-depth: 0

      - name: Configure AWS credentials (OIDC)
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: arn:aws:iam::ACCOUNT_ID:role/hugo-deploy
          aws-region: ${{ env.AWS_REGION }}

      - name: Cache Hugo modules
        uses: actions/cache@v4
        with:
          path: |
            ~/.cache/hugo_cache
            /tmp/hugo_cache
          key: ${{ runner.os }}-hugo-${{ hashFiles('go.sum') }}
          restore-keys: |
            ${{ runner.os }}-hugo-

      - name: Build with Hugo
        run: hugo --minify

      - name: Deploy to S3
        run: hugo deploy --target production --maxDeletes 100

      - name: Invalidate CloudFront cache
        run: |
          aws cloudfront create-invalidation \
            --distribution-id ${{ vars.CLOUDFRONT_DISTRIBUTION_ID }} \
            --paths "/*"
```

**Step 4:** Create section index template.

```markdown
---
title: "{{SECTION_TITLE}}"
description: "{{SECTION_DESCRIPTION}}"
weight: {{WEIGHT}}
---

{{SECTION_CONTENT}}
```

**Step 5:** Commit.

```bash
git add hugo-repo/templates/
git commit -m "add hugo.toml, GitHub Actions, and section index templates"
```

---

## Summary of Commits

1. `add hugo-repo plugin scaffolding and marketplace registration` (Task 1)
2. `add hugo-fundamentals skill` (Task 2)
3. `add hugo-module-mounts skill` (Task 3)
4. `add hugo-themes skill` (Task 4)
5. `add hugo-content-authoring skill` (Task 5)
6. `add hugo-data-templates skill` (Task 6)
7. `add hugo-github-actions skill` (Task 7)
8. `add hugo-s3-deployment skill` (Task 8)
9. `add hugo-site-architect and hugo-build-doctor agents` (Task 9)
10. `add hugo-init, hugo-serve, hugo-deploy, hugo-add-section commands` (Task 10)
11. `add hugo.toml, GitHub Actions, and section index templates` (Task 11)
