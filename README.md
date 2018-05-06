1. use "dep" dependency manager to get proper version of used client go and its transitive dependencies like: `dep ensure`
2. For authorization used config with certificates, in this case its **minikube** username

How program works
=============
1.  To run program user should provide path to kubernetes config
`SomeApp.exe -kubeconfig=C:\Users\Wizzard\.kube\config`

2.  program run several steps, every step in namespace `default`:

    2.1. It shows all roleBinding objects for namespace
     
    2.2. It creates role with name `deployment-manager`
    
    2.3. It binds role to user with name `someEmployee`
    
    2.4. It shows all roleBinding again to indicate that role was binded    

3. after that, to use this `someEmployee` and its role you need:

    3.1. create certificates and sign it with kubernetes CA cert, **NOTE** that CN= name should be same as user name in RoleBinding!
                     
          openssl.exe genrsa -out someEmployee.key 2048
          openssl.exe req -new -key someEmployee.key -out someEmployee.csr -subj "/CN=someEmployee/O=test"
          openssl.exe x509 -req -in someEmployee.csr -CA ./ca.crt -CAkey ./ca.key -CAcreateserial -out someEmployee.crt -days 500
    3.2. create user specific context in kubectl config
          
          kubectl config set-credentials employee --client-certificate=./someEmployee.crt  --client-key=./someEmployee.key
          kubectl config set-context some-employee-context --cluster=minikube --namespace=default --user=employee
          kubectl --context=some-employee-context get pods
    3.3 `kubectl --context=some-employee-context get pods`  **OK**
    
    3.4 `kubectl --context=some-employee-context get secrets` **FORBIDDEN**
            

Restrictions
=====
1. program cant work without kubernetes config
2. program creates only one POC role and roleBinding to show that program itself works, Role verbs selected just for demonstration
3. everything happensin `default` namespace                   