# Review requirements

The shared contract for manual and automatic code reviews, including replies in review threads.

## Scope and rule provenance

- Review the identified candidate against the actual comparison base, the merge base for a three-dot PR diff. Distinguish introduced, worsened, pre-existing and unresolved behavior. An old bug exposed by a new caller can be worsened; show that changed path. Keep unrelated pre-existing problems outside the findings.
- Use the active design, milestone and recorded decisions. Cite the requirement behind a spec finding. A later milestone is not a missing requirement in this PR. Conflicting or ambiguous requirements need a named decision, not an invented violation. Without requirements, state that spec coverage is unavailable.
- A design's **Concurrency and failures** section, or its decision rows, records what each operation guarantees and which windows are accepted limits. A finding that breaks a listed guarantee is a defect. A finding inside an accepted limit cites the row and is not a defect, unless it shows an impact the row didn't weigh. A race with an actor the design doesn't list still needs its trigger and impact, and names the row the design lacks.
- Tracked project and component guidance (the `AGENTS.md` files and `docs/agents/`) supplies shared conventions. Personal agent files outside the repo guide their owner's agent; a personal preference does not become an author obligation until the team or task adopts it. Distinguish behavior contracts, shared rules, agreed requirements and heuristics.
- When an automated reviewer supplies rules from the merge base, use that version. A PR's changes to review policy are proposals until adopted. Context files from the PR head explain the change but do not replace the supplied rules. A task may explicitly adopt proposed requirements.

## Responsibilities and affected callers

- The Go service owns application policy and authorization. Repositories own queries and persistence. Controllers map the service result to HTTP. Keep saved preferences, observed state and cached metadata distinct.
- For affected flows, check resource ownership and event recipients, path containment, partial writes, interrupted scans, stale events, concurrent updates and restart recovery. For API changes, trace supported client/server combinations, omitted/null/default values and generated consumers.
- When a PR adds a refusal, error, guard or new precondition to shared code, list every caller and check how each one handles it, including callers outside the flow the PR is about. State a coverage gap if the callers cannot all be checked.
- Find what else changes the same data from the data, not from the changed function: search for everything that reads or writes the same file path, table or field, older APIs and background jobs included. A guard on one writer leaves the others unguarded.

## Validate the premise

- A correctness finding needs a real caller, trigger, observable impact and source evidence. Trace all prerequisites, input shape, existing guards and the exact failure path. Actively look for what prevents the proposed trigger before confirming it.
- Check that a negative test reaches the intended rule instead of stopping at an unrelated prerequisite. A permissive fake, mismatched mock return type or missing generated input can invalidate a reproduction; establish that the failure comes from the behavior being reviewed.
- A test that shows one operation waits for another by checking that nothing changed for a fixed time also passes when the other operation hasn't started. It needs a barrier that shows the other operation is blocked.
- Separate source evidence from checks actually run. A decisive source trace can establish a bug without runtime reproduction. Report commands and outcomes only when observed; passing unrelated tests do not prove the finding safe.
- A mocked adapter does not establish real datastore, kernel, device or released-client behavior. Name the real and fake parts of integration, VM and hardware checks, and what remains unproven. Unavailable tools are a verification gap, not a reason to invent results.
- Unsettled premises stay unverified. Put genuine questions in the host's notes or discussion format, with the decision owner and missing information. Do not publish an uncertain assertion as an established defect. Agreement between reviewers is not proof; a supported finding from one reviewer remains valid.

## Simplification

- Look for a smaller implementation of the current requirements: independently mutable copies of the same state, redundant caches or lifecycle flags, duplicated policy, unused extension points and wrappers with no useful boundary or vocabulary.
- Before proposing another lock, map, worker or retry layer, check whether an existing owner can make the decision. A broad workflow with repeated mode switches may need smaller concrete paths; repeated policy may need one shared owner.
- Name the exact state, branch, parameter, type or indirection to remove. Explain the replacement flow, affected callers, maintenance cost removed and required behavior preserved. Fewer lines alone do not prove simplicity.
- Trace the smaller alternative through applicable authorization, errors, cancellation, ordering, retry and restart behavior. An adapter, transaction, generation check or cancellation owner can earn its place with one caller. Do not remove a safeguard without preserving its purpose.
- A heuristic design judgment is optional. A required correction needs a cited shared rule, substantive requirement or demonstrated defect. An unproved replacement stays unverified. Do not demand a broad rewrite or create a finding quota.

## Finding quality and priority

- Use a short concrete title. Explain the trigger, impact, smallest useful correction and focused regression check. Name the affected location. Keep proposed fixes distinct from fixes actually applied. Write plainly, without filler or claims of certainty beyond the evidence.
- Rank impact and urgency independently of confidence. Use the host's finding schema and categories. For Kritika, P0 is critical active risk, P1 is a serious correctness, security or compatibility problem, and P2 is an ordinary defect or substantive requirement correction. Do not add unsupported fields or categories.
- Separate required corrections, optional improvements and human decisions. Kritika supports P0 through P2; do not promote our manual review's P3 nits into P2 findings. Keep minor style, personal preferences and heuristic suggestions without substantive impact in optional notes or the manual review. Skip mechanical changes already enforced by tooling.
- Deduplicate by underlying cause, even when the same bug breaks both behavior and a requirement. Group repeated optional test or style feedback by rule. The existence of a test file, or a PR description calling tests mechanical, does not establish its coverage.
- Authentication, destructive filesystem operations, migrations, OS updates and CI changes need explicit verification gaps and human scrutiny. They do not automatically make every finding high priority. Agents do not replace required human approval.

## Rereviews and outcomes

- Use prior candidate and finding IDs when supplied. Recheck fixes against the current head, then inspect the new delta and affected callers. State incremental coverage; do not present it as another complete review.
- Do not repost a resolved finding or reopen an agreed tradeoff without new evidence. A similar title is not proof that two findings share a cause; check the actual failure path before deduplicating.
- Distinguish the reviewer's validation status, the human's posting decision and the author's observed response. Record fixed, incorrect premise, intentional, optional declined or deferred only with supporting evidence; otherwise the author outcome is not observed. Acceptance, silence and model agreement do not establish correctness or accuracy metrics.

## Code comments

Do not add code comments. Make the code self-explanatory through naming and structure instead.
