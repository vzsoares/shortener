resource "aws_iam_role" "iam" {
  name = "shortener_lambda_executor_${var.stage}"
  assume_role_policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : "sts:AssumeRole",
          "Principal" : {
            "Service" : "lambda.amazonaws.com"
          },
          "Effect" : "Allow",
        }
      ]
    }
  )
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.iam.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_basic_sqs_execution" {
  role       = aws_iam_role.iam.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaSQSQueueExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_smm_readonly" {
  role       = aws_iam_role.iam.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess"
}

# dynamo policy

data "aws_iam_policy_document" "url-table-full-acesss-policy-doc" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:*"]
    resources = ["arn:aws:dynamodb:*:*:table/${var.dynamodb_table_name}"]
  }
}

resource "aws_iam_policy" "url-table-full-acesss-policy" {
  name        = "ShortenerUrlsTableFullAccess-${var.stage}"
  description = "Dynamodb shortener urls table access"
  policy      = data.aws_iam_policy_document.url-table-full-acesss-policy-doc.json
}

resource "aws_iam_role_policy_attachment" "urls-table-only" {
  role       = aws_iam_role.iam.name
  policy_arn = aws_iam_policy.url-table-full-acesss-policy.arn
}
