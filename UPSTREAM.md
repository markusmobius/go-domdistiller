# Upstream Reference

## README Format and Attribution

The September 29 documentation update applies the approved nine-section README
format and the full requirements in [AGENTS.md](AGENTS.md), including verified
creator acknowledgments. `go-domdistiller` remains 1.0.0; runtime sources,
dependencies, tags and historical benchmark evidence are unchanged. The
Chromium license link uses the complete pinned notice because the inherited
LICENSE-domdistiller.txt is empty.

## Released Suite Benchmark

The [2026-09-29 shared FAST report](https://github.com/markusmobius/content-extractor-benchmark/blob/49c426d6135df81b7d492bea7e6aec8e6d77d80c/go_rust_shared_performance_2026_09_29.json)
is authoritative for the current six-engine README comparison. Its published
LF-byte SHA-256 is
`382f869ae8f91c493c2c623c90a42a574f711ea584aead27387e907f3a24c523`.
It uses Go/Rust Readability 0.6.0/0.6.5, DomDistiller 1.0.0/1.0.1 and
Trafilatura 2.2.6/2.2.6 on 2,659 saved development pages. Source/build pins,
all four measured passes after one warmup, metadata scores and differences
remain in the report. Timings use `all_passes`, not a fastest-pass selection.

DomDistiller extraction is Go 6.159 / Rust 3.400 ms/page, with pagination off.
Parsing is a separate Go 11.283 / Rust 6.386 ms/page charge for the whole shared
suite, including Trafilatura's separate noscript tree when needed. This is not
a standalone reader measurement or an old/new version speedup. The Go release
remains 1.0.0; runtime source, dependencies and benchmark data are unchanged.
All six READMEs follow [AGENTS.md](AGENTS.md).

## Historical Suite Benchmark: 2026-09-23

The [2026-09-23 JSON](https://github.com/markusmobius/content-extractor-benchmark/blob/d433ab637f0a56c0926aa3698f470794a553472f/go_rust_shared_performance_2026_09_23.json)
records the historical comparison below, not the current README measurements.
Its published-file SHA-256 (LF line endings) is `7d7be9839f1652606cb91850af5134b188f2508be25623df889372dab4a06cc6`.
Read text scores at `quality[worker][engine].evaluations[corpus].overall.f1`,
selected timings at `overall`, and all-four timings at `all_passes`.

Go pins are Readability 0.6.0 (`db6ab179f951f80ad176850001aaf486e2fbd367`),
DomDistiller 1.0.0 (`25b8d046ffb4053bf68345d6fa59bc9ae1961ad8`) and Trafilatura
2.2.2 (`f4684e100869274311107325e3b72e47cc78db20`). Rust uses Readability 0.6.3,
DomDistiller 1.0.1 and Trafilatura 2.2.4; exact commits, dependency graphs and
binary hashes are in the build receipts. Rust-Trafilatura is now public;
the recorded source commits remain available. No Go version or tag changes.

| Implementation | Author Sets Exact / 1,290 | Author-Unit F1 | Titles Exact / 2,364 | Dates Exact / 1,530 |
| --- | ---: | ---: | ---: | ---: |
| go-readabilityV2-0.6.0 | 640 | 56.38767% | 1,247 | 763 |
| rust-readability-0.6.3 | 640 | 56.38767% | 1,247 | 763 |
| go-domdistiller-1.0.0 | 0 | 0.00000% | 1,106 | 0 |
| rust-domdistiller-1.0.1 | 0 | 0.00000% | 1,106 | 0 |
| go-trafilatura-2.2.2 | 695 | 58.80923% | 1,228 | 1,227 |
| rust-trafilatura-2.2.4 | 696 | 58.86640% | 1,227 | 1,227 |

Metadata uses only nonempty supplied annotations; unannotated is not negative,
and missing output is not filled by another engine. All six scored-output
digests match the preceding September 22 report. Trafilatura's two Go/Rust
differences concern one title and one author, not extracted text.

The full 2,659-page development run used seed 20260922, one warmup and four
measured passes. Passes 1 and 3 were selected by combined extraction time for
every row (5,318 observations each); all-four means retain 10,636 observations.
Worker order is balanced per page; three-engine order is a partial six-pass
block. Go uses `GOMAXPROCS=1`, `GOGC=100`, without forced collection. Native
timers exclude file reads and IPC; parsing and extraction stay separate.
The 26,590-response audit passed with no recorded sleep and AC power throughout.
The Go worker uses the unchanged Trafilatura 2.2.2 module graph; DomDistiller's
pseudo-version resolves to the exact 1.0.0 tag commit. This is not a standalone
reader comparison or unseen holdout test. Historical scores remain separate.

## Source and Scope

Go-DomDistiller is the native Go port of Chromium DOM Distiller, with retained
Boilerpipe-derived components. The main branch additionally incorporates the
Readability-inspired improvements described in [IMPROVEMENTS.md](IMPROVEMENTS.md).
The separate stable branch is the closer port of the original Java engine.

The benchmark uses release 1.0.0 at
`25b8d046ffb4053bf68345d6fa59bc9ae1961ad8`, also identified by development
module version `v0.0.0-20240926050704-25b8d046ffb4`. Source and module versions
are not changed by this documentation refresh. Server-side extraction omits
browser-computed CSS and layout checks, as documented in [README.md](README.md#limitations).
This is not a claim of exact browser or Java-output equivalence.

The [LICENSE](LICENSE), [DOM Distiller license](LICENSE-domdistiller.txt),
[Boilerpipe license](LICENSE-boilerpipe.txt), and
[Boilerpipe notice](NOTICE-boilerpipe.txt) remain unchanged.