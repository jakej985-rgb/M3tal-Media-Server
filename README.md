# M3TAL Media Server

Autonomous, self-healing media automation platform.

## 🏗 Architecture (Source Services)
| Service | Type | Path |
|---------|------|------|
| api | CLI binary | [cmd/api](./cmd/api) |
| m3tal | CLI binary | [cmd/m3tal](./cmd/m3tal) |
| api | Go package | [internal/api](./internal/api) |
| auth | Go package | [internal/auth](./internal/auth) |
| config | Go package | [internal/config](./internal/config) |
| containers | Go package | [internal/containers](./internal/containers) |
| engine | Go package | [internal/engine](./internal/engine) |
| health | Go package | [internal/health](./internal/health) |
| orchestrator | Go package | [internal/orchestrator](./internal/orchestrator) |
| preflight | Go package | [internal/preflight](./internal/preflight) |
| store | Go package | [internal/store](./internal/store) |
| system | Go package | [internal/system](./internal/system) |
| dashboard | Deploy artifact | [deploy/dashboard](./deploy/dashboard) |
| stack | Deploy artifact | [deploy/stack](./deploy/stack) |
| github.com/jakej985-rgb/m3tal-core | Go module | [go.mod](./go.mod) |

## 🐳 Infrastructure Stacks
### M3tal Stack
- **m3tal-dashboard** (ghcr.io/jakej985-rgb/m3tal-godash:debug)

### Routing Stack
- **traefik** (traefik:latest)
  - Ports: 127.0.0.1:8081:8080
- **cloudflared** (cloudflare/cloudflared:latest)

## ⚙️ Environment Configuration
| Variable | Default |
|----------|---------|
| DASHBOARD\_PORT | `8082` |
| DASHBOARD\_EXPOSE\_MODE | `local` |
| HTTP\_PORT | `8080` |
| STATE\_DIR | `./state` |
| LOG\_LEVEL | `info` |
| DASHBOARD\_SECRET | `change\_me\_immediately` |
| API\_TOKEN | `change\_me\_api\_token` |
| ADMIN\_PASSWORD | `admin\_pass` |
| NETWORK\_NAME | `m3tal` |
| LOCAL\_IP | `127.0.0.1` |
| DOMAIN | `localhost` |
| VPN\_USER | `user` |
| VPN\_PASSWORD | `password` |
| BASE\_STORAGE\_PATH | `./data` |
| MEDIA\_PATH | `./data/media` |
| CONFIG\_PATH | `./data/config` |
| DOWNLOADS\_PATH | `./data/downloads` |
| PUID | `1000` |
| PGID | `1000` |
| TZ | `America/Denver` |
| TRAEFIK\_WEB\_PORT | `80` |
| TRAEFIK\_WEBHTTPS\_PORT | `443` |
| TRAEFIK\_DASHBOARD\_PORT | `8080` |
| DEBUG\_MODE | `false` |
| METRICS\_ENABLED | `true` |

## 🚀 Deployment

```bash
python m3tal.py init
```
