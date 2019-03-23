package kujo

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestConvert(t *testing.T) {
	input, err := os.Open("testdata/convert-input.yaml")
	if err != nil {
		t.Errorf("Expected no error opening the input file, got '%s'", err)
	}
	defer input.Close()

	generated, err := SuffixJobs(input)
	if err != nil {
		t.Errorf("Did not expect error, got '%s'", err)
	}

	output, err := os.Open("testdata/convert-output.yaml")
	if err != nil {
		t.Errorf("Expected no error opening the input file, got '%s'", err)
	}
	defer output.Close()

	fixture, err := ioutil.ReadAll(output)
	if err != nil {
		t.Errorf("Did not expect error, got '%s'", err)
	}

	if string(fixture) != string(generated) {
		t.Errorf("Expected generated output\n%s\nto match fixture\n%s", string(generated), string(fixture))
	}
}
