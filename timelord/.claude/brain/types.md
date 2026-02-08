# Type Definitions

This document lists all exported types with their fields and methods.

## Package: cmd

### NamespaceInfo (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| Name | string | Yes | - |
| State | string | Yes | - |
| Retention | string | Yes | - |
| Description | string | Yes | - |

---

### NamespaceListResult (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| Namespaces | []NamespaceInfo | Yes | - |
| Count | int | Yes | - |

---

### NamespaceCreateResult (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| Name | string | Yes | - |
| Retention | string | Yes | - |
| Message | string | Yes | - |

---

### ScaffoldResult (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| Type | string | Yes | - |
| Name | string | Yes | - |
| Files | []string | Yes | - |
| Message | string | Yes | - |

---

### TemplateData (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| Name | string | Yes | - |
| PackageName | string | Yes | - |
| TaskQueue | string | Yes | - |

---

### ReplayResult (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| Success | bool | Yes | - |
| HistoryFile | string | Yes | - |
| WorkflowName | string | Yes | - |
| EventCount | int | Yes | - |
| Error | string | Yes | - |
| Suggestions | []string | Yes | - |

---

### ValidationResult (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| Valid | bool | Yes | - |
| Path | string | Yes | - |
| Errors | []ValidationError | Yes | - |
| Warnings | []ValidationError | Yes | - |

---

### ValidationError (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| File | string | Yes | - |
| Line | int | Yes | - |
| Type | string | Yes | - |
| Description | string | Yes | - |
| Suggestion | string | Yes | - |

---

### WorkflowListResult (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| Workflows | []WorkflowSummary | Yes | - |
| Count | int | Yes | - |

---

### WorkflowSummary (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| WorkflowID | string | Yes | - |
| RunID | string | Yes | - |
| Type | string | Yes | - |
| Status | string | Yes | - |
| StartTime | string | Yes | - |
| CloseTime | string | Yes | - |

---

### WorkflowDescription (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| WorkflowID | string | Yes | - |
| RunID | string | Yes | - |
| Type | string | Yes | - |
| Status | string | Yes | - |
| StartTime | string | Yes | - |
| CloseTime | string | Yes | - |
| HistoryLength | int | Yes | - |
| TaskQueue | string | Yes | - |
| Memo | map[string]interface{} | Yes | - |
| SearchAttributes | map[string]interface{} | Yes | - |

---

### DiagnosisResult (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| WorkflowID | string | Yes | - |
| Status | string | Yes | - |
| Issues | []DiagnosisIssue | Yes | - |
| Suggestions | []string | Yes | - |
| LastEvent | string | Yes | - |
| PendingItems | []string | Yes | - |

---

### DiagnosisIssue (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| Severity | string | Yes | - |
| Type | string | Yes | - |
| Description | string | Yes | - |

---

### ClusterStatus (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| Healthy | bool | Yes | - |
| Status | string | Yes | - |
| Address | string | Yes | - |
| Message | string | Yes | - |
| Error | string | Yes | - |

---

### ClusterInfo (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| ServerVersion | string | Yes | - |
| ClusterID | string | Yes | - |
| ClusterName | string | Yes | - |
| HistoryShardCount | int | Yes | - |
| PersistenceStore | string | Yes | - |
| VisibilityStore | string | Yes | - |
| Features | map[string]bool | Yes | - |

---

### ClusterMetrics (struct)

#### Fields

| Name | Type | Exported | Description |
| --- | --- | --- | --- |
| ActiveNamespaces | int | Yes | - |
| WorkflowsRunning | int | Yes | - |
| TaskQueuesActive | int | Yes | - |
| FrontendLatencyP99 | float64 | Yes | - |
| HistoryLatencyP99 | float64 | Yes | - |

---

