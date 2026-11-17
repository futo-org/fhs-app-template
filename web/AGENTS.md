# SvelteKit rules · extends [core](../docs/agents/core.md)

Governs `web/`. Use Svelte 5, SvelteKit, Tailwind and strict TypeScript. The production Go server serves the static SPA and handles `/api`.

## Skills to load

When `modern-web-guidance` is available, load it before writing UI, CSS or browser API code. Use its search and retrieval commands to read guidance for the change.

## Vocabulary

| Core | Here |
| --- | --- |
| transport | components, routes and dialogs |
| logic | plain TypeScript services and Svelte rune state in managers |
| adapter | a generated API client and helpers around browser APIs |
| shape | SDK response types and a `type Props` for each component |

## Layout

Grow these directories under `src/lib/` as the app needs them.

```text
components/    presentational components with typed props
modals/        dialog components with an onClose result
services/      plain TypeScript rules and action handlers
managers/      .svelte.ts modules for shared state using runes
utils/         pure shared functions and browser adapters
constants.ts   shared constants named per core
```

## Idioms

Use Svelte 5 runes for new state and components.

```svelte
<script lang="ts">
  type Props = {
    albumIds: number[];
    selectedIds?: number[];
  };

  let { albumIds, selectedIds = $bindable([]) }: Props = $props();
  const selected = $derived(new Set(selectedIds));
  const unselectedIds = $derived(albumIds.filter((id) => !selected.has(id)));
</script>
```

* Derive collections with `$derived`. Use immutable updates for arrays. Keep values on screen during refresh and ignore responses for superseded requests.
* Define shared actions once so menus, buttons and shortcuts use the same label, availability and handler.
* In localized apps, use the translation catalog for user facing copy and follow core's key naming rules.
* Generate clients from [the OpenAPI spec](../server/doc/openapi.json) when the app makes API calls. Use the generated DTOs and keep request details in services.
* Reuse existing components, token values and icons. Use Tailwind utilities in the markup. In flex layouts, use `min-h-0 grow` for a body that needs to scroll.
* Return a dialog result through `onClose(result)`. When the app uses a modal manager, await its result.
* Keep the SPA build compatible with the static adapter. Put backend routes in the Go API.

## Web platform

* Use native `<dialog>` with `showModal()` for modals and `popover` for light dismiss menus. Check focus entry, Escape and return to the trigger. Use `@starting-style` and `transition-behavior: allow-discrete` when animating open and close.
* When the app supports light and dark themes, set `color-scheme` and use shared tokens. Allow the system preference. Sync a stored theme preference across tabs through the `storage` event.
* Render untrusted text as text. When rich HTML is required, use `setHTML()` where supported and a maintained sanitizer for fallback. Apply the content security policy at the server that serves the built app and verify the production build's scripts and styles load.
* Use semantic HTML, keyboard navigation and accessible labels. Give loading, empty and error states explicit views.

## Errors

Handle errors at the UI boundary for every async user action. Share error formatting, keep feedback specific to the action and suppress cancellation errors. Report a failed refresh without clearing data already on screen.

## Tests

Use the app's installed harness. Vitest fits services and managers, Testing Library fits component behavior, and Playwright fits smoke tests in a running app.

| Tier | What to exercise |
| --- | --- |
| Unit | shared rules and state with the SDK mocked |
| Transport | rendered components, user actions and emitted callbacks |
| Integration | queries and data correctness in the server's SQLite suite |
| End to end | visible flows against the dev stack |

Name suites for the component or service and write each case as one behavior. Assert SDK arguments and calls prevented by a guard. In browser tests, assert what the user can see and do through roles, labels and text.

## Checklist

Run `mise run web:check` for the frozen lockfile install, Svelte checks and production build. Run any app tests that cover the change. Verify user flows against the running app.

## Derived artifacts

Run `mise run open-api` when the server DTOs or routes change. Update SDK consumers in the same change when the app uses a generated client. Update the translation catalog when localized copy changes.

## Overrides

Use the generated client directly from services as the typed API adapter. The server integration tier owns datastore correctness.
