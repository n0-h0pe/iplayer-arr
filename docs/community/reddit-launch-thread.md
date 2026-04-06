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

---

## Thread Comments and Responses

### Comment: "iPlayarr already exists"

> I'm sorry to burst your bubble, but Iplayarr has existed for a while. The fact you said "nothing existed" for this, despite your applications having 2 character difference in name is not a good look.

**Response posted:**
> iPlayarr wraps the get_iplayer executable, it's not a native implementation.
>
> It also doesn't tackle the series and episode numbering problem this is built around, which is the core difference.
>
> You're right that iPlayarr has been around a while though, and I should have been clearer. Nothing existed that does *this specifically*, not that nothing existed at all.

### Comment: "Why not contribute to iPlayarr?"

> strich: I appreciate the work you've put in. But maybe it would have been better to contribute to the existing solution rather than this?

**Response posted:**
> That's a fair question, and I did look at it early on. The problem is the two projects share a goal but not much else. iPlayarr is built around get_iplayer as its engine, so contributing the native HLS downloader, the Newznab/SABnzbd API layer, and the episode resolution chain would have meant replacing most of the core. At that point you're not really contributing, you're proposing a rewrite of someone else's project, which didn't feel right either.
>
> Building separately meant Nikorag's project stays intact for the people it works well for, and this one exists for a different use case.

### Comment: VPN feedback (strich)

> That's fair. Well the addition of the VPN alone is a great feature. I have given it a go and spent the past 30 mins setting it up. Unfortunately my VPN provider appears to be blocked by BBC so I cannot test further at the moment. Only suggestion right now is that the VPN setup was a fair bit of work - You'd probably attract more attention if that became part of the setup in the GUI.

**Response posted:**
> Really glad you gave it a go, sorry the VPN provider got in the way. That's a known pain point with some providers - BBC's geo-detection has got stricter over time.
>
> The VPN setup UX is fair feedback and honestly something I've been thinking about. The hotio base handles the WireGuard config under the hood but surfacing that through the setup wizard rather than environment variables is a solid suggestion. I'll log it as a feature request.

---

## Action Items from Thread

- [ ] GitHub issue: VPN configuration surfaced through setup wizard UI
- [ ] Clarify "nothing existed" wording in original post going forward
- [ ] Monitor iPlayarr for feature updates (season offsetting, Lunr search improvements)
