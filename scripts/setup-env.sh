#!/usr/bin/env sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔧 Environment Setup"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

cd "$PROJECT_DIR"

# Check if .env.dev exists
if [ ! -f ".env.dev" ]; then
    echo "📝 Creating .env.dev from env.example..."
    cp env.example .env.dev
    echo "✅ .env.dev created"
    echo ""
    echo "💡 Edit .env.dev with your development values"
else
    echo "ℹ️  .env.dev already exists"
fi

# Check if .env.prod exists
if [ ! -f ".env.prod" ]; then
    echo "📝 Creating .env.prod from env.example..."
    cp env.example .env.prod
    # Update for production
    sed -i '' 's/ENV=development/ENV=production/' .env.prod 2>/dev/null || \
    sed -i 's/ENV=development/ENV=production/' .env.prod
    sed -i '' 's/LOG_LEVEL=info/LOG_LEVEL=info/' .env.prod 2>/dev/null || \
    sed -i 's/LOG_LEVEL=debug/LOG_LEVEL=info/' .env.prod
    sed -i '' 's/DB_SSLMODE=disable/DB_SSLMODE=require/' .env.prod 2>/dev/null || \
    sed -i 's/DB_SSLMODE=disable/DB_SSLMODE=require/' .env.prod
    echo "✅ .env.prod created (production defaults)"
    echo ""
    echo "⚠️  IMPORTANT: Update .env.prod with production secrets!"
else
    echo "ℹ️  .env.prod already exists"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Environment Setup Complete"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📋 Next Steps:"
echo ""
echo "Development:"
echo "  1. Edit .env.dev (optional, defaults work)"
echo "  2. docker-compose up -d"
echo ""
echo "Production:"
echo "  1. Edit .env.prod with production values"
echo "  2. Create Kubernetes secrets:"
echo "     kubectl create secret generic portfolio-secrets \\"
echo "       --from-env-file=.env.prod -n portfolio"
echo "  3. Apply ConfigMap:"
echo "     kubectl apply -f k8s/configmap.yaml"
echo "  4. Deploy:"
echo "     kubectl apply -k k8s/"
echo ""

