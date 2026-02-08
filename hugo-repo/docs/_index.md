---
title: "Hugo Repo"
description: "Hugo site generation with skills, agents, commands, and templates for documentation sites"
weight: 10
---

# Hugo Repo

Comprehensive Hugo site management with 7 skills, 2 agents, 4 commands, and 4 templates covering site scaffolding, content authoring, theme customization, and deployment.

## Components

| Type | Name | Description |
|------|------|-------------|
| Skill | hugo-fundamentals | Hugo site initialization, directory structure, and development workflow |
| Skill | hugo-content-authoring | Shortcodes, render hooks, taxonomies, and page bundles |
| Skill | hugo-data-templates | Data-driven site content with YAML/JSON/TOML files |
| Skill | hugo-themes | Theme installation, customization, and creation |
| Skill | hugo-module-mounts | Mounting multiple directories into Hugo's content tree |
| Skill | hugo-github-actions | GitHub Pages deployment with CI/CD workflows |
| Skill | hugo-s3-deployment | AWS S3 deployment with CloudFront and OIDC auth |
| Agent | hugo-site-architect | Designs and scaffolds Hugo sites for complex repositories |
| Agent | hugo-build-doctor | Diagnoses and troubleshoots Hugo build failures |
| Command | /hugo-init | Scaffold a Hugo site for the current repository |
| Command | /hugo-serve | Start the Hugo development server |
| Command | /hugo-deploy | Set up or trigger Hugo deployment |
| Command | /hugo-add-section | Mount a new directory as a content section |
| Template | hugo.toml.tmpl | Hugo site configuration template |
| Template | github-pages.yml.tmpl | GitHub Actions workflow for GitHub Pages |
| Template | s3-deploy.yml.tmpl | GitHub Actions workflow for S3 deployment |
| Template | _index.md.tmpl | Section index page template |

## Documentation

- [Getting Started](tutorials/getting-started/) — Scaffold a Hugo site and deploy to GitHub Pages
- [Set Up Hugo Site](howto/set-up-hugo-site/) — Configure module mounts for your repo
- [Add Content Section](howto/add-content-section/) — Mount a new directory into the site
- [Deploy to GitHub Pages](howto/deploy-to-github-pages/) — CI/CD with GitHub Actions
- [Deploy to S3](howto/deploy-to-s3/) — AWS S3 with CloudFront
- [Customize Theme](howto/customize-theme/) — Override and extend Hugo themes
- [Skills Reference](reference/skills/) — All 7 skills specification
- [Agents Reference](reference/agents/) — Both agents specification
- [Commands Reference](reference/commands/) — All 4 commands specification
- [Templates Reference](reference/templates/) — All 4 templates specification
- [Architecture](explanation/architecture/) — How 17 components fit together
- [Design Decisions](explanation/design-decisions/) — Why 7 skills, module mounts, separate deployments
- [Component Interaction](explanation/component-interaction/) — How skills, agents, commands, and templates work together
