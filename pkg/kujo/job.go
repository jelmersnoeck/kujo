package kujo

import (
	"encoding/json"

	v1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// JobSlice goes over a set of unstructured objects and returns a new slice with
// only jobs in them.
func JobSlice(uList []unstructured.Unstructured) ([]v1.Job, error) {
	var jobList []v1.Job
	for _, un := range uList {
		if un.GetKind() == "Job" && validObjectKind(un) {
			var job v1.Job
			data, err := un.MarshalJSON()
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(data, &job); err != nil {
				return nil, err
			}
			jobList = append(jobList, job)
		}
	}

	return jobList, nil
}
