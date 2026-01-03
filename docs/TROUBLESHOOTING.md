# Troubleshooting Guide - Portfolio Dashboard

## Yaygın Sorunlar ve Çözümleri

### 1. Pod Başlatılamıyor

#### Sorun: ImagePullBackOff

**Belirtiler**:
```
NAME                                  READY   STATUS             RESTARTS   AGE
portfolio-dashboard-xxxxx-xxxxx      0/1     ImagePullBackOff   0          2m
```

**Nedenler**:
- Image registry'de bulunamıyor
- Image tag yanlış
- Registry authentication hatası

**Çözümler**:

```bash
# Image'ın var olduğunu kontrol et
docker images | grep portfolio-dashboard

# kind kullanıyorsan image'ı load et
kind load docker-image portfolio-dashboard:local

# minikube kullanıyorsan
eval $(minikube docker-env)
docker build -t portfolio-dashboard:local .

# Image tag'ini kontrol et
kubectl -n portfolio describe pod <pod-name> | grep Image
```

#### Sorun: CrashLoopBackOff

**Belirtiler**:
```
NAME                                  READY   STATUS             RESTARTS   AGE
portfolio-dashboard-xxxxx-xxxxx      0/1     CrashLoopBackOff    3          2m
```

**Nedenler**:
- Container crash ediyor
- Health check başarısız
- Configuration hatası

**Çözümler**:

```bash
# Logları kontrol et
kubectl -n portfolio logs <pod-name>

# Önceki logları gör
kubectl -n portfolio logs <pod-name> --previous

# Pod'u describe et
kubectl -n portfolio describe pod <pod-name>

# Health check endpoint'ini test et
kubectl -n portfolio exec -it <pod-name> -- wget -qO- http://localhost/healthz
```

### 2. Service Erişilemiyor

#### Sorun: Port-Forward Çalışmıyor

**Belirtiler**:
```bash
$ kubectl -n portfolio port-forward svc/portfolio-dashboard 8080:80
Error: unable to forward port 8080: address already in use
```

**Çözümler**:

```bash
# Port'u kullanan process'i bul
lsof -ti:8080

# Process'i kill et
kill -9 $(lsof -ti:8080)

# Farklı port kullan
kubectl -n portfolio port-forward svc/portfolio-dashboard 8081:80
```

#### Sorun: Service Endpoint Bulunamıyor

**Belirtiler**:
```bash
$ kubectl -n portfolio get endpoints
NAME                  ENDPOINTS
portfolio-dashboard   <none>
```

**Nedenler**:
- Pod'lar Ready değil
- Service selector yanlış
- Pod label'ları eşleşmiyor

**Çözümler**:

```bash
# Pod label'larını kontrol et
kubectl -n portfolio get pods --show-labels

# Service selector'ını kontrol et
kubectl -n portfolio get svc portfolio-dashboard -o yaml | grep selector

# Pod'ların Ready olduğunu kontrol et
kubectl -n portfolio get pods
```

### 3. Deployment Başarısız

#### Sorun: Rollout Stuck

**Belirtiler**:
```bash
$ kubectl -n portfolio rollout status deploy/portfolio-dashboard
Waiting for deployment "portfolio-dashboard" rollout to finish: 0 of 1 updated replicas are available...
```

**Çözümler**:

```bash
# Deployment'ı describe et
kubectl -n portfolio describe deploy portfolio-dashboard

# ReplicaSet'leri kontrol et
kubectl -n portfolio get rs

# Pod'ları kontrol et
kubectl -n portfolio get pods

# Rollout'u rollback et
kubectl -n portfolio rollout undo deploy/portfolio-dashboard
```

### 4. Health Check Başarısız

#### Sorun: Liveness Probe Failing

**Belirtiler**:
```
Liveness probe failed: HTTP probe failed with statuscode: 500
```

**Çözümler**:

