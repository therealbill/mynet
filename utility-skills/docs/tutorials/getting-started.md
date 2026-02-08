---
title: "Getting Started with Utility Skills"
description: "Learn to handle nested code blocks with the markdown-nested-codeblocks skill"
weight: 1
---

This tutorial walks you through the most common problem the utility-skills plugin solves: writing documentation that contains code examples inside markdown code fences. You will trigger the **markdown-nested-codeblocks** skill, observe its guidance, and apply the k+1 backtick rule to produce correctly rendered markdown.

## Prerequisites

- The **utility-skills** plugin is installed in your Claude Code environment.
- You have a markdown file open or are about to write one that includes code examples.

## The Problem

Suppose you are writing a README that teaches readers how to use a Go code block in markdown. You need to show them the literal markdown syntax -- triple backticks and all -- inside your documentation. The moment you put triple backticks inside triple backticks, the rendering breaks.

## Step 1: Try the Naive Approach

Open a markdown file and write the following:

````markdown
```markdown
Here is how you write a Go code block:

```go
fmt.Println("Hello")
```
```
````

Preview the file. The second set of triple backticks closes the outer block prematurely. Everything after it renders as plain text instead of staying inside the code fence.

## Step 2: Ask Claude for Help

With the utility-skills plugin installed, ask Claude a question such as:

> "How do I nest code blocks in markdown?"

The **markdown-nested-codeblocks** skill activates automatically when it detects questions about nested fences, backtick escaping, or the k+1 rule. You do not need to invoke it by name.

## Step 3: Observe the Guidance

Claude explains the k+1 rule:

> If your content contains a run of k backticks, wrap it in a fence of k+1 backticks, or switch the outer fence to tildes.

Since the inner content uses 3 backticks, the outer fence needs 4 backticks.

## Step 4: Apply the Fix

Replace the broken markdown with a 4-backtick outer fence:

`````markdown
````markdown
Here is how you write a Go code block:

```go
fmt.Println("Hello")
```
````
`````

Preview the file again. The inner triple-backtick fences now render as literal text inside the outer code block. The markdown displays exactly as intended.

## Checkpoint

Verify the following before moving on:

- The outer fence uses exactly 4 backticks on both the opening and closing lines.
- The closing fence matches the opening fence in character and length.
- The inner triple-backtick fences appear as visible syntax in the rendered output.

## Step 5: Go Deeper -- Content with 4 Backticks

What if the content you are documenting already uses 4 backticks? The rule scales. Content containing a run of 4 backticks needs an outer fence of 5 backticks:

``````markdown
`````markdown
````markdown
Here is how you write a Go code block:

```go
fmt.Println("Hello")
```
````
`````
``````

Each layer adds one backtick. The rule works at any depth.

## Step 6: The Tilde Alternative

Instead of counting backticks, you can switch the outer fence to tildes. Since inner content almost always uses backticks, tildes avoid the counting problem entirely:

`````markdown
~~~markdown
Here is how you write a Go code block:

```go
fmt.Println("Hello")
```
~~~
`````

Tildes and backticks are interchangeable for fencing in CommonMark and GitHub Flavored Markdown. Use whichever produces the clearest result.

## Summary

The **markdown-nested-codeblocks** skill teaches one essential rule:

- **The k+1 rule:** If content contains a run of k backticks, wrap it in a fence of k+1 backticks.
- **The tilde shortcut:** Switch the outer fence to tildes to avoid counting entirely.
- **Any depth:** The rule scales to arbitrary nesting levels.

The skill activates proactively whenever Claude writes documentation containing code examples, preventing broken rendering before it occurs.

## See Also

- [Nest Code Blocks Correctly]({{< ref "howto/nest-code-blocks-correctly" >}}) -- step-by-step guide for applying the k+1 rule
- [Skill Reference]({{< ref "reference/skills" >}}) -- technical specification for markdown-nested-codeblocks
- [Architecture]({{< ref "explanation/architecture" >}}) -- why this is a skill rather than an agent
