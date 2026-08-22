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
