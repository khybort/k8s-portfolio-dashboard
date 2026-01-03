# Portfolio Dashboard

A modern, full-stack portfolio management system built with Kubernetes, Go, React, and microservices architecture. This project demonstrates production-ready practices including containerization, orchestration, event-driven architecture, and comprehensive testing.

## 🚀 Features

### Core Functionality
- **Portfolio Management**: Manage personal portfolio information, bio, and social links
- **Article Management**: Create, edit, and publish blog articles with markdown support
- **Project Showcase**: Display GitHub projects with descriptions, technologies, and live demos
- **Admin Panel**: Full-featured admin interface for content management
- **Public Portfolio**: Beautiful, responsive public-facing portfolio site

### Technical Features
- **Microservices Architecture**: Separated backend, frontend, and authentication services
- **Kubernetes Orchestration**: Full Kubernetes deployment with Ingress, Services, and StatefulSets
- **Event-Driven Architecture**: Kafka integration for asynchronous event processing
- **Caching Layer**: Redis for performance optimization
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT-based authentication with refresh tokens
- **Modern UI/UX**: React 18 + TypeScript + Tailwind CSS with dark mode
- **Comprehensive Testing**: Playwright E2E tests for all admin operations
- **CI/CD Ready**: GitHub Actions workflow with security scanning

## 📋 Table of Contents

- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Development](#development)
- [Deployment](#deployment)
- [Testing](#testing)
- [API Documentation](#api-documentation)
- [Configuration](#configuration)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## 🏗️ Architecture

The system follows a microservices architecture with three main services:

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend  │────▶│   Backend   │────▶│ Auth Service│
│   (React)   │     │    (Go)     │     │    (Go)     │
└─────────────┘     └─────────────┘     └─────────────┘
                            │
                            ▼
                    ┌─────────────┐
                    │ PostgreSQL  │
                    └─────────────┘
                            │
        ┌───────────────────┴───────────────────┐
        ▼                                       ▼
┌─────────────┐                         ┌─────────────┐
│    Redis    │                         │    Kafka    │
│   (Cache)   │                         │  (Events)   │
└─────────────┘                         └─────────────┘
```

### Services

1. **Frontend Service** (React + TypeScript)
   - Public portfolio site
   - Admin panel for content management
   - Responsive design with dark mode

2. **Backend Service** (Go + Gin)
   - RESTful API for portfolio, articles, and projects
   - Kafka event publishing
   - Redis caching
   - GORM database operations

3. **Auth Service** (Go + Gin)
   - User authentication and authorization
   - JWT token generation and validation
   - User management

### Infrastructure

- **PostgreSQL**: Primary database for all services
- **Redis**: Caching and session management
- **Kafka**: Event streaming for asynchronous operations
- **Nginx Ingress**: Routing and load balancing

For detailed architecture documentation, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

## 🛠️ Tech Stack

### Frontend
- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool
- **Tailwind CSS** - Styling
- **Framer Motion** - Animations
- **React Query** - Data fetching
- **React Router** - Routing
- **Playwright** - E2E testing

### Backend
- **Go 1.23** - Programming language
- **Gin** - Web framework
- **GORM** - ORM
- **Zap** - Logging
- **JWT** - Authentication
- **bcrypt** - Password hashing

### Infrastructure
- **Kubernetes** - Container orchestration
- **Docker** - Containerization
- **PostgreSQL** - Database
- **Redis** - Cache
- **Kafka** - Event streaming
- **Nginx Ingress** - Load balancer

## 📦 Prerequisites

### Required
- **Docker** (20.10+) or **OrbStack** or **Docker Desktop**
- **Kubernetes** cluster (kind, minikube, or Docker Desktop Kubernetes)
- **kubectl** (1.28+)
- **Node.js** (18+) and **npm** (for local frontend development)
- **Go** (1.23+) (for local backend development)

### Optional
- **TablePlus** or **pgAdmin** (for database management)
- **k9s** (for Kubernetes management)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd k8s-portfolio-dashboard
```

### 2. Start the Project

```bash
./scripts/start.sh
```

This script will:
- Create Kubernetes namespace
- Deploy all services
- Add local domain entries to `/etc/hosts`
- Start port-forwarding
- Deploy Ingress controller (if needed)

### 3. Access the Application

- **Public Portfolio**: http://localhost:5173
- **Admin Panel**: http://localhost:5173/admin/login
  - Email: `admin@portfolio.com`
  - Password: `Admin123!`
- **Backend API**: http://localhost:8080
- **Auth API**: http://localhost:8081

### 4. Stop the Project

```bash
./scripts/stop.sh
```

This will stop port-forwarding and remove local domain entries, but keep resources running.

### 5. Complete Cleanup

```bash
./scripts/cleanup.sh
```

This will remove all resources and namespace.

For detailed setup instructions, see [docs/SETUP.md](docs/SETUP.md).

## 📁 Project Structure

```
k8s-portfolio-dashboard/
├── backend/                 # Go backend service
│   ├── cmd/
│   │   ├── server/         # Main server application
│   │   └── seed/           # Database seeder
│   ├── internal/
│   │   ├── api/            # HTTP handlers
│   │   ├── model/           # Data models
│   │   ├── repository/     # Database layer
│   │   ├── service/         # Business logic
│   │   └── config/         # Configuration
│   ├── migrations/          # Database migrations
│   └── Dockerfile
├── frontend/                # React frontend
│   ├── src/
│   │   ├── admin/          # Admin panel
│   │   ├── public/         # Public portfolio site
│   │   └── services/       # API clients
│   ├── tests/              # E2E tests
│   └── Dockerfile
├── auth-service/           # Authentication service
│   ├── cmd/
│   │   ├── server/         # Main server
│   │   └── seed/           # User seeder
│   ├── internal/
│   └── migrations/
├── k8s/                    # Kubernetes manifests
│   ├── namespace.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   └── kustomization.yaml
├── infrastructure/         # Infrastructure components
│   ├── postgresql/
│   ├── redis/
│   └── kafka/
├── scripts/                # Utility scripts
│   ├── start.sh           # Start project
│   ├── stop.sh            # Stop project
│   └── cleanup.sh         # Cleanup everything
└── docs/                   # Documentation
    ├── ARCHITECTURE.md
    ├── SETUP.md
    ├── DEPLOYMENT.md
    ├── DEVELOPMENT.md
    ├── TESTING.md
    └── TROUBLESHOOTING.md
```

## 💻 Development

### Local Development Setup

#### Backend Development

```bash
cd backend
go mod download
go run cmd/server/main.go
```

#### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

#### Auth Service Development

```bash
cd auth-service
go mod download
go run cmd/server/main.go
```

### Environment Variables

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for detailed environment variable configuration.

### Database Migrations

Migrations are automatically applied on service startup. For manual migration:

```bash
# Backend
kubectl exec -it <backend-pod> -n portfolio -- /app/migrate up

# Auth Service
kubectl exec -it <auth-pod> -n portfolio -- /app/migrate up
```

### Database Seeding

Seeders run automatically via Kubernetes Jobs. For manual seeding:

```bash
# Backend seeder
kubectl create job --from=cronjob/seed-backend manual-seed-backend -n portfolio

# Auth seeder
kubectl exec -it <auth-pod> -n portfolio -- /app/seed
```

For detailed development guide, see [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md).

## 🚢 Deployment

### Kubernetes Deployment

#### Using Kustomize

```bash
kubectl apply -k k8s/
```

#### Manual Deployment

```bash
# Create namespace
kubectl apply -f k8s/namespace.yaml

# Deploy infrastructure
kubectl apply -f infrastructure/postgresql/
kubectl apply -f infrastructure/redis/
kubectl apply -f infrastructure/kafka/

# Deploy services
kubectl apply -f backend/k8s/
kubectl apply -f frontend/k8s/
kubectl apply -f auth-service/k8s/

# Deploy ingress
kubectl apply -f k8s/ingress.yaml
```

### Docker Compose (Local Development)

```bash
docker-compose up -d
```

### Production Considerations

- Use proper secrets management (Kubernetes Secrets, Vault)
- Configure TLS/SSL for Ingress
- Set up monitoring and logging (Prometheus, Grafana, ELK)
- Configure resource limits and requests
- Enable horizontal pod autoscaling
- Set up backup strategies for databases

For detailed deployment guide, see [docs/DEPLOYMENT.md](docs/DEPLOYMENT.md).

## 🧪 Testing

### Running Tests

#### E2E Tests (Playwright)

```bash
cd frontend
npm install
npx playwright install
npx playwright test
```

#### Test Coverage

- **Authentication**: Login, logout, protected routes
- **Articles CRUD**: Create, read, update, delete, search
- **Projects CRUD**: Create, read, update, delete, search
- **Portfolio**: View, update, cancel

### Test Reports

```bash
npx playwright show-report
```

For detailed testing guide, see [docs/TESTING.md](docs/TESTING.md).

## 📚 API Documentation

### Public Endpoints

#### Articles
- `GET /api/v1/articles` - List articles
- `GET /api/v1/articles/:id` - Get article by ID
- `GET /api/v1/articles/slug/:slug` - Get article by slug

#### Projects
- `GET /api/v1/projects` - List projects
- `GET /api/v1/projects/:id` - Get project by ID

#### Portfolio
- `GET /api/v1/portfolio` - Get portfolio information

### Admin Endpoints (Require Authentication)

#### Articles
- `POST /api/v1/admin/articles` - Create article
- `PUT /api/v1/admin/articles/:id` - Update article
- `DELETE /api/v1/admin/articles/:id` - Delete article

#### Projects
- `POST /api/v1/admin/projects` - Create project
- `PUT /api/v1/admin/projects/:id` - Update project
- `DELETE /api/v1/admin/projects/:id` - Delete project

#### Portfolio
- `PUT /api/v1/admin/portfolio` - Update portfolio

### Authentication Endpoints

- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/refresh` - Refresh access token
- `GET /api/v1/auth/verify` - Verify token

For detailed API documentation, see [docs/API.md](docs/API.md) (if available).

## ⚙️ Configuration

### Environment Setup

The project supports multiple environments (development, production) with clear separation of configuration.

#### Quick Setup

```bash
# Setup environment files
./scripts/setup-env.sh
```

This creates:
- `.env.dev` - Development configuration
- `.env.prod` - Production configuration

#### Environment Files

- **Development** (`.env.dev`): Used with Docker Compose
- **Production** (`.env.prod`): Used with Kubernetes (via ConfigMap/Secrets)
- **Example** (`env.example`): Template with all available variables

### Environment Variables

#### Backend Service

| Variable | Description | Dev Default | Prod Default |
|----------|-------------|-------------|--------------|
| `ENV` | Environment name | `development` | `production` |
| `SERVER_PORT` | Server port | `8080` | `8080` |
| `LOG_LEVEL` | Log level | `debug` | `info` |
| `DB_HOST` | Database host | `postgresql` | `postgresql.portfolio.svc.cluster.local` |
| `DB_PASSWORD` | Database password | `password` | From Secret |
| `REDIS_HOST` | Redis host | `redis` | `redis.portfolio.svc.cluster.local` |
| `KAFKA_BROKERS` | Kafka brokers | `kafka:9092` | `kafka.portfolio.svc.cluster.local:9092` |
| `AUTH_SERVICE_URL` | Auth service URL | `http://auth-service:8081` | `http://auth-service:80` |
| `JWT_SECRET` | JWT secret | `dev-secret-key` | From Secret |

#### Frontend Service

| Variable | Description | Dev Default | Prod Default |
|----------|-------------|-------------|--------------|
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` | `https://api.portfolio.local` |
| `VITE_AUTH_URL` | Auth service URL | `http://localhost:8081` | `https://auth.portfolio.local` |

**Note**: Frontend variables must be set at build time (Vite requirement).

#### Auth Service

| Variable | Description | Dev Default | Prod Default |
|----------|-------------|-------------|--------------|
| `AUTH_DB_NAME` | Database name | `auth_db` | `auth_db` |
| `AUTH_DB_PASSWORD` | Database password | `password` | From Secret |
| `JWT_SECRET` | JWT secret | `dev-secret-key` | From Secret |
| `JWT_ACCESS_EXPIRY` | Access token expiry | `15m` | `15m` |
| `JWT_REFRESH_EXPIRY` | Refresh token expiry | `168h` | `168h` |

### Configuration Methods

#### Development (Docker Compose)

Uses `.env.dev` file:

```bash
# Copy example
cp env.example .env.dev

# Edit if needed
# Then start
docker-compose up -d
```

#### Production (Kubernetes)

Uses ConfigMap and Secrets:

**1. Create Secrets**:
```bash
kubectl create secret generic portfolio-secrets \
  --from-literal=db-password='your-secure-password' \
  --from-literal=auth-db-password='your-secure-password' \
  --from-literal=jwt-secret='your-strong-random-jwt-secret' \
  --from-literal=redis-password='your-redis-password' \
  -n portfolio
```

**2. Apply ConfigMap**:
```bash
kubectl apply -f k8s/configmap.yaml
```

**3. Deploy**:
```bash
kubectl apply -k k8s/
```

For detailed environment configuration guide, see [docs/ENVIRONMENT.md](docs/ENVIRONMENT.md).

## 🔧 Troubleshooting

### Common Issues

#### Services Not Starting
- Check pod logs: `kubectl logs -n portfolio <pod-name>`
- Verify secrets exist: `kubectl get secrets -n portfolio`
- Check resource limits

#### Database Connection Issues
- Verify PostgreSQL is running: `kubectl get pods -n portfolio -l app=postgresql`
- Check connection string in config
- Verify network policies

#### Authentication Failures
- Verify auth service is running
- Check JWT secret configuration
- Verify token expiration settings

#### Ingress Not Working
- Check ingress controller: `kubectl get pods -n ingress-nginx`
- Verify ingress resource: `kubectl get ingress -n portfolio`
- Check `/etc/hosts` entries

For detailed troubleshooting guide, see [docs/TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md).

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go code style guidelines
- Use TypeScript for frontend code
- Write tests for new features
- Update documentation
- Follow commit message conventions

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [React](https://react.dev/)
- [Tailwind CSS](https://tailwindcss.com/)
- [Playwright](https://playwright.dev/)
- [Kubernetes](https://kubernetes.io/)

## 📞 Support

For issues and questions:
- Open an issue on GitHub
- Check the [Troubleshooting Guide](docs/TROUBLESHOOTING.md)
- Review the [Documentation](docs/)

## 🗺️ Roadmap

- [ ] Add more comprehensive API documentation
- [ ] Implement real-time updates with WebSockets
- [ ] Add analytics and monitoring
- [ ] Support for multiple portfolios
- [ ] Enhanced search functionality
- [ ] Image upload and management
- [ ] Comment system for articles
- [ ] RSS feed generation

---

**Built with ❤️ using Kubernetes, Go, and React**
