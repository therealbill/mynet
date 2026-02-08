---
title: "Design Decisions"
description: "Why the hugo-repo plugin uses 7 skills, separate deployment targets, module mounts over symlinks, and commands alongside agents"
weight: 2
---

# Design Decisions

This page explains the reasoning behind the major architectural choices in the hugo-repo plugin.

## Why 7 skills instead of fewer

The plugin could have used a single monolithic "hugo" skill covering all Hugo topics. Instead, it uses 7 separate skills, each focused on a specific domain. The reasons for this choice:

**Precision of activation.** Each skill has its own trigger phrases. When a user asks about shortcodes, only the hugo-content-authoring skill activates, not the deployment or module mount knowledge. This prevents irrelevant information from diluting the response.

**Manageable scope.** Each skill covers one topic thoroughly rather than many topics superficially. The hugo-module-mounts skill, for example, covers mount precedence, standard mount requirements, debugging commands, and common mistakes -- depth that would be impractical in a combined skill.

**Independent versioning.** Skills can be updated independently. If Hugo changes its deployment behavior, only hugo-github-actions and hugo-s3-deployment need updates, while the other five skills remain unchanged.

**Composability.** Complex tasks naturally activate multiple skills. Setting up a Hugo site for a monorepo might activate hugo-fundamentals, hugo-module-mounts, and hugo-themes simultaneously. The skill system handles this combination automatically.

The 7 skills map to natural boundaries in Hugo's feature set:

| Skill | Boundary |
|-------|----------|
| hugo-fundamentals | Core Hugo concepts (needed by everyone) |
| hugo-content-authoring | Content features beyond basic markdown |
| hugo-data-templates | Data-driven site patterns |
| hugo-themes | Visual presentation layer |
| hugo-module-mounts | Multi-directory content aggregation |
| hugo-github-actions | GitHub Pages deployment |
| hugo-s3-deployment | AWS S3 deployment |

These boundaries mirror how users think about Hugo. A user working on shortcodes is not simultaneously thinking about S3 bucket policies.

## Why separate deployment skills for GitHub Pages and S3

GitHub Pages and S3 deployment could have been a single "hugo-deployment" skill. They are separate because:

**Different prerequisite knowledge.** GitHub Pages deployment requires understanding GitHub Actions, permissions, and Pages configuration. S3 deployment requires understanding AWS IAM, OIDC federation, S3 bucket policies, and CloudFront. A user deploying to GitHub Pages does not need AWS knowledge, and vice versa.

**Different workflows.** The GitHub Pages workflow uses `actions/deploy-pages` with artifact upload. The S3 workflow uses `hugo deploy` with OIDC authentication and CloudFront invalidation. These are fundamentally different deployment mechanisms that happen to deploy the same Hugo output.

**Different troubleshooting paths.** When GitHub Pages deployment fails, the diagnostic process involves checking Pages settings, permissions, and baseURL. When S3 deployment fails, the diagnostic process involves checking IAM roles, trust policies, and bucket permissions. Combining these would make the hugo-build-doctor agent's diagnostic work less focused.

**Template alignment.** Each deployment target has its own template (`github-pages.yml.tmpl` and `s3-deploy.yml.tmpl`). The skill separation mirrors the template separation.

## Why module mounts over symlinks

Hugo's module mount system is the plugin's core mechanism for multi-directory content aggregation. Alternative approaches considered and rejected:

**Symlinks.** Symlinks from `content/` to scattered docs directories would achieve the same virtual content tree. However:

- Symlinks do not work reliably across operating systems (particularly Windows)
- Git's symlink handling varies by configuration
- CI environments may not follow symlinks by default
- Hugo's module mounts are the officially supported mechanism for this use case

**Build-time copying.** A pre-build script could copy docs directories into `content/`. However:

- This creates duplicate files that can drift from the originals
- The copy step adds build complexity and failure modes
- File watching during development does not work across the copy boundary
- Deleted source files leave orphans in the copy destination

**Hugo's union file system.** Module mounts create a virtual union file system without any of these drawbacks. Hugo resolves mounts at build time, watches the original source files, and handles precedence when paths conflict. This is the mechanism Hugo was designed to provide for exactly this use case.

The one requirement of module mounts -- that standard mounts must be explicitly included when custom mounts are defined -- is a known friction point. The plugin addresses this by including all standard mounts in the `hugo.toml.tmpl` template and by documenting the requirement prominently in the hugo-module-mounts skill.

## Why commands alongside agents

The plugin has both commands (deterministic, user-invoked actions) and agents (autonomous, multi-step workers). These overlap in capability: the `/hugo-init` command and the hugo-site-architect agent can both set up a Hugo site. The separation exists because they serve different interaction patterns.

**Commands are for known tasks.** When a user knows exactly what they want to do -- scaffold a site, start a server, add a section -- commands provide a direct path. The user types `/hugo-init` and gets a predictable process. There is no ambiguity about what will happen.

**Agents are for open-ended work.** When a user describes a situation rather than a task -- "we have a monorepo and need a documentation site" -- the agent determines the approach. The hugo-site-architect agent decides how to analyze the repository, what hierarchy to propose, and what configuration to generate.

**Commands use agents.** The `/hugo-init` command delegates its repository analysis to the hugo-site-architect agent. Commands are the structured entry point; agents are the execution engine. This means commands benefit from the agent's analysis capabilities while providing a predictable interface.

**Agents handle what commands cannot.** Migration from another static site generator, architectural design for an unusual repository structure, or integrating multiple Hugo features for a specific project -- these tasks are too variable for a fixed command process. Agents handle them by adapting their approach to the situation.

The division follows a general principle: if the task has a fixed, predictable process, it is a command. If the task requires analysis and judgment, it is an agent's responsibility.

## Related

- {{< ref "explanation/architecture" >}} -- Overview of how the 17 components fit together
- {{< ref "explanation/component-interaction" >}} -- How these design decisions play out in practice
- {{< ref "reference/skills" >}} -- Full skill specifications showing the domain boundaries
- {{< ref "reference/agents" >}} -- Full agent specifications showing tool and process differences
