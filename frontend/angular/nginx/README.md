# Nginx Configuration for Angular Frontend

## Overview

This directory contains nginx configuration files for serving the Angular application in both Docker and Devspace deployments.

## Architecture

### Two-Tier Nginx Configuration

Nginx uses a two-tier configuration structure:

- **Main config (`/etc/nginx/nginx.conf`)**: Provided by the base `nginx:alpine` image. Contains global settings like `user`, `worker_processes`, `events`, and `http` blocks. This file includes all configurations from `/etc/nginx/conf.d/*.conf`.

- **Server block (`/etc/nginx/conf.d/default.conf`)**: Our custom configuration that defines how the Angular app is served. This includes:
  - Static file serving from `/usr/share/nginx/html`
  - Angular SPA fallback routing (`try_files $uri $uri/ /index.html`)
  - Security headers (X-Frame-Options, X-Content-Type-Options, X-XSS-Protection)
  - Gzip compression
  - Cache headers for static assets
  - Health check endpoint at `/healthz`

## Files

- **`default.conf`**: The server block configuration. This is the only file that needs to be customized for the Angular app.
- **`nginx.conf`** (deprecated): Previously used but misleading name. This file is no longer used. Use `default.conf` instead.

## Deployment

### Docker

The Dockerfile copies `default.conf` to `/etc/nginx/conf.d/default.conf` inside the container:

```dockerfile
COPY default.conf /etc/nginx/conf.d/default.conf
```

### Devspace

The `devspace.yaml` syncs `default.conf` directly to the running pod:

```yaml
- path: ./frontend/angular/dev-nginx/default.conf:/etc/nginx/conf.d/default.conf
```

## Best Practices

✅ **Do NOT override the main `nginx.conf`** unless you need to change global settings (worker processes, log format, etc.). For typical SPA deployments, only `default.conf` is needed.

✅ **Keep server blocks in `conf.d/`** - This follows nginx's modular design and allows the base image to manage global settings.

✅ **Use `default.conf` naming** - This is the conventional name nginx includes by default.

## Why This Approach?

The base nginx:alpine image already provides an optimized `nginx.conf` that:
- Sets appropriate worker processes (auto)
- Configures logging and error handling
- Includes the `conf.d/` directory

By only providing a server block in `conf.d/default.conf`, we:
- Reduce maintenance burden
- Can upgrade the base image without merging config changes
- Follow Docker and nginx community best practices
- Keep our configuration simple and focused

## References

- [Nginx Configuration Structure](https://nginx.org/en/docs/ngx_core_module.html)
- [Docker Nginx Best Practices](https://docs.docker.com/engine/examples/nginx/)
- [Angular Deployment with Nginx](https://angular.io/guide/deployment)

## TODO: Migrate to Helm/Kubernetes ConfigMap

For production Kubernetes deployments using Helm charts, the nginx configuration should be managed via a ConfigMap rather than static files. The Helm chart in `deployment/fullstack-blueprint-app/` contains a ConfigMap template that can embed the `default.conf` content directly:

```yaml
# templates/frontend-nginx-configmap.yaml (or similar)
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ include "fullstack-blueprint.fullname" . }}-frontend-nginx-config
data:
  default.conf: |
    server {
        listen 80;
        server_name localhost;
        root /usr/share/nginx/html;
        index index.html;

        # Enable gzip compression
        gzip on;
        gzip_vary on;
        gzip_min_length 1024;
        gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;

        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;

        # Serve static files with caching
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
            try_files $uri =404;
        }

        # Angular SPA fallback - all other routes return index.html
        location / {
            try_files $uri $uri/ /index.html;
        }

        # Health check endpoint
        location /healthz {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }
```

This approach:
- Keeps configuration versioned with the Helm chart
- Allows environment-specific overrides via Helm values
- Follows Kubernetes best practices for immutable infrastructure
- Eliminates the need for Devspace sync or Docker COPY of config files

**Migration steps:**
1. Move the `default.conf` content into the Helm chart's ConfigMap template
2. Update the frontend nginx deployment to mount the ConfigMap as a volume at `/etc/nginx/conf.d/default.conf`
3. Remove the `COPY` directive from Dockerfile (or keep for standalone Docker builds)
4. Update devspace.yaml to use Helm instead of direct file sync for production-like deployments
