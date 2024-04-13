# Overview

## k8s-playground-server 

To build the container run:

```
docker build --load --tag johnwesonga/k8s-playground-server:latest .
```

This will build an arm64 image with the tag `johnwesonga/k8s-playground-server:latest` and load it into your local docker registry.

To run the docker container:
```
docker run -p 8080:8080 johnwesonga/k8s-playground-server:latest
```


### Kubernetes minikube

1) Set the environment variables with ```eval $(minikube docker-env)```
2) Build the image with the Docker daemon of Minikube (e.g., ```docker build -t my-image .```)
3) Set the image in the pod specification like the build tag (e.g., my-image)
4) Set the imagePullPolicy to Never, otherwise Kubernetes will try to download the image.

**Important note**: You have to run ```eval $(minikube docker-env)``` on each terminal you want to use, since it only sets the environment variables for the current shell session.


Deploy to minikube using kustomize:

```cd k8s-playground-server```

```kustomize build k8s/ | kubectl apply -f - ```

Running ingress on minikube apple silicon https://stackoverflow.com/questions/75204589/minikube-ingress-on-macos-appears-to-be-working-but-never-connects-times-out-no