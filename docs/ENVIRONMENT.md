# Environment Configuration Guide

This guide explains how to configure environment variables for different deployment environments (development, production).

## Overview

The project supports multiple environments:
- **Development**: Local development with Docker Compose
- **Production**: Kubernetes deployment

Environment variables are managed through:
- `.env.dev` - Development environment (Docker Compose)
- `.env.prod` - Production environment (Kubernetes)
- ConfigMaps - Non-sensitive configuration in Kubernetes
- Secrets - Sensitive data in Kubernetes

## Environment Files

### Development (.env.dev)

Used for local development with Docker Compose.

**Location**: Project root directory

**Key Features**:
- Local service names (postgresql, redis, kafka)
- Development-friendly defaults
- Debug logging enabled
- SSL disabled

**Example**:
```bash
ENV=development
SERVER_PORT=8080
DB_HOST=postgresql
DB_PASSWORD=password
LOG_LEVEL=debug
```

### Production (.env.prod)

Used for production Kubernetes deployment.

**Location**: Project root directory (not committed to git)

**Key Features**:
- Kubernetes service DNS names
- Production-ready defaults
- SSL enabled
- Info logging

**Example**:
```bash
ENV=production
DB_HOST=postgresql.portfolio.svc.cluster.local
DB_SSLMODE=require
LOG_LEVEL=info
```

## Environment Variable Reference

### Backend Service

| Variable | Description | Default (Dev) | Default (Prod) |
|----------|-------------|---------------|----------------|
| `ENV` | Environment name | `development` | `production` |
| `SERVER_PORT` | Server port | `8080` | `8080` |
| `SERVER_HOST` | Server host | `0.0.0.0` | `0.0.0.0` |
| `LOG_LEVEL` | Log level | `debug` | `info` |
| `DB_HOST` | Database host | `postgresql` | `postgresql.portfolio.svc.cluster.local` |
| `DB_PORT` | Database port | `5432` | `5432` |
| `DB_USER` | Database user | `portfolio` | `portfolio` |
| `DB_PASSWORD` | Database password | `password` | From Secret |
| `DB_NAME` | Database name | `portfolio` | `portfolio` |
| `DB_SSLMODE` | SSL mode | `disable` | `require` |
| `REDIS_HOST` | Redis host | `redis` | `redis.portfolio.svc.cluster.local` |
| `REDIS_PORT` | Redis port | `6379` | `6379` |
| `REDIS_PASSWORD` | Redis password | `` | From Secret |
| `REDIS_DB` | Redis database | `0` | `0` |
| `KAFKA_BROKERS` | Kafka brokers | `kafka:9092` | `kafka.portfolio.svc.cluster.local:9092` |
| `AUTH_SERVICE_URL` | Auth service URL | `http://auth-service:8081` | `http://auth-service:80` |
| `JWT_SECRET` | JWT secret key | `dev-secret-key` | From Secret |

### Auth Service

| Variable | Description | Default (Dev) | Default (Prod) |
|----------|-------------|---------------|----------------|
| `ENV` | Environment name | `development` | `production` |
| `AUTH_SERVICE_PORT` | Service port | `8081` | `8081` |
| `AUTH_SERVICE_HOST` | Service host | `0.0.0.0` | `0.0.0.0` |
| `AUTH_DB_NAME` | Database name | `auth_db` | `auth_db` |
| `AUTH_DB_HOST` | Database host | `postgresql-auth` | `postgresql.portfolio.svc.cluster.local` |
| `AUTH_DB_PORT` | Database port | `5432` | `5432` |
| `AUTH_DB_USER` | Database user | `portfolio` | `portfolio` |
| `AUTH_DB_PASSWORD` | Database password | `password` | From Secret |
| `AUTH_DB_SSLMODE` | SSL mode | `disable` | `require` |
| `JWT_SECRET` | JWT secret key | `dev-secret-key` | From Secret |
| `JWT_ACCESS_EXPIRY` | Access token expiry | `15m` | `15m` |
| `JWT_REFRESH_EXPIRY` | Refresh token expiry | `168h` | `168h` |
| `LOG_LEVEL` | Log level | `debug` | `info` |

### Frontend Service

| Variable | Description | Default (Dev) | Default (Prod) |
|----------|-------------|---------------|----------------|
| `ENV` | Environment name | `development` | `production` |
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` | `https://api.portfolio.local` |
| `VITE_AUTH_URL` | Auth service URL | `http://localhost:8081` | `https://auth.portfolio.local` |

**Note**: Frontend environment variables must be set at build time (Vite requirement).

## Configuration Methods

### 1. Docker Compose (Development)

Uses `.env.dev` file:

```yaml
services:
  backend:
    env_file:
      - .env.dev
    environment:
      ENV: development
```

**Setup**:
```bash
# Copy example file
cp .env.example .env.dev

# Edit .env.dev with your values
# Then start services
docker-compose up -d
```

### 2. Kubernetes (Production)

Uses ConfigMap and Secrets:

#### ConfigMap

Non-sensitive configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: portfolio-config
data:
  ENV: "production"
  SERVER_PORT: "8080"
  DB_HOST: "postgresql"
  # ... other non-sensitive vars
