data "aws_acm_certificate" "issued" {
  domain   = var.issued_certificate_domain
  statuses = ["ISSUED"]
  types    = ["AMAZON_ISSUED"]
}

data "aws_cloudfront_cache_policy" "s3_cache" {
  name = "Managed-CachingOptimized"
}
data "aws_cloudfront_cache_policy" "no_cache_policy" {
  name = "Managed-CachingDisabled"
}
data "aws_cloudfront_response_headers_policy" "cors_policy" {
  name = "Managed-SimpleCORS"
}
data "aws_cloudfront_origin_request_policy" "origin_policy_all" {
  name = "Managed-AllViewer"
}

locals {
  bucket_origin_id = "shortenerbucketstatic"
  api_origin_id = "public-api"
}

resource "aws_cloudfront_distribution" "s3_distribution" {
  origin {
    domain_name = var.bucket_regional_domain_name
    origin_id   = local.bucket_origin_id
    custom_origin_config {
      http_port              = "80"
      https_port             = "443"
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }
  origin {
    domain_name = var.api_cloudfront_origin_domain
    origin_path = "/shortener/v1/public-api/url"
    origin_id   = local.api_origin_id
    custom_origin_config {
      http_port              = "80"
      https_port             = "443"
      origin_protocol_policy = "https-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  enabled             = true
  comment             = "Shortner distribution ${var.stage}"
  default_root_object = "index.html"

  aliases = [var.cloudfront_alias]

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = local.bucket_origin_id

    cache_policy_id        = data.aws_cloudfront_cache_policy.s3_cache.id
    compress               = true
    viewer_protocol_policy = "redirect-to-https"
  }

  ordered_cache_behavior {
    path_pattern     = "p-*"
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = local.api_origin_id

    response_headers_policy_id = data.aws_cloudfront_response_headers_policy.cors_policy.id
    # This legacy config is to allow api-gateway usage
    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  price_class = "PriceClass_100"

  restrictions {
    geo_restriction {
      restriction_type = "none"
      locations        = []
    }
  }

  tags = {
    Stage = var.stage
  }

  viewer_certificate {
    acm_certificate_arn      = data.aws_acm_certificate.issued.arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.2_2021"
  }
}
