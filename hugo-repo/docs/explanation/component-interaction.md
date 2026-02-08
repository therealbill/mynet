---
title: "Component Interaction"
description: "How skills inform agents, commands use templates, and agents handle what commands cannot in the hugo-repo plugin"
weight: 3
---

# Component Interaction

The 17 components of the hugo-repo plugin interact in defined patterns. This page explains the three primary interaction patterns: skills informing agents, commands using templates, and agents extending beyond what commands can do.

## Skills inform agents

Skills are passive knowledge. Agents are active workers. When an agent encounters a Hugo task, the relevant skills provide the domain knowledge the agent needs to make correct decisions.

### Example: hugo-site-architect and hugo-module-mounts

When the hugo-site-architect agent sets up a Hugo site for a repository, it needs to generate the `hugo.toml` configuration with module mounts. The hugo-module-mounts skill provides the knowledge that:

- Custom mounts replace Hugo's default mounts, so all standard mounts must be explicitly included
- Every mounted directory needs an `_index.md` file to appear as a proper section
- Parent sections also need `_index.md` files
- Mount precedence follows the order of entries in the configuration (first mount wins)
- Source paths must be relative to the repository root

Without this skill knowledge, the agent might generate a configuration that defines custom mounts but omits the standard mounts -- a common and breaking mistake that causes the `content/`, `static/`, and `layouts/` directories to stop working.

### Example: hugo-build-doctor and hugo-github-actions

When the hugo-build-doctor agent diagnoses a CI deployment failure, the hugo-github-actions skill provides knowledge about:

- The required workflow permissions (`contents: read`, `pages: write`, `id-token: write`)
- Hugo Extended edition being necessary for SCSS processing
- Content-aware path filtering and why changes might not trigger a rebuild
- The difference between local and CI Hugo versions
- Module caching keys and potential cache staleness

This knowledge lets the agent systematically check the most likely failure causes rather than guessing.

### How skill activation works

Skills activate based on pattern matching against the user's input and the current task context. Multiple skills can activate simultaneously. When the hugo-site-architect agent is setting up a site with a custom theme and module mounts, three skills might be active:

- hugo-fundamentals (site initialization, configuration structure)
- hugo-module-mounts (mount configuration, standard mounts requirement)
- hugo-themes (theme installation as Hugo module, template lookup order)

The agent draws on all active skills to make informed decisions. This composability is why the plugin uses 7 separate skills rather than a single combined skill.

## Commands use templates

Commands follow a defined process, and templates provide the scaffolding for their output. Each command that generates files uses one or more templates.

### /hugo-init uses hugo.toml.tmpl

The `/hugo-init` command generates the site's configuration file using the `hugo.toml.tmpl` template. The command:

1. Scans the repository to determine the site title, base URL, and which directories to mount
2. Fills in the template's placeholders (`{{BASE_URL}}`, `{{SITE_TITLE}}`, `{{SITE_DESCRIPTION}}`)
3. Adds custom mount entries for each discovered docs directory
4. Writes the completed configuration to `hugo.toml`

The template encodes the plugin's recommended configuration: Goldmark renderer with unsafe mode disabled, table of contents from levels 2-4, HTML and RSS output formats, and pagination at 20 items per page. These defaults come from the template, not from the command's logic.

### /hugo-init and /hugo-add-section use _index.md.tmpl

Both commands create section index pages. The `_index.md.tmpl` template provides the front matter structure with `{{SECTION_TITLE}}`, `{{SECTION_DESCRIPTION}}`, and `{{WEIGHT}}` placeholders.

For `/hugo-init`, the template is used for every section index page the command creates: the home page, parent sections for grouped content, and verification that each mounted directory has its own index.

For `/hugo-add-section`, the template is used when the mounted directory does not already have an `_index.md`, and when parent sections need to be created.

### /hugo-deploy uses deployment templates

The `/hugo-deploy` command uses one of two templates depending on the deployment target:

- `github-pages.yml.tmpl` for `/hugo-deploy pages`
- `s3-deploy.yml.tmpl` for `/hugo-deploy s3`

The command customizes the template by:

- Adding path filter entries for all mounted docs directories
- Setting the Hugo version
- For S3: adding the deployment configuration section to `hugo.toml`

The templates encode complex workflow configurations -- OIDC authentication, module caching, artifact upload -- that would be error-prone to generate from scratch each time.

## Agents handle what commands cannot

Commands are designed for predictable, well-defined tasks. Agents handle situations that are too variable for a fixed process.

### Migration from other static site generators

The `/hugo-init` command can scaffold a Hugo site for a repository, but it assumes the repository does not already have a static site generator. Migrating from Jekyll, MkDocs, or Docusaurus requires the hugo-site-architect agent because:

- The agent must understand the source SSG's content model and map it to Hugo's
- URL preservation requires analyzing the source configuration's permalink patterns
- Content files may need front matter translation (Jekyll's layout names differ from Hugo's)
- The source's navigation structure needs to be recreated in Hugo's model

No fixed command process can handle this variability. The agent analyzes the specific source and target, plans the migration, and executes it step by step.

### Build troubleshooting

The `/hugo-serve` command starts the development server, but when the server fails or content renders incorrectly, the hugo-build-doctor agent takes over. Troubleshooting is inherently investigative:

- The agent reads the exact error message to determine the class of problem
- It checks the most likely cause first based on the error pattern
- It reads the relevant files to verify its hypothesis
- It provides a specific fix with file path and line reference

This diagnostic process cannot be a command because the path through it depends on what the agent finds at each step.

### Complex architecture design

When a repository has an unusual structure -- nested monorepos, mixed documentation and code, multiple output sites from one repo -- the hugo-site-architect agent designs a custom architecture. It might:

- Propose a content hierarchy that does not follow the standard patterns
- Recommend data-driven navigation instead of menu configuration
- Suggest using multiple Hugo configurations for different site outputs
- Design a mount structure that handles conflicting paths with explicit precedence

This design work requires judgment that adapts to the specific repository, which is the defining characteristic of agent work.

## Interaction summary

| Pattern | Mechanism | Example |
|---------|-----------|---------|
| Skill informs agent | Skills provide domain knowledge that agents use for decisions | hugo-module-mounts tells hugo-site-architect about the standard mounts requirement |
| Command uses template | Commands fill template placeholders with project-specific values | `/hugo-init` fills `hugo.toml.tmpl` with discovered configuration |
| Agent extends command | Agents handle variable tasks that commands cannot cover with fixed processes | hugo-site-architect handles SSG migration that `/hugo-init` cannot |
| Multiple skills compose | Several skills activate together for complex tasks | fundamentals + mounts + themes all inform a site setup task |
| Command invokes agent | Commands delegate analysis to agents | `/hugo-init` uses hugo-site-architect for repository analysis |

## Related

- [Overview of the four component types](../../explanation/architecture/)
- [Why the plugin is structured this way](../../explanation/design-decisions/)
- [See these interactions in practice during a full setup](../../tutorials/getting-started/)
