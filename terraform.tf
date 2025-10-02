aws_region = "eu-west-1"db_host = "votre-rds.rds.amazonaws.com"db_name = "postgres"db_user = "postgres"db_p = ""vpc_id = "vpc-xxxxx"subnet_ids = ["subnet-xxxxx", "subnet-yyyyy"]rds_security_group_id = "sg-xxxxx"
terraform {
  required_version = ">= 1.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = "~> 2.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# Variables
variable "aws_region" {
  description = "Région AWS"
  type        = string
  default     = "eu-west-1"
}

variable "db_host" {
  description = "Endpoint de la base RDS PostgreSQL"
  type        = string
}

variable "db_name" {
  description = "Nom de la base de données"
  type        = string
  default     = "postgres"
}

variable "db_user" {
  description = "Utilisateur de la base de données"
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "Mot de passe de la base de données"
  type        = string
  sensitive   = true
}

variable "db_port" {
  description = "Port PostgreSQL"
  type        = string
  default     = "5432"
}

variable "vpc_id" {
  description = "ID du VPC où se trouve RDS"
  type        = string
}

variable "subnet_ids" {
  description = "IDs des subnets privés pour les Lambdas"
  type        = list(string)
}

variable "rds_security_group_id" {
  description = "ID du security group de RDS"
  type        = string
}

# Security Group pour les Lambdas
resource "aws_security_group" "lambda_sg" {
  name        = "postgres-test-lambda-sg"
  description = "Security group pour les Lambdas de test PostgreSQL"
  vpc_id      = var.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "postgres-test-lambda-sg"
  }
}

# Règle pour autoriser les Lambdas à accéder à RDS
resource "aws_security_group_rule" "lambda_to_rds" {
  type                     = "ingress"
  from_port                = var.db_port
  to_port                  = var.db_port
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.lambda_sg.id
  security_group_id        = var.rds_security_group_id
  description              = "Allow Lambda to access RDS PostgreSQL"
}

# IAM Role pour la Lambda Worker
resource "aws_iam_role" "lambda_worker_role" {
  name = "postgres-test-worker-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

# Policy pour le Worker
resource "aws_iam_role_policy" "lambda_worker_policy" {
  name = "postgres-test-worker-policy"
  role = aws_iam_role.lambda_worker_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:*:*:*"
      },
      {
        Effect = "Allow"
        Action = [
          "ec2:CreateNetworkInterface",
          "ec2:DescribeNetworkInterfaces",
          "ec2:DeleteNetworkInterface",
          "ec2:AssignPrivateIpAddresses",
          "ec2:UnassignPrivateIpAddresses"
        ]
        Resource = "*"
      }
    ]
  })
}

# Lambda Layer pour psycopg2
resource "aws_lambda_layer_version" "psycopg2" {
  filename            = "psycopg2-layer.zip"
  layer_name          = "psycopg2-layer"
  compatible_runtimes = ["python3.11", "python3.12"]
  source_code_hash    = filebase64sha256("psycopg2-layer.zip")

  description = "psycopg2-binary pour connexion PostgreSQL"
}

# Archive du code Worker
data "archive_file" "lambda_worker" {
  type        = "zip"
  output_path = "${path.module}/lambda_worker.zip"
  
  source {
    content  = file("${path.module}/lambda_worker.py")
    filename = "lambda_function.py"
  }
}

# Lambda Worker
resource "aws_lambda_function" "worker" {
  filename         = data.archive_file.lambda_worker.output_path
  function_name    = "postgres-test-worker"
  role            = aws_iam_role.lambda_worker_role.arn
  handler         = "lambda_function.lambda_handler"
  source_code_hash = data.archive_file.lambda_worker.output_base64sha256
  runtime         = "python3.11"
  timeout         = 900  # 15 minutes
  memory_size     = 512

  layers = [aws_lambda_layer_version.psycopg2.arn]

  vpc_config {
    subnet_ids         = var.subnet_ids
    security_group_ids = [aws_security_group.lambda_sg.id]
  }

  environment {
    variables = {
      DB_HOST     = var.db_host
      DB_NAME     = var.db_name
      DB_USER     = var.db_user
      DB_PASSWORD = var.db_password
      DB_PORT     = var.db_port
    }
  }

  tags = {
    Name = "postgres-test-worker"
  }
}

# IAM Role pour Step Functions
resource "aws_iam_role" "step_functions_role" {
  name = "postgres-test-step-functions-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "states.amazonaws.com"
      }
    }]
  })
}

