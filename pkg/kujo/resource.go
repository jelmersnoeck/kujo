package kujo

import (
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

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
			if validObjectKind(un) {
				result = append(result, un)
			}
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

func validObjectKind(un unstructured.Unstructured) bool {
	if len(un.Object) == 0 {
		return false
	}

	versions, ok := validObjectKinds[un.GetKind()]
	if !ok {
		return false
	}

	unVersion := un.GetAPIVersion()
	for _, version := range versions {
		if version == unVersion {
			return true
		}
	}

	return false
}
