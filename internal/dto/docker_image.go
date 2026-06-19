package dto

import "strings"

type CreateImageRequest struct {
	// Docker Image ID
	ImageTag string `json:"image_tag,omitempty"`
	Lang     string `json:"lang,omitempty"`
}

func (r *CreateImageRequest) Sanitize() {
	r.ImageTag = strings.TrimSpace(r.ImageTag)
	r.Lang = strings.TrimSpace(r.Lang)
}

func (r *CreateImageRequest) Validate() error {
	r.Sanitize()
	var v ValidationErrors

	if r.ImageTag == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "image_tag",
			Message: "image_tag is required",
		})
	} else if !strings.Contains(r.ImageTag, ":") {
		// Assuming the required format is image:tag like ubuntu:general
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "image_tag",
			Message: "image_tag must be in a valid format containing a tag (e.g., name:general)",
		})
	}

	if r.Lang == "" {
		v.Violations = append(v.Violations, FieldViolation{
			Field:   "lang",
			Message: "lang is required",
		})
	}

	if len(v.Violations) > 0 {
		return &v
	}
	return nil
}

type DockerImageResponse struct {
	ID       string `json:"id"`
	ImageTag string `json:"image_tag"`
	Lang     string `json:"lang"`
}
