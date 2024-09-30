## API Gateway
resource "aws_apigatewayv2_api" "gateway" {
  name          = "${var.gateway_name}-${var.stage}"
  protocol_type = "HTTP"
  cors_configuration {
    allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }
  disable_execute_api_endpoint = true

  tags = {
    Terraform = "true"
    stage     = var.stage
  }
}

resource "aws_apigatewayv2_stage" "gateway" {
  api_id = aws_apigatewayv2_api.gateway.id

  name        = "v1"
  auto_deploy = true
}

resource "aws_apigatewayv2_api_mapping" "mapping" {
  api_id          = aws_apigatewayv2_api.gateway.id
  domain_name     = var.gateway_api_mapping_domain
  stage           = aws_apigatewayv2_stage.gateway.id
  api_mapping_key = "shortener/v1"
}
