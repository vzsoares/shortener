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

# front cloudfront
data "aws_acm_certificate" "issued" {
  domain   = "zenhalab.com"
  statuses = ["ISSUED"]
  types    = ["AMAZON_ISSUED"]
}

resource "aws_cloudfront_origin_access_control" "s3_origin_control" {
  name                              = "s3_origin_control"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "s3_distribution" {
  origin {
    domain_name              = module.front_bucket.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.s3_origin_control.id
    origin_id                = "shortenerbucket"
  }
  origin {
    domain_name = "api${local.stage == "dev" ? "-dev" : ""}.zenhalab.com"
    origin_path = "/shortener/v1/public-api/url"
    origin_id   = "public-api"
    custom_origin_config {
      http_port              = "80"
      https_port             = "443"
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  enabled             = true
  comment             = "Shortner distribution"
  default_root_object = "index.html"

  aliases = ["s${local.stage == "dev" ? "-dev" : ""}.zenhalab.com"]

  default_cache_behavior {
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "shortenerbucket"

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
    forwarded_values {
      query_string = true

      cookies {
        forward = "all"
      }
    }
  }

  ordered_cache_behavior {
    path_pattern     = "/p-*"
    allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "public-api"

    forwarded_values {
      query_string = true

      cookies {
        forward = "all"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
    compress               = true
  }

  price_class = "PriceClass_100"

  restrictions {
    geo_restriction {
      restriction_type = "none"
      locations        = []
    }
  }

  tags = {
    Environment = local.stage
  }

  viewer_certificate {
    acm_certificate_arn      = data.aws_acm_certificate.issued.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }
}
