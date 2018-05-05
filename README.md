1. use "dep" dependency manager to get proper version of used client go and its transitive dependencies like: `dep ensure`
2. For authorization used config with certificates, in this case its **minikube** username
3. for authentication CA cert is needed in Windows its in `{minikube_home}/.minikube/ca.crt`