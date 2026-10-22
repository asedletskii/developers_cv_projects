package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LayoutSnapshot struct {
	order                  []AnalyticsSectionID
	hashes                 map[AnalyticsSectionID]string
	technologies           []string
	technologyByFolded     map[string]string
	technologyCatalogKnown bool
}

func NewLayoutSnapshot(
	order []AnalyticsSectionID,
	hashes map[AnalyticsSectionID]string,
	technologySets ...[]string,
) (LayoutSnapshot, error) {
	if len(order) == 0 {
		return LayoutSnapshot{}, fmt.Errorf("%w: layout order is empty", ErrInvalidLayoutSnapshot)
	}
	seen := make(map[AnalyticsSectionID]struct{}, len(order))
	normalizedOrder := make([]AnalyticsSectionID, 0, len(order))
	for _, rawSection := range order {
		section, err := ParseAnalyticsSectionID(string(rawSection))
		if err != nil || section == AnalyticsSectionTechStack {
			return LayoutSnapshot{}, fmt.Errorf("%w: invalid section %q", ErrInvalidLayoutSnapshot, rawSection)
		}
		if _, ok := seen[section]; ok {
			return LayoutSnapshot{}, fmt.Errorf("%w: duplicate section %q", ErrInvalidLayoutSnapshot, section)
		}
		seen[section] = struct{}{}
		normalizedOrder = append(normalizedOrder, section)
	}

	normalizedHashes := make(map[AnalyticsSectionID]string, len(hashes))
	for _, section := range normalizedOrder {
		if _, ok := hashes[section]; !ok {
			return LayoutSnapshot{}, fmt.Errorf("%w: missing hash for section %q", ErrInvalidLayoutSnapshot, section)
		}
	}
	for section, hash := range hashes {
		if _, ok := seen[section]; !ok {
			return LayoutSnapshot{}, fmt.Errorf("%w: hash for hidden section %q", ErrInvalidLayoutSnapshot, section)
		}
		if len(hash) != 16 || strings.Trim(hash, "0123456789abcdef") != "" {
			return LayoutSnapshot{}, fmt.Errorf("%w: invalid hash for section %q", ErrInvalidLayoutSnapshot, section)
		}
		normalizedHashes[section] = hash
	}

	technologies := make([]string, 0)
	technologyByFolded := make(map[string]string)
	if len(technologySets) > 0 {
		for _, value := range technologySets[0] {
			value = strings.TrimSpace(value)
			if value == "" {
				continue
			}
			folded := strings.ToLower(value)
			if _, exists := technologyByFolded[folded]; exists {
				continue
			}
			technologyByFolded[folded] = value
			technologies = append(technologies, value)
		}
	}
	return LayoutSnapshot{
		order: normalizedOrder, hashes: normalizedHashes,
		technologies: technologies, technologyByFolded: technologyByFolded,
		technologyCatalogKnown: len(technologySets) > 0,
	}, nil
}

func RestoreLayoutSnapshot(
	order []AnalyticsSectionID,
	hashes map[AnalyticsSectionID]string,
	technologySets ...[]string,
) (LayoutSnapshot, error) {
	if len(technologySets) > 1 {
		return LayoutSnapshot{}, fmt.Errorf("%w: multiple technology catalogs", ErrInvalidLayoutSnapshot)
	}
	snapshot, err := NewLayoutSnapshot(order, hashes, technologySets...)
	if err != nil {
		return LayoutSnapshot{}, err
	}
	if len(technologySets) == 1 {
		stored := technologySets[0]
		if len(stored) != len(snapshot.technologies) {
			return LayoutSnapshot{}, fmt.Errorf("%w: technology catalog is not normalized", ErrInvalidLayoutSnapshot)
		}
		for i, technology := range stored {
			if technology != snapshot.technologies[i] {
				return LayoutSnapshot{}, fmt.Errorf("%w: technology catalog is not normalized", ErrInvalidLayoutSnapshot)
			}
		}
	}
	return snapshot, nil
}

func (s LayoutSnapshot) Order() []AnalyticsSectionID {
	return append([]AnalyticsSectionID(nil), s.order...)
}

func (s LayoutSnapshot) Hashes() map[AnalyticsSectionID]string {
	result := make(map[AnalyticsSectionID]string, len(s.hashes))
	for key, value := range s.hashes {
		result[key] = value
	}
	return result
}

func (s LayoutSnapshot) IndexOf(section AnalyticsSectionID) int {
	for i, value := range s.order {
		if value == section {
			return i
		}
	}
	return -1
}

func (s LayoutSnapshot) HashFor(section AnalyticsSectionID) (string, bool) {
	value, ok := s.hashes[section]
	return value, ok
}

func (s LayoutSnapshot) Technologies() []string {
	return append([]string(nil), s.technologies...)
}

func (s LayoutSnapshot) NormalizeTechnology(value string) (string, bool) {
	canonical, ok := s.technologyByFolded[strings.ToLower(strings.TrimSpace(value))]
	return canonical, ok
}

func (s LayoutSnapshot) HasTechnologyCatalog() bool {
	return s.technologyCatalogKnown
}

func (s LayoutSnapshot) MarshalJSON() ([]byte, error) {
	type legacySnapshot struct {
		Order  []AnalyticsSectionID          `json:"order"`
		Hashes map[AnalyticsSectionID]string `json:"hashes"`
	}
	type snapshotWithTechnologies struct {
		Order        []AnalyticsSectionID          `json:"order"`
		Hashes       map[AnalyticsSectionID]string `json:"hashes"`
		Technologies []string                      `json:"technologies"`
	}
	if !s.technologyCatalogKnown {
		return json.Marshal(legacySnapshot{Order: s.Order(), Hashes: s.Hashes()})
	}
	return json.Marshal(snapshotWithTechnologies{
		Order: s.Order(), Hashes: s.Hashes(), Technologies: s.Technologies(),
	})
}
