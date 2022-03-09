package kubernetes

import (
	"context"
	"log"
	"net"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func InitKubernetes() *kubernetes.Clientset {
	// kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	// log.Println("INFO\t", "Using kubeconfig ", kubeconfig)
	// Load kubeconfig
	// config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("INFO\t", "Loading Service account")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	log.Println("INFO\t Load Service account successfully")

	//Load clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INFO\t Creating Clientset")

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
					Env: []v1.EnvVar{
						{
							Name: "hihi",
							ValueFrom: &v1.EnvVarSource{
								FieldRef: &v1.ObjectFieldSelector{
									APIVersion: "v1",
									FieldPath:  "spec.nodeName",
								},
							},
						},
					},
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
			NodeSelector:  map[string]string{nodeSelector: "true"},
		},
	}
}
func GetFreePort(hostIP string) (int, error) {
	host := hostIP + ":0"
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
func MakePod(clientset *kubernetes.Clientset, namespace, podName, nodeSelector, cpuLimit, cpuRequest, memoryLimit, memoryRequest string) (podname, nodeselector, hostIP, hostName string, hostPort int) {
	//make pod spec
	pod := MakePodSpec(namespace, podName, nodeSelector, cpuLimit, cpuRequest, memoryLimit, memoryRequest)
	// create pod
	_, err := clientset.CoreV1().Pods(pod.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	// wait until pod is create
	time.Sleep(300 * time.Millisecond)

	// log.Printf("%v \n", podCreate.Status.HostIP)
	// Check pod

	podSpec, err := clientset.CoreV1().Pods(pod.Namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	port, err := GetFreePort(podSpec.Status.HostIP)
	if err != nil {
		log.Println(err)
	}
	var zero int64 = 0
	clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{GracePeriodSeconds: &zero})
	return podName, nodeSelector, podSpec.Status.HostIP, podSpec.Spec.NodeName, port

}
