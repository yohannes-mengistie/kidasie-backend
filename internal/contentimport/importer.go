package contentimport

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Importer struct {
	pool *pgxpool.Pool
}

func NewImporter(pool *pgxpool.Pool) *Importer {
	return &Importer{
		pool: pool,
	}
}

func (i *Importer) Import(
	ctx context.Context,
	document *Document,
) error {
	if document == nil {
		return errors.New("content document is required")
	}

	if err := document.Validate(); err != nil {
		return fmt.Errorf("validate content document: %w", err)
	}

	tx, err := i.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return fmt.Errorf("begin content import transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	liturgyID, err := upsertLiturgy(
		ctx,
		tx,
		document,
	)
	if err != nil {
		return err
	}

	if err := replaceSections(
		ctx,
		tx,
		liturgyID,
		document.Sections,
	); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit content import: %w", err)
	}

	return nil
}

func upsertLiturgy(
	ctx context.Context,
	tx pgx.Tx,
	document *Document,
) (int64, error) {
	const query = `
		INSERT INTO liturgies (
			slug,
			name,
			name_am,
			status,
			published_at
		)
		VALUES ($1, $2, $3, 'draft', NULL)
		ON CONFLICT (slug) DO UPDATE
		SET
			name = EXCLUDED.name,
			name_am = EXCLUDED.name_am,
			status = 'draft',
			published_at = NULL,
			updated_at = NOW()
		RETURNING id
	`

	var liturgyID int64

	err := tx.QueryRow(
		ctx,
		query,
		document.Slug,
		document.Name,
		document.NameAm,
	).Scan(&liturgyID)
	if err != nil {
		return 0, fmt.Errorf("upsert liturgy: %w", err)
	}

	return liturgyID, nil
}

func replaceSections(
	ctx context.Context,
	tx pgx.Tx,
	liturgyID int64,
	sections []Section,
) error {
	const deleteQuery = `
		DELETE FROM sections
		WHERE liturgy_id = $1
	`

	if _, err := tx.Exec(
		ctx,
		deleteQuery,
		liturgyID,
	); err != nil {
		return fmt.Errorf("delete existing sections: %w", err)
	}

	for i := range sections {
		sectionID, err := insertSection(
			ctx,
			tx,
			liturgyID,
			sections[i],
		)
		if err != nil {
			return fmt.Errorf(
				"import section order %d: %w",
				sections[i].Order,
				err,
			)
		}

		if err := insertVerses(
			ctx,
			tx,
			sectionID,
			sections[i].Verses,
		); err != nil {
			return fmt.Errorf(
				"import section order %d: %w",
				sections[i].Order,
				err,
			)
		}
	}

	return nil
}

func insertSection(
	ctx context.Context,
	tx pgx.Tx,
	liturgyID int64,
	section Section,
) (int64, error) {
	const query = `
		INSERT INTO sections (
			liturgy_id,
			sort_order,
			title,
			title_am,
			audio_url,
			audio_duration_ms,
			audio_size_bytes,
			audio_mime_type,
			audio_sha256
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9
		)
		RETURNING id
	`

	var (
		audioURL        any
		audioDurationMs any
		audioSizeBytes  any
		audioMIMEType   any
		audioSHA256     any
	)

	if section.Audio != nil {
		audioURL = section.Audio.URL
		audioDurationMs = section.Audio.DurationMs
		audioSizeBytes = section.Audio.SizeBytes
		audioMIMEType = section.Audio.MIMEType
		audioSHA256 = section.Audio.SHA256
	}

	var sectionID int64

	err := tx.QueryRow(
		ctx,
		query,
		liturgyID,
		section.Order,
		section.Title,
		section.TitleAm,
		audioURL,
		audioDurationMs,
		audioSizeBytes,
		audioMIMEType,
		audioSHA256,
	).Scan(&sectionID)
	if err != nil {
		return 0, fmt.Errorf("insert section: %w", err)
	}

	return sectionID, nil
}

func insertVerses(
	ctx context.Context,
	tx pgx.Tx,
	sectionID int64,
	verses []Verse,
) error {
	const query = `
		INSERT INTO verses (
			section_id,
			sort_order,
			text_geez,
			text_am,
			text_en,
			role,
			start_ms,
			end_ms,
			source_page,
			source_part,
			source_kind,
			source_note,
			source_needs_review
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12,
			$13
		)
	`

	for i := range verses {
		verse := verses[i]

		var textEn any
		if strings.TrimSpace(verse.TextEn) != "" {
			textEn = verse.TextEn
		}

		if _, err := tx.Exec(
			ctx,
			query,
			sectionID,
			verse.Order,
			verse.TextGeez,
			verse.TextAm,
			textEn,
			verse.Role,
			verse.StartMs,
			verse.EndMs,
			verse.SourcePage,
			optionalString(verse.SourcePart),
			optionalString(verse.SourceKind),
			optionalString(verse.SourceNote),
			verse.SourceNeedsReview,
		); err != nil {
			return fmt.Errorf(
				"insert verse order %d: %w",
				verse.Order,
				err,
			)
		}
	}

	return nil
}

func optionalString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	return value
}
