#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Ingress Test Script"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Test Ingress Controller via Port 8888
echo "1️⃣ Testing via Port-Forward (8888):"
curl -s -H "Host: portfolio.local" http://localhost:8888 | head -3
echo ""
echo ""

# Test via NodePort
echo "2️⃣ Testing via NodePort (30080):"
curl -s -H "Host: portfolio.local" http://localhost:30080 | head -3
echo ""
echo ""

# Test direct service
echo "3️⃣ Testing direct Frontend service:"
curl -s http://localhost:5173 | head -3
echo ""
echo ""

# Check /etc/hosts
echo "4️⃣ /etc/hosts kontrolü:"
grep portfolio.local /etc/hosts || echo "❌ Domain'ler /etc/hosts'ta yok"
echo ""

# Check ingress status
echo "5️⃣ Ingress durumu:"
kubectl get ingress portfolio-ingress -n portfolio
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "💡 Tarayıcıda kullanım:" 
echo "   - Browser extension: 'ModHeader' veya 'Header Editor'"
echo "   - Host header ekle: portfolio.local"
echo "   - URL: http://localhost:8888"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

