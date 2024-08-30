variable "stage" {
  type = string
}
variable "gateway_id" {
  type = string
}
variable "gateway_execution_arn" {
  type = string
}
variable "lambda_iam_arn" {
  type = string
}

locals {
  base_route = "engine"
}

module "lambda_function" {
  source = "../../../infra/modules/lambda"

  gateway_id            = var.gateway_id
  gateway_execution_arn = var.gateway_execution_arn
  lambda_iam_arn        = var.lambda_iam_arn
  stage                 = var.stage

  gateway_route_key = "ANY /${local.base_route}/{proxy+}"
  lambda_base_name  = "shortener-${local.base_route}-handler"
  filepath          = "${path.module}/../function"
  filename          = "function.zip"
  s3_prefix         = "build/lambda/shortener"
  s3_bucket         = "zenhalab-artifacts-${var.stage}"
}

