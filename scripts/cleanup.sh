#!/usr/bin/env sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧹 Portfolio Dashboard Tamamen Temizleniyor"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "⚠️  Bu işlem tüm portfolio namespace'ini ve kaynaklarını silecek!"
echo ""

read -p "Devam etmek istediğinize emin misiniz? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "❌ İşlem iptal edildi"
    exit 1
fi

cd "$PROJECT_DIR"

# Stop all port-forwards
echo "🔌 Port-forward'lar durduruluyor..."
pkill -f "kubectl port-forward.*portfolio" 2>/dev/null || true
pkill -f "kubectl port-forward.*ingress-nginx-controller" 2>/dev/null || true
sleep 1

# Remove domains from /etc/hosts
echo "🌐 /etc/hosts temizleniyor..."
if grep -q "portfolio.local" /etc/hosts 2>/dev/null; then
    sudo sed -i '' '/portfolio.local/d' /etc/hosts 2>/dev/null || \
    sudo sed -i '/portfolio.local/d' /etc/hosts 2>/dev/null || \
    echo "⚠️  /etc/hosts manuel olarak temizlenmeli"
    echo "✅ Domain'ler /etc/hosts'tan kaldırıldı"
fi

# Delete ingress resource
echo "🗑️  Ingress resource siliniyor..."
kubectl delete ingress portfolio-ingress -n portfolio 2>/dev/null || true

# Delete all resources in portfolio namespace
echo "🗑️  Portfolio namespace kaynakları siliniyor..."
kubectl delete all --all -n portfolio 2>/dev/null || true
kubectl delete ingress --all -n portfolio 2>/dev/null || true
kubectl delete job --all -n portfolio 2>/dev/null || true
kubectl delete secret --all -n portfolio 2>/dev/null || true
kubectl delete configmap --all -n portfolio 2>/dev/null || true

# Delete namespace
echo "🗑️  Portfolio namespace siliniyor..."
kubectl delete namespace portfolio 2>/dev/null || true

# Note: We keep ingress controller as it might be used by other projects
# To remove ingress controller completely, uncomment the following:
# echo "🗑️  Ingress controller siliniyor..."
# kubectl delete namespace ingress-nginx 2>/dev/null || true

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Temizlik Tamamlandı"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "ℹ️  Notlar:"
echo "   - Ingress controller hala mevcut (diğer projeler için)"
echo "   - Ingress controller'ı da silmek için:"
echo "     kubectl delete namespace ingress-nginx"
echo ""

