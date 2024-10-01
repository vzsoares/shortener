-include .env
export

.RECIPEPREFIX := >
.DEFAULT_GOAL := help

STAGE ?= dev
TERRAFORM ?= terraform

WORK_DIR_BASE=./infra/environments

deploy: ##@ Deploy to current $STAGE
> @echo "Deploying for stage: ${STAGE}"
> @make init
>  make apply

apply:
> cd ${WORK_DIR_BASE}/${STAGE} && ${TERRAFORM} apply \
  -var='issued_certificate_domain=zenhalab.com' \
  -var='cloudfront_alias=${FRONT_BASE_URL_DOMAIN}' \
  -var='api_cloudfront_origin_domain=${API_BASE_URL_DOMAIN}' \
  -var='front_bucket_name=${FRONT_BASE_URL_DOMAIN}' \
  -var='gateway_api_mapping_domain=${API_BASE_URL_DOMAIN}' \
  -var='gateway_api_name=${GATEWAY_API_NAME}' \
  -var='dynamodb_table_name=${DYNAMO_URL_TABLE_NAME}' \
  -var='artifacts_bucket_name=${ARTIFACTS_BUCKET_NAME}' \
  -auto-approve

init:
> cd ${WORK_DIR_BASE}/${STAGE} && ${TERRAFORM} init


