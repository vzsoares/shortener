-include .env
export

.RECIPEPREFIX := >
.DEFAULT_GOAL := help

STAGE ?= dev
TERRAFORM ?= terraform

WORK_DIR_BASE=./infra/environments

deploy: ##@ Deploy to current $STAGE
> @echo "Deploying for stage: ${STAGE}"
> @cd ${WORK_DIR_BASE}/${STAGE} && ${TERRAFORM} init && ${TERRAFORM} apply -auto-approve

