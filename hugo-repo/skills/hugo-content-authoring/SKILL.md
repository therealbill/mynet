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
` ``yaml
key: value
` ``
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
