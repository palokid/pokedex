package main

import (
	"testing"
	"time"

	"pokedex/internal/pokecache"
)

func TestCommandExplore(t *testing.T) {
	// 1. Setup a dummy config with a fresh cache
	// Use whatever package path you used for your cache
	conf := &config{
		Cache: pokecache.NewCache(5 * time.Minute),
	}

	// 2. Create raw mock JSON that matches your struct
	// We only need the pokemon_encounters array for this to pass
	mockJSON := []byte(`{
		"pokemon_encounters": [
			{
				"pokemon": {
					"name": "test-pikachu"
				}
			}
		]
	}`)

	// 3. Inject the mock data directly into the cache
	// This prevents the command from actually making an HTTP request!
	testArea := "test-zone"
	expectedURL := LocationAreasURL + testArea
	conf.Cache.Add(expectedURL, mockJSON)

	// 4. Table-driven test cases
	cases := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "fails with no arguments",
			args:        []string{},
			expectError: true,
		},
		{
			name:        "fails with too many arguments",
			args:        []string{"zone1", "zone2"},
			expectError: true,
		},
		{
			name:        "succeeds with valid cached argument",
			args:        []string{testArea},
			expectError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Call the function independently, passing our seeded config and test args
			err := commandExplore(conf, c.args...)

			if c.expectError && err == nil {
				t.Errorf("expected an error for %v args, but got nil", len(c.args))
			}
			if !c.expectError && err != nil {
				t.Errorf("did not expect an error, but got: %v", err)
			}
		})
	}
}
