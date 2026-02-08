---
title: "Skills Reference"
description: "Technical specifications for all 7 hugo-repo skills including version, trigger phrases, and coverage areas"
weight: 1
---

# Skills Reference

The hugo-repo plugin includes 7 skills. Each skill provides domain knowledge that Claude uses when responding to relevant questions. Skills are passive knowledge components -- they do not perform actions autonomously but inform Claude's responses and guide agents.

## hugo-fundamentals

| Field | Value |
|-------|-------|
| Version | 1.0.0 |
| Scope | Hugo site initialization, directory structure, configuration, content organization, development workflow |

**Trigger phrases**:

- Create a Hugo site
- Set up Hugo
- Start a new Hugo project
- Hugo directory structure
- Configure hugo.toml / hugo.yaml
- Hugo front matter
- Content types or archetypes
- Hugo server for local development
- Hugo content organization
- Section pages and page bundles
- Hugo basics
- URL structure in Hugo

**Topics covered**:

| Topic | Coverage |
|-------|----------|
| Site initialization | `hugo new site`, `hugo mod init` |
| Directory structure | archetypes, assets, content, data, i18n, layouts, static, themes |
| Configuration | `hugo.toml`, `hugo.yaml`, `hugo.json`, config directory splitting |
| Front matter | title, date, draft, weight, description, tags, categories, layout, aliases |
| Section pages | `_index.md` for list pages, section hierarchy |
| Leaf bundles | `index.md` for single pages with co-located resources |
| URL structure | Content path to URL mapping, `url` front matter, `[permalinks]` |
| Archetypes | Content templates for `hugo new`, lookup order |
| Development server | `hugo server`, `--buildDrafts`, `--navigateToChanged`, `--disableFastRender` |
| Production builds | `hugo`, `--environment`, `--minify`, output directory |
| Common mistakes | Spaces in filenames, missing `_index.md`, `index.md` vs `_index.md`, draft content |

---

## hugo-content-authoring

| Field | Value |
|-------|-------|
| Version | 1.0.0 |
| Scope | Shortcodes, render hooks, taxonomies, page bundles, content formatting |

**Trigger phrases**:

- Hugo shortcodes
- Creating custom shortcodes
- Hugo render hooks
- Custom rendering for links, images, headings, code blocks
- Hugo taxonomies (tags, categories, custom)
- Page bundles and page resources
- Content formatting beyond basic markdown
- Callout boxes, tabs, admonitions in Hugo
- Customizing how Hugo processes markdown

**Topics covered**:

| Topic | Coverage |
|-------|----------|
| Built-in shortcodes | `highlight`, `figure`, `ref`, `relref`, `gist` |
| Custom shortcodes | Named parameters, positional parameters, `.Inner`, `.Page` access |
| Shortcode patterns | Callout boxes, tabs, code blocks with filenames, badges |
| Render hooks | `render-link.html`, `render-image.html`, `render-heading.html`, `render-codeblock.html` |
| Taxonomies | Built-in (tags, categories), custom taxonomies, front matter usage |
| Page bundles | Leaf bundles with `index.md`, page resources, `.Resources` access |
| Table of contents | `.TableOfContents`, `[markup.tableOfContents]` configuration |

---

## hugo-data-templates

| Field | Value |
|-------|-------|
| Version | 1.0.0 |
| Scope | Data directory, YAML/JSON/TOML data files, `.Site.Data`, data-driven patterns |

**Trigger phrases**:

- Hugo data directory
- Loading YAML/JSON/TOML data files in templates
- `.Site.Data`
- Data-driven navigation or menus
- Rendering datasets from data files
- Generating content from data
- `range` to iterate over data collections
- `resources.GetRemote` for external data
- Organizing data files

**Topics covered**:

| Topic | Coverage |
|-------|----------|
| Data access | `.Site.Data`, file path to accessor mapping |
| File formats | YAML, JSON, TOML in data directory |
| Nested directories | Subdirectory access, hyphen-to-underscore conversion |
| Navigation from data | Data-driven menus with children, sorting by weight |
| Plugin registry | Data-driven registry rendering with status and keywords |
| Bibliography | Citation shortcode from data, DOI linking |
| Comparison matrix | Feature comparison tables from data |
| Changelog | Release-based changelog from structured data |
| Remote data | `resources.GetRemote`, `transform.Unmarshal` |
| Organization | Domain-based data directory structure |

---

## hugo-themes

| Field | Value |
|-------|-------|
| Version | 1.0.0 |
| Scope | Theme installation, customization, template override, asset pipeline, theme creation |

**Trigger phrases**:

