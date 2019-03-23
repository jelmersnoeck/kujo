package kujo

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	v1 "k8s.io/api/batch/v1"
	cv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// JobSlice goes over a set of unstructured objects and returns a new slice with
// only jobs in them.
func JobSlice(uList []unstructured.Unstructured) ([]v1.Job, error) {
	var jobList []v1.Job
	for _, un := range uList {
		if isJobResource(un) {
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

// HashedJobs goes over a list of jobs and creates a unique hash for said job's
// configuration. The list is returned as a map where the key represents the
// original namespace and name for the job so it can be mapped back to the
// original list of resources.
func HashedJobs(jobs []v1.Job, config map[string]string) (map[string]string, error) {
	hashedJobs := map[string]string{}
	for _, job := range jobs {
		ns := job.Namespace
		if ns == "" {
			ns = "default"
		}

		hj, err := hashedJobName(job, config)
		if err != nil {
			return nil, err
		}

		key := fmt.Sprintf("%s/%s", ns, job.Name)
		hashedJobs[key] = hj
	}

	return hashedJobs, nil
}

func hashedJobName(job v1.Job, config map[string]string) (string, error) {
	specData, err := json.Marshal(job.Spec)
	if err != nil {
		return "", err
	}

	hashes := []string{fmt.Sprintf("%x", sha256.Sum256([]byte(specData)))}
	hashes = append(hashes, jobVolumeHashes(job, config)...)
	hashes = append(hashes, jobContainerHashes(job, config)...)

	return encodeHashSlice(hashes)
}

func jobContainerHashes(job v1.Job, config map[string]string) []string {
	ns := job.Namespace
	if ns == "" {
		ns = "default"
	}
	hashes := []string{}
	for _, container := range job.Spec.Template.Spec.Containers {
		hashes = append(hashes, containerEnvHashes(ns, container, config)...)
		hashes = append(hashes, containerEnvFromHashes(ns, container, config)...)
	}
	return hashes
}

func containerEnvHashes(ns string, container cv1.Container, config map[string]string) []string {
	hashes := []string{}

	for _, env := range container.Env {
		if vf := env.ValueFrom; vf != nil {
			if cmr := vf.ConfigMapKeyRef; cmr != nil {
				key := fmt.Sprintf("ConfigMap/%s/%s", ns, cmr.LocalObjectReference.Name)
				if val, ok := config[key]; ok {
					hashes = append(hashes, val)
				}
			}

			if sr := vf.SecretKeyRef; sr != nil {
				key := fmt.Sprintf("Secret/%s/%s", ns, sr.LocalObjectReference.Name)
				if val, ok := config[key]; ok {
					hashes = append(hashes, val)
				}
			}
		}
	}

	return hashes
}

func containerEnvFromHashes(ns string, container cv1.Container, config map[string]string) []string {
	hashes := []string{}

	for _, env := range container.EnvFrom {
		if cmr := env.ConfigMapRef; cmr != nil {
			key := fmt.Sprintf("ConfigMap/%s/%s", ns, cmr.LocalObjectReference.Name)
			if val, ok := config[key]; ok {
				hashes = append(hashes, val)
			}
		}

		if sr := env.SecretRef; sr != nil {
			key := fmt.Sprintf("Secret/%s/%s", ns, sr.LocalObjectReference.Name)
			if val, ok := config[key]; ok {
				hashes = append(hashes, val)
			}
		}
	}

	return hashes
}

func jobVolumeHashes(job v1.Job, config map[string]string) []string {
	if job.Spec.Template.Spec.Volumes == nil {
	}

	ns := job.Namespace
	hashes := []string{}
	for _, volume := range job.Spec.Template.Spec.Volumes {
		if volume.ConfigMap != nil {
			key := fmt.Sprintf("ConfigMap/%s/%s", ns, volume.ConfigMap.LocalObjectReference.Name)
			if val, ok := config[key]; ok {
				hashes = append(hashes, val)
			}
		}

		if volume.Secret != nil {
			key := fmt.Sprintf("Secret/%s/%s", ns, volume.Secret.SecretName)
			if val, ok := config[key]; ok {
				hashes = append(hashes, val)
			}
		}
	}

	return hashes
}

func encodeHashSlice(hashes []string) (string, error) {
	joinedString := strings.Join(hashes[:], "")
	return encodeHash(fmt.Sprintf("%x", sha256.Sum256([]byte(joinedString))))
}

// encodeHash extracts the first 40 bits of the hash from the hex string
// (1 hex char represents 4 bits), and then maps vowels and vowel-like hex
// characters to consonants to prevent bad words from being formed (the theory
// is that no vowels makes it really hard to make bad words). Since the string
// is hex, the only vowels it can contain are 'a' and 'e'.
// We picked some arbitrary consonants to map to from the same character set as GenerateName.
// See: https://github.com/kubernetes/apimachinery/blob/dc1f89aff9a7509782bde3b68824c8043a3e58cc/pkg/util/rand/rand.go#L75
// If the hex string contains fewer than ten characters, returns an error.
func encodeHash(hex string) (string, error) {
	if len(hex) < 10 {
		return "", fmt.Errorf("the hex string must contain at least 10 characters")
	}
	enc := []rune(hex[:10])
	for i := range enc {
		switch enc[i] {
		case '0':
			enc[i] = 'g'
		case '1':
			enc[i] = 'h'
		case '3':
			enc[i] = 'k'
		case 'a':
			enc[i] = 'm'
		case 'e':
			enc[i] = 't'
		}
	}
	return string(enc), nil
}
