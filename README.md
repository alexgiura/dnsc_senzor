# DNSC Senzor API

API HTTP pentru ingestia alertelor de rețea (JSON) și persistarea lor în fișier JSON Lines (`.jsonl`).

## Pornirea proiectului

### Varianta 1: Go local

Din directorul `backend`:

```bash
cd backend
go run ./cmd
```

Serverul ascultă pe portul din variabila `SERVER_PORT` (implicit **8080**), conform fișierului `.env` din rădăcina repository-ului.

### Varianta 2: Docker Compose

Din rădăcina repository-ului (unde se află `docker-compose.yml`):

```bash
docker compose up --build
```

Portul host ↔ container este `${SERVER_PORT:-8080}` (vezi `.env`).

## Unde se salvează datele

| Mod de rulare | Fișier / locație |
|----------------|------------------|
| **Go local** | `backend/data/network_alerts.jsonl` — calea relativă `data/network_alerts.jsonl` este rezolvată față de directorul `backend/` (vezi `NETWORK_ALERTS_STORAGE_PATH` în `.env`). |
| **Docker** | În container: `/root/data/network_alerts.jsonl`, mapat la **`./backend/data`** pe host (volum din `docker-compose.yml`). |

Fiecare alertă primită la `POST` este adăugată ca **o linie JSON** (format JSON Lines).

## Documentație API (Swagger)

După ce serverul rulează, deschide în browser:

**http://localhost:8080/swagger**

Acolo poți consulta schema OpenAPI și poți trimite cereri de probă (inclusiv `POST /api/v1/network-alerts`). Specificația brută YAML este disponibilă la:

**http://localhost:8080/openapi.yaml**

## Exemplu: POST alertă de rețea

**URL:** `http://localhost:8080/api/v1/network-alerts`  
**Header:** `Content-Type: application/json`  
**Metodă:** `POST`

**Corp (body) exemplu:**

```json
{
  "agent_id": "test-postman",
  "exported_at": "2026-04-01T12:00:00Z",
  "event": {
    "timestamp": "2026-04-01T12:00:00Z",
    "protocol": "TCP",
    "src_ip": "10.0.0.2",
    "src_port": 443,
    "dst_ip": "10.0.0.2",
    "dst_port": 80,
    "watchlist_match": "src",
    "direction": "outbound",
    "tcp_flags": "S",
    "packet_size": 100
  }
}
```

**Răspuns la succes:** `201 Created` cu corp `{"status":"created"}`.

**Cu curl:**

```bash
curl -sS -X POST "http://localhost:8080/api/v1/network-alerts" \
  -H "Content-Type: application/json" \
  -d '{
  "agent_id": "test-postman",
  "exported_at": "2026-04-01T12:00:00Z",
  "event": {
    "timestamp": "2026-04-01T12:00:00Z",
    "protocol": "TCP",
    "src_ip": "10.0.0.2",
    "src_port": 443,
    "dst_ip": "10.0.0.2",
    "dst_port": 80,
    "watchlist_match": "src",
    "direction": "outbound",
    "tcp_flags": "S",
    "packet_size": 100
  }
}'
```

## Health check

- `GET http://localhost:8080/healthz` sau `GET http://localhost:8080/health` → răspuns text `OK`.
