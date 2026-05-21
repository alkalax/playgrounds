terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>6.46.0"
    }
  }
}

provider "aws" {}

data "archive_file" "example" {
  type        = "zip"
  output_path = "${path.module}/lambda_function.zip"

  source {
    content = <<EOF
def lambda_handler(event, context):
  return {
    'statusCode': 200,
    'body': 'Hello world!'
  }
EOF

    filename = "lambda_function.py"
  }
}

resource "aws_lambda_function" "example" {
  filename      = data.archive_file.example.output_path
  function_name = "example_lambda_function"
  handler       = "lambda_function.lambda_handler"
  runtime       = "python3.12"
  role          = aws_iam_role.example.arn

  source_code_hash = data.archive_file.example.output_base64sha256
}

resource "aws_iam_role" "example" {
  name               = "lambda_execution_role"
  assume_role_policy = data.aws_iam_policy_document.example.json
}

resource "aws_iam_role_policy_attachment" "example" {
  role       = aws_iam_role.example.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "example" {
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}