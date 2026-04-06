# Changelog

All notable changes to iplayer-arr will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2026-04-06

### Fixed

- **BBC daily soaps now reach Sonarr correctly.** EastEnders, Casualty, Holby City, Doctors, Coronation Street, Neighbours and any other BBC show whose iPlayer subtitle is just a date were silently broken end-to-end:
  - Newznab releases were emitted as `EastEnders.S01E7307.06042026.1080p...` because Tier 2 used `parent_position` as the episode number. Sonarr's parser interpreted this as season 1 episode 7307, found no matching episode in TVDB, and rejected every release.
  - The Newznab `tvsearch` filter compared Sonarr's TVDB-style `season`/`ep` against iPlayer's internal `Series`/`EpisodeNum`, which are both 0 for these shows. Every interactive search returned an empty RSS feed.
- `internal/newznab/titles.go`: added Tier 1.5 — when the iPlayer subtitle parses as a bare date (DD/MM/YYYY, DD-MM-YYYY, DD.MM.YYYY) and an air date is available, the release title is now generated in date form: `EastEnders.2026.04.06.1080p.WEB-DL.AAC.H264-iParr`. Sonarr's daily-series parser maps these to the correct `S{season}E{episode}` automatically. No per-show overrides required.
- `internal/newznab/search.go`: `handleTVSearch` now recognises Sonarr's daily-series query convention (`season=YYYY&ep=MM/DD`) and filters by air date instead of integer season/episode. The standard integer compare remains for normal numbered shows.
- `internal/bbc/ibl.go`: `IBLResult.AirDate` is now normalised to canonical `YYYY-MM-DD` at parse time via the new `normaliseAirDate` helper, regardless of whether BBC iBL returned `"6 Apr 2026"` or `"2026-04-09"`. Both `Search` and `ListEpisodes` paths covered. Downstream consumers (filters, title generation, RSS pubDate) can rely on a single shape.

### Tests

- +8 new tests bringing the suite from 84 to 109:
  - `bbc/ibl_test.go`: `TestListEpisodesNormalisesLooseAirDate`, `TestSearchNormalisesLooseAirDate`
  - `newznab/titles_test.go`: `TestGenerateTitleSubtitleIsBareDate`, `TestGenerateTitleSubtitleDateAlternateSeparators`, `TestGenerateTitleNumberedShowNotPromoted` (regression guard)
  - `newznab/handler_test.go`: `TestHandleTVSearchDailyMatchByDate`, `TestHandleTVSearchDailyMismatchByDate`, `TestHandleTVSearchStandardSEStillWorks` plus shared `fakeBBCSearchServer` helper. Closes the long-standing gap where the only handler test covered the `caps` endpoint.

### Verified end-to-end

- Sonarr `/api/v3/release?episodeId=49265` (live EastEnders S42E54): now returns 3 releases (1080p / 720p / 540p), all mapped to `S42E54`, `rejected: false`. Previously returned zero.
- Future episodes that haven't aired yet (e.g. S42E55, S42E56) correctly return zero items — the filter is precise, not just always returning the latest.
- Octonauts S1E1, In the Night Garden S1E1 (Tier 1 numbered shows): unchanged, still produce `S01E01.<episode-title>...` titles via Tier 1.

### Documentation

- Added Docker Hub and pkgbadge stats badges to README (b1bb865).

### CI

- Added weekly base image rebuild workflow that watches the hotio base image digest (0f3b805, c92dabd).
- Added multi-arch (amd64 + arm64) builds and Docker Hub publishing to release workflow (6f08605).

### Container images

```
docker pull ghcr.io/will-luck/iplayer-arr:1.0.1
docker pull willluck/iplayer-arr:1.0.1
```

## [1.0.0] - 2026-04-06

First stable release of iplayer-arr — a BBC iPlayer download manager that plugs into Sonarr as an indexer and download client.

### Added

- Full Sonarr integration via Newznab indexer and SABnzbd-compatible download API
- Built-in VPN support via hotio base image (WireGuard)
- Dashboard with download monitoring, history, and system health
- Setup wizard for guided Sonarr configuration
- Multi-arch images: `linux/amd64` and `linux/arm64`
- Published to both GHCR and Docker Hub
- Weekly automatic rebuild when the hotio base image updates

[1.0.1]: https://github.com/Will-Luck/iplayer-arr/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/Will-Luck/iplayer-arr/releases/tag/v1.0.0
