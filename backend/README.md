# kubejiangnan

A Tool for Clusters Managements

Tech Stack: 
-    Golang      version: go1.21.1 darwin/arm64
-   Kubernetes
-   Mac

## Web Framework
``` bash
go get -u github.com/gin-gonic/gin@v1.8.1 
```
> document: https://github.com/gin-gonic/gin


## Config Seperated
```bash
go get github.com/spf13/viper@v1.13.0
```
>document: https://github.com/spf13/viper

## Integrate Kubernetes
```bash
go get k8s.io/client-go@v0.20.4
```
>document: https://github.com/kubernetes/client-go


## APIs Development

### Pod Management APIs

- NameSpaceList APIs
- Pod Create 
- Pod Edit (Update/Upgrade)
- Pod View (Details, List)
- Pod Delete