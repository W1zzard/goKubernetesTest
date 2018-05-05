package main

import (
	"flag"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net"
)

// optional - local kubeconfig for testing
var kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig file")

func main() {

	// send logs to stderr so we can use 'kubectl logs'
	flag.Set("logtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()

	config, err := getConfig(*kubeconfig)
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

	// list pods
	pods, err := client.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("Failed to retrieve pods: %v", err)
		return
	}

	for _, p := range pods.Items {
		glog.V(3).Infof("Found pods: %s/%s", p.Namespace, p.Name)
	}
}

func getConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	return GetInClusterConfig()
}

func GetInClusterConfig() (*rest.Config, error) {
	host := "192.168.0.142"
	port := "8443"

	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ" +
		"2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1" +
		"bHQtdG9rZW4tcDlidjkiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1Y" +
		"mVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImFkMjYzN2QyLTUwMjktMTFlOC1iNGEyLTAwMTU1ZDAwMm" +
		"EwMSIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.ktc03NwATlU-r_7OypMwkSxwLqGDy3f4vEKsIm4a" +
		"CvKVoD_X-ecaHGs54KpJwyLt_-SBDm9j_Xqi2jKbHiHtdNKXw80ejB6WmmkHASNupsEGdqiaihVkCert7b5_zQtjswvbzJvBLf7lUCDK9aYKaY" +
		"V2AUEEFYI4CkxYSaEtuPdTnpKPoZatZVFC5RiJeBeEgJGfCAyMlBmB1BJr4XPaasXOIsaIEl8CruA2jn1-jBbNJJ0TXZN6M0xvp5lf2CMmcJDl" +
		"GPDuH8ekauuPmZfJ3z3jmEyZd5rjb_wJVswfu7AUCpa4xg5xn2uL2DPUZdllkw88s-2tkQoN0v-q9ZaVkA"
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
}