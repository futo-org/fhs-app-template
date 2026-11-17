# Agent rules: core

These rules apply across the FUTO Home app template. Stack specifics (layout, syntax, harness names, checklist commands) live in each component's `AGENTS.md`, listed in the root [AGENTS.md](../../AGENTS.md). Read the component file before writing code there. Where they differ, the component file wins, then the root file, then this one.

## How to work

1. Read the task and the surrounding code before editing. State the one idea this change delivers. A change with two ideas is two changes.
2. For any change touching more than one layer, write a two to five line plan (files, order) before the first edit.
3. If adding the feature feels awkward, the structure is wrong. Refactor first as its own change, then add the feature to the new shape. Refactor and feature never share a change.
4. Test first when the change needs a test (Tests, When to test): write it, run it, watch it fail on the current code, then write the code. A bug test drives the code as it is before the fix; if it needs a seam that does not exist, add the seam first as its own refactor. Code that only moves keeps its existing tests green and gets no new ones. Report the red then green run.
5. Prefer the change that deletes more than it adds. A net negative diff is the normal outcome of a refactor.
6. Run the affected tests while iterating. Run the full component checklist (lint, format, typecheck, tests, regenerated artifacts) once before reporting. Done means it is green and every scenario the change hits under Tests, When to test, has a test.
7. Do not add code comments. Make the code self-explanatory through naming and structure instead.
8. A failure you did not cause is reported, not fixed. Fixing it is a second idea.
9. Stop and ask only when the code cannot settle the decision: a new dependency, a schema change, a public API or breaking change. Everything else, decide, record the assumption in the PR, continue.
10. A user term that maps to stored state ("linked", "signed in") gets a definition in that state, confirmed by the developer. An assumption that changes what the user sees is a question, not an assumption. A choice question lists every state the code distinguishes.
11. Adding a dependency is an ask, not a decision. Prefer the standard library and what the project already depends on.
12. Use the vocabulary in the component's `CONTEXT.md` when it exists; when a term is settled in discussion, record it there.
13. Report tersely: what changed, what was verified, what is left. Verified means run where the behavior lives (device, real server); otherwise list the untested states. A UI change report includes a screenshot from the running app, before and after. Reasoning lives in names, types, structure, and tests, never in prose.

## Shape

- **Boundaries are sacred.** Business logic never touches a database client, filesystem, network, HTTP response, or external SDK directly. Each technology sits behind its own adapter, one technology per adapter. This is what makes unit tests cheap.
- Three layers: transport (controllers, handlers, UI) holds zero logic and maps in and out; logic holds the rules; adapters hold the I/O.
- Types flow from the data layer outward. Write only the types the compiler cannot derive.
- Shared helpers are pure, dependency free, generic: take thunks and options objects, return small result objects.
- Adapter inputs are canonical: one field per concept. Logic resolves client input against preferences and policy before the call, so an adapter never receives both `type` and `types` for the same column.
- Name a magic number or literal list (file extensions, mime types) only when the name earns its place: it is used in three or more places, counting tests and other constants, or two places apart must stay equal (a height another widget lays out around, a part of a sum). Names live in one constants file, never a module local constant. Anything else stays a literal where it is used; a width and height side by side need no name. A constant built from others is inlined as its expression, not its value. Membership lists are spelled out as a Set, not a regex. Enums are string valued, one home per project or package.

## Rule of three

A pattern seen three times becomes a shared helper in the same change, every call site replaced. If a call site resists the helper, the helper is wrong; adjust the helper, never special case the call site. Below three sightings, inline: a helper with one caller is premature. Replace on sight:

- hand rolled loops a helper covers (batching, chunking, find or fail)
- the same private helper in two modules
- scattered type checks → one type guard
- boolean parameters → options object or enum
- dead code, unused exports, unused imports

## Idioms

- Guard clause, early return, blank line after the guard.
- Options objects with defaults destructured at the top; rename with an underscore when normalizing.
- Return a small result object instead of a boolean when the caller branches on it.
- Sets for membership. Concurrent awaits for independent work. A plain loop with `continue` when awaiting inside.
- Higher order helpers take a thunk, not a value.
- When an app uses translations, use `snake_case` keys in its translation catalog. Use the word itself for short labels (`birthdays`, `select_all`) and a feature prefix for sentences (`memories_on_this_day_description`). Reuse the same key across clients.

## Naming

