package kujo

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// JobSlice goes over a set of unstructured objects and returns a new slice with
// only jobs in them.
func JobSlice(uList []unstructured.Unstructured) []unstructured.Unstructured {
	var jobList []unstructured.Unstructured
	for _, un := range uList {
		if un.GetKind() == "Job" && validObjectKind(un) {
			jobList = append(jobList, un)
		}
	}

	return jobList
}
