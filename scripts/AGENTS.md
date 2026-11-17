# Shell rules · extends [core](../docs/agents/core.md)

Governs shell scripts in the template, including scripts in mise tasks and GitHub workflows. Preserve the shell declared by the shebang or workflow. Use POSIX syntax for scripts that declare `sh`.

## Functions

Use lowercase names with underscores. Name each function for what it does.

* Predicates answer yes or no through the exit status. Prefix them with `is_` or `has_`.
* Actions change state or report failure. Start the name with a verb, such as `write_version`.
* Producers print a value or list for command substitution. Name them as a noun, such as `current_version`.
* Hooks called by another tool keep the name that tool requires.

Keep functional ShellCheck directives with the statements they cover.

## Tests

Test branches that affect observable script behavior. When the app uses Bats, mirror the script path under `tests/shell/`. Verify version bumps on a disposable copy of the repository. Use a local production build to check the embedded version.

## Checklist

Check syntax with the declared shell. Use ShellCheck and shfmt when installed. Run the affected script or workflow commands on disposable inputs and report any cases that need a hosted runner.
