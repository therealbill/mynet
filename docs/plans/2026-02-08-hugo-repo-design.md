# Hugo Repo Plugin Design

## Overview

**Plugin name:** `hugo-repo`
**Version:** 1.0.0
**Description:** Hugo site management for GitHub repositories — site creation, multi-directory content aggregation via module mounts, theme customization and creation, content authoring with custom shortcodes, data-driven content, and deployment to GitHub Pages or AWS S3.

The plugin provides expertise across the full Hugo lifecycle within a repository context: scaffolding a new site, organizing content from scattered `docs/` directories into a unified site using Hugo module mounts, selecting and customizing themes, authoring rich content with shortcodes and render hooks, leveraging Hugo's data directory for structured content, and deploying via GitHub Actions with content-aware path filtering.

**Target repo structures:**

- Plugin marketplaces with per-plugin documentation
- Monorepos with per-service/module documentation
- CLI tool collections with per-tool docs
- Research repos with data-driven content
- Any repo where documentation lives alongside code in multiple subdirectories

## Architecture

The plugin uses Hugo's native module mounts to map arbitrary repo directories into Hugo's content tree — no symlinks, no build-time copying. Each subdirectory's `docs/` folder is declared as a mount in `hugo.toml`, and Hugo treats the aggregated result as a single content tree. This is Hugo's official mechanism for exactly this problem.

Deployment supports two targets: GitHub Pages (standard, zero-infrastructure) and AWS S3 + CloudFront (for teams that need custom domains, auth, or multi-site support). Both use GitHub Actions with content-aware path filtering so documentation rebuilds only trigger when content, config, or theme files change — not on code-only commits.

## Components

### Skills (7)

#### 1. hugo-fundamentals

Core Hugo concepts for users new to Hugo.

**Covers:**

- Site initialization with `hugo new site`
- Directory structure: `content/`, `layouts/`, `static/`, `assets/`, `data/`, `themes/`
- Content front matter (YAML and TOML formats)
- Content types and archetypes
- The `hugo server` development workflow
- Configuration: `hugo.toml` vs `hugo.yaml` vs `hugo.json`
- Relationship between content organization and URL structure
- Section pages (`_index.md`) vs leaf bundles (`index.md`)
- Page bundles and resource organization

**Trigger phrases:** create a Hugo site, set up Hugo, new Hugo project, Hugo directory structure, Hugo front matter, Hugo configuration, hugo.toml, Hugo content organization.

#### 2. hugo-module-mounts

The central skill for multi-directory content aggregation. This is the plugin's core differentiator.

**Covers:**

- Hugo module initialization with `hugo mod init`
- The `module.mounts` configuration syntax
- Mapping arbitrary repo directories into Hugo's content tree
- Mount ordering and conflict resolution
- Preserving directory hierarchy in the generated site
- Patterns for common repo structures:
  - Monorepo: `services/*/docs/` mounted as `content/services/*/`
  - Marketplace: `plugins/*/docs/` mounted as `content/plugins/*/`
  - Tool collection: `tools/*/docs/` mounted as `content/tools/*/`
- Section `_index.md` files for mounted directories
- Navigation generation from mounted content
- Working example: mounting N subdirectory `docs/` folders with proper section structure

**Trigger phrases:** Hugo module mounts, mount multiple directories, aggregate docs directories, hugo mod init, module.mounts config, multi-directory Hugo site, monorepo documentation site.

#### 3. hugo-themes

Theme selection, customization, and creation.

**Emphasis:** Customizing existing themes (the common path). Theme creation from scratch covered at reference level.

**Covers:**

- Evaluating community themes (Hugo Themes directory, feature comparison)
- Installing themes as Hugo modules vs git submodules
- Theme override hierarchy: project layouts take precedence over theme layouts
- Overriding specific partials without forking the theme
- Customizing CSS and assets via the asset pipeline
- Theme configuration parameters
- Creating a theme from scratch (reference level):
  - `baseof.html` and the block system
  - Partial hierarchy and template lookup order
  - Asset pipeline (SCSS, PostCSS, minification)
  - Theme `theme.toml` metadata

