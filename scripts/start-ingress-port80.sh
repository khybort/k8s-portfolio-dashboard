#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 Ingress Controller Port 80 Başlatıcı"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "⚠️  Port 80 root yetkisi gerektirir."
echo ""

# Check if port 80 is already in use
if lsof -ti:80 > /dev/null 2>&1; then
    echo "❌ Port 80 zaten kullanımda. Mevcut process:"
    lsof -ti:80 | xargs ps -p
    echo ""
    echo "Port 80'i kullanan process'i durdurmak için:"
    echo "  sudo lsof -ti:80 | xargs sudo kill -9"
    exit 1
fi

# Kill existing port-forwards
pkill -f "kubectl port-forward.*ingress-nginx-controller" 2>/dev/null

echo "✅ Ingress controller port 80'de başlatılıyor..."
echo "   (Root yetkisi gerekecek)"
echo ""

# Start port-forward with sudo
sudo kubectl port-forward -n ingress-nginx svc/ingress-nginx-controller 80:80 > /tmp/pf-ingress-80.log 2>&1 &
PID=$!

sleep 3

# Check if it's running
if ps -p $PID > /dev/null; then
    echo "✅ Ingress controller port 80'de çalışıyor (PID: $PID)"
    echo ""
    echo "📍 Artık şu adresler çalışmalı:"
    echo "   http://portfolio.local"
    echo "   http://admin.portfolio.local"
    echo "   http://api.portfolio.local"
    echo "   http://auth.portfolio.local"
    echo ""
    echo "💡 Durdurmak için: sudo kill $PID"
else
    echo "❌ Port-forward başlatılamadı. Log:"
    cat /tmp/pf-ingress-80.log
    exit 1
fi

