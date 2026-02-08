---
title: "Nest Code Blocks Correctly"
description: "Apply the k+1 backtick rule to fix nested code block rendering"
weight: 1
---

This guide shows you how to write markdown that contains code fences inside other code fences without breaking the rendering. Follow these steps whenever you need to show markdown syntax examples, embed code block demonstrations, or document fenced code block usage.

## Goal

Produce markdown where inner code fences render as visible syntax inside an outer code block, rather than prematurely closing it.

## Steps

### 1. Identify the Inner Fences

Look at the content you want to wrap. Find the longest consecutive run of backticks (or tildes) in that content. This is your **k** value.

For example, if your content contains standard triple-backtick fences:

```text
```python
print("hello")
```
```

The longest backtick run is 3. So k = 3.

### 2. Apply the k+1 Rule

Add 1 to k. Use the result as the length of your outer fence.

- k = 3 means the outer fence uses 4 backticks.
- k = 4 means the outer fence uses 5 backticks.

### 3. Open the Outer Fence

Write the opening fence with k+1 backticks, followed by an optional language tag:

`````text
````markdown
`````

### 4. Write the Inner Content

Place your content -- including its original fences -- inside the outer fence. Do not modify or escape the inner backticks:

`````markdown
````markdown
```python
print("hello")
```
````
`````

### 5. Close with the Exact Same Fence

The closing fence must match the opening fence in both character type and length. If you opened with 4 backticks, close with exactly 4 backticks:

`````text
````
`````

A mismatch (opening with 4, closing with 3) leaves the block unclosed.

### 6. Alternative: Use Tildes Instead

If you prefer not to count backticks, switch the outer fence to tildes. Since inner content nearly always uses backticks, tildes create an unambiguous boundary:

`````markdown
~~~markdown
```python
print("hello")
```
~~~
`````

Tildes and backticks are equally valid fence characters in CommonMark and GitHub Flavored Markdown.

### 7. Verify the Output

Preview the markdown in your target renderer (GitHub, VS Code preview, Hugo, or any CommonMark-compliant tool). Confirm that:

- The outer code block displays as a single block.
- The inner fences appear as literal text, not as rendered code block boundaries.
- No raw text leaks outside the code block.

## Quick Reference

| Inner content has | Outer fence needed |
|---|---|
| ` ``` ` (3 backticks) | ```` (4 backticks) or `~~~` (3 tildes) |
| ```` (4 backticks) | ````` (5 backticks) or `~~~~` (4 tildes) |
| `~~~~` (4 tildes) | `~~~~~` (5 tildes) or ````` (5 backticks) |

## Common Pitfalls

- **Same-length fences:** Using 3 backticks for both inner and outer fences is the most frequent mistake. The parser closes the outer block at the first matching fence it encounters.
- **Mismatched closing fence:** Opening with 4 backticks but closing with 3. Always match the exact character and length.
- **Escaping inner backticks:** Backslash-escaping backticks inside code fences does not work. Use a longer outer fence instead.
