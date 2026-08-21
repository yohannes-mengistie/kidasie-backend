# Liturgy content extraction and alignment

This pipeline converts permitted slide-style liturgy PDFs into reviewed importer
JSON without committing source documents, recordings, or generated OCR drafts.

## Requirements

The local machine must provide:

- `pdfinfo` and `pdftoppm` from Poppler
- Tesseract OCR with the `amh` and `eng` language data
- Python 3 for serving the local alignment page

The current development machine satisfies these requirements.

## 1. Extract a review draft

The default source paths are:

- `source-material/liturgy.pdf`
- `source-material/audio/apostles.mp3`

Run:

```bash
make extract-slides
```

The command renders each permitted PDF page temporarily, runs Ethiopic and
English OCR, suggests speaker roles, fingerprints the audio, and writes:

```text
content/generated/apostles-slides.json
```

Every OCR page has `needs_review: true`. The generated directory is ignored by
Git because it can contain licensed liturgical text.

For a quick page-range check, run the Go command directly:

```bash
go run ./cmd/extractslides \
  -pdf source-material/liturgy.pdf \
  -out /tmp/slides-sample.json \
  -start 18 \
  -end 20
```

### Corrected slide drafts

The alignment page accepts both the raw extractor format and a curated slide
format. Curated pages may keep separate `text_geez`, `text_amharic`, and
`text_english` fields instead of copying text back into OCR fields.

The loader preserves the curated structure as follows:

- `people` is normalized to the backend role `congregation`.
- `assistant_deacon` is preserved as its own backend role.
- `parts` and separate people responses become individual audio segments.
- title-only headers and non-spoken instructions remain visible but are excluded
  from export by default.
- missing language or speaker values stay visible and block final export until
  reviewed.

Do not replace corrected text with OCR merely to satisfy the old draft shape.
The corrected JSON remains local and ignored by Git.

## 2. Review text, roles, and audio timing

Start the local alignment page:

```bash
make serve-aligner
```

Open:

```text
http://127.0.0.1:4173
```

Load the generated OCR draft and the matching MP3. The files stay in the local
browser and are not uploaded.

For every included segment:

1. Compare OCR text with the visible PDF.
2. Correct Geʽez and Amharic independently.
3. Confirm the participant role.
4. Split slides that contain multiple speakers.
5. Play the audio and mark each boundary.
6. Save review progress regularly.

The `Alt+N` shortcut records the current playback position as the next segment
boundary.

## 3. Export and import

The aligner blocks export until required metadata, roles, reviewed text, audio
checksum, and non-overlapping timings are complete. It downloads JSON compatible
with the transactional importer.

Import the reviewed file:

```bash
make import-content file=/path/to/apostles.json
```

The importer stores it as `draft`. A qualified church reviewer must approve the
text, roles, and timing before publication.

## Accuracy boundary

OCR and role detection are suggestions, not authoritative transcriptions.
Tesseract's Amharic language model recognizes Ethiopic characters but can alter
liturgical Geʽez words. Never publish generated text without line-by-line review.
