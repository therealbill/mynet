---
title: "Skills"
description: "Utility skills specifications"
weight: 1
---

## markdown-nested-codeblocks

| Field | Value |
|---|---|
| Name | Markdown Nested Code Blocks |
| Version | 1.1.0 |
| File | `skills/markdown-nested-codeblocks/SKILL.md` |

### Trigger Conditions

The skill activates when context matches any of the following:

- Writing markdown with code blocks inside code blocks
- Nested fences
- Triple-backtick examples inside documentation
- Fixing broken code block rendering
- Backtick escaping
- Fenced code block nesting
- The k+1 rule
- Showing markdown code blocks inside markdown
- Writing READMEs, documentation, or blog posts containing code examples inside markdown fences

### Activation

Proactive. The skill activates automatically when writing documentation that contains code examples inside markdown fences. It does not require explicit invocation.

### Core Rule (k+1)

Content containing a consecutive run of k backticks requires an outer fence of k+1 backticks. Alternatively, the outer fence may use tildes instead of backticks.

| Inner content | Outer fence options |
|---|---|
| ` ``` ` (3 backticks) | ```` (4 backticks) or `~~~` (3 tildes) |
| ```` (4 backticks) | ````` (5 backticks) or `~~~~` (4 tildes) |
| `~~~~` (4 tildes) | `~~~~~` (5 tildes) or ````` (5 backticks) |

### Fence Matching

The closing fence must use the exact same character (backtick or tilde) and the exact same length as the opening fence. A mismatch in character or length leaves the code block unclosed.

### Constraint

The rule applies unconditionally. There are no cases where same-length nested fences are acceptable, regardless of urgency or assumed parser tolerance.

### Compatibility

The k+1 rule is defined by the CommonMark specification and is supported by all CommonMark-compliant renderers, including:

- GitHub Flavored Markdown
- GitLab Flavored Markdown
- Hugo (Goldmark)
- VS Code Markdown Preview
- Pandoc

## See Also

- [Getting Started](../../tutorials/getting-started/) -- tutorial walkthrough
- [Nest Code Blocks Correctly](../../howto/nest-code-blocks-correctly/) -- step-by-step guide
- [Architecture](../../explanation/architecture/) -- design rationale
