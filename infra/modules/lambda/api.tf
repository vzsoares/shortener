resource "aws_apigatewayv2_integration" "gateway_integration" {
  api_id = var.gateway_id

  integration_uri    = aws_lambda_function.lambda.arn
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
}

resource "aws_apigatewayv2_route" "gateway_route" {
  api_id = var.gateway_id

  route_key = var.gateway_route_key
  target    = "integrations/${aws_apigatewayv2_integration.gateway_integration.id}"
}
