# Getting Started on Okteto with Go

[![Develop on Okteto](https://okteto.com/develop-okteto.svg)](https://cloud.okteto.com/deploy?repository=https://github.com/okteto/go-getting-started)

This example shows how to use the [Okteto CLI](https://github.com/okteto/okteto) to develop a Go Sample App directly in Kubernetes. The Go Sample App is deployed using Kubernetes manifests.

This is the application used for the [Getting Started on Okteto with Go](https://www.okteto.com/docs/samples/golang/) tutorial.

Two REST endpoints are specified:

/
    Returns "Hello Okteto!" for liveness check.
/pods
    /pods endpoint returns pods json defined in the okteto/k8s context along with their age in seconds and number of restarts.  
    /pods?sort=[name|age|restart] to sort list of pods by age, name, or restarts. 