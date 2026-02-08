---
title: "Agents Reference"
description: "Technical specifications for both hugo-repo agents including model, tools, trigger conditions, and capabilities"
weight: 2
---

# Agents Reference

The hugo-repo plugin includes 2 agents. Agents are autonomous components that perform multi-step work using their assigned tools. Unlike skills (which provide knowledge), agents take actions: reading files, running commands, writing configuration, and making changes to the repository.

---

## hugo-site-architect

| Field | Value |
|-------|-------|
| Model | sonnet |
| Color | blue |
| Tools | Read, Write, Edit, Bash, Glob, Grep, WebSearch, WebFetch |

### Description

Hugo site architecture and scaffolding specialist for repositories. Designs and scaffolds Hugo sites for repositories with complex content structures.

### Trigger conditions

- User asks to set up a Hugo site for a repository
- User asks to plan documentation structure for a Hugo site
- User asks to scaffold Hugo configuration with module mounts
- User asks to design site architecture for a monorepo or multi-project repo
- User asks to migrate from another static site generator (Jekyll, MkDocs, Docusaurus)
- User asks to integrate multiple Hugo features (mounts, themes, data, shortcodes) for a project

### Trigger examples

| User input | Agent response pattern |
|------------|----------------------|
| "Set up a Hugo site for this repo -- we have docs in each service directory" | Analyzes repo structure, identifies docs directories, generates Hugo configuration with module mounts |
| "We're using MkDocs but want to switch to Hugo for better performance" | Plans migration mapping MkDocs structure to Hugo content model, preserving URL scheme |
| "We just added a new microservice and need its docs in the Hugo site" | Adds module mount, creates section index, updates navigation |

### Process

1. Analyze repository structure -- find all directories containing content
2. Design the content hierarchy -- map directories to site sections
3. Generate Hugo configuration with module mounts
4. Scaffold directory structure and section index pages
5. Recommend and configure a theme
6. Set up deployment workflow

### Capabilities

| Capability | Details |
|------------|---------|
| Repository analysis | Scans for markdown files, docs directories, existing SSG configs |
| Content hierarchy design | Maps directories to sections with proper nesting |
| Configuration generation | `hugo.toml` with standard and custom mounts |
| Section scaffolding | Creates `_index.md` for all sections and parent sections |
| Theme selection | Recommends themes based on project type (Docsy, Book, Geekdoc, Ananke) |
| Deployment setup | GitHub Actions workflow with content-aware path filtering |
| Migration planning | Maps source SSG structure to Hugo equivalents |

### Constraints

- Does not create Hugo configuration without first analyzing the repository structure
- Does not omit standard mounts when adding custom mounts
- Does not hardcode content in templates when data files would be more maintainable
- Does not recommend themes without considering the project's actual needs
- Does not set up deployment without content-aware path filtering in repos with code alongside docs

---

## hugo-build-doctor

| Field | Value |
|-------|-------|
| Model | sonnet |
| Color | yellow |
| Tools | Read, Bash, Glob, Grep |

### Description

Hugo build diagnostic and troubleshooting specialist. Systematically identifies and fixes Hugo build problems.

### Trigger conditions

- Hugo build fails
- Hugo server shows errors or warnings
- Deployed site renders incorrectly
- Content is missing from the built site
- GitHub Actions deployment workflow fails
- User needs help interpreting Hugo error messages
- Builds are unexpectedly slow
- Assets (CSS, images, JS) are not loading or processing correctly

### Trigger examples

| User input | Agent response pattern |
|------------|----------------------|
| "Hugo build is failing with 'execute of template failed: template: index.html:12:14'" | Reads the referenced template file and line number to diagnose the template error |
| "I added docs for the new service but they don't show up on the site" | Checks mount configuration, section index pages, and build output for missing content |
| "The Hugo deploy workflow is failing in CI but works locally" | Compares CI environment with local setup, checking Hugo version, modules, and baseURL |

### Diagnostic approach

1. Read the exact error message or symptom description
2. Check the most likely cause first
3. Read relevant files (config, templates, content) to verify
4. Provide the specific fix with file path and line reference
5. Explain why it broke to prevent recurrence

### Common issue patterns

| Category | Issue | Typical cause |
|----------|-------|---------------|
| Template errors | "execute of template failed" | Template file error at the referenced line number |
| Template errors | "can't evaluate field X" | Template references a nonexistent variable or method |
| Template errors | "partial not found" | Missing partial file or incorrect path |
| Content issues | Pages not appearing | `draft: true`, missing `_index.md`, or mount not configured |
| Content issues | Wrong layout applied | Template lookup order, section vs single vs list mismatch |
| Content issues | Broken links | `ref` or `relref` shortcodes pointing to invalid content paths |
| Mount issues | Standard mounts missing | Custom mounts defined without including default mounts |
| Mount issues | Mount source path wrong | Paths must be relative to repository root |
| Mount issues | Missing section index | Every mounted directory needs `_index.md` |
| Build/deploy issues | SCSS fails in CI | Hugo Extended edition required |
| Build/deploy issues | Modules not resolving | `go.mod` and `go.sum` not committed |
| Build/deploy issues | baseURL wrong | Causes broken CSS/JS paths in production |

### Constraints

- Does not suggest multiple possible fixes without reading the files first
- Does not skip reading the actual error message
- Does not assume the issue without checking
- Does not recommend reinstalling Hugo or clearing caches as a first step
- Does not ignore the difference between local and CI environments when debugging deployment issues
