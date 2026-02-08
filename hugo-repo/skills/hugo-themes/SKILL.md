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
