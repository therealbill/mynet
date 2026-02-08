---
title: "How to write a tutorial"
description: "Create a learning-oriented tutorial using doc-tutorial-writer with the correct scope, checkpoints, and progressive skill-building."
weight: 2
doc_type: how-to
prerequisites:
  - "The diataxis-docs plugin installed"
  - "A clear understanding of what the tutorial should teach"
  - "Access to the source code or tool being documented"
est_time: "15 minutes"
roles: ["developer", "technical writer"]
stability: stable
---

# How to Write a Tutorial

Create a learning-oriented tutorial using the **doc-tutorial-writer** agent that guides beginners through building a working project with checkpoints and progressive skill-building.

**Goal:** Produce a tutorial that guarantees success -- if a reader follows every step exactly, they end up with a working result.

## Prerequisites

- The diataxis-docs plugin installed and available
- A clear understanding of what the learner should build
- Access to the source code or tool being documented
- Understanding of the difference between tutorials and how-to guides (see [Diataxis in Practice](../../explanation/diataxis-in-practice/))

## Steps

### 1. Define the learning outcome

Before invoking the agent, decide what the learner will build. A good tutorial outcome is:

- Concrete and demonstrable (a running application, a completed configuration, a working pipeline)
- Achievable in 90-120 minutes for onboarding tutorials
- Valuable enough that the learner feels accomplished

Write a one-sentence outcome: "The learner will build [specific thing] that [demonstrates value]."

### 2. List the minimal prerequisites

Identify only what the learner must have before starting. Keep prerequisites minimal:

- Required software and versions
- Required accounts or access
- Required prior knowledge (link to other tutorials if needed)

Do not list nice-to-have knowledge. The tutorial teaches everything else.

### 3. Invoke the doc-tutorial-writer agent

Pass the learning outcome, prerequisites, and source code path to the agent:

```
Write an onboarding tutorial for <project> at <source-path>.
Learning outcome: The learner will build <specific thing>.
Prerequisites: <list>.
Target time: 90 minutes.
Place the output at <docs-path>/tutorials/<filename>.md
```

The agent reads the source code, designs a progressive learning journey, and writes the complete tutorial.

### 4. Verify the tutorial structure

Check that the generated tutorial contains:

- **Frontmatter** with title, description, doc_type, prerequisites, est_time, roles, stability
- **"What you'll build"** section describing the concrete outcome
- **Sequential steps** each introducing one new concept
- **Complete code and commands** at every step (no fragments or "fill in the blank")
- **Expected output** shown after every command
- **Checkpoints** every 3-5 steps with verification commands
- **"What you built"** summary at the end
- **"Next steps"** linking to how-to guides, reference, and explanations
- **Troubleshooting** section covering common setup issues

### 5. Test every step

Walk through the tutorial from start to finish on a clean environment. Verify:

- Every command runs without error
- Every expected output matches what the learner will see
- Checkpoints correctly validate progress
- The final outcome matches what was promised

Fix any steps that fail. A tutorial that breaks partway through is worse than no tutorial.

### 6. Remove choices and alternatives

Scan the tutorial for anywhere the learner is asked to make a decision. Tutorials follow one path only. Replace any "you can choose between X and Y" with the recommended option stated as the only option.

### 7. Add cross-links

Add links to related documentation:

- Reference pages for APIs or commands used in the tutorial
- How-to guides for tasks the learner might want to do next
- Explanation pages for concepts mentioned but not fully explained

Use Hugo ref shortcodes: `{{</* ref "reference/agents" */>}}`

## Verify It Works

A correctly written tutorial meets these criteria:

- A complete beginner following the steps exactly produces the promised outcome
- The tutorial takes approximately the stated time (within 20%)
- No step requires the learner to make choices or apply judgment
- Every code block is complete and copy-pasteable
- Checkpoints catch common mistakes before they cascade

## Troubleshooting

**The tutorial exceeds 120 minutes.**
Reduce the scope. Remove features that are not essential to demonstrating the core value. Move advanced features to a follow-up tutorial or how-to guide.

**The agent included explanations of why things work.**
This is a type violation. Tutorials show what to do, not why. Remove conceptual paragraphs and replace them with brief "What just happened?" boxes or link to an explanation page.

**The agent presented multiple approaches.**
Tutorials must follow one path. Edit the output to use only the recommended approach. Move alternative approaches to a how-to guide if they are worth documenting.

**Steps assume knowledge the learner does not have.**
Add the missing knowledge as a prerequisite or add earlier steps that establish it. Tutorials must not skip steps, even obvious ones.

## See Also

- [doc-tutorial-writer agent specification](../../reference/agents/)
- [Why tutorials are learning-oriented, not task-oriented](../../explanation/diataxis-in-practice/)
- [Restructuring existing docs that mix tutorial content with other types](../../howto/restructure-docs-to-diataxis/)
- [Validate the tutorial passes quality checks](../../howto/validate-doc-quality/)
