# =============================================================================
# Dev-only Dockerfile
#   docker build --target dev -t mini-app-web:dev .
#
# There is no production stage here: the production build is a static SPA
# (adapter-static) built and served by the server image — see
# ../server/Dockerfile (web-build stage)
# =============================================================================

# =============================================================================
# Stage: base — node with pnpm activated
# =============================================================================
FROM node:24-alpine AS base

RUN corepack enable && corepack prepare pnpm@10.33.0 --activate

WORKDIR /app

# =============================================================================
# Stage: dev — vite dev server; bind-mount the source dir over /app
# (node_modules and .svelte-kit are expected to be masked by named volumes,
# so dependencies are installed at container start)
# =============================================================================
FROM base AS dev

ENV NODE_ENV=development

EXPOSE 3000

CMD ["sh", "-c", "pnpm install && exec pnpm dev --host 0.0.0.0 --port 3000"]
