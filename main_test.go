package main

import (
	"strconv"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetPods(t *testing.T) {
	var results []Pod
	var err error

	results, err = GetPods("./okteto-kube.config", "ascarb")
	if err != nil {
		t.Errorf("Error was not nil.")
	}

	if len(results) == 0 {
		t.Errorf("Results array not populated.")
	} else {
		t.Logf("Found pods")
	}
	t.Logf("Pods: " + strconv.Itoa(len(results)))
}

func TestGetPodAge(t *testing.T) {
	creationTime := metav1.NewTime(time.Now().Add(-1 * time.Minute))

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testpod",
			Namespace: "ascarb",
			Labels: map[string]string{
				"app": "testapp",
			},
			CreationTimestamp: creationTime,
			Annotations: map[string]string{
				"test-annotation": "test-value",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "testcontainer",
					Image: "testimage:latest",
					Ports: []v1.ContainerPort{
						{
							ContainerPort: 8080,
						},
					},
				},
			},
			NodeSelector: map[string]string{
				"disk": "ssd",
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
		Status: v1.PodStatus{
			Phase: v1.PodSucceeded,
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:         "testcontainer",
					RestartCount: 2,
					Ready:        true,
				},
			},
		},
	}

	result := getPodAge(*pod)

	if result == 0 {
		t.Errorf("Age not calculated.")
	}
}

func TestGetRestartCount(t *testing.T) {
	creationTime := metav1.NewTime(time.Now().Add(-1 * time.Minute))

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testpod",
			Namespace: "ascarb",
			Labels: map[string]string{
				"app": "testapp",
			},
			CreationTimestamp: creationTime,
			Annotations: map[string]string{
				"test-annotation": "test-value",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "testcontainer",
					Image: "testimage:latest",
					Ports: []v1.ContainerPort{
						{
							ContainerPort: 8080,
						},
					},
				},
			},
			NodeSelector: map[string]string{
				"disk": "ssd",
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
		Status: v1.PodStatus{
			Phase: v1.PodSucceeded,
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:         "testcontainer",
					RestartCount: 2,
					Ready:        true,
				},
			},
		},
	}

	result := getPodRestarts(*pod)

	if result != 2 {
		t.Errorf("Restart count was not 2.")
	}
}

func TestSort(t *testing.T) {

	var pod1 Pod
	var pod2 Pod
	var pod3 Pod

	pod1.Age = 3500
	pod1.Name = "aaa-testpod1"
	pod1.RestartCount = 5

	pod2.Age = 100
	pod2.Name = "testpod2"
	pod2.RestartCount = 0

	pod3.Age = 2000
	pod3.Name = "b-pod"
	pod3.RestartCount = 100

	pods := []Pod{pod1, pod2, pod3}

	result1 := sortPods("name", pods)
	result2 := sortPods("age", pods)
	result3 := sortPods("restart", pods)

	if len(result1) < 3 {
		t.Errorf("FAIL")
	}

	if len(result2) < 3 {
		t.Errorf("FAIL")
	}

	if len(result3) < 3 {
		t.Errorf("FAIL")
	}

}
