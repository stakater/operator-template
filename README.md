# operator-boilerplate
Boilerplate project for kick-starting your operator development

## Description
For effective development of operator. 
This project includes some common utils to handle reconciliation and caching.

## Getting Started
You’ll need a Kubernetes/Openshift cluster to run against.
### Running locally
1. Install Instances of Custom Resources:
```sh
kubectl apply -f config/samples/
```
2. Login to your development cluster, preferably CRC or SNO
3. Run the operator
```sh
make install run
```

### Running on the cluster with OLM
1. Add semver using `VERSION` variable in `MakeFile`:

```sh
VERSION?= x.x.x eg 0.1.0
```

2. Specify docker repo baseURL by modifying `IMAGE_TAG_BASE` variable in `MakeFile`:

3. Build and push all images:

```sh
make publish
```
4. Modify `catalog.yaml` and `subscription.yaml` with the same `VERSION` in `MakeFile`.
```sh
kind: CatalogSource
spec:
  image: stakaterdockerhubpullroot/operator-boilerplate-catalog:v{VERSION}
  
kind: Subscription
spec:
  startingCSV: operator-boilerplate.v{VERSION}
```

4. Deploy the `catalogs` and `subscription` manifests.
```sh
oc apply -f misc/olm
```

### Status updater
A small example on how to update conditions and event in the status of the resource.
Look into [resourcewatcher_types.go](api%2Fv1alpha1%2Fresourcewatcher_types.go) for required status field
Look into [statusupdater_controller.go](controllers%2Fstatusupdater_controller.go) ```SetupWithManager``` where metadata 
changes are ignored to allow frequent updates in status fields without triggering controller recocile

### Resource watcher
A small example on how to watch for externally manged resource by indexing all watchers.
For each change all watchers in the same namespace as will be reconciled

### Advanced topics
1. [Channels and upgrade strategy](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/doc/design/how-to-update-operators.md)
2. [Advanced indexing & caching](https://github.com/kubernetes-sigs/controller-runtime/blob/master/designs/use-selectors-at-cache.md)
3. [API design best-practices](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
4. [In depth information about caching related details](https://medium.com/@timebertt/kubernetes-controllers-at-scale-clients-caches-conflicts-patches-explained-aa0f7a8b4332)