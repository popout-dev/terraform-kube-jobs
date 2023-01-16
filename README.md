# terraform-kube-jobs

A helper chart for applying and destroying terraform resources with helm

## How to use

This chart should ideally be a dependency of a service chart so that it can create and destroy terraform resources. This chart takes in terraform files as a config map and allows for state to be tracked as a kubernetes secret.

```
helm repo add terraform-jobs  https://raw.githubusercontent.com/popout-dev/terraform-kube-jobs/gh-pages
```