**Trigger phrases:** Hugo theme, install theme, customize theme, override theme layout, Hugo partials, create Hugo theme, theme override, Hugo CSS, Hugo assets.

#### 4. hugo-content-authoring

Rich content features beyond basic markdown — shortcodes, render hooks, taxonomies.

**Covers:**

- Built-in shortcodes (`figure`, `highlight`, `ref`, `relref`, `gist`)
- Creating custom shortcodes:
  - Named vs positional parameters
  - Inner content (paired shortcodes)
  - Accessing page and site variables from shortcodes
  - Common patterns: callout boxes, API endpoint cards, citation blocks, tabs, code with filename headers
- Markdown render hooks:
  - Custom rendering for links, images, headings, code blocks
  - Adding attributes, classes, or wrapper elements
- Page resources and page bundles for co-located assets
- Taxonomies (tags, categories, custom taxonomies)
- Content summaries and descriptions
- Table of contents configuration
- Multilingual content (brief coverage, pointer to Hugo docs for depth)

**Trigger phrases:** Hugo shortcode, custom shortcode, render hook, Hugo taxonomy, page bundle, content formatting, Hugo callout, Hugo tabs component, Hugo code block customization.

#### 5. hugo-data-templates

Using Hugo's `data/` directory for structured content.

**Covers:**

- Loading YAML, JSON, and TOML data files
- Accessing data in templates via `.Site.Data`
- Using `range` to iterate over data collections
- Data-driven patterns:
  - Navigation and sidebar menus from data files
  - Plugin/service registries rendered as directory pages
  - Research datasets and bibliographies
  - Feature comparison matrices
  - Changelog rendering from structured data
  - Team/contributor pages
- Generating list pages from data files
- Remote data with `resources.GetRemote`
- Data file organization for large datasets
- Combining data with content (data supplements markdown pages)

**Trigger phrases:** Hugo data directory, Hugo data templates, data-driven content, Hugo .Site.Data, render data file, Hugo YAML data, data-driven navigation, Hugo bibliography, Hugo dataset.

#### 6. hugo-github-actions

GitHub Actions workflow for building and deploying Hugo sites to GitHub Pages.

**Covers:**

- Standard Pages deployment workflow (peaceiris/actions-hugo + peaceiris/actions-gh-pages or the official GitHub Pages action)
- Content-aware path filtering with `on.push.paths`:
  - Trigger on: `content/**`, `layouts/**`, `static/**`, `assets/**`, `data/**`, `themes/**`, `hugo.toml`, `go.mod`, `go.sum`, and each mounted docs directory
  - Skip on: code-only commits
- Caching Hugo modules and resources for faster builds
- Hugo extended edition for SCSS/PostCSS support
- PR preview deployments
- Environment-based configuration (baseURL per environment)
- Branch protection and deployment status checks
- Hugo version pinning and update strategy

**Trigger phrases:** Hugo GitHub Actions, deploy Hugo to GitHub Pages, Hugo CI/CD, Hugo deployment workflow, Hugo path filter, Hugo preview deployment, Hugo build cache.

#### 7. hugo-s3-deployment

S3 static site hosting as an alternative or primary deployment target.

**Covers:**

- S3 bucket configuration for static website hosting
- CloudFront distribution setup for CDN and HTTPS
- Hugo's built-in `hugo deploy` command and its `[deployment]` config section
- GitHub Actions workflow for S3 deployment:
  - AWS credential management via OIDC (preferred over long-lived access keys)
  - IAM role configuration for GitHub Actions
  - S3 sync with cache headers
  - CloudFront cache invalidation after deploy
- Environment-based deployment (staging vs production buckets)
- Custom domain and DNS configuration
- Cost considerations vs GitHub Pages
- When to choose S3 over Pages (auth requirements, multiple sites, custom headers, size limits)

**Trigger phrases:** Hugo S3 deployment, deploy Hugo to AWS, Hugo CloudFront, hugo deploy command, S3 static site, Hugo AWS GitHub Actions, Hugo S3 OIDC.