# Policy pour Step Functions
resource "aws_iam_role_policy" "step_functions_policy" {
  name = "postgres-test-step-functions-policy"
  role = aws_iam_role.step_functions_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "lambda:InvokeFunction"
        ]
        Resource = aws_lambda_function.worker.arn
      },
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogDelivery",
          "logs:GetLogDelivery",
          "logs:UpdateLogDelivery",
          "logs:DeleteLogDelivery",
          "logs:ListLogDeliveries",
          "logs:PutResourcePolicy",
          "logs:DescribeResourcePolicies",
          "logs:DescribeLogGroups"
        ]
        Resource = "*"
      }
    ]
  })
}

# CloudWatch Log Group pour Step Functions
resource "aws_cloudwatch_log_group" "step_functions_logs" {
  name              = "/aws/stepfunctions/postgres-test"
  retention_in_days = 7

  tags = {
    Name = "postgres-test-step-functions-logs"
  }
}

# Step Function avec Map pour paralléliser
resource "aws_sfn_state_machine" "postgres_test" {
  name     = "postgres-test-parallel"
  role_arn = aws_iam_role.step_functions_role.arn

  definition = jsonencode({
    Comment = "Test PostgreSQL avec Lambdas parallèles"
    StartAt = "PrepareWorkers"
    States = {
      PrepareWorkers = {
        Type = "Pass"
        Parameters = {
          "num_workers.$"        = "$.num_workers"
          "inserts_per_worker.$" = "$.inserts_per_worker"
          "table_name.$"         = "$.table_name"
          "workers.$"            = "States.Array(States.ArrayRange(0, $.num_workers, 1))"
        }
        ResultPath = "$.config"
        Next       = "MapWorkers"
      }
      MapWorkers = {
        Type       = "Map"
        ItemsPath  = "$.config.workers"
        MaxConcurrency = 0  # Illimité, parallélisation maximale
        ItemProcessor = {
          ProcessorConfig = {
            Mode = "INLINE"
          }
          StartAt = "InvokeWorker"
          States = {
            InvokeWorker = {
              Type     = "Task"
              Resource = "arn:aws:states:::lambda:invoke"
              Parameters = {
                "FunctionName" = aws_lambda_function.worker.arn
                "Payload" = {
                  "worker_id.$"         = "States.Format('worker-{}', $)"
                  "num_inserts.$"       = "$.config.inserts_per_worker"
                  "table_name.$"        = "$.config.table_name"
                }
              }
              Retry = [
                {
                  ErrorEquals     = ["Lambda.ServiceException", "Lambda.AWSLambdaException", "Lambda.SdkClientException"]
                  IntervalSeconds = 2
                  MaxAttempts     = 3
                  BackoffRate     = 2
                }
              ]
              End = true
            }
          }
        }
        ResultPath = "$.results"
        Next       = "AggregateResults"
      }
      AggregateResults = {
        Type = "Pass"
        Parameters = {
          "message"                   = "Test terminé"
          "num_workers.$"             = "$.config.num_workers"
          "inserts_per_worker.$"      = "$.config.inserts_per_worker"
          "total_expected_inserts.$"  = "States.MathAdd(States.MathMultiply($.config.num_workers, $.config.inserts_per_worker), 0)"
          "worker_results.$"          = "$.results[*].Payload.body"
          "execution_arn.$"           = "$$.Execution.Id"
          "start_time.$"              = "$$.Execution.StartTime"
        }
        End = true
      }
    }
  })

  logging_configuration {
    log_destination        = "${aws_cloudwatch_log_group.step_functions_logs.arn}:*"
    include_execution_data = true
    level                  = "ALL"
  }

  tags = {
    Name = "postgres-test-parallel"
  }
}

# Outputs
output "worker_function_name" {
  description = "Nom de la fonction Lambda Worker"
  value       = aws_lambda_function.worker.function_name
}

output "step_function_arn" {
  description = "ARN de la Step Function"
  value       = aws_sfn_state_machine.postgres_test.arn
}

output "step_function_name" {
  description = "Nom de la Step Function"
  value       = aws_sfn_state_machine.postgres_test.name
}

output "security_group_id" {
  description = "ID du security group des Lambdas"
  value       = aws_security_group.lambda_sg.id
}

output "execution_command" {
  description = "Commande pour lancer un test"
  value       = "aws stepfunctions start-execution --state-machine-arn ${aws_sfn_state_machine.postgres_test.arn} --input '{\"num_workers\": 50, \"inserts_per_worker\": 2000, \"table_name\": \"test_data\"}'"
}

-----------

