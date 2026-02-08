# Reference Template

Use this template when creating reference documentation. Read this before writing.

## Frontmatter

```yaml
---
title: "[API/CLI/Config] Reference"
summary: "Complete specification of [thing] including parameters, return values, errors, and limits."
doc_type: reference
stability: stable
version: "2.1.0"
---
```

## Page Structures

### API Endpoint Reference

```markdown
# [Resource] API

## [METHOD] /path/to/endpoint

[One-sentence description of what this endpoint does — no advice.]

### Request

**Headers:**

| Header | Required | Description |
|--------|----------|-------------|
| `Authorization` | Yes | Bearer token |
| `Content-Type` | Yes | `application/json` |

**Parameters:**

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `name` | string | Yes | — | The resource name. 3-50 characters, alphanumeric and hyphens. |
| `limit` | integer | No | 20 | Maximum items to return. Range: 1-100. |

**Request body:**

\`\`\`json
{
  "name": "example",
  "config": {
    "key": "value"
  }
}
\`\`\`

### Response

**Success (200):**

\`\`\`json
{
  "id": "abc-123",
  "name": "example",
  "created_at": "2024-01-15T09:30:00Z"
}
\`\`\`

**Response fields:**

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier. Format: UUID v4. |
| `name` | string | The resource name as provided. |
| `created_at` | string | ISO 8601 timestamp of creation. |

### Errors

| Code | HTTP Status | Condition |
|------|-------------|-----------|
| `INVALID_NAME` | 400 | Name contains invalid characters or exceeds length limit. |
| `DUPLICATE_NAME` | 409 | A resource with this name already exists. |
| `RATE_LIMITED` | 429 | More than 100 requests per minute. Retry after `Retry-After` header. |

### Limits

- Request body: 1 MB maximum
- Rate limit: 100 requests/minute per API key
- Name length: 3-50 characters

### Example

\`\`\`bash
curl -X POST https://api.example.com/resources \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "my-resource"}'
\`\`\`

### See Also

- [How to create resources](../how-to/create-resources.md)
- [Understanding the resource lifecycle](../explanation/resource-lifecycle.md)
```

### CLI Command Reference

```markdown
# [command-name]

[One sentence: what this command does.]

## Synopsis

\`\`\`
command-name [flags] <required-arg> [optional-arg]
\`\`\`

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `<required-arg>` | Yes | The target to operate on. |
| `[optional-arg]` | No | Additional parameter. Default: current directory. |

## Flags

| Flag | Short | Type | Default | Description |
|------|-------|------|---------|-------------|
| `--output` | `-o` | string | `text` | Output format. Values: `text`, `json`, `yaml`. |
| `--verbose` | `-v` | bool | `false` | Enable detailed output. |
| `--force` | `-f` | bool | `false` | Skip confirmation prompts. |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Resource not found |

## Examples

\`\`\`bash
# Basic usage
command-name my-target

# With JSON output
command-name --output json my-target
\`\`\`

## See Also

- [`related-command`](./related-command.md)
- [How to accomplish task using this command](../how-to/task.md)
```

### Configuration Reference

```markdown
# Configuration Reference

## File location

`~/.config/tool/config.yaml` or `$TOOL_CONFIG` environment variable.

## Options

### `option_name`

| | |
|---|---|
| **Type** | string |
| **Default** | `"default_value"` |
| **Required** | No |
| **Environment** | `TOOL_OPTION_NAME` |
| **Since** | v1.2.0 |

Description of what this option controls.

Values: `"value1"` | `"value2"` | `"value3"`
```

## The "No Advice" Rule

This is the most important rule in reference documentation. Reference states facts. It never advises.

### Examples of violations and fixes

| Violation | Fix |
|-----------|-----|
| "You should validate inputs before calling this function." | "Input validation is the caller's responsibility." |
| "We recommend using the async version for better performance." | "An async version is available as `sendEmailAsync`." |
| "It's best practice to set a timeout." | "Accepts an optional `timeout` parameter (default: 30s)." |
| "Don't forget to close the connection." | "The connection must be explicitly closed via `close()`." |
| "For most use cases, the default is fine." | "Default: `auto`. Other values: `manual`, `disabled`." |

### Where advice belongs

If you want to say "you should" — write a how-to guide for that task. If you want to explain why — write an explanation page. Reference only states what exists and how it behaves.

## Completeness Checklist

Reference must document EVERY public surface:

- [ ] All public functions/methods with full signatures
- [ ] All parameters with types, defaults, and constraints
- [ ] All return values with types and possible values
- [ ] All error codes with conditions and HTTP status codes
- [ ] All configuration options with types, defaults, and valid ranges
- [ ] All CLI commands with all flags and arguments
- [ ] All environment variables
- [ ] Rate limits, size limits, and other constraints
- [ ] Deprecation notices with migration paths
- [ ] Version/stability markers
