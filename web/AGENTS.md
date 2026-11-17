# SvelteKit web app

* Use Svelte 5 runes and strict TypeScript. Reuse existing components and styles. Put rules shared by multiple routes in `src/lib`.
* Consume server DTOs through an SDK generated from [the OpenAPI spec](../server/doc/openapi.json) when adding API calls. Keep request details out of components. Handle loading, empty, and error states when rendering asynchronous data.
* Components and routes should work with semantic elements and keyboard navigation. Test what the user can see and do through roles, labels, and text.
* Preserve the static SPA build. The production Go server serves the output and handles requests under `/api`.
* Run `mise run web:check` for a frozen lockfile install, Svelte checks, and a production build.
