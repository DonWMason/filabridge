# FilaBridge — Project Instructions

## GitHub lifecycle

When Don asks to get a change "into GitHub" (or similar), drive the full lifecycle rather than stopping after a commit. The standard flow for this repo:

1. Branch off `main` (never commit directly to the default branch).
2. Commit (end messages with the `Co-Authored-By` trailer).
3. Push the branch.
4. Open a PR against `main` with a clear problem/root-cause/fix body.
5. Merge it (squash, delete the branch) — **when Don asks to merge**.
6. Sync local `main`.
7. Release: tag `vX.Y.Z` on `main` and push the tag. **Do not edit `CHANGELOG.md` by hand** — the tag push triggers `.github/workflows/release.yml` (binaries + GitHub release), which generates the changelog entry from commit messages since the previous tag and commits it back to `main` as `chore(release): update changelog for vX.Y.Z`; a manual entry produces a duplicate section. The same tag push triggers `.github/workflows/docker-build.yml` (multi-arch image, `:vX.Y.Z` + `:latest`, ~16 min). Version comes from the git tag; there is no version constant in the code. Use semver (bug fix → patch bump, feature → minor bump). Afterwards, sync local `main` to pick up the changelog commit.

Merging and pushing release tags are outward-facing/hard-to-reverse — confirm before those steps unless Don has just asked for them in context. Flag anything risky and recommend a real end-to-end verification against the live instance before tagging.
