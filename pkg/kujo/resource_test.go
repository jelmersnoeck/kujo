package kujo

import (
	"os"
	"testing"
)

func TestResourcesFromReader(t *testing.T) {
	tcs := map[string]struct {
		fixture     string
		err         error
		resourceLen int
	}{
		"without data in the reader": {
			fixture: "testdata/no-config.yaml",
		},
		"with a list of data": {
			fixture:     "testdata/full-config.yaml",
			resourceLen: 4,
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
			if err != tc.err {
				t.Errorf("Expected err to be '%s', got '%s'", tc.err, err)
			}

			// there's an error, we shouldn't process further
			if err != nil {
				return
			}

			if len(rs) != tc.resourceLen {
				t.Errorf("Expected %d resources, got %d", tc.resourceLen, len(rs))
			}
		})
	}
}
