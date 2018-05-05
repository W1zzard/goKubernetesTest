package main

import (
	"flag"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// optional - local kubeconfig for testing
var kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig file")

func main() {

	// send logs to stderr so we can use 'kubectl logs'
	flag.Set("logtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()

	config, err := getConfig("C:/Users/Wizzard/.kube/config")
	if err != nil {
		glog.Errorf("Failed to load client config: %v", err)
		return
	}

	// build the Kubernetes client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorf("Failed to create kubernetes client: %v", err)
		return
	}

	getSomething(client)

	/*// list pods
	pods, err := client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("Failed to retrieve pods: %v", err)
		return
	}

	for _, p := range pods.Items {
		glog.V(3).Infof("Found pods: %s/%s", p.Namespace, p.Name)
	}*/
}

func getConfig(kubeconfig string) (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func getSomething(client *kubernetes.Clientset)  () {
	bindings, err := client.RbacV1().ClusterRoleBindings().List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("Failed to get binding list: %v", err)
		return
	}
	for _, p := range bindings.Items {
		glog.V(3).Infof("Found name: %s  SUBJECTS: %s", p.Name, p.Subjects)
	}

}

/*func GetInClusterConfig() (*rest.Config, error) {
	host := "192.168.0.149"
	port := "8443"

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tNDJsYjkiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImUyZTE3MjQ2LTUwNGMtMTFlOC05NDA3LTAwMTU1ZDAwMmEwMyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.UMgBhV6aGFnLIcY8MhcGsnaTpJPji5VvtFI_lSd9RxNwNAvIQN4WKS92l_ryaaoKbQkqlWbXCc5HEXYQ4OBoZEpgoPXvSd6eujn1td-WLSXLVwvhwYUudg8Y5YgsL4M0k9vewftm9SOFYh2v5MsfhKpcaXOeWmJ0axtxJHcKUvuhiwjb6ZO_a8aCzmoiJhZL0gbtKP7PqBy4XjRNkamqEVDudRgoVR4fbF_CdwfDoKK9WS2oq_ocLD5UP_cqZOezwBzzoLMjNtDfJ7LzuOT5Wz6P78Ll-tyOro36UdAtN8uqNkVn3vXKRXTOzAAL37blDc7ISrod_vGm3I2kZ5R-EA"
	tlsClientConfig := rest.TLSClientConfig{}
	rootCAFile := "C:/minikube_home/.minikube/ca.crt"
	if _, err := certutil.NewPool(rootCAFile); err != nil {
		glog.Errorf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host:            "https://" + net.JoinHostPort(host, port),
		BearerToken:     string(token),
		TLSClientConfig: tlsClientConfig,
	}, nil
}*/