---

### Agents (2)

#### 1. hugo-site-architect

**Model:** sonnet | **Color:** blue
**Tools:** Read, Write, Edit, Bash, Glob, Grep, WebSearch, WebFetch

The primary agent for complex Hugo tasks requiring analysis of repository structure and multi-step planning.

**Responsibilities:**

- Initial site scaffolding: analyze repo structure, identify all directories containing documentation or content, recommend mount configuration, generate `hugo.toml`, scaffold Hugo directory structure
- Content reorganization planning: when documentation exists but needs restructuring for Hugo
- Section hierarchy design: determining how mounted directories map to site navigation and URL structure
- Multi-component integration: configuring theme + shortcodes + data + mounts to work together for a specific repo's needs
- Migration from other static site generators (Jekyll, MkDocs, Docusaurus) to Hugo

**Trigger description:** This agent should be used when the user asks to set up a Hugo site for a repository, plan documentation structure, scaffold Hugo configuration, design site architecture for a monorepo, migrate from another static site generator, or integrate multiple Hugo features (mounts, themes, data, shortcodes) for a specific project. Use for complex, multi-step Hugo tasks that require analyzing the repository and making architectural decisions.

#### 2. hugo-build-doctor

**Model:** sonnet | **Color:** yellow
**Tools:** Read, Bash, Glob, Grep

Diagnostic agent for Hugo build problems and deployment failures.

**Responsibilities:**

- Interpreting Hugo error messages (often cryptic, referencing internal template execution)
- Diagnosing build failures: broken mounts, missing theme components, template errors, front matter problems, shortcode failures
- Debugging deployment workflow errors: GitHub Actions failures, S3 permission issues, CloudFront invalidation problems
- Identifying rendering issues: content not appearing, wrong layout applied, missing sections, broken navigation
- Performance diagnosis: slow builds, unnecessary rebuilds, large asset processing

**Trigger description:** This agent should be used when Hugo build fails, hugo server shows errors, the site renders incorrectly, deployment workflows fail, content is missing from the built site, or the user needs help interpreting Hugo error messages. Also applies when builds are slow or assets are not processing correctly.

---

### Commands (4)

#### 1. /hugo-init

Analyzes the current repository structure and scaffolds a Hugo site.

**Behavior:**

1. Scan repo for directories containing markdown documentation
2. Identify existing `docs/` directories at any depth
3. Present findings and recommend site structure
4. Generate `hugo.toml` with appropriate module mounts
5. Scaffold Hugo directory structure (`layouts/`, `static/`, `assets/`, `data/`)
6. Create root `content/_index.md`
7. Optionally recommend and install a theme
8. Create `.gitignore` entries for Hugo build artifacts (`public/`, `resources/`)

#### 2. /hugo-serve

Starts the Hugo development server with appropriate configuration.

**Behavior:**

1. Verify Hugo is installed (suggest installation if not)
2. Resolve Hugo modules if `go.mod` exists
3. Start `hugo server` with flags appropriate to context:
   - `--buildDrafts` for development
   - `--disableFastRender` if content issues are suspected
   - `--navigateToChanged` for live editing
4. Report the local URL and any warnings

#### 3. /hugo-deploy

Sets up or triggers deployment to the configured target.

**Behavior:**

1. Check if a deployment workflow already exists
2. If not, ask: GitHub Pages or S3?
3. Generate the appropriate GitHub Actions workflow from template
4. For S3: generate the `[deployment]` section in `hugo.toml`
5. If workflow exists, validate it and explain how to trigger
6. For S3: verify AWS credential configuration

#### 4. /hugo-add-section

Adds a new content section to the site from an existing directory.

**Behavior:**

1. Accept a directory path (e.g., `new-plugin/docs/`)
2. Add the mount entry to `hugo.toml`
3. Create the section's `_index.md` with appropriate front matter (title, weight for ordering)
4. Update data-driven navigation if configured
5. Verify the new section appears with `hugo server`

---

### Templates (4)

