package kujo

import (
	"io"
	"log"
	"strconv"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const annKey = "kujo.sphc.io"

// ResourcesFromReader takes a reader object and parses the data into a slice
// of unstructured resources. The reader should either contain JSON or YAML
// objects.
// The objects are filtered based on their type. Only ConfigMap, Job and Secret
// resources are returned.
func ResourcesFromReader(rdr io.Reader) ([]unstructured.Unstructured, error) {
	decoder := yaml.NewYAMLOrJSONDecoder(rdr, 1024)

	var result []unstructured.Unstructured
	var err error
	for err == nil {
		var un unstructured.Unstructured
		if err = decoder.Decode(&un); err == nil {
			result = append(result, un)
		}
	}

	if err == io.EOF {
		err = nil
	}

	return result, err
}

// validObjectKinds is a map of data which represents the items we're looking
// for in a list of unstructured objects. It's mapped as Kind: []apiVersions.
var validObjectKinds = map[string][]string{
	"Job":       []string{"batch/v1"},
	"Secret":    []string{"v1"},
	"ConfigMap": []string{"v1"},
}

func isJobResource(un unstructured.Unstructured) bool {
	if un.GetKind() == "Job" {
		unVersion := un.GetAPIVersion()
		for _, version := range validObjectKinds["Job"] {
			if version == unVersion {
				ann := un.GetAnnotations()
				if val, ok := ann[annKey]; ok {
					pb, err := strconv.ParseBool(val)
					if err != nil {
						log.Println(err)
						return false
					}

					return pb
				}
			}
		}
	}

	return false
}
