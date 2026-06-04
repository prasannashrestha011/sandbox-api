package models

import (
	"time"
)

type Lab struct {
	ID          string
	Title       string
	Description string
	Lang        string
	ContainerID string
	CreatedByID string // from req
	CreatedBy   User   //from res
	Chapters    []Chapter
	Tags        []Tag
	Difficulty  string
	IsPublic    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type Chapter struct {
	ID          string
	LabID       string
	Title       string
	Description string
	OrderIndex  int
	Exercises   []Exercise
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Exercise struct {
	ID             string
	ChapterID      string
	Title          string
	Description    string
	StarterCode    string
	ExpectedOutput string
	Hints          string
	OrderIndex     int
	Solution       string
	MaxAttempts    int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Tag struct {
	ID   string
	Name string
}

type LabEnrollment struct {
	UserID      string
	LabID       string
	Status      string // enrolled | in_progress | completed
	ProgressPct int
	EnrolledAt  time.Time
	CompletedAt *time.Time
}
