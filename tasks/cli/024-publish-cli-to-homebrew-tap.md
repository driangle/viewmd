---
title: "Publish CLI to Homebrew tap"
id: "024"
status: completed
priority: medium
type: chore
tags: ["release", "homebrew"]
touches: ["ci"]
created: "2026-03-07"
---

# Publish CLI to Homebrew tap

## Description

Publish the Go CLI as a Homebrew formula in the existing tap at [driangle/homebrew-tap](https://github.com/driangle/homebrew-tap), so users can install viewmd via `brew install driangle/tap/viewmd`.

## Tasks

- [x] Set up cross-platform Go binary builds (macOS amd64/arm64, Linux amd64) via GitHub Actions or GoReleaser
- [x] Create a GitHub release workflow that builds binaries and attaches them as release assets
- [x] Write the Homebrew formula (`Formula/viewmd.rb`) in the `driangle/homebrew-tap` repo
- [x] Formula should download the correct binary for the user's platform and architecture
- [x] Include SHA256 checksums for each binary in the formula
- [x] Add a CI step or script to auto-update the formula on new releases (update URL, version, SHA256)
- [x] Test `brew install driangle/tap/viewmd` and `brew upgrade` on macOS
- [x] Document installation via Homebrew in the project README

## Acceptance Criteria

- `brew tap driangle/tap && brew install viewmd` installs the CLI successfully
- Formula supports macOS (Intel + Apple Silicon) at minimum
- New releases automatically update the formula (or a clear manual process is documented)
- `viewmd --version` (or equivalent) reports the correct version after install
- SHA256 checksums are verified during install
