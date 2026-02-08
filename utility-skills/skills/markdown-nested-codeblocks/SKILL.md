---
name: Markdown Nested Code Blocks
description: Fix broken markdown rendering when code blocks contain code blocks
when_to_use: when writing markdown that contains code examples, documenting code syntax, creating tutorials with code snippets, or fixing broken code block rendering in markdown files
version: 1.1.0
languages: markdown
---

# Markdown Nested Code Blocks

## Overview

When markdown contains code blocks inside other code blocks (nested fences), same-length fences break rendering. The k+1 rule solves this: wrap content containing k backticks in a fence of k+1 backticks, or switch the outer fence to tildes.

**Critical:** This rule applies ALWAYS - even under time pressure, even for "quick" documentation, even when you think "it'll probably work." Broken rendering is worse than taking 10 extra seconds to count backticks.

## The k+1 Rule

**If your content contains a run of k backticks, wrap it in a fence of k+1 backticks (or use tildes).**

```
Inner uses:  ```  (3 backticks)
Outer needs: ```` (4 backticks = 3+1)  OR  ~~~ (tildes)
```

**Rules:**
- Content has ``` (3 backticks) → Use ```` (4 backticks) or ~~~ (3 tildes) for outer fence
- Content has ```` (4 backticks) → Use ````` (5 backticks) or ~~~~ (4 tildes) for outer fence
- Close with the **exact same** fence you opened with (same character, same length)
- You don't need to escape the inner backticks; the longer outer fence makes them literal text

**Mnemonic:** Count inner backticks, add one for outer. Always.

## Common Scenario: Documentation with Code Examples

### ❌ Broken (Same Length Fences)

This markdown breaks because both fences use 3 backticks:

````markdown
```markdown
# Example

```go
func main() {
    fmt.Println("Hello")
}
```
```
````

**Problem:** First ``` after opening closes the outer block prematurely. The Go code isn't escaped.

### ✅ Fixed (k+1 Backticks)

Use 4 backticks for outer fence when inner has 3:

`````markdown
````markdown
# Example

```go
func main() {
    fmt.Println("Hello")
}
```
````
`````

### ✅ Alternative (Tildes)

Switch outer fence to tildes:

````markdown
~~~markdown
# Example

```go
func main() {
    fmt.Println("Hello")
}
```
~~~
````

## Quick Reference

| Inner Content | Outer Fence Options |
|---------------|-------------------|
| ` ``` ` (3 backticks) | ```` (4 backticks) or `~~~` (3 tildes) |
| ` ```` ` (4 backticks) | ````` (5 backticks) or `~~~~` (4 tildes) |
| ` ~~~ ` (3 tildes) | ```` (4 backticks) or `~~~~` (4 tildes) |

## Detection Checklist

Nested code blocks are needed when you're:

- Writing markdown that shows markdown syntax
- Creating tutorials that include code examples
- Documenting how to use code blocks
- Showing examples of configuration files with code in them

## Common Mistakes

### Using Same-Length Fences
```markdown
# ❌ WRONG
```markdown
Content with ```code``` blocks
```

# ✅ CORRECT
````markdown
Content with ```code``` blocks
````
```

### Forgetting to Match Closing Fence
```markdown
# ❌ WRONG - Opening is 4 backticks, closing is 3
````markdown
Content
```

# ✅ CORRECT - Both are 4 backticks
````markdown
Content
````
```

### Trying to Escape Inner Backticks
```markdown
# ❌ WRONG - Don't escape, use longer fence
```markdown
Content with \`\`\`code\`\`\` blocks
```

# ✅ CORRECT - Longer outer fence makes inner backticks literal
````markdown
Content with ```code``` blocks
````
```

## When to Use Tildes vs Longer Backticks

**Use tildes when:**
- Inner content uses backticks (more readable contrast)
- You want to avoid counting backtick length

**Use k+1 backticks when:**
- Consistency with existing codebase style
- Inner content uses tildes
- Both work - choose based on readability

## Red Flags - STOP and Apply k+1 Rule

If you're thinking any of these thoughts, **STOP** and apply the k+1 rule:

- "I'll just use triple backticks, that's standard"
- "No time to count backticks"
- "This will probably render fine"
- "GitHub will figure it out"
- "It's simple markdown, same fences should work"
- "I'll fix it later if it breaks"

**All of these lead to broken rendering.** Take 10 seconds now to count backticks and add one. Fixing broken documentation later wastes more time.

## Real-World Impact

Broken nested code blocks cause:
- Documentation that doesn't render on GitHub/GitLab
- Code examples visible as plain text instead of highlighted
- Tutorials that confuse readers
- README files that look unprofessional

Applying k+1 rule ensures clean rendering across all markdown parsers.