```bash
# Health endpoint'ini manuel test et
kubectl -n portfolio exec -it <pod-name> -- wget -qO- http://localhost/healthz

# Nginx config'i kontrol et
kubectl -n portfolio exec -it <pod-name> -- cat /etc/nginx/conf.d/default.conf

# Nginx loglarını kontrol et
kubectl -n portfolio exec -it <pod-name> -- tail -f /var/log/nginx/error.log
```

### 5. Resource Limitleri

#### Sorun: OOMKilled

**Belirtiler**:
```
NAME                                  READY   STATUS      RESTARTS   AGE
portfolio-dashboard-xxxxx-xxxxx      0/1     OOMKilled    3          2m
```

**Çözümler**:

```bash
# Resource kullanımını kontrol et
kubectl -n portfolio top pod

# Resource limitlerini artır
kubectl -n portfolio edit deploy portfolio-dashboard
# resources.limits.memory değerini artır
```

### 6. Network Sorunları

#### Sorun: DNS Resolution Failed

**Belirtiler**:
```
Error: dial tcp: lookup service-name: no such host
```

**Çözümler**:

```bash
# DNS'i test et
kubectl -n portfolio run -it --rm debug --image=busybox --restart=Never -- nslookup portfolio-dashboard.portfolio.svc.cluster.local

# Service'i kontrol et
kubectl -n portfolio get svc portfolio-dashboard
```

### 7. Image Build Sorunları

#### Sorun: Docker Build Failed

**Belirtiler**:
```
ERROR: failed to solve: failed to fetch...
```

**Çözümler**:

```bash
# Docker daemon'ın çalıştığını kontrol et
docker info

# Build cache'i temizle
docker builder prune

# Network bağlantısını kontrol et
docker pull nginx:1.25-alpine

# Build'i verbose modda çalıştır
docker build --progress=plain -t portfolio-dashboard:local .
```

### 8. Kafka Sorunları

#### Sorun: KAFKA_PROCESS_ROLES is not set

**Belirtiler**:
```
Error: environment variable "KAFKA_PROCESS_ROLES" is not set
```

**Nedenler**:
- Kafka image versiyonu uyumsuz (latest yerine 7.5.0 kullanılmalı)
- Environment variable eksik

**Çözümler**:

```bash
# Kafka deployment'ı kontrol et
kubectl -n portfolio get deployment kafka -o yaml | grep image

# Image'ı 7.5.0'a güncelle
kubectl -n portfolio set image deployment/kafka kafka=confluentinc/cp-kafka:7.5.0

# Veya deployment.yaml'ı düzenle
# image: confluentinc/cp-kafka:7.5.0
# env:
# - name: KAFKA_AUTO_CREATE_TOPICS_ENABLE
#   value: "true"
```

#### Sorun: Kafka Pod CrashLoopBackOff

**Belirtiler**:
```
NAME      READY   STATUS             RESTARTS   AGE
kafka-0   0/1     CrashLoopBackOff   3          2m
```

**Çözümler**:

```bash
# Kafka loglarını kontrol et
kubectl -n portfolio logs kafka-0 --tail 50

# Zookeeper'ın çalıştığını kontrol et
kubectl -n portfolio get pods | grep zookeeper

# Kafka'yı yeniden başlat
kubectl -n portfolio delete pod kafka-0
```

### 9. Database Connection Sorunları

#### Sorun: CreateContainerConfigError

**Belirtiler**:
```
Error: couldn't find key db-host in Secret portfolio/portfolio-secrets
```

**Çözümler**:

```bash
# Secret'ı kontrol et
kubectl -n portfolio get secret portfolio-secrets -o yaml

# Secret'ı yeniden oluştur
kubectl -n portfolio delete secret portfolio-secrets
kubectl -n portfolio create secret generic portfolio-secrets \
  --from-literal=db-password=password \
  --from-literal=jwt-secret=your-secret-key-change-in-production

# Deployment'ı yeniden başlat
kubectl -n portfolio rollout restart deployment/backend
```

#### Sorun: Database "auth_db" does not exist

**Belirtiler**:
```
FATAL: database "auth_db" does not exist (SQLSTATE 3D000)
```

