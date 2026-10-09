# FilaBridge — Project Instructions

## GitHub lifecycle

When Don asks to get a change "into GitHub" (or similar), drive the full lifecycle rather than stopping after a commit. The standard flow for this repo:

1. Branch off `main` (never commit directly to the default branch).
2. Commit (end messages with the `Co-Authored-By` trailer).
3. Push the branch.
4. Open a PR against `main` with a clear problem/root-cause/fix body.
5. Merge it (squash, delete the branch) — **when Don asks to merge**.
6. Sync local `main`.
7. Release: add a `## [vX.Y.Z]` entry to `CHANGELOG.md`, then tag `vX.Y.Z` on `main` and push the tag — this triggers `.github/workflows/release.yml` (builds binaries + Docker image). Version comes from the git tag; there is no version constant in the code. Use semver (bug fix → patch bump).

Merging and pushing release tags are outward-facing/hard-to-reverse — confirm before those steps unless Don has just asked for them in context. Flag anything risky and recommend a real end-to-end verification against the live instance before tagging.
