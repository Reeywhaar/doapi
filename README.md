# doctl

A minimal HTTP server and CLI for managing DigitalOcean DNS records.

## Setup

Create a `.env` file:

```
DO_TOKEN=dop_v1_...
```

## Build & deploy

## Run

```bash
docker run -d --name doapi -p 8080:8080 --env-file .env ghcr.io/reeywhaar/doapi:latest
```

## HTTP API

### List all domains and their records

```
GET /dns
```

```bash
curl http://localhost:8080/dns
```

### Create a DNS record

```
POST /dns
```

```bash
curl -X POST http://localhost:8080/dns \
  -H 'Content-Type: application/json' \
  -d '{"domain":"example.com","type":"TXT","name":"_acme-challenge","data":"abc123","ttl":3600}'
```

| Field    | Description                                            |
| -------- | ------------------------------------------------------ |
| `domain` | Domain name (e.g. `example.com`)                       |
| `type`   | Record type (`TXT`, `A`, `CNAME`, etc.)                |
| `name`   | Subdomain or `@` for root                              |
| `data`   | Record value                                           |
| `ttl`    | Time-to-live in seconds (optional, defaults to `1800`) |

Returns the created record JSON (status 201).

### Delete a DNS record

```
DELETE /dns
```

```bash
curl -X DELETE http://localhost:8080/dns \
  -H 'Content-Type: application/json' \
  -d '{"domain":"example.com","id":12345678}'
```

The record `id` is returned by the create or list endpoints. Returns 204 on success.

## CLI

```bash
# List all domains and their records
docker exec doapi /doapi dns list

# Create a DNS record (ttl optional, defaults to 1800)
docker exec doapi /doapi dns create example.com TXT _acme-challenge abc123
docker exec doapi /doapi dns create example.com TXT _acme-challenge abc123 30

# Delete a record by ID
docker exec doapi /doapi dns delete example.com 12345678
```

## Project structure

```
├── main.go               # entry point — wires client → api or cli
├── internal/client.go    # DO API client and types
├── api/handler.go        # HTTP handlers
├── cli/cli.go            # CLI commands
├── Dockerfile
└── build.sh
```
