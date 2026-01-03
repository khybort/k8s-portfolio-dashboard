# Deployment Guide

This guide covers deployment strategies for the Portfolio Dashboard, including Kubernetes deployment, Docker Compose, and production considerations.

## Quick Start

### Start Project
```bash
./scripts/start.sh
```

This will:
- Create namespace if needed
- Deploy all resources
- Add domains to `/etc/hosts`
- Start port-forwards
- Deploy ingress controller if needed

### Stop Project (Keep Resources)
```bash
./scripts/stop.sh
```

This will:
- Stop all port-forwards
- Remove domains from `/etc/hosts`
- Delete ingress resource (namespace-specific)

### Complete Cleanup
```bash
./scripts/cleanup.sh
```

This will:
- Stop all port-forwards
- Remove domains from `/etc/hosts`
- Delete all resources in portfolio namespace
- Delete portfolio namespace
- **Note**: Ingress controller is kept (may be used by other projects)

## Manual Deployment

### 1. Deploy Infrastructure
```bash
kubectl apply -k k8s/
```

### 2. Configure Local DNS
Add to `/etc/hosts`:
```
127.0.0.1 portfolio.local admin.portfolio.local api.portfolio.local auth.portfolio.local
```

### 3. Start Port-Forwards
```bash
kubectl port-forward -n portfolio svc/frontend 5173:80 &
kubectl port-forward -n portfolio svc/backend 8080:80 &
kubectl port-forward -n portfolio svc/auth-service 8081:80 &
```

### 4. Deploy Ingress Controller (Optional)
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
```

### 5. Start Ingress Port-Forward
```bash
kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 8888:80 &
```

## Access Points

### Port-Forward (Recommended)
- Public: http://localhost:5173
- Admin: http://localhost:5173/admin/login
- Backend API: http://localhost:8080
- Auth API: http://localhost:8081

### Ingress (Requires /etc/hosts)
- Public: http://portfolio.local (or http://localhost:8888 with Host header)
- Admin: http://admin.portfolio.local (or http://localhost:8888 with Host header)
- Backend API: http://api.portfolio.local
- Auth API: http://auth.portfolio.local

## Resource Management

### Project-Specific Resources (Namespace: portfolio)
- All deployments, services, ingress, jobs, secrets, configmaps
- These are created/removed with the project

### Shared Resources
- Ingress controller (namespace: ingress-nginx)
- Not removed by stop/cleanup scripts (may be used by other projects)
- To remove manually: `kubectl delete namespace ingress-nginx`

## Troubleshooting

### Port Already in Use
```bash
# Find process using port
lsof -ti:5173

# Kill process
kill -9 $(lsof -ti:5173)
```

### /etc/hosts Permission Denied
```bash
sudo sh -c 'echo "127.0.0.1 portfolio.local" >> /etc/hosts'
```

### Ingress Not Working
1. Check ingress controller is running:
   ```bash
   kubectl get pods -n ingress-nginx
   ```

2. Check ingress resource:
   ```bash
   kubectl get ingress -n portfolio
   ```

3. Test with curl:
   ```bash
   curl -H "Host: portfolio.local" http://localhost:8888
   ```
