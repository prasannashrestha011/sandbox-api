package models

import (
	"time"
)

type Lab struct {
	ID          string
	Title       string
	Description string
	Lang        string
	Exercises   []Exercise
	ContainerID string
	CreatedByID string
	CreatedBy   User
	Tags        []Tag
	Difficulty  string
	IsPublic    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Exercise struct {
	ID             string
	LabID          string
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
