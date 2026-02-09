# k8s-fullstack-blueprint

**Production-oriented cloud-native full-stack blueprint for Go + Angular + Kubernetes**

A repository structure and documentation template for building and deploying cloud-native full-stack applications using Go (backend), Angular (frontend), Docker, and Kubernetes.

> **Note:** This is a **blueprint/template** project, not an actual application. It provides structural foundation and documentation for teams to build upon.

---

## Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────┐
│                     Kubernetes Cluster                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │
│  │   Go API    │  │ Angular SPA │  │   PostgreSQL    │  │
│  │  Service    │  │  Service    │  │   Database      │  │
│  │  (DevSpace) │  │  (DevSpace) │  │  (DevSpace)     │  │
│  └─────────────┘  └─────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────┘
          │                    │                 │
          └────────────────────┼─────────────────┘
                               │
                     ┌─────────▼─────────┐
                     │    DevSpace CLI   │
                     │  (build + sync)   │
                     └───────────────────┘
```

### Component Responsibilities

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Backend API** | Go (Golang) | REST API server, business logic, data access |
| **Frontend UI** | Angular | Progressive Web App with reactive UI, client-side routing |
| **Orchestration** | Kubernetes | Container orchestration: deployments, services, ingress |
| **Containerization** | Docker | Consistent build artifacts, environment isolation |

### Design Principles

- **Infrastructure as Code (IaC):** All infrastructure in version-controlled YAML manifests
- **Immutable Infrastructure:** Docker images built once, deployed everywhere
- **Separation of Concerns:** Clear boundaries between backend, frontend, infrastructure
- **Environment Parity:** Local development mirrors production using DevSpace and K8s
- **Security First:** Secrets management, non-root containers, network policies, RBAC

---

## Security Considerations

> **CRITICAL:** The default credentials amd APIs in this blueprint are for **development and testing only**. Never use them in production.

### Database Credentials

- **Development:** The default PostgreSQL password is `password`. This is intentionally weak for local testing.
- **Production:** You MUST provide strong, randomly generated passwords via Kubernetes Secrets or external secret management (e.g., HashiCorp Vault, AWS Secrets Manager, Azure Key Vault).
- **Helm Values:** In [`deployment/fullstack-blueprint-app/values.yaml`](deployment/fullstack-blueprint-app/values.yaml), both `backend.config.DATABASE_URL` and `postgres.config.PASSWORD` are set to empty strings with warnings. You must override these with secure values via Helm `--set` flags or custom values files.

### Authentication

- **APIs:** This blueprint uses unauthenticated APIs for local development
- **Production:**: In production, implement JWT/Other authentication and use an API Gateway/Ingress with JWT/Other validation.

### Local Development

When running locally with DevSpace or port-forwarding:
- Use the default credentials only in isolated development environments
- Never expose these services to public networks
- Clear browser cache and avoid using production data in development

### Production Deployment Requirements

- [ ] Generate strong random passwords for all services
- [ ] Store credentials in a secure secrets manager
- [ ] Configure Kubernetes Secrets (do not commit secrets to version control)
- [ ] Set appropriate RBAC and network policies
- [ ] Enable TLS/SSL for database connections
- [ ] Rotate credentials regularly
- [ ] Audit secret access logs
- [ ] API Authentication and API Gateway Ingress to Lock Down External APIs

---

## Directory Structure

```
k8s-fullstack-blueprint/
├── backend/go/          # Go backend service
│   ├── main.go          # Application entry point
│   ├── go.mod           # Go module definition
│   ├── Dockerfile       # Multi-stage container build
│   ├── api/             # HTTP layer (routes, middleware, controllers)
│   ├── service/         # Business logic layer
│   └── db/              # Database layer (models, repositories, connections)
├── frontend/angular/    # Angular frontend application
│   ├── src/             # Application source code
│   ├── nginx/           # Nginx configuration for container
│   ├── Dockerfile       # Multi-stage container build
│   ├── angular.json     # Angular CLI configuration
│   ├── package.json     # NPM dependencies
│   └── tsconfig.json    # TypeScript configuration
├── deployment/
│   └── fullstack-blueprint-app/  # Helm chart for Kubernetes
│       ├── Chart.yaml    # Chart metadata
│       ├── values.yaml   # Default configuration values
│       └── templates/    # Helm templates (deployments, services, ingress, etc.)
└── devspace.yaml        # DevSpace configuration (hot-reload on K8s)
```

---

## Appointment Booking Feature (POC)

The appointment booking feature is a **proof-of-concept (POC)** implementation demonstrating:
- Service and staff management
- Appointment scheduling with availability checking
- Customer-facing booking flow
- Dashboard with basic insights

**Important:** This POC currently returns **all appointments** to any user. In a production environment, you must implement:

- **JWT Authentication:** Extract user identity from JWT tokens
- **Authorization:** Filter appointments by `user_id` to ensure users can only access their own appointments
- **Role-based Access:** Differentiate between customer and staff/administrator views

Backend endpoints to secure:
- `GET /api/appt_booking/appointments` → filter by authenticated user ID
- `POST /api/appt_booking/appointments` → associate with authenticated user ID

See [`backend/go/api/appt_booking/appointment_handler.go`](backend/go/api/appt_booking/appointment_handler.go) for implementation details.

---

## Local Development Setup

### Prerequisites

- **Docker Desktop** - [Download](https://www.docker.com/products/docker-desktop)
- **Go** 1.22+ - [Install](https://go.dev/dl/)
- **Node.js** 20+ & **npm** - [Install](https://nodejs.org/)
- **Angular CLI** (optional) - `npm install -g @angular/cli`
- **DevSpace** CLI - [Install](https://devspace.sh/cli/docs/installation)
- **kubectl** - [Install](https://kubernetes.io/docs/tasks/tools/)
- **Helm** (v3+) - [Install](https://helm.sh/docs/intro/install/)

### Quick Start with DevSpace

The fastest way to run the full stack with hot-reload on Kubernetes:

```bash
# 1. Ensure your Kubernetes cluster is running & correct namespace
kubectl cluster-info
devspace use namespace {devspace_namespace} 