Templates provide starting-point files that commands and agents populate with repo-specific values.

#### 1. hugo.toml.tmpl

Base site configuration with module mount patterns, common settings, and commented examples for customization.

#### 2. github-pages.yml.tmpl

GitHub Actions workflow for Pages deployment with content-aware path filtering, Hugo module caching, and extended edition support.

#### 3. s3-deploy.yml.tmpl

GitHub Actions workflow for S3 deployment with OIDC authentication, S3 sync, and CloudFront invalidation.

#### 4. _index.md.tmpl

Section page template with standard front matter (title, description, weight) for new content sections.

---

## Plugin Structure

```
hugo-repo/
├── .claude-plugin/plugin.json
├── agents/
│   ├── hugo-site-architect.md
│   └── hugo-build-doctor.md
├── skills/
│   ├── hugo-fundamentals/SKILL.md
│   ├── hugo-module-mounts/SKILL.md
│   ├── hugo-themes/SKILL.md
│   ├── hugo-content-authoring/SKILL.md
│   ├── hugo-data-templates/SKILL.md
│   ├── hugo-github-actions/SKILL.md
│   └── hugo-s3-deployment/SKILL.md
├── commands/
│   ├── hugo-init.md
│   ├── hugo-serve.md
│   ├── hugo-deploy.md
│   └── hugo-add-section.md
└── templates/
    ├── hugo.toml.tmpl
    ├── github-pages.yml.tmpl
    ├── s3-deploy.yml.tmpl
    └── _index.md.tmpl
```

## Component Interaction

Commands lean on skills for domain knowledge and may invoke agents for complex operations:

- `/hugo-init` uses hugo-fundamentals for scaffolding and hugo-module-mounts for analyzing the repo and generating mount config
- `/hugo-deploy` uses either hugo-github-actions or hugo-s3-deployment depending on the user's target choice
- `/hugo-add-section` uses hugo-module-mounts for the mount entry and hugo-content-authoring for the section page
- `/hugo-serve` uses hugo-fundamentals for the correct server invocation

The hugo-site-architect agent draws on all skills when planning a site setup. The hugo-build-doctor agent primarily uses hugo-fundamentals, hugo-module-mounts, and hugo-themes knowledge for diagnosis.

## Marketplace Registration

```json
{
  "name": "hugo-repo",
  "source": "./hugo-repo",
  "description": "Hugo site management for GitHub repositories — multi-directory content aggregation, theme customization, data-driven content, and deployment to GitHub Pages or AWS S3",
  "version": "1.0.0",
  "keywords": ["hugo", "static-site", "github-pages", "s3", "documentation", "module-mounts"]
}
```

## Key Design Decisions

1. **Module mounts over symlinks or build scripts.** Hugo's native module mount system is the official, cross-platform mechanism for mapping arbitrary directories into the content tree. No duplication, no fragile symlinks, declarative config.

2. **Content-aware path filtering for CI.** In repos where Hugo docs live alongside code, rebuilding the site on every code commit is wasteful. Path-filtered GitHub Actions workflows only trigger when content, config, or theme files change.

3. **OIDC over long-lived AWS keys.** For S3 deployment, the plugin recommends GitHub Actions OIDC federation for AWS credentials rather than storing access keys as secrets. More secure, no key rotation burden.

4. **Theme customization over creation.** Most users select and customize existing themes. Full theme creation is covered at reference level but the skills emphasize the override/extend workflow.

5. **Separate content-authoring from themes.** Shortcodes, render hooks, and taxonomies are content features independent of theme choice. Keeping them in their own skill prevents the themes skill from becoming a catch-all.

6. **Two focused agents, not one generalist.** Site architecture (planning/building) and build diagnosis (debugging/fixing) are distinct workflows with different tool needs and reasoning patterns. Keeping them separate gives clearer trigger boundaries.

7. **Skills cover full spectrum.** From hugo-fundamentals for first-time users through hugo-module-mounts and hugo-data-templates for advanced patterns. Similar to the gnu-make plugin's fundamentals-through-advanced structure.
