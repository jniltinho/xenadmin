---
name: create-release
description: Cut a release of the XCP-ng Center desktop app (Linux, macOS and Windows) on GitHub with the gh CLI — CHANGELOG, checks, the linux-v* tag, the release workflow (tarball, AppImage, macOS universal zip, Windows zip, SHA256SUMS) and its verification. Use when asked to release, tag or publish a version of the desktop app.
---

# Releasing the XCP-ng Center desktop app (GitHub, gh CLI)

Adapted from the `create-release` skill of github.com/jniltinho/ad-manager.

The repository is **github.com/jniltinho/xenadmin**. A release is: bump `XenAdmin.Desktop/CHANGELOG.md`, commit,
push the branch, wait for its CI, push the tag. The tag `linux-vX.Y.Z` starts `.github/workflows/linux.yml`, which runs
the unit tests and the net481 build of the core again, builds the self-contained Linux tarball, the AppImage, the universal macOS app and the Windows x64 app,
smoke-tests the AppImage under Xvfb and publishes a GitHub **pre-release** with the notes of that version's section in
`XenAdmin.Desktop/CHANGELOG.md` and one `SHA256SUMS`.

> The version comes from the **git tag** (`Makefile`: `git describe --tags --match 'linux-v*'`; the workflow passes
> `VERSION=X.Y.Z` to `make appimage`): nothing to bump in a file. The root `CHANGELOG.md` belongs to the Windows
> (WinForms) app and is not touched.

Release only when the user asks for it.

## Facts

```bash
gh auth status                                  # logged in, with access to jniltinho/xenadmin
REPO=jniltinho/xenadmin
BRANCH=feature/linux-desktop                    # the desktop line, until it is merged into development
LAST=$(git describe --tags --abbrev=0 --match 'linux-v*' 2>/dev/null || git rev-list --max-parents=0 HEAD)
```

Versioning (0.x, tags `linux-v`): a user-visible feature bumps the **minor** (`linux-v0.1.0` → `linux-v0.2.0`);
fixes, docs or CI only bump the **patch** (`linux-v0.2.0` → `linux-v0.2.1`).

## Process

### 1. Prerequisites

```bash
git checkout "$BRANCH" && git pull --ff-only
git status                                      # clean, besides the release changes
```

### 2. What changed

```bash
git log "$LAST"..HEAD --oneline
```

### 3. `XenAdmin.Desktop/CHANGELOG.md`

Keep a Changelog, in English. Move the items of `## [Unreleased]` into a new section **above** the previous one —
the workflow takes the release notes from it by its exact heading (`## [X.Y.Z]`):

```markdown
## [Unreleased]

—

## [0.2.0] — YYYY-MM-DD

**One or two lines of summary.**

### Added
### Changed
### Fixed
### Security
```

Omit empty sections.

### 4. Check, commit, push, wait for the CI

```bash
make test && make test-win && make appimage && make smoke && make macos && make windows
git add -A && git diff --cached --name-only        # nothing from _output/, dist/, .env, screenshots
git commit -m "release: linux-vX.Y.Z"
git push origin "$BRANCH"
SHA=$(git rev-parse HEAD)
until RUN=$(gh run list --repo "$REPO" --workflow linux.yml --commit "$SHA" --event push --limit 1 --json databaseId -q '.[0].databaseId') && [ -n "$RUN" ]; do sleep 5; done
gh run watch --repo "$REPO" --exit-status "$RUN"
```

Integration tests need the lab servers and run locally (`XCP_*` variables, see `XenAdmin.Desktop/README.md`); run
them before a minor release. Never tag a commit whose CI is red.

### 5. Tag (starts the release run)

```bash
git tag -a linux-vX.Y.Z -m "Release linux-vX.Y.Z"
git push origin linux-vX.Y.Z
until RUN=$(gh run list --repo "$REPO" --workflow linux.yml --commit "$SHA" --event push --json databaseId,headBranch -q '[.[] | select(.headBranch=="linux-vX.Y.Z")][0].databaseId') && [ -n "$RUN" ]; do sleep 5; done
gh run watch --repo "$REPO" --exit-status "$RUN"    # ~10 min
```

If a job fails: `gh run view --repo "$REPO" "$RUN" --log-failed`, fix on the branch, and release a new patch version
(do not move a published tag).

### 6. Verify

```bash
gh release view linux-vX.Y.Z --repo "$REPO" --json name,isPrerelease,assets -q '.name, .isPrerelease, (.assets[].name)'
```

Expect: `xcp-ng-center-X.Y.Z-linux-x64.tar.gz`, `xcp-ng-center-X.Y.Z-x86_64.AppImage`,
`xcp-ng-center-X.Y.Z_macos_universal.zip`, `xcp-ng-center-X.Y.Z-win-x64.zip` and `SHA256SUMS`.

### 6b. Format the GitHub release notes (rustXterm style)

The workflow publishes with the raw CHANGELOG section. Keep the CHANGELOG in Keep a Changelog, but give the **GitHub
release** the release-notes style of `github.com/jniltinho/rustXterm/releases`: a `# Release linux-vX.Y.Z` title,
emoji sections, a `## 📦 Downloads` list and a `**Full Changelog**` compare link. Map `### Added` → `## ✨ New Features`,
`### Changed`/`### Fixed` → `## 🔧 Improvements`, `### Security` → `## 🔒 Security` (keep the bullets; omit empty
sections; drop the one-line summary). Build the Downloads list from the published assets and the compare link from
the previous tag.

```markdown
# Release linux-vX.Y.Z

## ✨ New Features
- …(from ### Added)

## 🔧 Improvements
- …(from ### Changed / ### Fixed)

## 📦 Downloads
- `xcp-ng-center-X.Y.Z-x86_64.AppImage` — desktop app for Linux (any x86_64 desktop with X11 or XWayland; `chmod +x` and run)
- `xcp-ng-center-X.Y.Z-linux-x64.tar.gz` — the same app as a self-contained folder (no .NET install needed)
- `xcp-ng-center-X.Y.Z_macos_universal.zip` — desktop app for macOS (universal: Apple Silicon and Intel; unsigned, run `xattr -dr com.apple.quarantine "XCP-ng Center.app"` once)
- `xcp-ng-center-X.Y.Z-win-x64.zip` — desktop app for Windows 10/11 x64 (the same Avalonia app; unzip and run `xcp-ng-center.exe`)
- `SHA256SUMS` — SHA-256 checksums of the files above

**Full Changelog**: https://github.com/jniltinho/xenadmin/compare/$LAST...linux-vX.Y.Z
```

Write it to a file in the scratchpad and apply: `gh release edit linux-vX.Y.Z --repo "$REPO" --notes-file <file>`
(never commit it).

## Guardrails

- Never push a tag before the tests and the CI of that commit are green.
- Never commit secrets, `.env`, `_output/` (lab ISOs, answerfiles, screenshots) or `dist/`.
- The tag is the version; to fix a published release, cut a new patch.
- Releases stay **pre-releases** until the app reaches parity with the Windows app.
