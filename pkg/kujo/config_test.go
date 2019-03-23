package kujo

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHashedConfig(t *testing.T) {
	tcs := map[string]struct {
		fixture string
		config  map[string]string
		err     error
	}{
		"without data in the reader": {
			fixture: "testdata/no-config.yaml",
			config:  map[string]string{},
		},
		"without config provided": {
			fixture: "testdata/job-only.yaml",
			config:  map[string]string{},
		},
		"with a configmap provided": {
			fixture: "testdata/job-configmap.yaml",
			config: map[string]string{
				"default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
		},
		"with a secret provided": {
			fixture: "testdata/job-secret.yaml",
			config: map[string]string{
				"default/mysecret": "8ebf17fe046d11996943eee66edbe0a487cb0a7b75f34d3d469ab58649530fbd",
			},
		},
		"with configmap and secret provided": {
			fixture: "testdata/full-config.yaml",
			config: map[string]string{
				"default/mysecret":        "8ebf17fe046d11996943eee66edbe0a487cb0a7b75f34d3d469ab58649530fbd",
				"default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
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
				t.Errorf("Expected no errors getting the resources, got '%s'", err)
			}

			cm, err := HashedConfig(rs)
			if tc.err != err {
				t.Errorf("Expected error '%s', got '%s'", tc.err, err)
			}

			// we've got an error, don't run further tests
			if err != nil {
				return
			}

			if !cmp.Equal(tc.config, cm) {
				t.Errorf("Expected config to equal\n\n%s\n\ngot\n\n%s", tc.config, cm)
			}
		})
	}
}
