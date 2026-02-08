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
