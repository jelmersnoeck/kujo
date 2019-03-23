package kujo

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/batch/v1"
	cv1 "k8s.io/api/core/v1"
	mv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestJobSlice(t *testing.T) {
	tcs := map[string]struct {
		fixture  string
		jobCount int
	}{
		"without data": {
			fixture: "testdata/no-config.yaml",
		},
		"with no job in the config": {
			fixture: "testdata/config-only.yaml",
		},
		"with jobs in the config": {
			fixture:  "testdata/full-config.yaml",
			jobCount: 1,
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
				t.Errorf("Expected no error loading the resources, got '%s'", err)
			}

			sl, err := JobSlice(rs)
			if err != nil {
				t.Errorf("Expected no error getting the job slice, got '%s'", err)
			}

			if len(sl) != tc.jobCount {
				t.Errorf("Expected %d jobs, got %d", tc.jobCount, len(sl))
			}
		})
	}
}

func TestHashedJobs(t *testing.T) {
	tcs := map[string]struct {
		jobs   []v1.Job
		config map[string]string
		list   map[string]string
	}{
		"without extra config": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
				},
			},
			list: map[string]string{
				"default/foo": "k86kg7tt2c",
			},
		},
		"with config but not linked": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
				},
			},
			config: map[string]string{
				"ConfigMap/default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
			list: map[string]string{
				"default/foo": "k86kg7tt2c",
			},
		},
		"with configmap linked but not configured": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Volumes: []cv1.Volume{
									{
										Name: "linked-configmap",
										VolumeSource: cv1.VolumeSource{
											ConfigMap: &cv1.ConfigMapVolumeSource{
												LocalObjectReference: cv1.LocalObjectReference{
													Name: "perl-job-config",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			list: map[string]string{
				"default/foo": "b6k5k6d596",
			},
		},
		"with existing configmap linked": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Volumes: []cv1.Volume{
									{
										Name: "linked-configmap",
										VolumeSource: cv1.VolumeSource{
											ConfigMap: &cv1.ConfigMapVolumeSource{
												LocalObjectReference: cv1.LocalObjectReference{
													Name: "perl-job-config",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			config: map[string]string{
				"ConfigMap/default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
			list: map[string]string{
				"default/foo": "b6k5k6d596",
			},
		},
		"with existing secret linked": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Volumes: []cv1.Volume{
									{
										Name: "linked-configmap",
										VolumeSource: cv1.VolumeSource{
											Secret: &cv1.SecretVolumeSource{
												SecretName: "mysecret",
											},
										},
									},
									{
										Name: "linked-configmap",
										VolumeSource: cv1.VolumeSource{
											ConfigMap: &cv1.ConfigMapVolumeSource{
												LocalObjectReference: cv1.LocalObjectReference{
													Name: "perl-job-config",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			config: map[string]string{
				"Secret/default/mysecret":           "8ebf17fe046d11996943eee66edbe0a487cb0a7b75f34d3d469ab58649530fbd",
				"ConfigMap/default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
			list: map[string]string{
				"default/foo": "d9dkck4b5b",
			},
		},
		"with existing configmap and secret linked": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Volumes: []cv1.Volume{
									{
										Name: "linked-configmap",
										VolumeSource: cv1.VolumeSource{
											Secret: &cv1.SecretVolumeSource{
												SecretName: "mysecret",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			config: map[string]string{
				"Secret/default/mysecret":           "8ebf17fe046d11996943eee66edbe0a487cb0a7b75f34d3d469ab58649530fbd",
				"ConfigMap/default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
			list: map[string]string{
				"default/foo": "cg8g5t8m84",
			},
		},
		"without env linked": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Containers: []cv1.Container{
									{
										Env: []cv1.EnvVar{
											{
												Name: "my-cm-env-var",
												ValueFrom: &cv1.EnvVarSource{
													ConfigMapKeyRef: &cv1.ConfigMapKeySelector{
														LocalObjectReference: cv1.LocalObjectReference{
															Name: "perl-job-config",
														},
													},
												},
											},
											{
												Name: "my-secret-env-var",
												ValueFrom: &cv1.EnvVarSource{
													SecretKeyRef: &cv1.SecretKeySelector{
														LocalObjectReference: cv1.LocalObjectReference{
															Name: "mysecret",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			list: map[string]string{
				"default/foo": "hh9ktgc5mg",
			},
		},

		"with env configmap": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Containers: []cv1.Container{
									{
										Env: []cv1.EnvVar{
											{
												Name: "my-cm-env-var",
												ValueFrom: &cv1.EnvVarSource{
													ConfigMapKeyRef: &cv1.ConfigMapKeySelector{
														LocalObjectReference: cv1.LocalObjectReference{
															Name: "perl-job-config",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			list: map[string]string{
				"default/foo": "fffg62dhc8",
			},
			config: map[string]string{
				"Secret/default/mysecret":           "8ebf17fe046d11996943eee66edbe0a487cb0a7b75f34d3d469ab58649530fbd",
				"ConfigMap/default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
		},

		"with env secret": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Containers: []cv1.Container{
									{
										Env: []cv1.EnvVar{
											{
												Name: "my-secret-env-var",
												ValueFrom: &cv1.EnvVarSource{
													SecretKeyRef: &cv1.SecretKeySelector{
														LocalObjectReference: cv1.LocalObjectReference{
															Name: "mysecret",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			list: map[string]string{
				"default/foo": "2c22g5g4tt",
			},
			config: map[string]string{
				"Secret/default/mysecret":           "8ebf17fe046d11996943eee66edbe0a487cb0a7b75f34d3d469ab58649530fbd",
				"ConfigMap/default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
		},

		"with multiple envs": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Containers: []cv1.Container{
									{
										Env: []cv1.EnvVar{
											{
												Name: "my-cm-env-var",
												ValueFrom: &cv1.EnvVarSource{
													ConfigMapKeyRef: &cv1.ConfigMapKeySelector{
														LocalObjectReference: cv1.LocalObjectReference{
															Name: "perl-job-config",
														},
													},
												},
											},
											{
												Name: "my-secret-env-var",
												ValueFrom: &cv1.EnvVarSource{
													SecretKeyRef: &cv1.SecretKeySelector{
														LocalObjectReference: cv1.LocalObjectReference{
															Name: "mysecret",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			list: map[string]string{
				"default/foo": "dh98674kcb",
			},
			config: map[string]string{
				"Secret/default/mysecret":           "8ebf17fe046d11996943eee66edbe0a487cb0a7b75f34d3d469ab58649530fbd",
				"ConfigMap/default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
		},

		"with no linked envfroms": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Containers: []cv1.Container{
									{
										EnvFrom: []cv1.EnvFromSource{
											{
												ConfigMapRef: &cv1.ConfigMapEnvSource{
													LocalObjectReference: cv1.LocalObjectReference{
														Name: "perl-job-config",
													},
												},
											},
											{
												SecretRef: &cv1.SecretEnvSource{
													LocalObjectReference: cv1.LocalObjectReference{
														Name: "mysecret",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			list: map[string]string{
				"default/foo": "mt8f6tcg88",
			},
		},

		"with configmap envfrom": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Containers: []cv1.Container{
									{
										EnvFrom: []cv1.EnvFromSource{
											{
												ConfigMapRef: &cv1.ConfigMapEnvSource{
													LocalObjectReference: cv1.LocalObjectReference{
														Name: "perl-job-config",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			list: map[string]string{
				"default/foo": "6hk9ttb952",
			},
			config: map[string]string{
				"Secret/default/mysecret":           "8ebf17fe046d11996943eee66edbe0a487cb0a7b75f34d3d469ab58649530fbd",
				"ConfigMap/default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
		},

		"with secret envfrom": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Containers: []cv1.Container{
									{
										EnvFrom: []cv1.EnvFromSource{
											{
												SecretRef: &cv1.SecretEnvSource{
													LocalObjectReference: cv1.LocalObjectReference{
														Name: "mysecret",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			list: map[string]string{
				"default/foo": "mkh9h2ftf6",
			},
			config: map[string]string{
				"Secret/default/mysecret":           "8ebf17fe046d11996943eee66edbe0a487cb0a7b75f34d3d469ab58649530fbd",
				"ConfigMap/default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
		},

		"with multiple envfroms": {
			jobs: []v1.Job{
				{
					ObjectMeta: mv1.ObjectMeta{
						Name: "foo",
					},
					Spec: v1.JobSpec{
						Template: cv1.PodTemplateSpec{
							Spec: cv1.PodSpec{
								Containers: []cv1.Container{
									{
										EnvFrom: []cv1.EnvFromSource{
											{
												ConfigMapRef: &cv1.ConfigMapEnvSource{
													LocalObjectReference: cv1.LocalObjectReference{
														Name: "perl-job-config",
													},
												},
											},
											{
												SecretRef: &cv1.SecretEnvSource{
													LocalObjectReference: cv1.LocalObjectReference{
														Name: "mysecret",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			list: map[string]string{
				"default/foo": "h9gfkdb5hm",
			},
			config: map[string]string{
				"Secret/default/mysecret":           "8ebf17fe046d11996943eee66edbe0a487cb0a7b75f34d3d469ab58649530fbd",
				"ConfigMap/default/perl-job-config": "6b01af86bab978c892006d41097f29c7b040d459e6613fad29293c1d2c624046",
			},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			hashedJobs, err := HashedJobs(tc.jobs, tc.config)
			if err != nil {
				t.Errorf("Expected no error hashing the jobs, got '%s'", err)
			}

			if !cmp.Equal(tc.list, hashedJobs) {
				t.Errorf("Expected hashes to equal\n\n%s\n\ngot\n\n%s", tc.list, hashedJobs)
			}
		})
	}
}
