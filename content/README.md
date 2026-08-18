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
- `congregation`
- `chanter`
- `reader`
- `rubric`

Section and verse orders must start at 1 and remain consecutive. Verse audio
ranges cannot overlap. If an audio object is present, every audio metadata field
is required, including a lowercase 64-character SHA-256 checksum.
