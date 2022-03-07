#!/usr/bin/env bash
# Paulo Aleixo Campos
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
function shw_info { echo -e '\033[1;34m'"$1"'\033[0m'; }
function error { echo "ERROR in ${1}"; exit 99; }
trap 'error $LINENO' ERR
PS4='████████████████████████${BASH_SOURCE}@${FUNCNAME[0]:-}[${LINENO}]>  '
set -o errexit
set -o pipefail
set -o nounset
shopt -s inherit_errexit
set -o xtrace


show_usage_and_exit() {
  cat <<EOT
  USAGE: 
     # Environment should contain these 2 env-vars defined:
     export AZURE_STORAGE_ACCOUNT_NAME=zzzzzz
     export AZURE_STORAGE_PRIMARY_ACCOUNT_KEY=yyyyyyyy
     $0 
EOT
  exit 1
}

main() {
  if [[ "${AZURE_STORAGE_ACCOUNT_NAME:-x}" == "x" ]] 
  then
    show_usage_and_exit
  fi

  AZURE_STORAGE_ACCOUNT_NAME="${AZURE_STORAGE_ACCOUNT_NAME}" 
  AZURE_STORAGE_PRIMARY_ACCOUNT_KEY="${AZURE_STORAGE_PRIMARY_ACCOUNT_KEY}"
  cd "${__dir}"


  # REQUIRES these binaries already installed and available in $PATH
  #   docker
  #   kubectl
  #   helm
  #   kind


  # compile golang code and create docker container (local image)
  cd golang/az-storage-enum
  ./docker-build.sh
  cd "${__dir}"

  # create k8s cluster with kind 
  kind create cluster
    # From this point on, kubectl was configured by kind to point to then newly-created cluster

  # load the local docker-image into kind, so its available to be deployed in k8s
  kind load docker-image az-storage-enum:1.0.0

  # Deploy docker-image using helm-chart.
  # Helm deployment command will take ENV-vars and pass them into the chart, to create a secret, which is read into pod env-vars
  helm upgrade --install --atomic --wait --timeout 30s \
    --reset-values \
    --debug \
    ht-release \
    k8s/ht-chart \
    --set secret.AZURE_STORAGE_ACCOUNT_NAME="${AZURE_STORAGE_ACCOUNT_NAME}" \
    --set secret.AZURE_STORAGE_PRIMARY_ACCOUNT_KEY="${AZURE_STORAGE_PRIMARY_ACCOUNT_KEY}"

  # The pod deployed will list azure-storage account containers and its blob (and their's contents)
  # If pod executes correctly, it's logs will show the blobs info, and the pod will stay "sleeping" for 1h (ugly but functional)
  # If the pod cannot read the blobs (ex: credentials error) then it will restart continuously
  kubectl logs --selector=app.kubernetes.io/instance=ht-release --tail=99

#  # In the end, delete kind cluster
#  kind delete cluster

  shw_info "== Finished successfully without errors"
  shw_info "To stop-and-destroy the kind cluster, manually run 'kind delete cluster'"
}
main "${@}"




