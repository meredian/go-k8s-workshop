APP?=austin
RELEASE?=0.0.1

NAMESPACE?=workshop
REGISTRY?=registry.k8s.gromnsk.ru
PROJECT=gitlab.k8s.gromnsk.ru/${NAMESPACE}/${APP}
GOOS?=linux
GOARCH?=amd64

CONTAINER_IMAGE?=${REGISTRY}/${NAMESPACE}/${APP}

.PHONY: build
build: clean
	@echo "+ $@"
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -a -installsuffix cgo \
		-o bin/${GOOS}-${GOARCH}/${APP} cmd/tracking/main.go
	docker build --pull -t $(CONTAINER_IMAGE):$(RELEASE) .

clean:
	rm -rf bin/${GOOS}-${GOARCH}/${APP}

.PHONY: push
push: build
	@echo "+ $@"
	docker push $(CONTAINER_IMAGE):$(RELEASE)

.PHONY: deploy
deploy: push
	helm upgrade ${APP} -f charts/values.yaml charts --namespace ${NAMESPACE} --version=${RELEASE} -i --wait
