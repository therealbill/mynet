---
title: "Architecture"
description: "Why nested codeblocks break and how the skill addresses it"
weight: 1
---

This document explains why nested code fences break in markdown, how the k+1 rule resolves the problem, and why the markdown-nested-codeblocks component is implemented as a skill rather than an agent.

## Why Nested Fences Break

Markdown parsers match fence boundaries by looking for a line that starts with the same fence character at the same length as the opening fence. When inner content uses the same fence length as the outer fence, the parser encounters a matching boundary before reaching the intended closing line. It closes the outer block prematurely, and everything after that point renders as raw, unformatted text.

This is not a bug. The parser is following the specification correctly. The problem is that the author used identical delimiters for two different structural boundaries, creating an ambiguity that the parser resolves by closing at the first match.

## The k+1 Solution

The CommonMark specification (and GitHub Flavored Markdown, which extends it) allows code fences of any length greater than or equal to 3 characters. A code fence closes only when it encounters a line with the exact same character at the exact same length. No other combination of character or length will close it.

The k+1 rule exploits this: if the inner content contains a consecutive run of k backticks, the outer fence uses k+1 backticks. The inner fences, being shorter, cannot match the outer boundary. They become literal text inside the code block.

This approach scales to arbitrary depth. Content with 4 backticks uses a 5-backtick outer fence. Content with 5 backticks uses a 6-backtick outer fence. There is no practical limit.

## The Tilde Alternative

Tildes function identically to backticks for code fencing in CommonMark. A tilde fence closes only when it encounters another tilde fence of the same length. Backtick fences and tilde fences do not interact.

Since inner content almost always uses backticks, switching the outer fence to tildes avoids the need to count backtick lengths entirely. The two fence types operate in separate namespaces, so there is no ambiguity regardless of what the inner content contains. This makes tildes the simpler choice when the goal is clarity over convention.

## Why a Skill, Not an Agent

The Claude Code plugin system distinguishes between agents and skills. Agents are autonomous entities that use tools, maintain state, and execute multi-step workflows. Skills are knowledge injections -- focused rules or techniques that any agent benefits from knowing.

The markdown-nested-codeblocks component is a skill because it encapsulates a single, self-contained rule. It does not require tool access, multi-step reasoning, or autonomous decision-making. Any agent that writes markdown -- whether it is a documentation agent, a code generation agent, or a general-purpose assistant -- benefits from knowing the k+1 rule. Making it a skill means the knowledge is available across the entire system without duplicating it in every agent definition.

## Proactive Activation

The skill does not wait for the user to ask about nesting. It activates whenever Claude writes markdown that will contain code examples inside fences. This is a deliberate design choice: the most common failure mode is not that users cannot find the rule, but that the rule is not applied in the moment when documentation is being written. Proactive activation prevents broken rendering before it occurs, rather than fixing it after the fact.

This is the key advantage of encoding the k+1 rule as a skill with broad trigger conditions. A narrowly triggered skill that only responds to explicit questions about "backtick nesting" would miss the majority of cases where the rule is needed -- cases where the user is focused on writing content, not on markdown syntax.
