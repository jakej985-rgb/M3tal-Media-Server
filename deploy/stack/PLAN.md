# m3tal-stack Plan (Docker Infrastructure)

## Purpose
Runs all containers.

NO logic allowed.

---

## Goals
- Clean Docker Compose setup
- Consistent volumes
- No cross-filesystem moves

---

## Structure

```
m3tal-stack/
├── *-compose.yml
├── .env
├── services/
└── volumes/
```

---

## Tasks

### 1. Standardize volumes

ALL containers MUST use:

```yaml
volumes:
  - /mnt:/mnt
```

---

### 2. Fix media paths

Inside apps:

* Radarr → `/mnt/media/movies`
* Sonarr → `/mnt/media/tv`
* Downloads → `/mnt/downloads`

---

### 3. Networking

```yaml
networks:
  m3tal:
    driver: bridge
```

All containers join this network.

---

### 4. Remove logic

❌ No scripts
❌ No agents
❌ No automation

---

### 5. Validate imports

* No copying between folders
* Moves should be instant (same filesystem)

---

## Done When

* Containers see same `/mnt`
* Imports are fast (no copy)
* Stack runs independently
