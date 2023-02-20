package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//Pod String struct of k8s pod info.
type Pod struct {
	Name         string
	Age          int
	RestartCount int
}

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

//GetPods Returns a string array, or error, of the names of pods for the kubeconfig passed in.
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
		//pod.Age = podInfo.
		podList = append(podList, pod)
	}

	var podStringList []string
	for _, podInfo := range pods.Items {
		fmt.Println(podInfo.Name)
		podStringList = append(podStringList, podInfo.Name)
	}
	return podStringList, err
}

func getPodAge(podInf v1.Pod) int {

}
