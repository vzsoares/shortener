locals {
  region     = "us-east-1"
  stage      = "dev"
  account_id = "355738159777"
}

module "urls-table" {
  source = "../../services/urls-table"
  stage  = local.stage
}

module "role" {
  source    = "../../services/lambda-iam-role"
  table-arn = module.urls-table.arn
}

module "api_gateway" {
  source       = "../../services/gateway"
  stage        = local.stage
  gateway_name = "shortener"
}

module "engine-lambda" {
  source = "../../../apps/engine/infra"

  stage                 = local.stage
  gateway_id            = module.api_gateway.id
  gateway_execution_arn = module.api_gateway.execution_arn
  lambda_iam_arn        = module.role.iam_role_arn
}

