package main

import (
	"errors"
	"fmt"
	"time"
)

const (
	_humanDate   = "Jan 2006"    // for display
	_machineDate = time.DateOnly // for YAML (text sorts chronologically)
)

type metadata struct {
	Title      string `yaml:"title"`
	RawCreated string `yaml:"created"`
	Hidden     bool   `yaml:"hidden"`

	// for normal posts
	Description string `yaml:"description"`
	RawUpdated  string `yaml:"updated"`
	HideHome    bool   `yaml:"hide_home"`
	HideLicense bool   `yaml:"hide_license"`

	// for external posts
	Link string `yaml:"link"`
	Via  string `yaml:"via"`
}

func (m metadata) Validate() error {
	desc := m.Description != ""
	link := m.Link != ""
	if !desc && !link {
		return errors.New("posts require either description or external link")
	}
	if desc && link {
		return errors.New("posts cannot have both description and external link")
	}
	if link && m.Via == "" {
		return errors.New("external links must have via text")
	}
	if !m.Hidden && m.RawCreated == "" {
		return errors.New("all displayed posts must have a created date")
	}
	if m.RawUpdated != "" && m.RawUpdated <= m.RawCreated {
		return errors.New("updated date is before created date")
	}
	for _, d := range []string{m.RawCreated, m.RawUpdated} {
		if d == "" {
			continue
		}
		if _, err := time.Parse(_machineDate, d); err != nil {
			return fmt.Errorf("parse %s: %w", d, err)
		}
	}
	return nil
}

func (m metadata) Created() string {
	return machineToHuman(m.RawCreated)
}

func (m metadata) Updated() string {
	return machineToHuman(m.RawUpdated)
}

func (m metadata) Compare(other metadata) int {
	// Reverse chronological sort.
	if m.RawCreated > other.RawCreated {
		return -1
	}
	if m.RawCreated == other.RawCreated {
		return 0
	}
	return 1
}

func machineToHuman(date string) string {
	if date == "" {
		return ""
	}
	d, err := time.Parse(_machineDate, date)
	if err != nil {
		// checked in metadata.Validate, unreachable here
		panic(err.Error())
	}
	return d.Format(_humanDate)
}
