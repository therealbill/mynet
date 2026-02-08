---
title: "Deploy to GitHub Pages"
description: "How to deploy a Hugo site to GitHub Pages using the /hugo-deploy command and GitHub Actions"
weight: 3
---

# Deploy to GitHub Pages

This guide covers setting up automated deployment of a Hugo site to GitHub Pages using the `/hugo-deploy` command and the hugo-github-actions skill.

## When to use this

Use this when you want to host your Hugo site on GitHub Pages with automated deployments triggered by content changes. GitHub Pages is the simpler deployment target -- free, minimal setup, and tightly integrated with GitHub.

If you need custom authentication, custom HTTP headers, or sites larger than 1GB, see {{< ref "howto/deploy-to-s3" >}} instead.

## Prerequisites

- An existing Hugo site in a GitHub repository
- Repository settings access (for GitHub Pages configuration)
- Hugo version 0.142.0 or later pinned in the workflow

## Steps

### 1. Run the `/hugo-deploy` command

```
/hugo-deploy pages
```

The command checks for an existing deployment workflow. If none exists, it generates `.github/workflows/hugo-deploy.yml` using the `github-pages.yml.tmpl` template.

### 2. Review the generated workflow

The workflow includes:

- **Trigger conditions**: Runs on push to master/main when Hugo-related files change
- **Content-aware path filtering**: Only rebuilds when content, layouts, static files, configuration, or mounted docs directories change
- **Hugo Extended edition**: Required for SCSS/PostCSS processing
- **Module caching**: Caches Hugo modules to speed up builds
- **Minified build**: Runs `hugo --minify` for production output
- **GitHub Pages deployment**: Uses the official `actions/deploy-pages` action

Check that all mounted docs directories appear in the path filter:

```yaml
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
  # Ensure your mounted directories are listed
  - 'plugin-a/docs/**'
  - 'plugin-b/docs/**'
```

If any mounted directory is missing from the filter, changes to that directory will not trigger a rebuild.

### 3. Configure GitHub Pages

In your GitHub repository:

1. Go to Settings
2. Select Pages in the left sidebar
3. Under Source, select "GitHub Actions"

This tells GitHub to use the workflow output rather than serving from a branch.

### 4. Push and verify

Commit the workflow file and push to your main branch:

```bash
git add .github/workflows/hugo-deploy.yml
git commit -m "Add Hugo GitHub Pages deployment workflow"
git push
```

The first push triggers the workflow. Monitor progress in the Actions tab of your GitHub repository.

After successful deployment, your site is available at `https://USERNAME.github.io/REPO/`.

### 5. Verify the deployed site

Check that:

- The site loads at the expected URL
- All sections from mounted directories appear
- CSS and JavaScript assets load correctly
- Internal links work (the `baseURL` is set automatically from GitHub Pages configuration)

## Updating the workflow

When you add a new content section with `/hugo-add-section`, the command automatically updates the path filter in the workflow. If you add directories manually, remember to add their path to the filter.

To update the pinned Hugo version, change the `HUGO_VERSION` environment variable in the workflow:

```yaml
env:
  HUGO_VERSION: '0.142.0'
```

## Troubleshooting

- **404 on deployed site**: The `baseURL` is set from `steps.pages.outputs.base_url` in the workflow. Verify GitHub Pages is configured with Source set to "GitHub Actions".
- **CSS/JS not loading**: Ensure the baseURL includes a trailing slash and uses HTTPS.
- **Build succeeds but site is empty**: Confirm content files do not all have `draft: true`. The production build excludes drafts.
- **Module resolution fails in CI**: Ensure `go.mod` and `go.sum` are committed to the repository.
- **SCSS compilation fails**: The workflow uses Hugo Extended edition (`hugo_extended_`). Verify the download URL in the workflow is correct.
- **Workflow does not trigger**: Check that the changed file's path matches one of the path filter entries.

## Related

- {{< ref "howto/deploy-to-s3" >}} -- Alternative deployment to AWS S3
- {{< ref "reference/commands" >}} -- Full `/hugo-deploy` command specification
- {{< ref "reference/templates" >}} -- The `github-pages.yml.tmpl` template specification
- {{< ref "explanation/design-decisions" >}} -- Why separate deployment skills exist
