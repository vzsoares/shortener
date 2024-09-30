-include .env
export

.RECIPEPREFIX := >
.DEFAULT_GOAL := help

STAGE ?= dev
TERRAFORM ?= terraform

WORK_DIR_BASE=./infra/environments

deploy: ##@ Deploy to current $STAGE
> @echo "Deploying for stage: ${STAGE}"
#> @make init
>  make apply

apply:
> cd ${WORK_DIR_BASE}/${STAGE} && ${TERRAFORM} apply \
  -var='issued_certificate_domain=zenhalab.com' \
  -var='cloudfront_alias=${FRONT_BASE_URL_DOMAIN}' \
  -var='api_cloudfront_origin_domain=${API_BASE_URL_DOMAIN}'

#-auto-approve

init:
> cd ${WORK_DIR_BASE}/${STAGE} && ${TERRAFORM} init


