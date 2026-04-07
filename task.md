# Task List: Resource Monitoring System

## 1. Project Initialization
- [x] Create monorepo directory structure (`/backend`, `/frontend`, `/k8s/helm`).
- [x] Create base `docker-compose.yml` for local development (TimescaleDB, Backend, Frontend).
- [x] Use .env files for configuration.

## 2. Database Setup (TimescaleDB)
- [x] Define PostgreSQL/TimescaleDB schema for Users (roles).
- [x] Define TimescaleDB schema for Metrics (hypertable).

## 3. Backend Implementation (Go)
- [x] Initialize Go module and install dependencies (Chi/Gin, pgx, jwt, validator).
- [x] Implement database connection and migrations.
- [x] Implement Authentication & Authorization (JWT, BCrypt, User vs Admin roles).
- [x] Implement Metrics API (Receiving metrics from agents).
- [x] Implement System Metrics Collector (Mock or real system stats).
- [x] Apply OWASP Security Practices (Validation, Escaping, Auth).
- [x] Write Backend Unit & Integration Tests.

## 4. Frontend Implementation (Vue 3)
- [x] Initialize Vue 3 + Vite project.
- [x] Setup Brutalist Design System (CSS variables, reset, core components).
- [x] Implement Authentication views (Login, Register).
- [x] Implement Admin/User Dashboard views (Fetching and displaying metrics).
- [x] Secure routes based on user roles.
- [x] Write Frontend Unit tests (Vitest).

## 5. Containerization & Deployment
- [x] Write `Dockerfile` for Backend (multi-stage, scratch/alpine).
- [x] Write `Dockerfile` for Frontend (Nginx static serving).
- [x] Create Helm Chart for deploying the entire application to Kubernetes.
