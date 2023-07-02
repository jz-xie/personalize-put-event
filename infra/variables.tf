variable "aws_region" {
  description = "AWS Region to be used that was defined in AWS CLI configuration"
  type        = string
  default     = "us-east-1"
}

variable "project" {
  description = "Name of your project/product/application"
  type        = string
  default     = "recommendation-engine"
}

variable "package_name" {
  description = "Name of your the go package that contains the lambda function"
  type        = string
  default     = "main"
}
