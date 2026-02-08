# Diataxis Classification Signals

Use this reference when classifying existing documentation pages or deciding which type a new page should be.

## Quick Classification Matrix

| Signal | Tutorial | How-to | Reference | Explanation |
|--------|----------|--------|-----------|-------------|
| **Orientation** | Learning | Task | Information | Understanding |
| **Answers** | "How do I learn?" | "How do I do X?" | "What is the spec?" | "Why/how does it work?" |
| **Reader is** | Beginner studying | Practitioner working | Practitioner checking | Anyone thinking |
| **Author mood** | Teaching | Directing | Describing | Discussing |
| **Verb mood** | Imperative | Imperative | Indicative | Indicative |
| **Code examples** | Complete, working | Complete, working | Syntax only | Illustrative |

## Classification Signals by Type

### Tutorial Signals

Strong indicators:

- Title contains "Getting started", "Your first", "Build a", "Quickstart"
- Walks through building something from zero to working result
- Has sequential steps that depend on each other
- Shows expected output after each step
- Has checkpoints or "verify your progress" sections
- Addresses a complete beginner audience
- One path, no choices presented

Weak/ambiguous indicators:

- "Step 1, Step 2..." (how-tos also have steps)
- Uses imperative mood (how-tos also use imperative)
- Has code examples (all types can)

### How-to Guide Signals

Strong indicators:

- Title starts with "How to..." or "Configure...", "Set up...", "Deploy..."
- Solves one specific, practical problem
- Assumes the reader already has a working environment
- Has a verification/testing section at the end
- Has a troubleshooting section
- References prerequisites (often tutorials or other how-tos)
- Under 2000 words

Weak/ambiguous indicators:

- Numbered steps (tutorials also have these)
- "Run this command" style instructions (tutorials also have these)

### Reference Signals

Strong indicators:

- Consistent, repeating structure across items (every function documented the same way)
- Parameter tables with types, defaults, constraints
- Error code tables with conditions
- No opinions, advice, or "you should" language
- Could be auto-generated from source
- Organized by API surface, not by task
- Includes version/stability markers

Weak/ambiguous indicators:

- Has code examples (but reference examples show syntax, not workflows)
- Has tables (all types can use tables)

### Explanation Signals

Strong indicators:

- Title contains "Understanding", "Architecture", "Concepts", "How X works", "Why we..."
- Discusses design decisions and trade-offs
- Compares alternatives (what was considered, what was chosen, why)
- Uses analogies or mental models
- Has architectural diagrams
- Discusses history or evolution of the design
- Uses discursive language ("because", "the reason", "this trades X for Y")
- Links outward to how-tos and reference for practical application

Weak/ambiguous indicators:

- Discusses "how something works" (could be a tutorial if it's step-by-step)
- Has diagrams (reference can also have diagrams)

## Mixed-Type Detection

A page is mixed when it contains strong signals from multiple types. Common mixtures:

### Tutorial + Explanation (most common)

Symptom: Tutorial steps interrupted by paragraphs explaining how/why the underlying system works.

Fix: Extract the "how it works" content into an explanation page. Add a link: "Want to understand how this works? See [Understanding X](...)."

### How-to + Reference

Symptom: A guide that lists every API option and parameter while also walking through a task.

Fix: Extract the complete parameter/option listing into a reference page. Keep only the parameters needed for the specific task in the how-to.

### How-to + Explanation

Symptom: A guide that spends 500+ words explaining why before getting to the steps.

Fix: Extract the "why" into an explanation page. The how-to gets a one-sentence link: "For background on why this approach is recommended, see [Understanding X](...)."

### Reference + Explanation

Symptom: API documentation that includes design rationale, recommendations, or "you should" language alongside specifications.

Fix: Strip all advisory language from reference. Move rationale to an explanation page.

## Quality Match Ratings

When classifying, rate how well a page matches its type:

- **Strong** — Page cleanly matches one type with no significant signals from other types
- **Moderate** — Page mostly matches one type but has minor contamination (1-2 sentences of wrong-type content)
- **Weak** — Page nominally fits a type but has significant structural issues (missing checkpoints in tutorial, missing troubleshooting in how-to)
- **Mixed** — Page contains substantial content from two or more types and should be split