**Çözümler**:

```bash
# PostgreSQL'de database oluştur
kubectl exec -it postgresql-0 -n portfolio -- psql -U portfolio -c "CREATE DATABASE auth_db;"

# Veya migration job'ını çalıştır
kubectl apply -f k8s/jobs/migrate-auth.yaml
```

### 10. Local Image Sorunları

#### Sorun: ImagePullBackOff (Local Images)

**Belirtiler**:
```
NAME      READY   STATUS             RESTARTS   AGE
backend   0/1     ImagePullBackOff   0          2m
```

**Çözümler**:

```bash
# Image'ları build et
docker build -t portfolio-backend:latest ./backend
docker build -t portfolio-auth:latest ./auth-service
docker build -t portfolio-frontend:latest ./frontend

# Deployment'da imagePullPolicy: Never olduğundan emin ol
kubectl -n portfolio get deployment backend -o yaml | grep imagePullPolicy

# imagePullPolicy: Never olmalı
```

### 11. Kubernetes Cluster Sorunları

#### Sorun: kubectl Connection Refused

**Belirtiler**:
```
The connection to the server localhost:8080 was refused
```

**Çözümler**:

```bash
# Context'i kontrol et
kubectl config current-context

# Context'leri listele
kubectl config get-contexts

# Context'i değiştir
kubectl config use-context <context-name>

# Cluster'ın çalıştığını kontrol et
kubectl cluster-info
```

#### Sorun: Node Not Ready

**Belirtiler**:
```
NAME       STATUS     ROLES                  AGE   VERSION
node-1     NotReady   control-plane,master   5m    v1.24.0
```

**Çözümler**:

```bash
# Node'u describe et
kubectl describe node <node-name>

# kind için cluster'ı yeniden başlat
kind delete cluster --name portfolio
kind create cluster --name portfolio

# minikube için
minikube delete
minikube start
```

### 9. Kafka Sorunları

#### Sorun: KAFKA_PROCESS_ROLES is not set

**Belirtiler**:
```
Error: environment variable "KAFKA_PROCESS_ROLES" is not set
```

**Çözümler**:

```bash
# Kafka deployment'ı kontrol et
kubectl -n portfolio get deployment kafka -o yaml | grep image

# Image'ı 7.5.0'a güncelle (deployment.yaml'da)
# image: confluentinc/cp-kafka:7.5.0
```

#### Sorun: Kafka Pod CrashLoopBackOff

**Belirtiler**:
```
NAME      READY   STATUS             RESTARTS   AGE
kafka-0   0/1     CrashLoopBackOff   3          2m
```

**Çözümler**:

```bash
# Kafka loglarını kontrol et
kubectl -n portfolio logs kafka-0 --tail 50

# Zookeeper'ın çalıştığını kontrol et
kubectl -n portfolio get pods | grep zookeeper
```

### 10. Database Connection Sorunları

#### Sorun: CreateContainerConfigError

**Belirtiler**:
```
Error: couldn't find key db-host in Secret portfolio/portfolio-secrets
```

**Çözümler**:

```bash
# Secret'ı yeniden oluştur
kubectl -n portfolio delete secret portfolio-secrets
kubectl -n portfolio create secret generic portfolio-secrets \
  --from-literal=db-password=password \
  --from-literal=jwt-secret=your-secret-key-change-in-production
```

#### Sorun: Database "auth_db" does not exist

**Belirtiler**:
```
FATAL: database "auth_db" does not exist (SQLSTATE 3D000)
```

**Çözümler**:

```bash
# PostgreSQL'de database oluştur
kubectl exec -it postgresql-0 -n portfolio -- psql -U portfolio -c "CREATE DATABASE auth_db;"
```

### 11. Local Image Sorunları

#### Sorun: ImagePullBackOff (Local Images)

**Çözümler**:

```bash
# Image'ları build et
docker build -t portfolio-backend:latest ./backend
docker build -t portfolio-auth:latest ./auth-service
docker build -t portfolio-frontend:latest ./frontend

# Deployment'da imagePullPolicy: Never olduğundan emin ol
```

