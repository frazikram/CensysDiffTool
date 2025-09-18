# Host Diff Tool

In security and network monitoring, understanding change over time is critical.  
A port that was closed yesterday might be open today; a new vulnerability might appear on a service that was previously considered safe.

The **Host Diff Tool** provides a simple way to ingest host snapshots at different points in time, compare them, and highlight meaningful changes.

---

## Features

- Upload JSON snapshots of host data.
- View a history of snapshots for a given host.
- Compare any two snapshots to see changes in:
  - Ports
  - Services
  - Software versions
  - Vulnerabilities
  - TLS configurations
- Generate structured diff reports.
- Web-based UI for interaction.

---

## Project Structure

```
/packages
  /api     → Go backend (REST API + snapshot storage)
  /ui      → React/TypeScript frontend
```

- **Backend (API)**: Handles snapshot ingestion, storage (SQLite), and diff computation.
- **Frontend (UI)**: Provides a simple interface for uploading, viewing, and comparing snapshots.

---

## Running the Project

### Run Manually 

#### Backend (Go API)
```bash
cd packages/api
go run cmd/server/main.go
```

Runs on **http://localhost:8080**

#### Frontend (React/TS UI)
```bash
cd packages/ui
npm install
npm start
```

Runs on **http://localhost:5173**

---

## API Endpoints

- `GET /health` – Health check
- `POST /snapshots` – Upload a new snapshot (JSON)
- `GET /snapshots/{ip}` – Get all snapshots for a host
- `POST /compare` – Compare two snapshots and return a diff

---

## Assumptions

- Each snapshot JSON is well-formed and follows the agreed schema.
- Host identity is based on **IP address**.
- Snapshots are relatively small and can be stored in SQLite (no distributed DB needed at this scale).
- No authentication required (development/demo context).
- TLS certificate data and vulnerabilities are stored as provided, without external validation.

---

## Testing

### Manual Testing
- Start both API and UI.  
- Upload multiple snapshot JSON files via UI.  
- Confirm that:
  - Snapshots appear in history.  
  - Selecting two snapshots shows differences clearly.  
  - Invalid JSON upload fails gracefully.  

### Automated Testing
- **API**: Run unit tests  
  ```bash
  cd packages/api
  go test ./...
  ```
- **UI**: Currently no automated tests are implemented.

---

## AI Techniques

- **Diff Computation Logic**:  
  Inspired by change detection algorithms, we implemented a structured comparison of snapshots by mapping services to unique keys (`port/protocol`).  
  - Detects **added, removed, and changed** services.  
  - Tracks software version drift, TLS differences, and vulnerability changes.  

This deterministic approach is lightweight but borrows from principles used in AI for **state comparison and anomaly detection**.

---

## Future Enhancements

If given more time, we would extend the system with:

- **Authentication & Multi-User Support**: Secure uploads and user-specific histories.  
- **Search & Filtering**: Query snapshots by port, protocol, or vulnerability.  
- **Visualization**: Graphs or timelines showing service history.  
- **AI-powered Anomaly Detection**: Use ML to detect unusual changes (e.g., suspicious new services).  
- **Scalability**: Migrate storage from SQLite to Postgres or a distributed DB.  
- **Integration with Security Feeds**: Auto-enrich snapshots with live CVE data.  
- **Exportable Reports**: PDF/CSV export of diffs for compliance or audits.  

---

## Getting Started

1. Clone the repo  
   ```bash
   git clone <your_repo_url>
   cd hostdiff
   ```

2. Run with manually as described above.  

3. Upload sample snapshots from the `examples/` folder (or create your own).  

4. Compare two snapshots and view the diff report.



