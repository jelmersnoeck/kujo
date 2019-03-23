package kujo

import (
	"os"
	"testing"
)

func TestJobSlice(t *testing.T) {
	tcs := map[string]struct {
		fixture  string
		jobCount int
	}{
		"without data": {
			fixture: "testdata/no-config.yaml",
		},
		"with no job in the config": {
			fixture: "testdata/config-only.yaml",
		},
		"with jobs in the config": {
			fixture:  "testdata/full-config.yaml",
			jobCount: 1,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(tc.fixture)
			if err != nil {
				t.Errorf("Expected no error opening the file '%s', got '%s'", tc.fixture, err)
			}
			defer f.Close()

			rs, err := ResourcesFromReader(f)
			if err != nil {
				t.Errorf("Expected no error loading the resources, got '%s'", err)
			}

			sl, err := JobSlice(rs)
			if err != nil {
				t.Errorf("Expected no error getting the job slice, got '%s'", err)
			}

			if len(sl) != tc.jobCount {
				t.Errorf("Expected %d jobs, got %d", tc.jobCount, len(sl))
			}
		})
	}
}
