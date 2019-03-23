package kujo

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// SuffixJobs takes a list of Kubernetes resources and goes over all the jobs.
// It matches the job's configuration with ConfigMap and Secret objects and
// calculates a unique hash from all three configurations to determine a unique
// name for each job.
// Once done, it replaces all the job names from the input with the newly
// calculated job name and outputs the data byte slice filled with YAML.
func SuffixJobs(data io.Reader) ([]byte, error) {
	resourceList, err := ResourcesFromReader(data)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read the resources from input")
	}

	jobs, err := JobSlice(resourceList)
	if err != nil {
		return nil, errors.Wrap(err, "Could not filter out the jobs")
	}

	// no jobs in the resource list, return the original
	if len(jobs) == 0 {
		return marshalUnstructured(resourceList)
	}

	cm, err := HashedConfig(resourceList)
	if err != nil {
		return nil, errors.Wrap(err, "Could not calculate the config hashes")
	}

	jobHashes, err := HashedJobs(jobs, cm)
	if err != nil {
		return nil, errors.Wrap(err, "Could not calculate job hashes")
	}
	for i, rs := range resourceList {
		if isJobResource(rs) {
			ns := rs.GetNamespace()
			if ns == "" {
				ns = "default"
			}

			key := fmt.Sprintf("%s/%s", ns, rs.GetName())
			if hash, ok := jobHashes[key]; ok {
				resourceList[i].SetName(fmt.Sprintf("%s-%s", rs.GetName(), hash))
			}
		}
	}

	return marshalUnstructured(resourceList)
}

func marshalUnstructured(resourceList []unstructured.Unstructured) ([]byte, error) {
	var separator bool
	buf := bytes.NewBuffer([]byte{})
	for _, rs := range resourceList {
		out, err := yaml.Marshal(rs.Object)
		if err != nil {
			return nil, err
		}

		if separator {
			_, err = buf.WriteString("---\n")
			if err != nil {
				return nil, err
			}
		} else {
			separator = true
		}

		_, err = buf.Write(out)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
