-include .env
export

.RECIPEPREFIX := >
.DEFAULT_GOAL := help

DOCKER ?= docker
AWS ?= aws
SAM ?= sam
TERRAFORM ?= terraform
AWK ?= awk
GO ?= go
MAKE ?= make

STAGE ?= dev
DEBUG ?= DEBUG

BASE_TF_WORK_DIR=../../infra/environments

SAM_COMMON_AGRS=--hook-name terraform --beta-features
SAM_TERRAFORM_ROOT=--terraform-project-root-path ../../../

##@
##@ Build
##@

build: ##@ Build package
> @echo "Building for stage: ${STAGE}"
> @cd ./function && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc,${BUILD_TAG}
> @cd ./function && zip -X -r -j function.zip bootstrap

##@
##@ Run local SAM
##@

sam-build: ##@ Build for sam
> @make build BUILD_TAG=${DEBUG} && cd ${BASE_TF_WORK_DIR}/${STAGE} && ${SAM} build ${SAM_COMMON_AGRS} ${SAM_TERRAFORM_ROOT} ${EXTRA}

start-api-debug: ##@ start with listen and serve
> @cd ./function && STAGE=dev go run -tags=DEBUG, .

start: ##@ start local http api
> @npx nodemon -e go --signal SIGTERM --exec 'make' start-api-debug

start-sam: ##@ Start local api
       ##@ - if this is your first run then
       ##@ run make tf-init && make sam-build
       ##@ runs in DEBUG by default to disable do `$ make start DEBUG=""`
> @make watch-build-sam start-api-sam -j

start-api-sam:
> @cd ${BASE_TF_WORK_DIR}/${STAGE} && ${SAM} local start-api --debug --skip-pull-image ${SAM_COMMON_AGRS} --add-host=host.docker.internal:host-gateway

watch-build-sam:
> @npx nodemon -e go --signal SIGTERM --exec 'make' sam-build EXTRA=--skip-prepare-infra

##@
##@ Deployment
##@

deploy: ##@ Deploy to current $STAGE
> @echo "Deploying for stage: ${STAGE}"
> @make build
> @cd ${BASE_TF_WORK_DIR}/${STAGE} && ${TERRAFORM} init && ${TERRAFORM} apply -auto-approve

##@
##@ Terraform
##@

tf-%: ##@ run Terraform %
> @cd ${BASE_TF_WORK_DIR}/${STAGE} && ${TERRAFORM} $*

##@
##@ Misc commands
##@

test:
> @go test ./... -v -tags DEBUG,unit

help: ##@ (Default) This help menu
> @printf "\nUsage: STAGE=dev make <command>\n"
> @grep -F -h "##@" $(MAKEFILE_LIST) | grep -F -v grep -F | sed -e 's/\\$$//' | $(AWK) 'BEGIN {FS = ":*[[:space:]]*##@[[:space:]]*"}; \
	{ \
		if($$2 == "") \
			pass; \
		else if($$0 ~ /^#/) \
			printf "\n%s\n", $$2; \
		else if($$1 == "") \
			printf "     %-20s%s\n", "", $$2; \
		else \
      printf "\n    \033[38;2;156;207;216m%-20s\033[0m %s\n", $$1, $$2; \
	}'
