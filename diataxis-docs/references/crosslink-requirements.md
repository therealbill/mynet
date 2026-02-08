# Cross-Linking Requirements

Use this reference when validating or creating cross-links between documentation pages.

## Minimum Cross-Link Counts

Each documentation type must link to other types. These are minimums — more is better when relevant.

| From Type | Minimum Links | Must Link To |
|-----------|--------------|--------------|
| Tutorial | 3+ | How-tos (next tasks), Reference (API used), Explanations (deeper understanding) |
| How-to | 5+ | Prerequisites (tutorials/how-tos), Reference (API details), Related how-tos, Explanations (context) |
| Reference | 2+ | How-tos (showing practical usage of each API) |
| Explanation | 4+ | How-tos (practical application), Reference (specifications), Related explanations |

## Link Placement

### Tutorials

```markdown
## Next steps

Now that you've built [thing], you can:

- **[How-to]:** [Task that builds on what was learned] →
- **[How-to]:** [Another common next task] →
- **[Explanation]:** [Understand how the system works internally] →
- **[Reference]:** [Full API/CLI reference for what you used] →
```

Tutorials link OUTWARD at the end. Don't interrupt the learning flow with links mid-tutorial.

### How-to Guides

```markdown
## Prerequisites

- Completed [Tutorial: Getting Started](../tutorials/getting-started.md)
- [Tool] v2.0+ installed ([installation how-to](./install.md))

...

[Within steps, link to reference for details:]
Configure the `auth` section (see [auth configuration reference](../reference/config.md#auth)).

...

## See also

- [How to configure advanced auth](./advanced-auth.md) — for SSO and MFA
- [Understanding the auth system](../explanation/auth-architecture.md) — design rationale
- [Auth API reference](../reference/api/auth.md) — full endpoint specifications
```

How-to guides link to prerequisites at the top, reference inline where needed, and related content at the bottom.

### Reference

```markdown
## See Also

- [How to create a resource](../how-to/create-resource.md)
- [How to bulk import resources](../how-to/bulk-import.md)
```

Reference pages link to how-tos showing practical usage. Keep links minimal and relevant — reference is for looking things up, not for exploration.

### Explanations

```markdown
## Related

- **Do it:** [How to configure event-driven processing](../how-to/event-processing.md)
- **Look it up:** [Event API reference](../reference/api/events.md)
- **Related concept:** [Understanding message queues](./message-queues.md)
- **Get started:** [Tutorial: Building your first event handler](../tutorials/event-handler.md)
```

Explanations link to all other types — they're the conceptual hub connecting practical and reference content.

## Red Flags

### Missing links

- **Tutorial with no "Next steps"** — reader finishes and has nowhere to go
- **How-to with no prerequisites** — reader may not have the foundation to succeed
- **Reference with no how-to links** — reader knows the API exists but not how to use it
- **Explanation with no how-to links** — reader understands the concept but can't apply it

### Wrong-direction links

- **Reference linking to tutorials** — reference readers are looking things up, not learning
- **Tutorial linking to other tutorials mid-flow** — interrupts the learning path
- **How-to linking to explanations mid-step** — distracts from the task

### Orphaned pages

A page with no incoming links from other pages is an orphan. It's effectively invisible to readers navigating the documentation. Every page should be reachable from at least one other page.

## Validation Checklist

For each page, verify:

- [ ] Page links to at least the minimum number of other pages (see table above)
- [ ] Links go to pages of the correct types (tutorials link to how-tos, not other tutorials)
- [ ] Links are placed in the right location (prerequisites at top, "see also" at bottom)
- [ ] All internal links resolve (no broken links)
- [ ] Page is linked TO from at least one other page (no orphans)
- [ ] Links use relative paths (not absolute URLs to the same docs site)
- [ ] Link text describes the destination (not "click here")
