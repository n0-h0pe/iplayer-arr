# Reddit Launch Thread - Community Notes

## Original Post

Posted to Reddit introducing iplayer-arr. Key points covered:

- Newznab indexer + SABnzbd download client emulation
- 4-tier episode resolution chain (Full, Position, Date, Manual)
- Per-show overrides for TheTVDB alignment
- HLS download via ffmpeg, .mp4 + .srt output
- Real-time dashboard with SSE progress
- Setup wizard for Sonarr connection
- Built-in WireGuard VPN via hotio base image
- Written in Go, single binary, ~15 MB image

---

## Edit Added to Original Post

After iPlayarr comparison was raised in comments:

> **Edit:** A few people have pointed out iPlayarr (https://github.com/Nikorag/iplayarr) which does a similar job, so worth clarifying the differences since the names are close. And to be clear, I'm not dismissing Nikorag's work at all - if that's how the original post came across, that's on me.
>
> iPlayarr wraps the get_iplayer executable and builds its search index from BBC broadcast schedules. It's a solid project and supports both Sonarr and Radarr.
>
> iplayer-arr is a different architecture. It has no get_iplayer dependency, downloads directly from iPlayer's HLS streams via ffmpeg, and searches the full iPlayer catalogue via the BBC's own API rather than schedule data, so anything currently available to stream is searchable regardless of whether it's currently airing. The name was chosen deliberately to be distinct - the arr suffix places it in the same ecosystem without being a clone.
>
> The episode numbering system is also the core thing that makes this different from both iPlayarr and anything else out there. BBC iPlayer's metadata is inconsistent enough that a schedule-based approach can't solve it on its own.

---

## iPlayarr vs iplayer-arr - Full Comparison

### Architecture

| | iplayer-arr | iPlayarr |
|---|---|---|
| get_iplayer dependency | None | Required |
| Download method | Native HLS via ffmpeg | get_iplayer wrapper |
| Sonarr integration | Newznab + SABnzbd | Indexer + client |
| Radarr support | No | Yes |
| Search source | iPlayer catalogue API (ibl.api.bbci.co.uk) | BBC broadcast schedules |
| Episode numbering | 4-tier resolution chain | Basic |
| Per-show overrides | Yes | No |
| Built-in VPN | Yes (WireGuard via hotio) | No |

### Key Technical Notes

- iPlayarr uses `https://www.bbc.co.uk/schedules/` endpoints - confirmed via issue #143 bug report showing schedule fetch logs
- iplayer-arr uses `https://ibl.api.bbci.co.uk/ibl/v1/new-search` - full iPlayer catalogue search
- Practical difference: content available on iPlayer but not currently airing may not surface in iPlayarr search
- iPlayarr also recently added season offsetting (PR #205) - worth monitoring for feature parity

### VPN Compatibility

- Private Internet Access (PIA) tested and confirmed working
- Some providers are blocked by BBC geo-detection - not a project issue, a provider issue

---

## Thread Comments and Responses

### Comment: "iPlayarr already exists"

> I'm sorry to burst your bubble, but Iplayarr has existed for a while. The fact you said "nothing existed" for this, despite your applications having 2 character difference in name is