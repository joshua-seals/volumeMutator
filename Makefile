# Command line to get all certs in one command for TLS (you must actively enter prompt info)
# openssl req -x509 -nodes -newkey rsa:2048 -keyout server.rsa.key -out server.rsa.crt -days 3650

## Environment 
CERT_PATH := /helx/webhook/certs/
ENVIRONMENT := develop
BASE_IMAGE := helxplatform/volume-mutator
IMAGE_TAG := containers.renci.org/$(BASE_IMAGE)
VERSION := 0.0.1

## Kind Related
KIND_CLUSTER := helx-testing


build-arm:
	docker buildx build \
	--platform=linux/arm64 \
	--build-arg=BUILD_REF=$(ENVIRONMENT) \
	--build-arg=CERT_PATH=$(CERT_PATH) \
	--build-arg=BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
	--tag=$(IMAGE_TAG):$(VERSION) \
	--tag=$(BASE_IMAGE):$(VERSION) \
	.

build:
	docker buildx build \
	--platform=linux/amd64 \
	--build-arg=BUILD_REF=$(ENVIRONMENT) \
	--build-arg=CERT_PATH=$(CERT_PATH) \
	--build-arg=BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
	--tag=$(IMAGE_TAG):$(VERSION) \
	--tag=$(BASE_IMAGE):$(VERSION) \
	.


kind-up:
	kind create cluster --name $(KIND_CLUSTER)

kind-load:
	kind load docker-image $(BASE_IMAGE):$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kubectl apply -f k8s/volume-mutator.yml

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-default-CRB:
	kubectl create clusterrolebinding serviceaccounts-cluster-admin \
	--clusterrole=cluster-admin \
	--group=system:serviceaccounts