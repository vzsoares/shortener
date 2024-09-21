variable "gateway_id" {
  type = string
}

variable "filename" {
  type = string
}

variable "gateway_route_key" {
  type = string
}

variable "lambda_iam_arn" {
  type = string
}

variable "lambda_base_name" {
  type = string
}

variable "gateway_execution_arn" {
  type = string
}

variable "stage" {
  type = string
}

variable "filepath" {
  description = "path to ziped file"
  type        = string
}
variable "s3_prefix" {
  type = string
}
variable "s3_bucket" {
  type = string
}
