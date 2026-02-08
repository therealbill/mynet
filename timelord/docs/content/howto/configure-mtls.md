---
title: "Configure mTLS for Production"
weight: 2
---

# Configure mTLS for Production

Enable mutual TLS (mTLS) for secure communication between Temporal server and workers.

## Problem

You need to secure communication between your workers and the Temporal cluster for production deployment.

## Solution

Configure mTLS with certificates for both server and client authentication.

## Prerequisites

- Running Temporal cluster
- Certificate authority (CA) or ability to create one
- OpenSSL or similar tool
- kubectl access to cluster

## Steps

### 1. Generate Certificates

Create a CA and certificates for server and clients.

**Create CA:**

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

# Create server CSR
openssl req -new -key server.key \
  -out server.csr \
  -subj "/CN=temporal.example.com/O=YourOrg"

# Sign server certificate
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

# Create client CSR
openssl req -new -key client.key \
  -out client.csr \
  -subj "/CN=temporal-worker/O=YourOrg"

# Sign client certificate
openssl x509 -req -days 365 \
  -in client.csr \
  -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out client.crt
```

### 2. Create Kubernetes Secrets

```bash
# Server certificates
kubectl create secret generic temporal-server-tls \
  --namespace temporal \
  --from-file=ca.crt=ca.crt \
  --from-file=tls.crt=server.crt \
  --from-file=tls.key=server.key

# Client certificates (for workers)
kubectl create secret generic temporal-client-tls \
  --namespace your-app \
  --from-file=ca.crt=ca.crt \
  --from-file=tls.crt=client.crt \
  --from-file=tls.key=client.key
```

### 3. Configure Temporal Server

Update Helm values:

```yaml
# values-mtls.yaml
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
        client:
          serverName: temporal.example.com
          rootCaFiles:
            - /etc/temporal/tls/ca.crt
      internode:
        server:
          certFile: /etc/temporal/tls/tls.crt
          keyFile: /etc/temporal/tls/tls.key
          requireClientAuth: true
          clientCaFiles:
            - /etc/temporal/tls/ca.crt
        client:
          rootCaFiles:
            - /etc/temporal/tls/ca.crt

  # Mount the secret
  additionalVolumes:
    - name: tls-certs
      secret:
        secretName: temporal-server-tls

  additionalVolumeMounts:
    - name: tls-certs
      mountPath: /etc/temporal/tls
      readOnly: true
```

Upgrade the deployment:

```bash
helm upgrade temporal temporal/temporal \
  --namespace temporal \
  -f values-production.yaml \
  -f values-mtls.yaml
```

### 4. Configure Workers

Update your Go worker to use mTLS:

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "log"
    "os"

    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
)

func main() {
    // Load client certificate
    cert, err := tls.LoadX509KeyPair(
        "/etc/temporal/tls/tls.crt",
        "/etc/temporal/tls/tls.key",
    )
    if err != nil {
        log.Fatalf("Failed to load client cert: %v", err)
    }

    // Load CA certificate
    caCert, err := os.ReadFile("/etc/temporal/tls/ca.crt")
    if err != nil {
        log.Fatalf("Failed to load CA cert: %v", err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // Create TLS config
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
        ServerName:   "temporal.example.com",
    }

    // Create Temporal client with TLS
    c, err := client.Dial(client.Options{
        HostPort: "temporal-frontend:7233",
        ConnectionOptions: client.ConnectionOptions{
            TLS: tlsConfig,
        },
    })
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer c.Close()

    // Create and run worker
    w := worker.New(c, "your-queue", worker.Options{})
    // Register workflows and activities...
    w.Run(worker.InterruptCh())
}
```

### 5. Update Worker Deployment

Mount the client certificates:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: temporal-worker
spec:
  template:
    spec:
      containers:
      - name: worker
        volumeMounts:
        - name: tls-certs
          mountPath: /etc/temporal/tls
          readOnly: true
      volumes:
      - name: tls-certs
        secret:
          secretName: temporal-client-tls
```

### 6. Configure CLI (Optional)

For CLI access with mTLS:

```bash
export TEMPORAL_ADDRESS=temporal.example.com:7233
export TEMPORAL_TLS_CERT=/path/to/client.crt
export TEMPORAL_TLS_KEY=/path/to/client.key
export TEMPORAL_TLS_CA=/path/to/ca.crt
export TEMPORAL_TLS_SERVER_NAME=temporal.example.com

temporal operator cluster health
```

## Verification

- [ ] Workers connect successfully
- [ ] CLI works with certificates
- [ ] Non-TLS connections rejected
- [ ] Certificate expiration monitored

## Testing

```bash
# Test TLS connection
openssl s_client -connect temporal-frontend:7233 \
  -cert client.crt \
  -key client.key \
  -CAfile ca.crt

# Verify rejection without cert
temporal operator cluster health  # Should fail without TLS config
```

## Certificate Rotation

Plan for certificate renewal before expiration:

1. Generate new certificates
2. Update Kubernetes secrets
3. Rolling restart of Temporal server
4. Update and restart workers
5. Verify connectivity

## Troubleshooting

**Connection refused:**

- Verify certificates are mounted
- Check certificate CN matches server name
- Verify CA certificate is correct

**Certificate verification failed:**

- Check certificate chain
- Verify SAN includes server hostname
- Check certificate expiration

**Internode communication fails:**

- Ensure internode TLS is configured
- Verify all server pods have certificates

## Related

- `security-config` skill for authorization
- Production setup tutorial
