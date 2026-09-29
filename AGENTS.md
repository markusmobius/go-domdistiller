# Go-DomDistiller Maintenance

Read these instructions before changing README.md, UPSTREAM.md, CHANGELOG.md,
benchmark claims or releases. Inspect current files and git status; preserve
unrelated edits and follow explicit user constraints.

## Library Identity

This is the native Go DOM Distiller port, derived from Chromium and Boilerpipe.
Distinguish the main branch's improvements from the faithful stable branch.
Do not describe server-side extraction as computed browser layout, CSS or
JavaScript execution. Pagination is a supported independent feature. Removing
DomDistiller as a Trafilatura fallback does not remove this standalone library.
DomDistiller does clean and classify content; candidate reuse problems concern
which prepared input it receives, not an absence of its own cleanup.

## Document Roles

- README.md: purpose, scope, usage and concise current quality/speed. Preserve
  useful API/examples and clearly label historical comparisons and limitations.
- UPSTREAM.md: source/dependency pins, deliberate deviations, controlled evidence,
  reproduction commands, hashes, coverage and unresolved compatibility issues.
- CHANGELOG.md: dated/versioned user-visible changes and reasons, not worker
  incidents or an audit transcript. Label documentation-only changes explicitly.
- Release notes: the same release-facing changes, benchmark definitions and
  limitations as the README/changelog, with links to detailed evidence.
- AGENTS.md: durable instructions, not current results or work-in-progress.

## Six-Repository Benchmark Contract

Follow the complete
[benchmark maintenance guide](https://github.com/markusmobius/content-extractor-benchmark/blob/master/AGENTS.md).
Coordinate go-domdistiller, rust-domdistiller, go-readabilityV2, rust-readability,
go-trafilatura and rust-trafilatura, not only this library's language pair.

1. Every README has `## Current Quality and Speed` with the same six-engine
	comparison from one completed shared-suite report. Use measured version labels.
2. Columns: `Extractor`, `Go Version`, `Rust Version`, `Go ms/page`,
	`Rust ms/page`, `Go/Rust`. Rows: Readability, DomDistiller, Trafilatura FAST.
	Display milliseconds/page to three decimals and Go/Rust ratios to two decimals.
3. Read structured JSON and compute ratios from unrounded means. The current
	shared protocol uses all four measured passes after one warmup. Never mix
	dates, environments, modes, means/medians or selected/all-pass aggregates.
4. Report parsing separately, charged once per language/page. Include decoding,
	normalization, DOM construction and any separately required Trafilatura tree.
	State corpus/counts, options, hardware, toolchains and included/excluded work.
	DomDistiller pagination must be identified as on/off and by algorithm.
5. Keep non-FAST Trafilatura as a separate comparison with its own report. A
	paired language ratio is not an isolated version speedup or request latency.
6. Report named-corpus F1 percentages to five decimals, retaining errors in the
	denominators. Do not average different scoring definitions. Matching text
	scores do not prove byte-identical HTML or metadata; preserve known differences.
7. Link immutable reports/commits and preserve old artifacts unchanged. Clearly
	date old standalone timing/allocation tables; do not compare them directly
	to the current shared-input extraction boundary.
8. Documentation-only patches retain actual measured versions; do not relabel
	old measurements or rerun benchmarks just to change wording. Trafilatura's
	fallback selection frequency is not DomDistiller's standalone accuracy.

## Editing and Publication

1. Identify scope and evidence before editing. Documentation does not authorize
	algorithm, pagination, dependency or application-worker changes.
2. Update the three documents according to their roles and coordinate the common
	benchmark section across all six repositories. Use readable library prose.
3. Validate counts, percentages, units, ratios, versions, options and Markdown
	links. Compare common sections and run `git diff --check`. Use the existing
	Go test/vet and CI commands for relevant validation; disclose unrun gates.
4. Obtain authorization before commits, pushes, tags, versions or releases. Go
	documentation changes do not automatically require a new module version.
	Never move an already published tag or overwrite historical evidence.
5. Finalize documentation before an authorized release. Preserve runtime and
	dependency pins for documentation-only work. Verify ordinary module download
	with checksum checking enabled and verify the GitHub release page; a tag
	push alone does not create a release page.
6. Coordinated Rust README updates require attention to immutable crates.io
	archives. Obtain approval for a patch version, finish docs before packaging,
	inspect included files and verify the public source/doc bytes after publishing.
7. Hash committed/remote bytes, not assumed Windows checkout bytes: CRLF and LF
	differ. Preserve unrelated dirty files and keep credentials out of logs.
	Do not infer hosted CI success from a successful push or publication.
8. Report actual checks and publication state concisely. Detailed evidence and
	remaining limitations belong in UPSTREAM or the benchmark reports.