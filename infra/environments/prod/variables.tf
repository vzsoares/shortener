variable "issued_certificate_domain" {
  type = string
}
variable "cloudfront_alias" {
  type = string
}

variable "api_cloudfront_origin_domain" {
  type = string
}

variable "front_bucket_name" {
  type = string
}

variable "gateway_api_mapping_domain" {
  type = string
}

variable "gateway_api_name" {
  type = string
}

variable "dynamodb_table_name" {
  type = string
}

variable "artifacts_bucket_name" {
  type = string
}
