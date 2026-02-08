---
title: "Customize a Hugo Theme"
description: "How to customize a Hugo theme by overriding templates, modifying layouts, and adding custom CSS using the hugo-themes skill"
weight: 5
---

# Customize a Hugo Theme

This guide covers customizing a Hugo theme without modifying the theme's own files. The hugo-themes skill provides the knowledge Claude uses to guide theme customization.

## When to use this

Use this when you want to:

- Override specific templates or partials from a theme
- Add custom CSS or SCSS to a themed site
- Change the layout of specific sections
- Add site-wide elements (analytics, custom headers, footers)

## Steps

### 1. Identify what to change

Determine which part of the theme you want to modify. Common targets include:

- **Partials**: Header, footer, sidebar, navigation, head (for custom CSS/JS)
- **Layouts**: Single page layout, list page layout, section-specific layouts
- **Styles**: Colors, typography, spacing, responsive behavior

Ask Claude for help identifying the correct file:

```
I want to change the navigation in my Hugo site that uses the Book theme
```

Claude uses the hugo-themes skill to identify the relevant template files and the template lookup order.

### 2. Create the override file

Hugo's template lookup order checks the project's `layouts/` directory before the theme's. To override any theme template, create a file at the same path in your project.

For example, to override a theme's header partial:

```
layouts/partials/header.html
```

Copy the theme's original file first, then modify your copy. This preserves the original structure while letting you make targeted changes.

To find the original file's contents:

```bash
cat themes/your-theme/layouts/partials/header.html
```

If using Hugo modules instead of a `themes/` directory, the theme files are in the module cache. Use `hugo config` or browse the theme's repository to find the template source.

### 3. Override section-specific layouts

To change the layout for only one section, create a section-specific template. Hugo checks section-specific templates before `_default/` templates.

For a custom blog list page:

```
layouts/blog/list.html
```

For a custom single page layout in the docs section:

```
layouts/docs/single.html
```

### 4. Add custom CSS

Three approaches are available, in order of recommendation:

**Hugo Pipes (recommended)**: Place SCSS or CSS in `assets/` and process through Hugo's asset pipeline:

```html
<!-- layouts/partials/head-custom.html -->
{{ $style := resources.Get "scss/custom.scss" | toCSS | minify | fingerprint }}
<link rel="stylesheet" href="{{ $style.RelPermalink }}" integrity="{{ $style.Data.Integrity }}">
```

**Static file**: Place CSS in `static/css/custom.css` and reference it in a head partial override:

```html
<link rel="stylesheet" href="{{ "css/custom.css" | relURL }}">
```

**Theme parameters**: Some themes support CSS customization via configuration:

```toml
[params]
  customCSS = ["css/custom.css"]
```

### 5. Add custom shortcodes

Create shortcodes in `layouts/shortcodes/` to add reusable content components. The hugo-content-authoring skill covers shortcode creation in detail.

Example callout shortcode (`layouts/shortcodes/callout.html`):

```html
<div class="callout callout-{{ .Get "type" | default "info" }}">
  <strong>{{ .Get "title" | default "" }}</strong>
  {{ .Inner | markdownify }}
</div>
```

### 6. Verify changes

Start the development server to preview your customizations:

```bash
hugo server --buildDrafts
```

If changes do not appear, try disabling fast render:

```bash
hugo server --buildDrafts --disableFastRender
```

## Key principles

- **Never modify theme files directly**: All changes should be overrides in the project's `layouts/`, `assets/`, or `static/` directories.
- **Copy then modify**: When overriding a partial, copy the theme's version first to preserve its structure.
- **Template lookup order**: Project `layouts/` takes precedence over theme `layouts/`. Section-specific templates take precedence over `_default/`.
- **Hugo Pipes for CSS**: Use `assets/` and Hugo Pipes for CSS/SCSS processing with minification and fingerprinting, rather than placing unprocessed files in `static/`.

## Troubleshooting

- **Override not taking effect**: Verify the file path in `layouts/` exactly matches the theme's path. Check capitalization.
- **Broken layout after override**: Your override replaces the theme's file entirely. Make sure you copied the complete original before modifying.
- **SCSS compilation fails**: Hugo Extended edition is required for SCSS. Check your Hugo version with `hugo version` -- it should include "extended".
- **Theme update breaks overrides**: When updating a theme, review the changelog for template changes that affect your overrides.

## Related

- {{< ref "reference/skills" >}} -- hugo-themes and hugo-content-authoring skill specifications
- {{< ref "tutorials/getting-started" >}} -- Theme customization in the context of a full setup
- {{< ref "explanation/architecture" >}} -- How skills provide knowledge for theme work
