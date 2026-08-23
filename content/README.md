# Content import files

This directory contains JSON documents used by the content importer.

Run an import with:

```bash
make import-content file=content/example.json
```

An import validates the complete document and replaces the matching liturgy in
one PostgreSQL transaction. Imported content is always saved as `draft` and is
therefore hidden from the public mobile API until it is reviewed and published.

The example file contains development placeholders only. Do not treat it as
authoritative Ethiopian Orthodox Church content. Production files must contain
text and audio metadata approved by the responsible church reviewers.
Valid roles are:

- `priest`
- `assistant_priest`
- `deacon`
- `assistant_deacon`
- `congregation`
- `chanter`
- `reader`
- `rubric`

- `mixed`
Section and verse orders must start at 1 and remain consecutive. A text-only
verse must contain at least one of Geʽez, Amharic, or English. Its audio timing
may be omitted. When a section has audio, every verse requires a complete,
non-overlapping time range and every audio metadata field is required, including
a lowercase 64-character SHA-256 checksum.

## St. Mary text-first workflow

Prepare the importer document from the corrected slide source:

```bash
make prepare-st-mary
```

Importing replaces the matching content in one transaction and returns it to
draft status:

```bash
make import-st-mary
```

Normal publication is blocked while any source segment is marked as requiring
review:

```bash
make publish-content slug=st-mary
```

For local preview only, use the explicit review override:

```bash
make publish-st-mary
```

The combined local-preview workflow is `make integrate-st-mary`. Re-import and
review the content again when approved St. Mary audio and exact timing become
available.

## Flat anaphora preview workflow

The Apostles, Our Lord Jesus Christ, and St. Athanasius source files use a
flat slide-array format. Convert all three into importer documents with:

```bash
make prepare-anaphoras
```

The converter preserves explicit Geʽez, Amharic, and English fields, maps the
bilingual role labels to backend roles, and keeps duplicate source-page numbers
as separate spoken passages on the same reader page. When a source entry only
contains `ethiopic_text`, the source does not provide a safe Geʽez/Amharic
boundary. The converter keeps that combined text visible in the Amharic fallback
field and marks every imported passage as requiring review.

Import drafts with `make import-anaphoras`. For an authorized local or deployed
preview, `make publish-anaphoras` uses the explicit review override, and
`make integrate-anaphoras` performs both steps. These targets do not mean that
the transcription has received church review; re-import and publish normally
after the language separation and text corrections are approved.

### Separate combined Ethiopic fields

Use the matching permitted PDF to suggest separate Geʽez and Amharic fields:

```bash
make separate-ethiopic \
  file=content/generated/Anaphora-of-the-Apostles.json \
  pdf=source-material/KidasehawariatAugust2010.pdf \
  metadata_only=true
```

The command first uses PDF font, color, and page-position evidence. Without
`metadata_only=true`, it also runs slower OCR on pages whose metadata does not
prove a boundary. It writes a `-separated.json` draft and a
`-separated-report.json` review report beside the input.

Successful suggestions contain `geez_text`, `amharic_text`,
`original_ethiopic_text`, and `separation_confidence`. Uncertain entries keep
their original `ethiopic_text` unchanged. Generated suggestions are never
publication approval: compare every split with the visible PDF and obtain the
required church review before importing or publishing it.
