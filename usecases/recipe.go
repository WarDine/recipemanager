package usecases

import "domain"

type Recipe struct {
}

// Enforce interface
var _ domain.Recipe = (*Recipe)(nil)