| Pattern                         | Meaning                                      |
| ------------------------------- | -------------------------------------------- |
| `getX`, `getAllX`, `getXById`   | read                                         |
| `findOrFail`                    | read or throw the client error               |
| `requireX`                      | throws when the condition fails              |
| `checkX`, `isX`, `hasX`, `canX` | returns boolean or allowed set, never throws |
| `mapX`                          | entity → response shape                      |
| `handleX`                       | event, job, or UI handler with side effects  |
| `onX`                           | callback or hook                             |
| `newX`                          | test factory                                 |
| `xStub`                         | static test fixture                          |
| `sut`                           | subject under test                           |

Transport method names are the public operation IDs (`createTag`, `getAllTags`): name them for the SDK consumer. Logic method names are short (`create`, `getAll`) because the class or package already gives context.

## Errors

- Fail early with a short noun phrase: `Album not found`, `Email is not available`.
- Not found and no access produce the same client error, so existence never leaks.
- One global handler at the transport edge decides status code, log level, and whether the response is still writable. Logic code throws and moves on; transport code has no try/catch.
- Background work (fire and forget promise, goroutine, job) always runs through a wrapper that logs failures.
- Log through the injected logger: `debug` for expected failures (not found, client aborted), `error` only for the genuinely unknown.

## Risky changes

A change is risky when a mistake lands on users' data or on what they see, and tests on fixtures would pass anyway:

- deletes, retention, pruning, cleanup jobs
- filters and queries that decide what a user sees or owns
- migrations and transforms of stored data
- sync between devices or with the server
- auth, access, and sharing
- release, signing, and update paths

For a risky change:

1. Before editing, state what a wrong result looks like to a user and how many rows, files, or users the path touches at realistic scale.
2. Run the query or job against realistic data (a copy of real data, or a fixture built to production shape and volume) and report the counts before and after. A filter that passes every test can still hide most of a feature on live data.
3. A destructive path gets a dry run or a reversible step (trash, soft delete, backup) unless the project already has one.
4. Review it as big or risky, whatever its size.
5. The report names the risky area, the counts, and what was not checked at scale.

## Tests

A test exists to stop a failure a user would see. Name that failure before writing the test; if you cannot, do not write it.

### When to test

Write a test when the change:

1. Fixes a bug. The test replays the reported scenario and fails on the code before the fix.
2. Adds or changes a rule the user can observe: validation, access, a limit, ordering, pagination, dedupe, a state transition.
3. Adds a failure path the user sees: an error mapping, timeout, retry, partial failure, offline. One case per path the change adds.
4. Keeps state that persists or recurs. One case per state the user can reach: fresh, previously linked, after a failure, after restart, back to an earlier state (A, then B, then A), two at once, cancelled midway.
5. Adds or changes a query or an access rule: integration test or SQL snapshot.
6. Adds an input constraint at the edge: transport test.
7. Adds a shared helper: unit test of its contract.
8. Transforms stored data: a migration test from a realistic old row to the new one.
9. Removes access (a route, a permission): a test that it is now refused. Other removals delete their tests and add none.
10. Fixes a class of mistake seen before (a second sighting in review, CI, or a bug report): one test that stops the class, not this instance.

### When not to test

No new test when the change only:

- moves or renames code; existing tests stay green
- declares what the compiler, framework, or generator already checks: types, DI wiring, constants, enum values, DTO fields, generated clients
- changes layout, styling, or copy, or adds a view with no branching
- adds logs or metrics, unless they are the feature
- adds a mapper or getter with no branching; its caller's test covers it
- repeats a behavior an existing test proves at the same or a cheaper tier; extend that test or its table instead

### Tiers

Prove each behavior at the cheapest tier that can catch the bug.

| Tier        | What is real                         | What it proves                                       |
| ----------- | ------------------------------------ | ---------------------------------------------------- |
| Unit        | nothing, every adapter mocked        | branching, call ordering, which adapter got what     |
| Transport   | the framework pipeline, logic mocked | validation, auth, parsing, right method invoked      |
| Integration | real datastore plus chosen adapters  | query correctness, access rules, cross table effects |
| E2E         | everything                           | smoke only                                           |

### Writing one

