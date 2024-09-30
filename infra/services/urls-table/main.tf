variable "stage" {
  type = string
}

variable "dynamodb_table_name" {
  type = string
}

###### DynamoDB ######
resource "aws_dynamodb_table" "urls-table" {
  name         = var.dynamodb_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Rash"

  attribute {
    name = "Rash"
    type = "S"
  }

  ttl {
    attribute_name = "Ttl"
    enabled        = true
  }

  tags = {
    Terraform = "true"
    Stage     = var.stage
  }

  lifecycle {
    prevent_destroy = true
  }
}

output "arn" {
  value = aws_dynamodb_table.urls-table.arn
}

