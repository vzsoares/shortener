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
> @GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap -tags ${BUILD_TAG},

##@
##@ Run local
##@

css-watch:
> @npx tailwindcss -i ./src/styles.css -o ./dist/styles.css --watch

watch-build:
> @npx nodemon -e go,go.html --exec go run .

##@
##@ Misc commands
##@

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
