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