# 2. Start DevSpace (Terminal 1)
devspace dev

# 3. Build and sync Angular frontend (Terminal 2)
./frontend/angular/devspace_angular.sh   # Linux/macOS
# or
./frontend/angular/devspace_angular.ps1  # Windows PowerShell
```
- NOTE: To set back namespace: use 
-     devspace use namespace default

DevSpace will:
- Deploy the application's helm chart
- Build Docker images from source
- Patch running pods to use local images
- Sync source code for live reload (Angular: `frontend/angular/dist/browser`)
- Port-forward services to localhost
- Open browser to `http://localhost:4200`

**Access services:**
- Frontend: http://localhost:4200
- Backend API: http://localhost:8080

**Access database:**
```bash
kubectl port-forward svc/fullstack-postgres 5432:5432 -n {namespace}
```
Then connect with:
- Server: localhost:5432
- Username: postgres
- Password: password
- Database: k8s_blueprint
    
**Stop:** `Ctrl+C` to stop DevSpace (pods remain running). Use `devspace purge` to delete all resources.

**Note:** If UI changes don't appear, clear browser cache: F12 → Right-click reload button → "Empty Cache and Hard Reload".

---

## Kubernetes Deployment (Helm)

### Prerequisites

- `helm` CLI (v3+)
- `kubectl` configured with cluster access
- Kubernetes cluster with ingress controller (optional)

### 1. Build and Push Docker Images

```bash
# Backend
cd backend/go
docker build -t your-registry/backend:latest .
docker push your-registry/backend:latest

# Frontend
cd ../angular
docker build -t your-registry/frontend:latest .
docker push your-registry/frontend:latest
```

### 2. Deploy with Helm

```bash
# Basic installation
helm install fullstack ./deployment/fullstack-blueprint-app

# With custom values
helm install fullstack ./deployment/fullstack-blueprint-app -f values-prod.yaml

# Or with overrides
helm install fullstack ./deployment/fullstack-blueprint-app \
  --set backend.image.repository=your-registry/backend \
  --set backend.image.tag=latest \
  --set frontend.image.repository=your-registry/frontend \
  --set frontend.image.tag=latest
```

### 3. Verify Deployment

```bash
helm list
kubectl get pods -n backend
kubectl get pods -n frontend
kubectl get pods -n database
```

### 4. Access the Application

**With Ingress:**
```bash
kubectl get ingress -n frontend
# Access via configured hostname/IP
```

**Without Ingress (port-forwarding):**
```bash
# Frontend
kubectl port-forward svc/fullstack-frontend 8080:80 -n frontend
# Open: http://localhost:8080

# Backend API
kubectl port-forward svc/fullstack-backend 8081:8080 -n backend
# API: http://localhost:8081

# PostgreSQL Database (for local development/verification)
kubectl port-forward svc/fullstack-postgres 5432:5432 -n database
# Then connect with:
#   Server: localhost:5432
#   Username: {Username}
#   Password: {Password}
#   Database: k8s_blueprint
```

### 5. Update and Rollback

```bash
# Update images
helm upgrade fullstack ./deployment/fullstack-blueprint-app \
  --set backend.image.tag=v2.0 \
  --set frontend.image.tag=v2.0

# Monitor rollout
kubectl rollout status deployment/fullstack-backend -n backend
kubectl rollout status deployment/fullstack-frontend -n frontend

# Rollback if needed
helm rollback fullstack 1
```

### 6. Uninstall

```bash
helm uninstall fullstack

# Note: PersistentVolumeClaims are NOT automatically deleted
# To remove PostgreSQL data PVC:
kubectl delete pvc -l app.kubernetes.io/instance=fullstack -n database
```

---

## Notes

- **Authentication:** This blueprint uses unauthenticated APIs for local development. In production, implement JWT authentication and use an API Gateway/Ingress with JWT validation.
- **Network Policies:** Included in the Helm chart for security hardening. Requires a NetworkPolicy controller (Calico, Cilium, etc.). Disable by deleting NetworkPolicy resources if not supported.
- **Versions:** Always verify and use the latest stable releases for Kubernetes, Go, Angular, Docker, Node.js, and PostgreSQL.

---

**Maintained by:** Adrian B
**Last Updated:** 2026-02-09
**Blueprint Version:** 1.0.0
