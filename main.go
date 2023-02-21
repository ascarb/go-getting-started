package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Pod String struct of k8s pod info.
type Pod struct {
	Name         string
	Age          int
	RestartCount int
}

type ByName []Pod
type ByAge []Pod
type ByRestartCount []Pod

const (
	Age     = 0
	Name    = 1
	Restart = 2
)

func main() {
	fmt.Println("Starting hello-okteto server...")
	//http.HandleFunc("/", helloServer)
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	panic(err)
	// }

	results, err := GetPods("/Users/adam/code/go-getting-started/okteto-kube.config", "ascarb")
	if err != nil {
		panic(err)
	}

	fmt.Println(len(results))

}

func helloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Okteto!")
}

// GetPods Returns a string array, or error, of the names of pods for the kubeconfig passed in.
func GetPods(kubeConfigPath string, namespace string) ([]string, error) {
	//TODO:read path and namespace from config file.
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("ERROR: Kube config not found!!!")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset.")
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("ERROR: error getting pod list. " + err.Error())

	}

	var podList []Pod
	for _, podInfo := range pods.Items {
		var pod Pod
		pod.Name = podInfo.Name
		pod.Age = getPodAge(podInfo)
		pod.RestartCount = getPodRestarts(podInfo)
		podList = append(podList, pod)
	}

	var podStringList []string
	for _, podInfo := range pods.Items {
		fmt.Println(podInfo.Name)
		podStringList = append(podStringList, podInfo.Name)
	}
	return podStringList, err
}

func getPodAge(podInfo v1.Pod) int {
	creationTime := podInfo.ObjectMeta.CreationTimestamp.Time
	age := int(time.Since(creationTime).Seconds())

	return age

}

func getPodRestarts(podInfo v1.Pod) int {
	restarts := 0

	for _, containerStatus := range podInfo.Status.ContainerStatuses {
		restarts += int(containerStatus.RestartCount)
	}

	return restarts
}

func sortPods(sortType string, pods []Pod) []Pod {
	if sortType == "name" {
		sort.Sort(ByName(pods))
	}
	if sortType == "age" {
		sort.Sort(ByAge(pods))
	}
	if sortType == "restart" {
		sort.Sort(ByRestartCount(pods))
	}
	return pods
}

// Sorting interface functions
func (p ByName) Len() int {
	return len(p)
}

func (p ByName) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByName) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}

func (p ByAge) Len() int {
	return len(p)
}

func (p ByAge) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByAge) Less(i, j int) bool {
	return p[i].Age < p[j].Age
}

func (p ByRestartCount) Len() int {
	return len(p)
}

func (p ByRestartCount) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByRestartCount) Less(i, j int) bool {
	return p[i].RestartCount < p[j].RestartCount
}
