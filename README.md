# 🔍 Host Diff Tool (Go + React Monorepo)

The **Host Diff Tool** helps security engineers and network admins track changes on a host over time.  
It compares snapshots of a host’s open services and highlights what changed (new ports, closed ports, or updated software versions).

This repository is a **monorepo** containing both:
- **API** → Go backend
- **UI** → React frontend

---

## 📂 Repository Structure

```
host-diff-tool/
├── packages/
│   ├── api/    # Go backend
│   ├── ui/     # React frontend
│
├── shared/     # Optional shared schemas/utilities
├── docker-compose.yml
└── README.md
```

---

## ✨ Features

- Upload JSON snapshots of a host
- Store and view snapshot history
- Compare any two snapshots to see:
  - ✅ New services
  - ❌ Removed services
  - 🔄 Changed services (e.g., software upgrades)

---

## 📂 Snapshot Format

Snapshots are JSON files with this structure:

```json
{
  "timestamp": "2025-09-10T03:00:00Z",
  "ip": "203.0.113.45",
  "services": [
    {
      "port": 22,
      "protocol": "SSH",
      "software": {
        "vendor": "openssh",
        "product": "openssh",
        "version": "8.2p1"
      }
    }
  ]
}
```

---

## 📊 Example Diff Output

```json
{
  "new_services": [
    { "port": 443, "protocol": "HTTPS", "software": { "vendor": "nginx", "version": "1.24.0" } }
  ],
  "removed_services": [
    { "port": 21, "protocol": "FTP" }
  ],
  "changed_services": [
    { "port": 22, "protocol": "SSH", "old_version": "8.2p1", "new_version": "9.0p1" }
  ]
}
```

---

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/) (>= 1.21)
- [Node.js](https://nodejs.org/) (>= 18)
- (Optional) [Docker](https://www.docker.com/) for containerized setup

---

### Clone & Install

```bash
git clone https://github.com/yourusername/host-diff-tool.git
cd host-diff-tool
```

---

### Running the Backend (Go API)

```bash
cd packages/api
go run main.go
```

Runs on [http://localhost:4000](http://localhost:4000)

---

### Running the Frontend (React UI)

```bash
cd packages/ui
npm install
npm run dev
```

Runs on [http://localhost:3000](http://localhost:3000)

---

### Running with Docker Compose

```bash
docker-compose up --build
```

- API → [http://localhost:4000](http://localhost:4000)  
- UI → [http://localhost:3000](http://localhost:3000)  

---

## 📡 API Endpoints

### `POST /snapshots`
Upload a new snapshot file.

**Request:**
```http
POST /snapshots
Content-Type: application/json
```

```json
{
  "timestamp": "2025-09-10T03:00:00Z",
  "ip": "203.0.113.45",
  "services": [
    {
      "port": 22,
      "protocol": "SSH",
      "software": {
        "vendor": "openssh",
        "product": "openssh",
        "version": "8.2p1"
      }
    }
  ]
}
```

**Response:**
```json
{ "message": "Snapshot stored successfully" }
```

---

### `GET /snapshots/{host}`
Retrieve all snapshots for a given host.

**Request:**
```http
GET /snapshots/203.0.113.45
```

**Response:**
```json
[
  {
    "timestamp": "2025-09-10T03:00:00Z",
    "ip": "203.0.113.45",
    "services": [ ... ]
  },
  {
    "timestamp": "2025-09-12T03:00:00Z",
    "ip": "203.0.113.45",
    "services": [ ... ]
  }
]
```

---

### `GET /diff?host=...&from=...&to=...`
Compare two snapshots of the same host.

**Request:**
```http
GET /diff?host=203.0.113.45&from=2025-09-10T03:00:00Z&to=2025-09-12T03:00:00Z
```

**Response:**
```json
{
  "new_services": [
    { "port": 443, "protocol": "HTTPS", "software": { "vendor": "nginx", "version": "1.24.0" } }
  ],
  "removed_services": [
    { "port": 21, "protocol": "FTP" }
  ],
  "changed_services": [
    { "port": 22, "protocol": "SSH", "old_version": "8.2p1", "new_version": "9.0p1" }
  ]
}
```

---

## 🏗 Tech Stack

- **Backend** → Go (net/http or Gin/Echo)
- **Frontend** → React + Vite
- **Storage** → JSON files (MVP), can upgrade to Postgres/Mongo
- **Monorepo** → `packages/` directory convention

---

## 🛠 Roadmap

- [ ] Basic Go API for snapshots/diffs  
- [ ] File-based snapshot storage  
- [ ] React UI for uploads + diffs  
- [ ] Database backend (Postgres/Mongo)  
- [ ] Export diffs (PDF/HTML)  
- [ ] Alerts for critical changes  

---

## 📜 License

MIT
