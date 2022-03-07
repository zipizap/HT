# Demo:

## REQUIRES:
- docker
- kind
- kubectl
- helm

## Git clone
```
git clone https://github.com/zipizap/HT.git
cd HT
```

## Build infra with terraform
```


# Login into your azure-account with:
az login
 Create infra: creates 1 storage-account, with 1 container, containaing 1 blob
cd terraform/
./terraform_do.sh init plan apply

export AZURE_STORAGE_ACCOUNT_NAME="$(terraform output AZURE_STORAGE_ACCOUNT_NAME| tr -d '"')"
export AZURE_STORAGE_PRIMARY_ACCOUNT_KEY="$(terraform output AZURE_STORAGE_PRIMARY_ACCOUNT_KEY | tr -d '"')"

cd ..
```


## From Docker to k8s pod
It will:
- Compile Golang program
- build docker image
- create k8s (kind) cluster
- deploy helm-chart with pod, using secret from env-vars

```
./compileGolang.buildDocker.deployKindHelmChart.sh
```

NOTE: The .sh scripts stop inmediately if there is an error in any command. And they are veeery verbose

After script executes successfully the kind cluster will have a pod in default namespace, with logs showing the storage-account cointainers and blobs and theirs contents


## Final cleanup
```
# Delete kind cluster
kind delete cluster

# Destroy infra
cd terraform
./terraform_do.sh destroy 
cd ..

```

