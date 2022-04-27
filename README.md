# kubviz
Visualize Kubernetes & DevSecOps Workflows. Tracks changes/events real-time across your entire K8s clusters, git repos, container registries, SBOM, Vulnerability foot print, etc. , analyzing their effects and providing you with the context you need to troubleshoot efficiently. Get the Observability you need, easily.


## Install kubviz using Helm:

#### Client Installation
```bash
helm repo add kubviz https://kube-tarian.github.io/kubviz/
helm repo update

helm upgrade -i kubviz-client kubviz/client -n kubviz
```
#### Agent Installation
```bash
helm upgrade -i kubviz-agent kubviz/agent -n kubviz
```