- Cases come from the requirement's words and the user states above, never from the code's branches. Write the reported scenario's case before reading the code.
- Suite is the unit under test, group is the method, case is one behavior stated as a sentence.
- Arrange, one action, assert. The action produces the state the test asserts; never arrange it. Assert the negative path too: what was _not_ called.
- Test through the public interface. Never read private fields or call private handlers. Never add a parameter or flag that only tests pass; inject the adapter instead.
- Assert the specific error type, never the generic one.
- Table tests for cases that differ by one input.
- A new test proves itself: once green, break the line it guards (invert the condition, drop the call), watch it fail, restore. A negative test fails on the rule it names, not on an earlier prerequisite.
- Adapters are the mock boundary: mock or fake adapters, never logic you own. Mock minimally: one resolved value per call that matters, ordered chains only when order is the point.
- A fake fails where the real adapter fails: it checks the arguments the real one checks (auth, ids), returns the adapter's own error rather than an already mapped one, and returns realistic data for every input the code branches on. An unexpected call fails loudly; a default error or zero value hides paths.
- Integration tests build state through the production flow (real service or adapter calls), never raw SQL or hand built rows.
- No sleeps. Drive time and races with fakes, completers, or an injected clock.
- Expected values are literals from the requirement, never the code's constants or a recomputation of its logic.
- Fixtures: static stubs for shapes, factories for values. The case sets every value it asserts on; the factory fills the rest.
- Rule of three covers test code: reuse the project's shared fakes and factories, and move the third copy of a double or setup into them.
- `sut` is the subject under test. Other setup names follow the file's convention; never stretch a name to fit another value.

## Commits and PRs

- Conventional commits, lowercase, no trailing period: `refactor: find or fail`, `fix(server): birthdate timezone`, `feat!: remove typeorm`. Scope is the top level directory or package; omit it when the change spans layers.
- Title names the thing, two to five words. Commit description one line or empty; PR description per PR descriptions below.
- One idea per PR. Large features ship as a chain of PRs, each revertable alone, refactors first.
- Squash merge, one commit per PR on master.
- Breaking change: `!` in the title plus a note where the project's changelog convention says.
- Derived artifacts (SQL snapshots, API spec, generated SDK, i18n keys) regenerate in the same PR as their source.

### PR descriptions

Plain words, short sentences, what changed for the user. Questions for the developer go in the reply, never the PR. End with `Fixes #n` or `Stacked on #n` when it applies, one line each. Include any notes required by the root `AGENTS.md`.

| Type                          | Description                                                                                           |
| ----------------------------- | ----------------------------------------------------------------------------------------------------- |
| `feat`                        | What the feature is, one or two sentences. Then a screenshot from the running app when it is visible. |
| `fix`                         | What the fix is, one or two sentences. Then two plain lines, `Before:` and `After:`.                  |
| `perf`                        | As `fix`, with measured numbers in `Before:` and `After:`.                                            |
| `refactor`                    | What the refactor is for and its benefit, one or two sentences.                                       |
| `chore`                       | What the chore is about, one or two sentences.                                                        |
| `docs`, `test`, `ci`, `build` | As `chore`.                                                                                           |

```markdown
Asset duration is null for a broken video, which kept the whole bucket from rendering.

Before: timeline blank past a broken video
After: timeline renders, the video shows no duration

Fixes #123
```

## Refactor playbook

1. Notice repetition while doing something else. Note it, finish the task.
2. Open a `refactor:` change scoped to that one pattern, touching every instance.
3. Write the helper's test, watch it fail, then write the helper as a generic, dependency free function. Defaults come from constants.
4. Replace all call sites. Expect net negative lines.
5. Tighten tests to assert the stricter behavior the refactor now guarantees.
6. Delete the old helpers and anything now unused.
7. Regenerate derived artifacts.
8. Then open the `feat:` change, which is now small.

Migrations at scale (ORM swap, test tier migration) follow the same loop: build the new harness, migrate one resource per PR, remove the old thing last with a breaking note.

## Checklist (component files fill in the commands)

- lint, format, typecheck, unit tests green
- derived artifacts regenerated
- migration generated if the schema changed
- translation catalog updated for new user facing strings in localized apps
- breaking change flagged in title and notes
- every scenario the change hits under Tests, When to test, has a test at its tier; nothing under When not to test got one
- the red then green run is in the report
- a risky change reports its counts at realistic scale

## In one paragraph

Write the smallest change that delivers one idea. Put logic in the logic layer, I/O in adapters, shapes in schemas, nothing in transport. When you see the same five lines twice you are already late: extract a generic helper and replace every instance now. Fail early with a plain message and let the edge handler do HTTP. Name things for the reader of the generated SDK and the reader of the test output. Test what a user would see break, first, with one focused test at the cheapest tier that catches it; a test that guards nothing is not written. Then delete more than you added.
