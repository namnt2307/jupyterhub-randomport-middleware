package kubernetes

import (
	"context"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func InitKubernetes() *kubernetes.Clientset {
	// kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// log.Println("Using kubeconfig ", kubeconfig)

	// Load kubeconfig
	log.Println("Loading service account")
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	//Load clientset
	log.Println("Load kubeconfig successfully \t Creating Clientset")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}
func MakePodSpec(namespace, podName, nodeSelector, cpuLimit, cpuRequest, memoryLimit, memoryRequest string) *v1.Pod {
	podResource := v1.ResourceRequirements{
		Limits: v1.ResourceList{
			"cpu":    resource.MustParse(cpuLimit),
			"memory": resource.MustParse(memoryLimit),
		},
		Requests: v1.ResourceList{
			"cpu":    resource.MustParse(cpuRequest),
			"memory": resource.MustParse(memoryRequest),
		},
	}
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            podName,
					Image:           "repo.bigdata.local/sleep:latest",
					ImagePullPolicy: v1.PullIfNotPresent,
					// Command: []string{
					// 	"sleep",
					// 	"60",
					// },
					Resources: podResource,
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
			NodeSelector:  map[string]string{nodeSelector: "true"},
		},
	}
}

func MakePod(clientset *kubernetes.Clientset, namespace, podName, nodeSelector, cpuLimit, cpuRequest, memoryLimit, memoryRequest string) string {
	//make pod spec
	pod := MakePodSpec(namespace, podName, nodeSelector, cpuLimit, cpuRequest, memoryLimit, memoryRequest)
	// create pod
	_, err := clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	// wait until pod is create
	time.Sleep(10 * time.Second)
	// log.Printf("%v \n", podCreate.Status.HostIP)
	// Check pod

	nodeIP, err := clientset.CoreV1().Pods(pod.Namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	var zero int64 = 0
	clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{GracePeriodSeconds: &zero})
	return nodeIP.Status.HostIP

}
