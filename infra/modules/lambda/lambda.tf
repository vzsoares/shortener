locals {
  full_file_path = "${var.filepath}/${var.filename}"
  function_name  = "${var.lambda_base_name}-${var.stage}"
  s3_prefix      = var.s3_prefix
  s3_bucket      = var.s3_bucket
  s3_key         = "${local.s3_prefix}/${var.lambda_base_name}/${local.function_name}.zip"
}

resource "aws_lambda_function" "lambda" {
  function_name = local.function_name
  role          = var.lambda_iam_arn
  handler       = "index.handler"

  source_code_hash  = filebase64sha256(local.full_file_path)
  s3_bucket         = local.s3_bucket
  s3_key            = local.s3_key
  s3_object_version = aws_s3_object.lambda_package.version_id

  runtime     = "provided.al2"
  memory_size = "128"

  timeout = 30
  publish = true

  environment {
    variables = {
      STAGE = var.stage
    }
  }

  tags = {
    datadog   = "monitored"
    Terraform = "true"
    stage     = var.stage
  }
  depends_on = [aws_s3_object.lambda_package]
}

resource "aws_s3_object" "lambda_package" {
  bucket = local.s3_bucket
  key    = local.s3_key
  source = local.full_file_path

  etag = filemd5(local.full_file_path)
}

resource "aws_lambda_permission" "lambda" {
  statement_id  = "AllowAPIGatewaySample"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${var.gateway_execution_arn}/*/*"
}

resource "null_resource" "sam_metadata_aws_lambda_function_lambda" {
  triggers = {
    resource_name        = "aws_lambda_function.lambda"
    resource_type        = "ZIP_LAMBDA_FUNCTION"
    original_source_code = var.filepath
    built_output_path    = local.full_file_path
  }
}
