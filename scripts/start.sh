#!/usr/bin/env sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 Portfolio Dashboard Başlatılıyor"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

cd "$PROJECT_DIR"

# Check if namespace exists
if ! kubectl get namespace portfolio > /dev/null 2>&1; then
    echo "📦 Namespace oluşturuluyor..."
    kubectl apply -f k8s/namespace.yaml
fi

# Deploy infrastructure (ingress will be created here)
echo "🔧 Infrastructure deploy ediliyor..."
kubectl apply -k k8s/

# Wait for services to be ready
echo "⏳ Servislerin hazır olması bekleniyor..."
kubectl wait --for=condition=ready pod -l app=frontend -n portfolio --timeout=120s || true
kubectl wait --for=condition=ready pod -l app=backend -n portfolio --timeout=120s || true
kubectl wait --for=condition=ready pod -l app=auth-service -n portfolio --timeout=120s || true

# Add domains to /etc/hosts if not already present
echo "🌐 /etc/hosts güncelleniyor..."
HOSTS_ENTRY="127.0.0.1 portfolio.local admin.portfolio.local api.portfolio.local auth.portfolio.local"
if ! grep -q "portfolio.local" /etc/hosts 2>/dev/null; then
    echo "$HOSTS_ENTRY" | sudo tee -a /etc/hosts > /dev/null
    echo "✅ Domain'ler /etc/hosts'a eklendi"
else
    echo "ℹ️  Domain'ler zaten /etc/hosts'ta mevcut"
fi

# Start port-forwards
echo "🔌 Port-forward'lar başlatılıyor..."
pkill -f "kubectl port-forward.*portfolio" 2>/dev/null || true
sleep 1

kubectl port-forward -n portfolio svc/frontend 5173:80 > /tmp/pf-frontend.log 2>&1 &
kubectl port-forward -n portfolio svc/backend 8080:80 > /tmp/pf-backend.log 2>&1 &
kubectl port-forward -n portfolio svc/auth-service 8081:80 > /tmp/pf-auth.log 2>&1 &

# Check if ingress controller exists, if not, deploy it
if ! kubectl get deployment ingress-nginx-controller -n ingress-nginx > /dev/null 2>&1; then
    echo "🌐 Ingress controller kuruluyor..."
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml > /dev/null 2>&1
    echo "⏳ Ingress controller'ın hazır olması bekleniyor..."
    sleep 15
fi

# Start ingress controller port-forward (port 8888 to avoid root requirement)
echo "🌐 Ingress controller port-forward başlatılıyor..."
pkill -f "kubectl port-forward.*ingress-nginx-controller" 2>/dev/null || true
sleep 1
kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 8888:80 > /tmp/pf-ingress.log 2>&1 &

sleep 3

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Portfolio Dashboard Başlatıldı"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📍 Erişim:"
echo "   Public: http://localhost:5173"
echo "   Admin: http://localhost:5173/admin/login"
echo "   Backend API: http://localhost:8080"
echo "   Auth API: http://localhost:8081"
echo ""
echo "🌐 Ingress (Port 8888):"
echo "   Public: http://portfolio.local (Host header ile: http://localhost:8888)"
echo "   Admin: http://admin.portfolio.local (Host header ile: http://localhost:8888)"
echo ""
echo "💡 Durdurmak için: ./scripts/stop.sh"
echo ""

