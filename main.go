package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/gin-gonic/gin"
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

func main() {
	fmt.Println("Starting hello-okteto server...")

	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		context.String(200, "Hello Okteto!")
	})

	//TODO: better error handling, and an actual api spec to follow against
	router.GET("/pods", func(context *gin.Context) {
		//TODO: Put the path to important config in a config file, and marshall into struct.
		pods, err := GetPods("./okteto-kube.config", "ascarb")
		if err != nil {
			log.Fatalf("ERROR: " + err.Error())
			context.JSON(http.StatusNotFound, "Error: Not Authorized or Not Found.")
		}
		var podResult []Pod
		if sortParam := context.Query("sort"); strings.ToLower(sortParam) != "" {
			switch sortParam {
			case "name":
				podResult = sortPods("name", pods)
			case "age":
				podResult = sortPods("age", pods)
			case "restart":
				podResult = sortPods("restart", pods)
			}

		} else {
			podResult = append(podResult, pods...)
		}

		context.JSON(http.StatusOK, podResult)

	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf(err.Error())
	}

}

// GetPods Returns a string array, or error, of the names of pods for the kubeconfig passed in.
func GetPods(kubeConfigPath string, namespace string) ([]Pod, error) {
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

	return podList, err
}

// getPodAge gets the age in seconds of the pod since creation to now.
func getPodAge(podInfo v1.Pod) int {
	creationTime := podInfo.ObjectMeta.CreationTimestamp.Time
	age := int(time.Since(creationTime).Seconds())

	return age

}

// getPodRestarts get the number of restarts as listed in the pod metadata container status.
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