- Hugo themes
- Installing a theme
- Customizing a theme
- Overriding theme layouts or partials
- Hugo template lookup order
- Creating a Hugo theme from scratch
- Hugo asset pipeline (SCSS, PostCSS)
- `baseof.html`
- Theme configuration parameters

**Topics covered**:

| Topic | Coverage |
|-------|----------|
| Module installation | `hugo mod get`, `[module.imports]` configuration |
| Template lookup order | Project layouts override theme layouts |
| Partial overrides | Copy-then-modify pattern |
| Section overrides | Section-specific templates before `_default/` |
| CSS customization | Static file, Hugo Pipes (SCSS), theme parameters |
| Hugo Pipes | `toCSS`, `minify`, `fingerprint`, `resources.Get` |
| Theme parameters | `[params]`, social links, navigation options |
| Theme creation | `hugo new theme`, scaffold structure, `baseof.html`, `single.html`, `list.html` |

---

## hugo-module-mounts

| Field | Value |
|-------|-------|
| Version | 1.0.0 |
| Scope | Module mount configuration, multi-directory content aggregation, union file system |

**Trigger phrases**:

- Hugo module mounts
- Mounting multiple directories into Hugo content tree
- Aggregating docs from scattered directories
- `hugo mod init`
- `module.mounts` in hugo.toml
- Unified documentation site from a monorepo
- Mapping subdirectory docs folders into a single Hugo site
- Hugo union file system

**Topics covered**:

| Topic | Coverage |
|-------|----------|
| Prerequisites | Go installation, `hugo mod init`, `go.mod` and `go.sum` |
| Standard mounts | content, static, layouts, data, assets, i18n, archetypes |
| Custom mounts | `source` and `target` configuration |
| Plugin marketplace pattern | Per-plugin docs directories mounted to content tree |
| Monorepo pattern | Per-service docs directories |
| CLI tool pattern | Per-tool docs directories |
| Section index pages | `_index.md` requirements for mounts and parent sections |
| Mount precedence | First mount wins for conflicting paths |
| Non-content mounts | Layouts, data, and static directory mounts |
| Adding mounts | Procedure for adding new mounts |
| Debugging | `hugo list all`, `hugo config mounts`, `hugo --verbose` |
| Common mistakes | Missing standard mounts, missing `_index.md`, absolute paths, uncommitted go files |

---

## hugo-github-actions

| Field | Value |
|-------|-------|
| Version | 1.0.0 |
| Scope | GitHub Pages deployment, GitHub Actions workflow, path filtering, module caching |

**Trigger phrases**:

- Deploying Hugo to GitHub Pages
- GitHub Actions workflow for Hugo
- Hugo CI/CD pipeline
- Content-aware path filtering for Hugo builds
- Caching Hugo modules in CI
- Hugo PR preview deployments
- Configuring GitHub Pages deployment

**Topics covered**:

| Topic | Coverage |
|-------|----------|
| Standard workflow | Complete GitHub Actions workflow for Hugo + GitHub Pages |
| Path filtering | Content-aware `paths` filter, per-mount path entries |
| Hugo Extended | Extended vs standard edition, SCSS requirement |
| Version pinning | `HUGO_VERSION` environment variable, deliberate updates |
| PR preview | Build-only workflow for pull request validation |
| GitHub Pages setup | Settings > Pages > Source > GitHub Actions |
| Module caching | `actions/cache` with `go.sum` hash key |
| Workflow permissions | `contents: read`, `pages: write`, `id-token: write` |
| Common issues | 404 errors, missing CSS/JS, empty site, module failures, SCSS failures |

---

## hugo-s3-deployment

| Field | Value |
|-------|-------|
| Version | 1.0.0 |
| Scope | AWS S3 deployment, CloudFront CDN, OIDC authentication, `hugo deploy` command |

**Trigger phrases**:

- Deploying Hugo to AWS S3
- S3 static site hosting
- Hugo with CloudFront
- `hugo deploy` command
- Hugo deployment configuration
- GitHub Actions with AWS OIDC for Hugo
- S3 bucket static website configuration
- CloudFront cache invalidation
- Choosing between GitHub Pages and S3

**Topics covered**:

| Topic | Coverage |
|-------|----------|
| S3 vs GitHub Pages | Cost, setup, domain, HTTPS, auth, headers, size, multi-site comparison |
| Hugo deploy config | `[deployment]`, targets, matchers, cache control, gzip |
| S3 bucket setup | Bucket creation, website hosting, bucket policy |
| OIDC authentication | IAM OIDC provider, trust policy, role creation, permissions |
| GitHub Actions workflow | Complete S3 deploy workflow with OIDC |
| CloudFront invalidation | Post-deploy cache invalidation |
| Environment-based deploy | Staging and production targets |
| Common issues | Access denied, OIDC failures, stale content, MIME types, 403 errors |
