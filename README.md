## Overview

### k8s-playground-server 

To build the container to run on Apple silicon:

```
docker build --load --tag johnwesonga/k8s-playground-server:latest .
```

To run the docker container:
```
docker run -p 8080:8080 johnwesonga/k8s-playground-server:latest
```