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

  #####################################
  ### other atributes for reference ###
  #  attribute {
  #    name = "Destination"
  #    type = "S"
  #  }
  #
  #  attribute {
  #    name = "CreatedAt"
  #    type = "N"
  #  }
  #
  #  attribute {
  #    name = "UpdatedAt"
  #    type = "N"
  #  }
  #
  #  attribute {
  #    name = "Version"
  #    type = "N"
  #  }
  #####################################
}

