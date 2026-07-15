# =============================================================================
# Multi-target Dockerfile
#   development: docker build --target dev -t mini-app:dev .
#   production:  docker build --build-context web=../web -t mini-app .
#                (prod is the final/default stage; it bundles the static web
#                 build from the named "web" context — docker-compose.yml
#                 wires it via additional_contexts)
# =============================================================================

# =============================================================================
# Stage: base — shared Go toolchain with warmed module cache
# =============================================================================
FROM golang:1.26-alpine AS base

WORKDIR /app

# Bake modules into the image (no cache mount) so the dev target
# doesn't re-download them at runtime
COPY go.mod go.sum ./
RUN go mod download

# =============================================================================
# Stage: dev — hot-reload via air; bind-mount the source dir over /app
#   docker run -v .:/app -p 3001:3001 mini-app:dev
# =============================================================================
FROM base AS dev

# Version matches the air pin in apps/templates/mise.toml
RUN go install github.com/air-verse/air@v1.65.3

# uid/gid 1000 matches the default first Linux user, so files air writes into
# the bind-mounted source dir (./tmp) stay owned by the host user. /go is
# chowned so the app user can use the module cache and install new modules
RUN addgroup -S -g 1000 app && adduser -S -u 1000 -G app app && \
    chown -R app:app /go

USER app

ENV APP_ENV=development

EXPOSE 3001

CMD ["air"]

# =============================================================================
# Stage: build — compile a static binary (pure-Go sqlite, no CGO)
# =============================================================================
FROM base AS build

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/mini-app .

# =============================================================================
# Stage: web-build — static SPA build (adapter-static) from the "web" context
# Only the needed sources are copied explicitly, so the host's node_modules
# and .svelte-kit never leak into the build
# =============================================================================
FROM node:24-alpine AS web-build

RUN corepack enable && corepack prepare pnpm@10.33.0 --activate

WORKDIR /web

COPY --from=web package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc ./
RUN pnpm install --frozen-lockfile

COPY --from=web vite.config.ts tsconfig.json ./
COPY --from=web src ./src
COPY --from=web static ./static

RUN pnpm build

# =============================================================================
# Stage: prod — slim runtime; migrations are embedded and run at startup,
# the web UI is served as static files from ./web
# =============================================================================
FROM alpine:3.22 AS prod

RUN apk add --no-cache ca-certificates tzdata tini

WORKDIR /app

COPY --from=build /out/mini-app ./mini-app
COPY --from=web-build /web/build ./web

# Config resolves the database to appdata/database.db relative to the workdir
RUN mkdir -p /app/appdata

ENV APP_ENV=production \
    SERVER_ADDR=0.0.0.0:3001

EXPOSE 3001

ENTRYPOINT ["tini", "--", "./mini-app"]
