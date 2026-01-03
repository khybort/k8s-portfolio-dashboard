#!/usr/bin/env sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🛑 Portfolio Dashboard Durduruluyor"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

cd "$PROJECT_DIR"

# Stop all port-forwards
echo "🔌 Port-forward'lar durduruluyor..."
pkill -f "kubectl port-forward.*portfolio" 2>/dev/null || true
pkill -f "kubectl port-forward.*ingress-nginx-controller" 2>/dev/null || true
sleep 1
echo "✅ Port-forward'lar durduruldu"

# Remove domains from /etc/hosts
echo "🌐 /etc/hosts temizleniyor..."
if grep -q "portfolio.local" /etc/hosts 2>/dev/null; then
    sudo sed -i '' '/portfolio.local/d' /etc/hosts 2>/dev/null || \
    sudo sed -i '/portfolio.local/d' /etc/hosts 2>/dev/null || \
    echo "⚠️  /etc/hosts manuel olarak temizlenmeli"
    echo "✅ Domain'ler /etc/hosts'tan kaldırıldı"
else
    echo "ℹ️  Domain'ler /etc/hosts'ta zaten yok"
fi

# Delete ingress resource (namespace-specific, safe to delete)
echo "🗑️  Ingress resource siliniyor..."
kubectl delete ingress portfolio-ingress -n portfolio 2>/dev/null || true
echo "✅ Ingress resource silindi"

# Note: We don't delete the namespace or ingress controller
# as they might be used by other projects or require manual cleanup
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Portfolio Dashboard Durduruldu"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "ℹ️  Notlar:"
echo "   - Namespace ve deployment'lar hala mevcut (tekrar başlatmak için)"
echo "   - Ingress controller hala mevcut (diğer projeler için)"
echo "   - Tamamen temizlemek için: ./scripts/cleanup.sh"
echo ""

