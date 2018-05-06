package main

import (
	"flag"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/api/rbac/v1"
)

// local kubeconfig for testing
var kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig file")

func main() {

	// send logs to stderr so we can use 'kubectl logs'
	flag.Set("logtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()
	if (*kubeconfig == "") {
		glog.Errorf("Please provide path to config  -kubeconfig=[path]")
		return
	}

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

	// show roleBindings already existed
	printRoleBinding(client)

	// add role and rolebinding
	role := createDeploymentRole()
	createRoleOnKub(client, role)
	bindRole(client, role, "someEmployee")

	// show roleBindings again to check that it was added
	printRoleBinding(client)


}

func getConfig(kubeconfig string) (*rest.Config, error) {
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

/*print list of role bindings*/
func printRoleBinding(client *kubernetes.Clientset)  () {
	bindings, err := client.RbacV1().RoleBindings("default").List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("Failed to get binding list: %v", err)
		return
	}
	for _, p := range bindings.Items {
		glog.V(3).Infof("Found name: %s  SUBJECTS: %s", p.Name, p.Subjects)
	}

}

/*Create role with parameters on Kubernetes*/
func createRoleOnKub(client *kubernetes.Clientset, role *v1.Role) (*v1.Role){

	roleCreationResult, err := client.RbacV1().Roles("default").Create(role)
	if err != nil {
		glog.Errorf("Failed to create role: %v", err)
		return nil
	}
	return roleCreationResult
}

/*Create role structure which will be created on kubernetes*/
func createDeploymentRole() (*v1.Role) {
	rule := v1.PolicyRule{
		Verbs:           []string{"get", "list", "watch", "create", "update", "patch", "delete"},
		APIGroups:       []string{"", "extensions", "apps"},
		Resources:       []string{"deployments", "replicasets", "pods"},
	}

	return &v1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-manager",
			Namespace: "default",
		},
		Rules: []v1.PolicyRule{rule},
	}
}

/*Bind some Role which is already on Kubernetes to user name provided*/
func bindRole(client *kubernetes.Clientset, role *v1.Role, userName string) {
	subject := v1.Subject{
		Kind:      "User",
		APIGroup:  "",
		Name:      userName,
		Namespace: "default",
	}

	_, err := client.RbacV1().RoleBindings("default").Create(&v1.RoleBinding{

		ObjectMeta: metav1.ObjectMeta{
			Name: role.Name + "-binding",
		},
		Subjects: []v1.Subject{subject},
		RoleRef: v1.RoleRef{
			APIGroup: "",
			Kind:     "Role",
			Name:     role.Name,
		},
	})
	if err != nil {
		glog.Errorf("Failed to bind role: %v", err)
	}
}

