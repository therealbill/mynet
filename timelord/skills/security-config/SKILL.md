---
name: security-config
description: This skill should be used when the user asks about "temporal mTLS", "temporal authorization", "temporal security", "secure temporal", "TLS temporal", "certificate configuration", "RBAC temporal", or needs guidance on securing Temporal clusters.
version: 1.0.0
---

# Temporal Security Configuration

Guidance for securing Temporal clusters with mTLS, authorization, and network policies.

## Security Layers

Temporal security involves multiple layers:

1. **Transport Security (mTLS)** - Encrypted, authenticated connections
2. **Authorization** - Role-based access control
3. **Namespace Isolation** - Multi-tenant separation
4. **Network Policies** - Infrastructure-level controls

## mTLS Configuration

### Certificate Requirements

Temporal requires certificates for:

| Component | Certificate Type | Purpose |
|-----------|-----------------|---------|
| Frontend | Server cert | Client connections |
| Internode | Server + Client | Service mesh |
| Worker | Client cert | Connect to frontend |
| CLI/SDK | Client cert | API access |

### Certificate Generation

**Create Certificate Authority:**

```bash
# Generate CA key
openssl genrsa -out ca.key 4096

# Generate CA certificate
openssl req -new -x509 -days 3650 -key ca.key \
  -out ca.crt \
  -subj "/CN=Temporal CA/O=YourOrg"
```

**Create Server Certificate:**

```bash
# Generate server key
openssl genrsa -out server.key 4096

# Create CSR
openssl req -new -key server.key \
  -out server.csr \
  -subj "/CN=temporal.example.com/O=YourOrg"

# Sign with CA (include SANs)
openssl x509 -req -days 365 \
  -in server.csr \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out server.crt \
  -extfile <(printf "subjectAltName=DNS:temporal.example.com,DNS:temporal-frontend,DNS:localhost")
```

**Create Client Certificate:**

```bash
# Generate client key
openssl genrsa -out client.key 4096

# Create CSR
openssl req -new -key client.key \
  -out client.csr \
  -subj "/CN=temporal-worker/O=YourOrg"

# Sign with CA
openssl x509 -req -days 365 \
  -in client.csr \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out client.crt
```

### Server Configuration

**Helm values for mTLS:**

```yaml
server:
  config:
    tls:
      frontend:
        server:
          certFile: /etc/temporal/tls/tls.crt
          keyFile: /etc/temporal/tls/tls.key
          requireClientAuth: true
          clientCaFiles:
            - /etc/temporal/tls/ca.crt
      internode:
        server:
          certFile: /etc/temporal/tls/tls.crt
          keyFile: /etc/temporal/tls/tls.key
          requireClientAuth: true
          clientCaFiles:
            - /etc/temporal/tls/ca.crt
```

### Worker Configuration

**Go SDK with mTLS:**

```go
import (
    "crypto/tls"
    "crypto/x509"
    "os"

    "go.temporal.io/sdk/client"
)

func createSecureClient() (client.Client, error) {
    cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
    if err != nil {
        return nil, err
    }

    caCert, err := os.ReadFile("ca.crt")
    if err != nil {
        return nil, err
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    return client.Dial(client.Options{
        HostPort: "temporal.example.com:7233",
        ConnectionOptions: client.ConnectionOptions{
            TLS: &tls.Config{
                Certificates: []tls.Certificate{cert},
                RootCAs:      caCertPool,
                ServerName:   "temporal.example.com",
            },
        },
    })
}
```

## Authorization

### Authorizer Plugin

Temporal supports pluggable authorization:

```go
type Authorizer interface {
    Authorize(ctx context.Context, claims *Claims, target *CallTarget) (Result, error)
}
```

### Claim Mapper

Map certificates to claims:

```go
type ClaimMapper interface {
    GetClaims(authInfo *AuthInfo) (*Claims, error)
}
```

### Common Authorization Patterns

| Pattern | Description |
|---------|-------------|
| Namespace-based | Users access specific namespaces |
| Role-based | Admin, operator, developer roles |
| Team-based | Team membership determines access |

### Configuration Example

```yaml
server:
  config:
    authorization:
      authorizer: default
      claimMapper: default
      permissionsClaimName: permissions
```

## Namespace Isolation

### Per-Namespace Security

Configure security per namespace:

```bash
# Create namespace with specific permissions
temporal operator namespace create \
  --namespace secure-ns \
  --retention 72h
```

### Multi-Tenant Patterns

| Pattern | Isolation Level | Use Case |
|---------|----------------|----------|
| Shared cluster | Namespace | Cost-effective |
| Cluster per tenant | Full | Maximum isolation |
| Hybrid | Mixed | Balance |

## Network Policies

### Kubernetes NetworkPolicy

**Restrict frontend access:**

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: temporal-frontend
  namespace: temporal
spec:
  podSelector:
    matchLabels:
      app: temporal-frontend
  policyTypes:
    - Ingress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              temporal-access: "true"
      ports:
        - port: 7233
```

**Internode communication:**

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: temporal-internode
  namespace: temporal
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: temporal
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app.kubernetes.io/name: temporal
```

## Security Checklist

### Production Readiness

- [ ] mTLS enabled for frontend
- [ ] mTLS enabled for internode
- [ ] Client certificates distributed securely
- [ ] Certificate rotation plan in place
- [ ] Authorization configured
- [ ] Network policies applied
- [ ] Secrets management (Vault, etc.)
- [ ] Audit logging enabled

### Certificate Management

- [ ] CA stored securely
- [ ] Certificates in Kubernetes secrets
- [ ] Expiration monitoring
- [ ] Rotation procedures documented

## Troubleshooting

### Connection Refused

- Verify certificates are mounted
- Check certificate CN matches server name
- Verify CA certificate is correct

### Certificate Verification Failed

- Check certificate chain
- Verify SAN includes hostname
- Check certificate expiration

### Authorization Denied

- Verify claims are mapped correctly
- Check namespace permissions
- Review authorizer logs

## Additional Resources

### Reference Files

For detailed security patterns, consult:

- **`references/certificate-rotation.md`** - Certificate lifecycle management
- **`references/authorization-patterns.md`** - RBAC implementation examples
