# stakater-operator-boilerplate
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
  image: stakaterdockerhubpullroot/stakater-operator-boilerplate-catalog:v{VERSION}
  
kind: Subscription
spec:
  startingCSV: stakater-operator-boilerplate.v{VERSION}
```

4. Deploy the `catalogs` and `subscription` manifests.
```sh
oc apply -f misc/olm
```

### Advanced topics
1. [Channels and upgrade strategy](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/doc/design/how-to-update-operators.md)
2. [Advanced indexing & caching](https://github.com/kubernetes-sigs/controller-runtime/blob/master/designs/use-selectors-at-cache.md)
3. [API design best-practices](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md)
