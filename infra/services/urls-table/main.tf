variable "stage" {
  type = string
}

###### DynamoDB ######
resource "aws_dynamodb_table" "urls-table" {
  name         = "shortener-urls"
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
}

output "arn" {
  value = aws_dynamodb_table.urls-table.arn
}

