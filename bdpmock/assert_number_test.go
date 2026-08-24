// SPDX-FileCopyrightText: 2026 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: MPL-2.0

package bdpmock

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// A snapshot is written to a file and read back, so its numbers are float64.
// The captured side is whatever the code under test produced — which is
// json.Number when the producer decoded with UseNumber to keep a literal
// exact. Both describe the same payload and must compare equal.
func TestUnifyNumbersHandlesJSONNumber(t *testing.T) {
	captured := map[string]any{
		"capacity": json.Number("490"),
		"lat":      json.Number("46.71602"),
		"name":     "Demo",
	}
	fromSnapshot := map[string]any{
		"capacity": float64(490),
		"lat":      float64(46.71602),
		"name":     "Demo",
	}

	a := unifyNumbersToFloat(captured).(map[string]any)
	b := unifyNumbersToFloat(fromSnapshot).(map[string]any)

	for _, k := range []string{"capacity", "lat", "name"} {
		if a[k] != b[k] {
			t.Errorf("%s: captured unified to %v (%T), snapshot to %v (%T)",
				k, a[k], a[k], b[k], b[k])
		}
	}
}

// Nested payloads are the normal case for station metadata, so the conversion
// has to reach through maps and slices rather than only the top level.
func TestUnifyNumbersHandlesNestedJSONNumber(t *testing.T) {
	captured := map[string]any{
		"metaData": map[string]any{
			"capacity": json.Number("245"),
			"netex":    map[string]any{"levels": json.Number("3")},
			"list":     []any{json.Number("1"), json.Number("2")},
		},
	}
	want := map[string]any{
		"metaData": map[string]any{
			"capacity": float64(245),
			"netex":    map[string]any{"levels": float64(3)},
			"list":     []any{float64(1), float64(2)},
		},
	}

	a := unifyNumbersToFloat(captured)
	b := unifyNumbersToFloat(want)

	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	if string(ja) != string(jb) {
		t.Errorf("nested json.Number was not unified:\n captured %s\n want     %s", ja, jb)
	}
}

type tagged struct {
	MediumUrl string    `json:"medium_url"`
	Skipped   string    `json:"-"`
	Absent    string    `json:"absent,omitempty"`
	When      time.Time `json:"when"`
}

// A struct is compared against the same value after a trip through a snapshot
// file, so it has to unify to the shape encoding/json gives — tag names, not Go
// field names.
func TestUnifyStructMatchesItsJSONRoundTrip(t *testing.T) {
	v := tagged{MediumUrl: "https://example/i.png", Skipped: "not serialised",
		When: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)}

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var roundTripped interface{}
	if err := json.Unmarshal(b, &roundTripped); err != nil {
		t.Fatal(err)
	}

	live := unifyNumbersToFloat(v)
	disk := unifyNumbersToFloat(roundTripped)

	if !reflect.DeepEqual(live, disk) {
		t.Errorf("a struct and its own round trip do not unify alike:\n  live = %#v\n  disk = %#v", live, disk)
	}
	m, ok := live.(map[string]interface{})
	if !ok {
		t.Fatalf("want a map, got %T", live)
	}
	if _, keyed := m["medium_url"]; !keyed {
		t.Errorf("keyed on the Go field name, not the json tag: %v", m)
	}
	if _, leaked := m["Skipped"]; leaked {
		t.Error(`a json:"-" field was included`)
	}
	if _, leaked := m["absent"]; leaked {
		t.Error("an empty omitempty field was included")
	}
	if _, isString := m["when"].(string); !isString {
		t.Errorf("time.Time unified to %T, want the marshalled string", m["when"])
	}
}
