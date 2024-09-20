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
  source = "../../services/lambda-iam-role"

  stage = local.stage
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

module "public-api-lambda" {
  source = "../../../apps/public-api/infra"

  stage                 = local.stage
  gateway_id            = module.api_gateway.id
  gateway_execution_arn = module.api_gateway.execution_arn
  lambda_iam_arn        = module.role.iam_role_arn
}

module "front_bucket" {
  source = "../../services/front-bucket"
  stage  = local.stage
}

module "cloudfront-distribution" {
  source = "../../services/cloudfront"

  stage                       = local.stage
  bucket_regional_domain_name = module.front_bucket.website_endpoint
}
