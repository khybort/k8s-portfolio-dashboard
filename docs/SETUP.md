# Setup Guide

This guide will help you set up the Portfolio Dashboard project on your local machine or in a Kubernetes cluster. - Portfolio Dashboard

## Gereksinimler

### Development Gereksinimleri

#### Temel Araçlar
- **Go**: v1.23+ ([kurulum](https://go.dev/doc/install))
- **Node.js**: v18+ ve npm v9+ ([kurulum](https://nodejs.org/))
- **Docker**: v20.10+ (veya OrbStack, Colima)
- **Docker Compose**: v2.0+ (Docker ile birlikte gelir)
- **Git**: v2.30+

#### Database & Infrastructure
- **PostgreSQL**: v14+ (Docker ile çalıştırılabilir)
- **Redis**: v7+ (Docker ile çalıştırılabilir)
- **Kafka**: v3.5+ (Docker ile çalıştırılabilir)
- **Zookeeper**: v3.8+ (Kafka için gerekli)

#### Kubernetes (Production/Staging)
- **kubectl**: v1.24+
- **Kubernetes Cluster**: 
  - kind v0.17+
  - minikube v1.28+
  - Docker Desktop Kubernetes
  - OrbStack (macOS)
  - Cloud provider (GKE, EKS, AKS)

### Opsiyonel Araçlar

- **kustomize**: v4.0+ (kubectl ile birlikte gelir)
- **make**: Build automation
- **air**: Go hot reload (development)
- **curl**: API test için
- **jq**: JSON parsing
- **Postman/Insomnia**: API testing

## Kurulum Yöntemleri

### Yöntem 1: Local Development (Docker Compose) - Önerilen

En hızlı ve kolay başlangıç yöntemi. Tüm servisler Docker Compose ile çalışır.

#### 1. Repository'yi Klonla

```bash
cd /Users/muhsin/Documents/gowit/k8s-portfolio-dashboard
```

#### 2. Environment Variables Ayarla

```bash
# .env dosyası oluştur
cp .env.example .env

# .env dosyasını düzenle
# Gerekli değişkenleri ayarla (database, redis, kafka, jwt secrets)
```

#### 3. Docker Compose ile Servisleri Başlat

```bash
# Tüm servisleri başlat (PostgreSQL, Redis, Kafka, Zookeeper)
docker-compose up -d

# Servis durumunu kontrol et
docker-compose ps
```

#### 4. Database Migration ve Seed

Docker Compose ile otomatik çalışır. Manuel çalıştırmak için:

```bash
# Backend dizinine git
cd backend

# Migration'ları çalıştır
go run cmd/migrate/main.go

# Seed'leri çalıştır (Portfolio, Articles, Projects)
go run cmd/seed/main.go

# Auth service migration ve seed
cd ../auth-service
go run cmd/migrate/main.go
go run cmd/seed/main.go
```

**Not**: Seeder'lar otomatik olarak:
- Portfolio bilgilerini (Muhsin Kılıç)
- 5 teknik makale
- 7 GitHub projesi
- Admin kullanıcı (admin@portfolio.com / Admin123!)

#### 5. Backend Servislerini Başlat

**Terminal 1 - Backend API**:
```bash
cd backend
go run cmd/server/main.go
# veya hot reload için
air
```

**Terminal 2 - Auth Service**:
```bash
cd auth-service
go run cmd/server/main.go
# veya hot reload için
air
```

#### 6. Frontend'i Başlat

**Terminal 3 - Frontend**:
```bash
cd frontend
npm install
npm run dev
```

#### 7. Erişim

- **Public Site**: `http://localhost:5173`
- **Admin Panel**: `http://localhost:5173/admin`
- **Backend API**: `http://localhost:8080`
- **Auth Service**: `http://localhost:8081`
- **API Docs**: `http://localhost:8080/swagger` (opsiyonel)

### Yöntem 2: Kubernetes Deployment

Production benzeri ortam için Kubernetes cluster'ında çalıştırma.

#### 1. Kubernetes Cluster Hazırlığı

##### Option A: kind (Kubernetes in Docker)

```bash
# kind kurulumu (macOS)
brew install kind

# Cluster oluştur (daha fazla kaynak için)
cat <<EOF | kind create cluster --name portfolio --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30080
    hostPort: 30080
  - containerPort: 30443
    hostPort: 30443
- role: worker
- role: worker
EOF

# Context'i kontrol et
kubectl config current-context
# Çıktı: kind-portfolio
```

##### Option B: minikube

```bash
# minikube kurulumu (macOS)
brew install minikube

# Cluster başlat (daha fazla kaynak için)
minikube start \
  --cpus=4 \
  --memory=8192 \
  --disk-size=20g \
  --driver=docker

# Context'i kontrol et
kubectl config current-context
# Çıktı: minikube
```

##### Option C: OrbStack (macOS - Önerilen)

```bash
# OrbStack kurulumu
# https://orbstack.dev/ adresinden indir

# Kubernetes'i etkinleştir
# OrbStack UI'dan Settings > Kubernetes > Enable

# Context'i kontrol et
kubectl config current-context
# Çıktı: orbstack
```

#### 2. Local Image'ları Build Et

```bash
# Image'ları build et
docker build -t portfolio-backend:latest ./backend
docker build -t portfolio-auth:latest ./auth-service
docker build -t portfolio-frontend:latest ./frontend
```

#### 3. Infrastructure Servislerini Deploy Et

**PostgreSQL**:
```bash
kubectl apply -f k8s/infrastructure/postgresql/deployment.yaml
kubectl apply -f k8s/infrastructure/postgresql/service.yaml
```

**Redis**:
```bash
kubectl apply -f k8s/infrastructure/redis/deployment.yaml
kubectl apply -f k8s/infrastructure/redis/service.yaml
```

**Kafka & Zookeeper**:
```bash
# Zookeeper önce
kubectl apply -f k8s/infrastructure/kafka/zookeeper.yaml

# Kafka (7.5.0 versiyonu kullanılıyor)
kubectl apply -f k8s/infrastructure/kafka/deployment.yaml
kubectl apply -f k8s/infrastructure/kafka/service.yaml
```

**Not**: Kafka için `confluentinc/cp-kafka:7.5.0` kullanılıyor (latest yerine).

#### 4. Secrets Oluştur

```bash
# Secrets oluştur
kubectl create secret generic portfolio-secrets \
  --from-literal=db-password=password \
  --from-literal=jwt-secret=your-secret-key-change-in-production \
  -n portfolio
```

#### 5. Database Migration ve Seed

**Manuel olarak** (database'e direkt bağlanarak):
```bash
# PostgreSQL'de auth_db oluştur
kubectl exec -it postgresql-0 -n portfolio -- psql -U portfolio -c "CREATE DATABASE auth_db;"

# Migration job'ları (opsiyonel - şu an binary'ler migrate/seed komutlarını desteklemiyor)
# kubectl apply -f k8s/jobs/migrate-backend.yaml
# kubectl apply -f k8s/jobs/migrate-auth.yaml
```

**Not**: Production'da init container veya Job kullanarak migration'ları çalıştırabilirsiniz.

#### 6. Application Servislerini Deploy Et

```bash
# Auth Service
kubectl apply -f k8s/auth-service/deployment.yaml
kubectl apply -f k8s/auth-service/service.yaml

# Backend API
kubectl apply -f k8s/backend/deployment.yaml
kubectl apply -f k8s/backend/service.yaml

# Frontend
kubectl apply -f k8s/frontend/deployment.yaml
kubectl apply -f k8s/frontend/service.yaml
```

**Not**: Deployment'larda `imagePullPolicy: Never` kullanılıyor (local image'lar için).

#### 6. Ingress Yapılandırması

```bash
# Ingress controller kur (eğer yoksa)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/cloud/deploy.yaml

# Ingress'i apply et
kubectl apply -f k8s/ingress/
```

#### 7. Erişim

**Port-Forward ile**:
```bash
# Frontend
kubectl port-forward -n portfolio svc/frontend 5173:80

# Backend API
kubectl port-forward -n portfolio svc/backend 8080:80

# Auth Service
kubectl port-forward -n portfolio svc/auth-service 8081:80
```

**Ingress ile** (eğer yapılandırıldıysa):
- Public Site: `http://portfolio.local`
- Admin Panel: `http://admin.portfolio.local`
- API: `http://api.portfolio.local`

## İlk Kurulum ve Yapılandırma

### 1. İlk Giriş

Seeder'lar otomatik olarak admin kullanıcı oluşturur:

- **Email**: `admin@portfolio.com`
- **Password**: `Admin123!`

1. Admin panel'e git: `http://localhost:5173/admin/login`
2. Yukarıdaki bilgilerle giriş yap
3. Dashboard'a yönlendirileceksin

### 2. Seed Data Kontrolü

Seeder'lar şunları oluşturur:
- Portfolio bilgileri (Muhsin Kılıç)
- 5 teknik makale
- 7 GitHub projesi

### 3. İlk İçerik Ekleme

Admin panel'den:
- Portfolio bilgilerini düzenle
- İlk makaleyi oluştur
- İlk GitHub projesini ekle

## Doğrulama

### Servis Durumlarını Kontrol Et

**Docker Compose**:
```bash
docker-compose ps
```

**Kubernetes**:
```bash
kubectl get pods -n portfolio
kubectl get svc -n portfolio
```

### Health Check'ler

```bash
# Backend health
curl http://localhost:8080/healthz

# Auth service health
curl http://localhost:8081/healthz

# Frontend (tarayıcıda aç)
open http://localhost:5173
```

### Database Bağlantısı

```bash
# PostgreSQL'e bağlan
docker exec -it portfolio-postgres psql -U portfolio -d portfolio

# Tabloları listele
\dt

# Çıkış
\q
```

### Redis Bağlantısı

```bash
# Redis CLI
docker exec -it portfolio-redis redis-cli

# Test
PING
# Çıktı: PONG

# Çıkış
exit
```

### Kafka Bağlantısı

```bash
# Kafka broker'a bağlan
docker exec -it portfolio-kafka kafka-topics.sh --bootstrap-server localhost:9092 --list

# Topic oluştur (test)
docker exec -it portfolio-kafka kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create \
  --topic test-topic \
  --partitions 3 \
  --replication-factor 1
```

**Kubernetes'te**:
```bash
# Kafka pod'una bağlan
kubectl exec -it kafka-0 -n portfolio -- kafka-topics.sh --bootstrap-server localhost:9092 --list
```

## Geliştirme Ortamı Yapılandırması

### Go Development

```bash
# Go modülü zaten hazır
cd backend

# Dependencies yükle
go mod download

# Air kurulumu (hot reload)
go install github.com/cosmtrek/air@latest

# .air.toml yapılandırması oluştur
air init
```

**Not**: Module path `github.com/portfolio/backend` ve `github.com/portfolio/auth-service` olarak yapılandırılmış.

### React Development

```bash
# Frontend dependencies
cd frontend
npm install

# Environment variables
cp .env.example .env.local

# Development server
npm run dev
```

### IDE Yapılandırması

**VS Code Extensions**:
- Go (golang.go)
- ESLint
- Prettier
- Docker
- Kubernetes

**GoLand/IntelliJ**:
- Go plugin
- Docker plugin
- Kubernetes plugin

## Sorun Giderme

### Port Çakışmaları

**Sorun**: Port zaten kullanılıyor

**Çözüm**:
```bash
# Port'u kullanan process'i bul
lsof -ti:8080 | xargs kill -9

# Veya docker-compose.yml'de port'ları değiştir
```

### Database Bağlantı Hatası

**Sorun**: Backend database'e bağlanamıyor

**Çözüm**:
```bash
# Database'in çalıştığını kontrol et
docker-compose ps postgres

# Connection string'i kontrol et
# .env dosyasındaki DB_HOST, DB_PORT, DB_USER, DB_PASSWORD
```

### Kafka Bağlantı Hatası

**Sorun**: Kafka'ya bağlanamıyor

**Çözüm**:
```bash
# Kafka ve Zookeeper'ın çalıştığını kontrol et
docker-compose ps kafka zookeeper

# Kafka loglarını kontrol et
docker-compose logs kafka
```

### Redis Bağlantı Hatası

**Sorun**: Redis'e bağlanamıyor

**Çözüm**:
```bash
# Redis'in çalıştığını kontrol et
docker-compose ps redis

# Redis'e manuel bağlan
docker exec -it portfolio-redis redis-cli PING
```

### Kubernetes Pod'ları Başlamıyor

**Sorun**: Pod'lar CrashLoopBackOff veya ImagePullBackOff

**Çözüm**:
```bash
# Pod loglarını kontrol et
kubectl logs -n portfolio <pod-name>

# Pod'u describe et
kubectl describe pod -n portfolio <pod-name>

# Image'ın var olduğunu kontrol et
docker images | grep portfolio

# Local image'lar için imagePullPolicy: Never kullanıldığından emin ol
# Deployment'larda imagePullPolicy: Never olmalı
```

### Kafka Hatası

**Sorun**: `KAFKA_PROCESS_ROLES is not set`

**Çözüm**: Kafka image'ı `confluentinc/cp-kafka:7.5.0` kullanılmalı (latest yerine). Deployment'da güncellenmiş olmalı.

## Sonraki Adımlar

- [Development Guide](./DEVELOPMENT.md) - Geliştirme ortamı detayları
- [Deployment Guide](./DEPLOYMENT.md) - Production deployment
- [API Documentation](./API.md) - API endpoint'leri
- [Auth Service](./AUTH.md) - Authentication detayları
- [Kafka Integration](./KAFKA.md) - Kafka kullanımı
- [Redis Integration](./REDIS.md) - Redis caching
