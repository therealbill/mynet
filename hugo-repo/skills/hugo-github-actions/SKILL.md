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
