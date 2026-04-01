# Stockyard Pasture

**Bookmark manager with full-text search — save links, tag them, search the page content**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9200:9200 -v pasture_data:/data ghcr.io/stockyard-dev/stockyard-pasture
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9200` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9200` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `PASTURE_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 100 bookmarks | Unlimited bookmarks |
| Price | Free | $2.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Operations & Teams

## License

Apache 2.0
