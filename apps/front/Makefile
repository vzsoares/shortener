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

build-go: ##@ Build package
> @echo "Building for stage: ${STAGE}"
> @GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap -tags ${BUILD_TAG},

build:
> @make rm-dist
> @make build-go
> @./bootstrap
> @make css-build-mini

css-build-mini:
> @npx tailwindcss -i ./src/css/styles.css -o ./dist/css/styles.css --minify
css-build:
> @npx tailwindcss -i ./src/css/styles.css -o ./dist/css/styles.css

##@
##@ Run local
##@

dev:
> @make serve-watch watch-build -j

serve-watch:
> @cd dist && npx live-server

watch-build:
> @npx nodemon -e tailwind.config.js,go.html,css,go --ignore "dist/*" --exec go run .


##@
##@ Deploy
##@

publish:
> @make sync-bucket
> @make invalidate-distribution

sync-bucket:
> $(eval BUCKET := $(shell [[ ${STAGE} = "dev" ]] && echo "s-dev.zenhalab.com" || echo "s.zenhalab.com"))
> @aws s3 sync --acl=public-read ./dist s3://${BUCKET}/

invalidate-distribution:
> $(eval DISTRI := $(shell [[ ${STAGE} = "dev" ]] && echo "E1WR5S0OYB2OKY" || echo "E1BLLPJRDAF46D"))
> @aws cloudfront create-invalidation --distribution-id=${DISTRI} --paths="/*" | cat

##@
##@ Misc commands
##@

rm-dist:
> @rm -rf ./dist

test:
> @go test ./... -v -tags DEBUG,

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
