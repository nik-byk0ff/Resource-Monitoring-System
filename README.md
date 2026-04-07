# Resource Monitoring System

A scalable, containerized, and lightweight resource monitoring system utilizing a Go backend, Vue 3 frontend, and PostgreSQL + TimescaleDB for time-series data storage.

## Architecture

- **Backend (Go)**: Built with `chi` for routing and `pgxpool` for high-performance PostgreSQL interaction. Exposes a secure REST API protected by JWT and Role-Based Access Control (RBAC).
- **Frontend (Vue 3)**: A responsive, Brutalist-style dashboard powered by Vite. Communicates with the backend API to fetch and stream near real-time metrics.
- **Database (TimescaleDB)**: A time-series optimized extension on top of standard PostgreSQL, utilizing **hypertables** for scalable and fast metric inserts/queries.
- **Agent Component**: A lightweight worker designed to gather local resource statistics and stream them directly into the central API.

---

## How Monitoring Works

This system utilizes a **Push-Based Architecture** (similar to Telegraf/InfluxDB or Datadog, as opposed to Prometheus' pull mechanism). 

1. **Authentication:** The Agent initializes by sending a login request `POST /api/auth/login` to the central backend using agent credentials, receiving a securely signed JSON Web Token (JWT).
2. **Metrics Collection:** The Agent calculates the host's system metrics (CPU utilization, Memory usage, Disk capacity). *Note: The current mock agent generates random stats for testing purposes, but can easily be replaced with `gopsutil` calls.*
3. **Data Ingestion:** The Agent periodically executes a `POST /api/metrics` carrying its `NodeName`, `Time`, and numeric statistics.
4. **Storage:** The backend validates the JWT and payload, writing the metrics to the TimescaleDB tracking hypertable.
5. **Visualization:** Operators log into the Frontend dashboard which automatically polls the system metrics and renders them by node.

---

## How to Set Monitoring on Other APIs / Servers

**Does it automatically monitor other APIs?**  
No. Because the architecture relies on *data push*, your other systems will not be recorded magically. They must actively dispatch data to this monitoring backend.

To monitor an external API, Node, or Server, you must use one of the following integration methods:

### Option 1: Deploy the Provided Agent (Recommended)
You can compile and deploy the `agent` application alongside your other API logic (often referred to as a **Sidecar** or **DaemonSet** pattern). 
The agent runs in the background and requires only environment variables to target your central registry.
```bash
# On the external server you wish to monitor
export API_URL="https://your-central-backend.com/api"
./agent
```

### Option 2: Implement a Custom Push Script
If you prefer not to use the Go agent, any application can hook into the system by sending simple REST commands.
For example, inside a Python API or Node.js application, you can orchestrate a periodic background task to collect hardware footprints and dispatch the HTTP request manually.

**Example Payload to `POST /api/metrics`**:
```json
{
  "node_name": "payment-api-server-1",
  "cpu_usage": 45.2,
  "memory_usage": 60.0,
  "disk_usage": 72.1
}
```
*(Requires `Authorization: Bearer <token>` in the HTTP headers)*

---

## Setup & Installation Instructions

### 1. Local Development via Docker Compose

This is the fastest method to execute and modify the application suite locally.

1. **Environment Setup:** Ensure the provided `.env` variables have been configured properly (e.g. database credentials, JWT secret).
2. **Run Containers:** Execute docker-compose from the project root.
    ```bash
    docker-compose up --build
    ```
    This automatically builds the Frontend (Node), Backend (Go), and provisions the TimescaleDB image, mounting the schema dump from `/db/init`.
3. **Verify:**
   - **Frontend GUI:** `http://localhost`
   - **Backend API:** `http://localhost:8080/api`

### 2. Production Deployment via Kubernetes (Helm)

The project includes a Helm chart configured to initialize the environment in a scalable Kubernetes cluster.

1. **Verify your Cluster:** Ensure `kubectl` is pointed to your active cluster.
2. **Review Values:** Look into `/k8s/helm/monitoring/values.yaml` to substitute deployment behaviors (such as image registry references, service node types or scaling limits).
3. **Deploy:**
   Navigate into the Helm directory and execute installation:
   ```bash
   cd k8s/helm
   helm install monitoring-release ./monitoring
   ```
4. **Scale:** Once the deployments stand up, you can easily control identical pods:
   ```bash
   kubectl scale deployment monitoring-release-backend --replicas=3
   ```

---

## Security

- **Injection Control:** Implemented via prepared statement bindings in Go `pgxpool`. 
- **Broken Auth:** Handled via properly salted `bcrypt` storage hashes and timed JWT token evaluations natively baked into Chi-Router Middleware.
- **Identification Flaws:** Basic structured validation across all API boundaries preventing payload manipulations.
