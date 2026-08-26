package domain

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	AnnouncementStatusDraft     = "draft"
	AnnouncementStatusPublished = "published"
	AnnouncementStatusArchived  = "archived"
)

var announcementSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Announcement struct {
	ID          int64              `json:"id"`
	Slug        string             `json:"slug"`
	Version     int64              `json:"version"`
	TitleAm     string             `json:"title_am"`
	TitleEn     string             `json:"title_en"`
	BodyAm      string             `json:"body_am"`
	BodyEn      string             `json:"body_en"`
	Kind        string             `json:"kind"`
	Action      AnnouncementAction `json:"action"`
	Priority    int                `json:"priority"`
	IsPinned    bool               `json:"is_pinned"`
	Status      string             `json:"-"`
	PublishedAt *time.Time         `json:"published_at,omitempty"`
	ExpiresAt   *time.Time         `json:"expires_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type AnnouncementAction struct {
	Type  string `json:"type"`
	Value string `json:"value,omitempty"`
}

type AnnouncementInput struct {
	Slug        string     `json:"slug"`
	TitleAm     string     `json:"title_am"`
	TitleEn     string     `json:"title_en"`
	BodyAm      string     `json:"body_am"`
	BodyEn      string     `json:"body_en"`
	Kind        string     `json:"kind"`
	ActionType  string     `json:"action_type"`
	ActionValue string     `json:"action_value"`
	Priority    int        `json:"priority"`
	IsPinned    bool       `json:"is_pinned"`
	Status      string     `json:"status"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

func (input *AnnouncementInput) NormalizeAndValidate() error {
	input.Slug = strings.TrimSpace(input.Slug)
	input.TitleAm = strings.TrimSpace(input.TitleAm)
	input.TitleEn = strings.TrimSpace(input.TitleEn)
	input.BodyAm = strings.TrimSpace(input.BodyAm)
	input.BodyEn = strings.TrimSpace(input.BodyEn)
	input.Kind = strings.TrimSpace(input.Kind)
	input.ActionType = strings.TrimSpace(input.ActionType)
	input.ActionValue = strings.TrimSpace(input.ActionValue)
	input.Status = strings.TrimSpace(input.Status)

	if !announcementSlugPattern.MatchString(input.Slug) {
		return fmt.Errorf("slug must use lowercase words separated by hyphens")
	}
	if input.TitleAm == "" || input.TitleEn == "" {
		return fmt.Errorf("Amharic and English titles are required")
	}
	if input.BodyAm == "" || input.BodyEn == "" {
		return fmt.Errorf("Amharic and English bodies are required")
	}
	if !oneOf(input.Kind, "general", "content", "audio", "app_update", "important") {
		return fmt.Errorf("unsupported announcement kind %q", input.Kind)
	}
	if !oneOf(input.ActionType, "none", "open_liturgy", "download_apk", "open_url") {
		return fmt.Errorf("unsupported action type %q", input.ActionType)
	}
	if input.ActionType == "none" && input.ActionValue != "" {
		return fmt.Errorf("action_value must be empty when action_type is none")
	}
	if input.ActionType != "none" && input.ActionValue == "" {
		return fmt.Errorf("action_value is required for %s", input.ActionType)
	}
	if input.ActionType == "open_liturgy" && !announcementSlugPattern.MatchString(input.ActionValue) {
		return fmt.Errorf("open_liturgy action_value must be a liturgy slug")
	}
	if input.ActionType == "open_url" || input.ActionType == "download_apk" {
		parsed, err := url.ParseRequestURI(input.ActionValue)
		if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
			return fmt.Errorf("%s action_value must be a public HTTPS URL", input.ActionType)
		}
	}
	if input.Priority < 0 || input.Priority > 100 {
		return fmt.Errorf("priority must be between 0 and 100")
	}
	if !oneOf(input.Status, AnnouncementStatusDraft, AnnouncementStatusPublished, AnnouncementStatusArchived) {
		return fmt.Errorf("unsupported status %q", input.Status)
	}
	return nil
}

func oneOf(value string, values ...string) bool {
	for _, candidate := range values {
		if value == candidate {
			return true
		}
	}
	return false
}
