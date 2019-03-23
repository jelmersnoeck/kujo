package kujo

import (
	"crypto/sha256"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// HashedConfig goes over a given set of unstructured objects and filters out
// the ConfigMap and Secret objects. It then hashes it's content and returns a
// map of hashes, where the key is in the `<namespace>/<name>` format.
func HashedConfig(uList []unstructured.Unstructured) (map[string]string, error) {
	uMap := map[string]string{}
	for _, un := range uList {
		switch un.GetKind() {
		case "ConfigMap", "Secret":
			ns := un.GetNamespace()
			if ns == "" {
				ns = "default"
			}

			key := fmt.Sprintf("%s/%s", ns, un.GetName())
			if _, ok := uMap[key]; !ok {
				hsh, err := hashUnstructured(un)
				if err != nil {
					return nil, err
				}

				uMap[key] = hsh
			}
		}
	}

	return uMap, nil
}

func hashUnstructured(obj unstructured.Unstructured) (string, error) {
	data, err := obj.MarshalJSON()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256([]byte(data))), nil
}
