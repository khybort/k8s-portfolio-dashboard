# Quick Start Guide

## 🚀 Hızlı Başlangıç

### 1. Docker Compose ile Local Development

```bash
# Tüm servisleri başlat (migration + seed otomatik)
make up
# veya
docker-compose up -d

# Servis durumunu kontrol et
docker-compose ps

# Logları izle
docker-compose logs -f backend
```

**Erişim**:
- Frontend (Admin): http://localhost:5173/admin/login
- Backend API: http://localhost:8080
- Auth Service: http://localhost:8081

### 2. İlk Giriş

Seeder'lar otomatik olarak admin kullanıcı oluşturur:

- **Email**: `admin@portfolio.com`
- **Password**: `Admin123!`

1. http://localhost:5173/admin/login adresine git
2. Yukarıdaki bilgilerle giriş yap
3. Dashboard'a yönlendirileceksin

### 3. Seed Data Kontrolü

Seeder'lar şunları oluşturur:

**Portfolio**:
- Muhsin Kılıç'ın bilgileri
- Social links (GitHub, LinkedIn, Phone)

**Makaleler** (5 adet):
- Real-Time Bidding Systems
- Kubernetes Best Practices
- React Performance
- Event-Driven Architecture
- Redis Caching

**Projeler** (7 adet):
- Real-Time Bidding Platform
- Campaign Management Dashboard
- Data Pipeline System
- SOAR Platform
- UAV Control System
- Energy Management Platform
- OpenAPI Generator

### 4. API Test

```bash
# Public articles
curl http://localhost:8080/api/v1/articles

# Public projects
curl http://localhost:8080/api/v1/projects

# Portfolio info
curl http://localhost:8080/api/v1/portfolio
```

### 5. Admin Panel Kullanımı

1. Login yap
2. **Articles** sayfasından makaleleri görüntüle/düzenle
3. **Projects** sayfasından projeleri yönet
4. **Portfolio** sayfasından portfolio bilgilerini güncelle

## 📋 Checklist

- [x] Backend API (Go)
- [x] Auth Service (Go + JWT)
- [x] Frontend (React + TypeScript)
- [x] Kafka Integration
- [x] Redis Caching
- [x] Docker Compose
- [x] Kubernetes Manifests
- [x] Database Migrations
- [x] Database Seeders
- [x] Documentation

## 🔗 Önemli Linkler

- [Architecture Docs](./docs/ARCHITECTURE.md)
- [API Documentation](./docs/API.md)
- [Setup Guide](./docs/SETUP.md)
- [Development Guide](./docs/DEVELOPMENT.md)