```

**Apply**:
```bash
kubectl apply -f k8s/configmap.yaml
```

#### Secrets

Sensitive data (passwords, keys):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: portfolio-secrets
type: Opaque
stringData:
  db-password: "your-secure-password"
  jwt-secret: "your-jwt-secret"
```

**Create from command line**:
```bash
kubectl create secret generic portfolio-secrets \
  --from-literal=db-password='your-password' \
  --from-literal=auth-db-password='your-password' \
  --from-literal=jwt-secret='your-jwt-secret' \
  --from-literal=redis-password='your-redis-password' \
  -n portfolio
```

**Create from file**:
```bash
kubectl create secret generic portfolio-secrets \
  --from-env-file=.env.prod \
  -n portfolio
```

#### Deployment Usage

```yaml
spec:
  containers:
  - name: backend
    envFrom:
    - configMapRef:
        name: portfolio-config
    env:
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: portfolio-secrets
          key: db-password
```

## Environment-Specific Configurations

### Development

**Characteristics**:
- Local service discovery
- Debug logging
- SSL disabled
- Development secrets
- Hot reload enabled

**Service URLs**:
- Backend: `http://localhost:8080`
- Auth: `http://localhost:8081`
- Frontend: `http://localhost:5173`

### Production

**Characteristics**:
- Kubernetes service discovery
- Info logging
- SSL enabled
- Secure secrets
- Optimized builds

**Service URLs**:
- Backend: `http://backend.portfolio.svc.cluster.local:80`
- Auth: `http://auth-service.portfolio.svc.cluster.local:80`
- Frontend: `http://frontend.portfolio.svc.cluster.local:80`

## Configuration Loading Order

1. **Environment variables** (highest priority)
2. **.env files** (`.env.dev` or `.env.prod` based on `ENV`)
3. **Default values** (lowest priority)

## Best Practices

### Security

1. **Never commit secrets**:
   - Add `.env.dev`, `.env.prod` to `.gitignore`
   - Use Kubernetes Secrets for production
   - Rotate secrets regularly

2. **Use strong secrets**:
   ```bash
   # Generate strong JWT secret
   openssl rand -base64 32
   
   # Generate database password
   openssl rand -base64 24
   ```

3. **Separate environments**:
   - Different secrets for dev/prod
   - Different databases
   - Different service accounts

### Organization

1. **Group related variables**:
   ```bash
   # Database
   DB_HOST=...
   DB_PORT=...
   
   # Redis
   REDIS_HOST=...
   REDIS_PORT=...
   ```

2. **Use descriptive names**:
   - `AUTH_SERVICE_URL` not `AUTH_URL`
   - `DB_PASSWORD` not `PASSWORD`

3. **Document defaults**:
   - Include in `.env.example`
   - Document in this guide

### Validation

1. **Validate on startup**:
   - Check required variables
   - Validate formats
   - Log warnings for missing optional vars

2. **Fail fast**:
   - Exit if critical vars missing
   - Don't use insecure defaults in production

## Troubleshooting

### Variables Not Loading

**Check**:
1. File exists and is readable
2. `ENV` variable is set correctly
3. File format is correct (no spaces around `=`)
4. No syntax errors

**Debug**:
```bash
# Check if variables are loaded
docker-compose config

# Check Kubernetes ConfigMap
kubectl get configmap portfolio-config -n portfolio -o yaml

# Check Kubernetes Secrets
kubectl get secret portfolio-secrets -n portfolio -o yaml
```

### Wrong Environment

**Symptoms**:
- Connecting to wrong database
- Using dev secrets in production

**Fix**:
1. Verify `ENV` variable
2. Check which `.env` file is loaded
3. Verify ConfigMap/Secret in Kubernetes

### Frontend Variables Not Working

**Issue**: Vite requires build-time variables

**Solution**:
1. Set variables before build
2. Use Docker build args
3. Rebuild image after changing variables

```dockerfile
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
```

## Migration Guide

### From Hardcoded to Environment Variables

1. **Identify hardcoded values**
2. **Create environment variables**
3. **Update code to use `os.Getenv()`**
4. **Add to `.env.example`**
5. **Update deployment configs**
6. **Test in both environments**

### Updating Existing Deployments

1. **Update ConfigMap**:
   ```bash
   kubectl edit configmap portfolio-config -n portfolio
   ```

2. **Update Secrets**:
   ```bash
   kubectl create secret generic portfolio-secrets \
     --from-literal=key=value \
     --dry-run=client -o yaml | kubectl apply -f -
   ```

3. **Restart pods**:
   ```bash
   kubectl rollout restart deployment/backend -n portfolio
   ```

## Examples

### Development Setup

```bash
# 1. Copy example file
cp .env.example .env.dev

# 2. Edit .env.dev (optional, defaults work)
# 3. Start services
docker-compose up -d
```

### Production Setup

```bash
# 1. Create secrets
kubectl create secret generic portfolio-secrets \
  --from-literal=db-password='secure-password' \
  --from-literal=jwt-secret='secure-jwt-secret' \
  -n portfolio

# 2. Apply ConfigMap
kubectl apply -f k8s/configmap.yaml

# 3. Deploy services
kubectl apply -k k8s/
```

## Additional Resources

- [Kubernetes ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/)
- [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/)
- [Docker Compose Environment Variables](https://docs.docker.com/compose/environment-variables/)
- [Vite Environment Variables](https://vitejs.dev/guide/env-and-mode.html)

