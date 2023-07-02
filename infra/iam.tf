data "aws_caller_identity" "current" {}


data "aws_iam_policy_document" "kinesis_read" {
  statement {
    effect = "Allow"
    actions = [
      "kinesis:GetRecords",
      "kinesis:GetShardIterator",
      "kinesis:DescribeStream",
    ]
    resources = [
      "arn:aws:kinesis:${var.aws_region}:${data.aws_caller_identity.current.account_id}:stream/ecommerce_unauth_dev",
      "arn:aws:kinesis:${var.aws_region}:${data.aws_caller_identity.current.account_id}:stream/ecommerce_unauth_prod"
    ]
  }
}

data "aws_iam_policy_document" "personalize_put" {
  statement {
    effect = "Allow"
    actions = [
      "personalize:PutEvents"
    ]
    resources = [
      "*"
    ]
  }
}

data "aws_iam_policy_document" "lambda_policy" {
  source_policy_documents = [
    data.aws_iam_policy_document.kinesis_read.json,
    data.aws_iam_policy_document.personalize_put.json
  ]
}

resource "aws_iam_policy" "event_tracker_lambda_policy" {
  name        = "${var.project}-personalize-upload"
  description = "Policy for writing real-time events to Personalize event tracker"
  policy      = data.aws_iam_policy_document.lambda_policy.json
  tags = {
    project     = var.project
    environment = "production"
  }
}

resource "aws_iam_role" "lambda_role" {
  name = "${var.project}-personalize-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_policy_attachment" {
  policy_arn = aws_iam_policy.event_tracker_lambda_policy.arn
  role       = aws_iam_role.lambda_role.name
}

resource "aws_iam_role_policy_attachment" "lambda_log_policy_attachment" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_role.name
}
