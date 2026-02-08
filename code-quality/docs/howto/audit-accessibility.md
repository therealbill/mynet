---
title: "Audit Web Accessibility"
description: "Use web-accessibility-checker for WCAG compliance audits"
weight: 3
---

# Audit Web Accessibility

Find and fix WCAG 2.2 accessibility violations in your web content using the web-accessibility-checker agent.

## Prerequisites

- Claude Code with the code-quality plugin installed
- A web project with HTML templates, JSX components, or rendered pages to audit

## Steps

### 1. Point web-accessibility-checker at Your Content

Activate the agent by describing what you want audited:

```
Audit this page for accessibility
```

```
Make sure this form is accessible
```

```
We need to be WCAG AA compliant before launch
```

For targeted audits, specify the files or components:

```
Check the checkout form in src/components/CheckoutForm.tsx for accessibility
```

web-accessibility-checker reads the specified files and evaluates them against WCAG 2.2 success criteria. It examines semantic HTML structure, keyboard navigation patterns, screen reader compatibility, color contrast, and form labeling.

### 2. Review Findings by WCAG Success Criterion

The agent reports each issue with:

- **WCAG success criterion** -- the specific rule violated (e.g., "1.1.1 Non-text Content", "2.1.1 Keyboard")
- **Level** -- A, AA, or AAA
- **File and line number** -- exact location in your source code
- **Description** -- what is wrong and why it matters for users
- **Fix** -- a code example showing the corrected markup

Findings are organized by WCAG level, with Level A violations listed first.

### 3. Fix Level A Violations First

Level A is the minimum accessibility standard. These violations affect the most users and represent the most severe barriers. Common Level A findings:

- Images missing `alt` text (1.1.1 Non-text Content)
- Non-semantic interactive elements -- `div` with click handler instead of `button` (4.1.2 Name, Role, Value)
- Missing form labels (1.3.1 Info and Relationships)
- Keyboard-inaccessible controls (2.1.1 Keyboard)

Fix all Level A issues before moving on.

### 4. Address Level AA Violations

Level AA is the standard compliance target for most organizations and legal requirements. Common Level AA findings:

- Insufficient color contrast -- normal text below 4.5:1 ratio, large text below 3:1 (1.4.3 Contrast Minimum)
- Missing focus indicators -- `outline: none` without a visible replacement (2.4.7 Focus Visible)
- Error messages not associated with their form fields via `aria-describedby` (3.3.1 Error Identification)
- Dynamic content updates not announced to screen readers via `aria-live` regions (4.1.3 Status Messages)

### 5. Consider Level AAA

Level AAA is the highest conformance level. Only pursue these fixes if your project specifically requires AAA compliance. Common Level AAA criteria include enhanced contrast ratios (7:1 for normal text) and extended audio descriptions. Most projects target AA.

### 6. Verification

After applying fixes, re-run the audit:

```
Re-audit the checkout form for accessibility
```

A clean audit confirms that the previously reported violations are resolved. New findings may appear if fixes introduced different markup patterns that need their own accessibility treatment.

For ongoing compliance, run accessibility audits whenever you modify UI components, add new pages, or change form structures.

## Troubleshooting

**False positives on dynamic content:**

web-accessibility-checker performs static analysis of your source code. It may flag dynamic content (modals, dropdowns, toast notifications) that is correctly managed at runtime through JavaScript. When a finding looks incorrect, verify the runtime behavior manually: tab through the component with a keyboard, test with a screen reader, and check that focus management works as intended. If the runtime behavior is correct, the finding is a false positive and can be noted as such.

**Components that conditionally render accessible attributes:**

If your components add `aria-label`, `role`, or other accessibility attributes conditionally based on props, web-accessibility-checker may flag the base case as non-compliant. Ensure that every code path produces accessible output, including default states and error states.

**Third-party component libraries:**

web-accessibility-checker audits your source code, not compiled output. If a third-party component renders accessible HTML at runtime but the JSX source does not contain explicit ARIA attributes, the agent may produce findings that do not reflect the actual rendered output. Check the component library's accessibility documentation before modifying wrapper code.

## See Also

- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- tutorial walkthrough of the code quality workflow
- [Agent Reference]({{< ref "reference/agents" >}}) -- web-accessibility-checker specification
- [Choosing the Right Agent]({{< ref "explanation/choosing-the-right-agent" >}}) -- when to use web-accessibility-checker versus other agents
