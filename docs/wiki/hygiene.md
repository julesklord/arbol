# Hygiene and Git Workflow

This project strictly follows the **FMG Development Bible**. Deviations will be met with silent judgment.

## Atomic Commits

The use of **Conventional Commits** is mandatory:
`<type>(<scope>): <subject>`

### Allowed Types

- `feat`: New functionality.
- `fix`: Bug correction.
- `docs`: Documentation changes.
- `style`: Visual changes (no logic).
- `refactor`: Code change that neither adds nor fixes anything.
- `chore`: Maintenance tasks, dependencies.

## Branch Workflow

- `main`: Production branch (linear history only).
- `feat/*`: Branches for new functionalities.
- `fix/*`: Branches for corrections.

**Banned:** `git push --force` to `main`. Don't even think about it.

## Verification Requirements

1. Fork the repo and create a topic branch.
2. Run `make build` and `make test` before opening a PR.
3. Include tests for new behaviors under `tests/` (bash bats or simple sh harness).
