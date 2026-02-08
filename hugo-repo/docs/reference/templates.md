---
title: "Templates Reference"
description: "Technical specifications for all 4 hugo-repo scaffolding templates including placeholders, default values, and usage context"
weight: 4
---

# Templates Reference

The hugo-repo plugin includes 4 templates used by commands to generate configuration and workflow files. Templates contain placeholders that are replaced with project-specific values during scaffolding.

---

## hugo.toml.tmpl

| Field | Value |
|-------|-------|
| Path | `templates/hugo.toml.tmpl` |
| Generated file | `hugo.toml` |
| Used by | `/hugo-init` command |

### Placeholders

| Placeholder | Description | Default |
|-------------|-------------|---------|
| `{{BASE_URL}}` | Site base URL | Detected from repository or set to placeholder |
| `{{SITE_TITLE}}` | Site title | Derived from repository name |
| `{{SITE_DESCRIPTION}}` | Site description for `[params]` | Empty string |
| `{{THEME}}` | Theme module path (commented out by default) | `github.com/theNewDynamic/gohugo-theme-ananke` |

### Structure

The template generates a complete Hugo configuration file with:

- **Base settings**: `baseURL`, `languageCode` (en-us), `title`
- **Theme import**: Commented-out `[module.imports]` block
- **Standard mounts**: content, static, layouts, data, assets, archetypes
- **Custom mount placeholders**: Commented-out examples for plugin docs directories
- **Site parameters**: `[params]` with description
- **Markup configuration**: Goldmark renderer with `unsafe = false`, table of contents levels 2-4
- **Output formats**: HTML and RSS for home and section pages, HTML for individual pages
- **Pagination**: `pagerSize = 20`

### Standard mounts included

| Source | Target |
|--------|--------|
| `content` | `content` |
| `static` | `static` |
| `layouts` | `layouts` |
| `data` | `data` |
| `assets` | `assets` |
| `archetypes` | `archetypes` |

---

## github-pages.yml.tmpl

| Field | Value |
|-------|-------|
| Path | `templates/github-pages.yml.tmpl` |
| Generated file | `.github/workflows/hugo-deploy.yml` |
| Used by | `/hugo-deploy pages` command |

### Configuration

| Setting | Value |
|---------|-------|
| Workflow name | Deploy Hugo site to GitHub Pages |
| Trigger branches | master, main |
| Hugo version | 0.142.0 (pinned via `HUGO_VERSION` env var) |
| Hugo edition | Extended (supports SCSS) |
| Build command | `hugo --minify --baseURL "${{ steps.pages.outputs.base_url }}/"` |
| Concurrency group | `pages` (cancel-in-progress: false) |

### Path filter entries

| Path pattern | Purpose |
|-------------|---------|
| `content/**` | Content files |
| `layouts/**` | Template files |
| `static/**` | Static assets |
| `assets/**` | Processed assets |
| `data/**` | Data files |
| `themes/**` | Theme files |
| `hugo.toml` | Hugo configuration |
| `hugo.yaml` | Hugo configuration (alternate) |
| `hugo.json` | Hugo configuration (alternate) |
| `config/**` | Split configuration directory |
| `go.mod` | Module dependencies |
| `go.sum` | Module checksums |
| `.github/workflows/hugo-deploy.yml` | The workflow file itself |

Custom mount paths are added as comments that the `/hugo-deploy` command uncomments and populates.

### Workflow permissions

| Permission | Level |
|-----------|-------|
| `contents` | read |
| `pages` | write |
| `id-token` | write |

### Jobs

**build**:

1. Install Hugo CLI (Extended edition)
2. Checkout repository (with submodules, full history)
3. Setup Pages (configure-pages action)
4. Install Node.js dependencies (if package-lock.json exists)
5. Cache Hugo modules (keyed on `go.sum` hash)
6. Build with Hugo (`--minify`)
7. Upload artifact (upload-pages-artifact action)

**deploy** (depends on build):

1. Deploy to GitHub Pages (deploy-pages action)

### Actions used

| Action | Version |
|--------|---------|
| `actions/checkout` | v4 |
| `actions/configure-pages` | v5 |
| `actions/cache` | v4 |
| `actions/upload-pages-artifact` | v3 |
| `actions/deploy-pages` | v4 |

---

## s3-deploy.yml.tmpl

| Field | Value |
|-------|-------|
| Path | `templates/s3-deploy.yml.tmpl` |
| Generated file | `.github/workflows/hugo-s3-deploy.yml` |
| Used by | `/hugo-deploy s3` command |

### Configuration

| Setting | Value |
|---------|-------|
| Workflow name | Deploy Hugo site to S3 |
| Trigger branches | master, main |
| Hugo version | 0.142.0 (pinned via `HUGO_VERSION` env var) |
| Hugo edition | Extended |
| AWS region | us-east-1 (default) |
| Build command | `hugo --minify` |
| Deploy command | `hugo deploy --target production --maxDeletes 100` |
| IAM role | `arn:aws:iam::ACCOUNT_ID:role/hugo-deploy` (requires substitution) |

### Path filter entries

Same as github-pages.yml.tmpl, with the workflow file path changed to `.github/workflows/hugo-s3-deploy.yml`.

### Workflow permissions

| Permission | Level |
|-----------|-------|
| `id-token` | write |
| `contents` | read |

### Job steps

1. Install Hugo CLI (Extended edition)
2. Checkout repository (with submodules, full history)
3. Configure AWS credentials via OIDC (`aws-actions/configure-aws-credentials@v4`)
4. Cache Hugo modules (keyed on `go.sum` hash)
5. Build with Hugo (`--minify`)
6. Deploy to S3 (`hugo deploy --target production --maxDeletes 100`)
7. Invalidate CloudFront cache (uses `vars.CLOUDFRONT_DISTRIBUTION_ID` repository variable)

### Actions used

| Action | Version |
|--------|---------|
| `actions/checkout` | v4 |
| `aws-actions/configure-aws-credentials` | v4 |
| `actions/cache` | v4 |

### Required setup

| Item | Type | Description |
|------|------|-------------|
| `ACCOUNT_ID` | Workflow edit | Replace placeholder with AWS account ID in the role ARN |
| `CLOUDFRONT_DISTRIBUTION_ID` | Repository variable | CloudFront distribution ID for cache invalidation |
| AWS OIDC provider | AWS resource | GitHub Actions OIDC identity provider in IAM |
| IAM role | AWS resource | Role with S3 and CloudFront permissions, trust policy for the repository |

---

## _index.md.tmpl

| Field | Value |
|-------|-------|
| Path | `templates/_index.md.tmpl` |
| Generated file | `_index.md` (in content section directories) |
| Used by | `/hugo-init` and `/hugo-add-section` commands |

### Placeholders

| Placeholder | Description |
|-------------|-------------|
| `{{SECTION_TITLE}}` | Title for the section, used in navigation and page heading |
| `{{SECTION_DESCRIPTION}}` | Description for metadata and list pages |
| `{{WEIGHT}}` | Sort order within the parent section (lower values appear first) |
| `{{SECTION_CONTENT}}` | Optional body content for the section index page |

### Generated front matter

```yaml
---
title: "{{SECTION_TITLE}}"
description: "{{SECTION_DESCRIPTION}}"
weight: {{WEIGHT}}
---
```

### Usage context

This template is used whenever a new section index page is needed:

- During `/hugo-init` for the home page and all discovered sections
- During `/hugo-add-section` for the mounted directory and any missing parent sections
- The hugo-site-architect agent uses this pattern when scaffolding sections