## Debug Komutları

### Genel Debug

```bash
# Tüm kaynakları listele
kubectl -n portfolio get all

# Pod'ların detaylı durumu
kubectl -n portfolio get pods -o wide

# Event'leri gör
kubectl -n portfolio get events --sort-by='.lastTimestamp'

# Resource kullanımı
kubectl -n portfolio top pod
kubectl -n portfolio top node
```

### Log Debug

```bash
# Tüm pod logları
kubectl -n portfolio logs -l app=portfolio-dashboard

# Belirli bir pod'un logları
kubectl -n portfolio logs <pod-name>

# Önceki container logları
kubectl -n portfolio logs <pod-name> --previous

# Logları takip et
kubectl -n portfolio logs -f <pod-name>

# Tüm container logları
kubectl -n portfolio logs <pod-name> --all-containers=true
```

### Network Debug

```bash
# Service endpoint'lerini kontrol et
kubectl -n portfolio get endpoints

# Service detaylarını gör
kubectl -n portfolio describe svc portfolio-dashboard

# Port-forward test
kubectl -n portfolio port-forward svc/portfolio-dashboard 8080:80 &
curl http://localhost:8080/healthz
```

### Container Debug

```bash
# Pod'a exec et
kubectl -n portfolio exec -it <pod-name> -- sh

# Nginx config'i kontrol et
kubectl -n portfolio exec <pod-name> -- cat /etc/nginx/conf.d/default.conf

# Nginx'i test et
kubectl -n portfolio exec <pod-name> -- nginx -t

# Dosya sistemini kontrol et
kubectl -n portfolio exec <pod-name> -- ls -la /usr/share/nginx/html/
```

## Performans Sorunları

### Yavaş Yükleme

**Çözümler**:

```bash
# Resource limitlerini kontrol et
kubectl -n portfolio describe pod <pod-name> | grep -A 5 Resources

# CPU throttling kontrolü
kubectl -n portfolio top pod

# Network latency
kubectl -n portfolio exec <pod-name> -- ping -c 3 google.com
```

### Yüksek Memory Kullanımı

**Çözümler**:

```bash
# Memory kullanımını izle
kubectl -n portfolio top pod

# Memory limitlerini artır
kubectl -n portfolio edit deploy portfolio-dashboard
```

## Yardım Alma

### Logları Toplama

```bash
# Tüm logları topla
kubectl -n portfolio logs -l app=portfolio-dashboard > logs.txt

# Pod describe çıktısı
kubectl -n portfolio describe pod <pod-name> > pod-describe.txt

# Event'ler
kubectl -n portfolio get events > events.txt
```

### GitHub Issue Açma

Issue açarken şunları ekle:
- Kubernetes versiyonu: `kubectl version`
- Cluster tipi: kind/minikube/Docker Desktop
- Hata mesajları: Log çıktıları
- Adımlar: Ne yaptığın ve ne beklediğin

## Önleyici Bakım

### Düzenli Kontroller

```bash
# Pod durumlarını kontrol et
kubectl -n portfolio get pods

# Resource kullanımını izle
kubectl -n portfolio top pod

# Event'leri kontrol et
kubectl -n portfolio get events

# Image'ları güncelle
docker pull nginx:1.25-alpine
```

### Monitoring Kurulumu

- Prometheus metrics
- Grafana dashboards
- Alerting rules

## Hızlı Referans

| Sorun | Komut |
|-------|-------|
| Pod logları | `kubectl -n portfolio logs <pod-name>` |
| Pod durumu | `kubectl -n portfolio get pods` |
| Service kontrol | `kubectl -n portfolio get svc` |
| Deployment rollback | `kubectl -n portfolio rollout undo deploy/portfolio-dashboard` |
| Port-forward | `kubectl -n portfolio port-forward svc/portfolio-dashboard 8080:80` |
| Health check | `curl http://localhost:8080/healthz` |

