provider "aws" {
  region = var.aws_region
}

locals {
  source_file   = "${path.module}/../src/${var.package_name}"
  zip_file_path = "${var.package_name}.zip"
}

resource "null_resource" "build" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command     = "GOOS=linux GOARCH=amd64 go build ${var.package_name}"
    working_dir = "${path.module}/../src"
  }
}

data "archive_file" "package" {
  type = "zip"

  source_file = local.source_file
  output_path = local.zip_file_path
  depends_on  = [null_resource.build]
}

resource "null_resource" "remove_exc" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command     = "rm ${var.package_name}"
    working_dir = "${path.module}/../src"
  }
}

data "aws_kinesis_stream" "dev_stream" {
  name = "ecommerce_unauth_dev"
}

resource "aws_lambda_function" "load_ecom_event" {
  function_name = "${var.project}-personalize-upload"

  runtime = "go1.x"
  handler = var.package_name

  filename         = local.zip_file_path
  source_code_hash = data.archive_file.package.output_base64sha256

  role = aws_iam_role.lambda_role.arn

  memory_size = 128
  timeout     = 10

}

resource "aws_lambda_event_source_mapping" "kinesis_mapping" {
  event_source_arn  = data.aws_kinesis_stream.dev_stream.arn
  function_name     = aws_lambda_function.load_ecom_event.function_name
  starting_position = "LATEST"
}


