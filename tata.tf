# Step Function avec Map pour paralléliser
resource "aws_sfn_state_machine" "postgres_test" {
  name     = "postgres-test-parallel"
  role_arn = aws_iam_role.step_functions_role.arn

  definition = jsonencode({
    Comment = "Test PostgreSQL avec Lambdas parallèles"
    StartAt = "PrepareWorkers"
    States = {
      PrepareWorkers = {
        Type = "Pass"
        Parameters = {
          "num_workers.$"        = "$.num_workers"
          "inserts_per_worker.$" = "$.inserts_per_worker"
          "table_name.$"         = "$.table_name"
          "workers.$"            = "States.ArrayRange(0, $.num_workers, 1)"
        }
        ResultPath = "$.config"
        Next       = "MapWorkers"
      }
      MapWorkers = {
        Type       = "Map"
        ItemsPath  = "$.config.workers"
        MaxConcurrency = 0
        ItemProcessor = {
          ProcessorConfig = {
            Mode = "INLINE"
          }
          StartAt = "InvokeWorker"
          States = {
            InvokeWorker = {
              Type     = "Task"
              Resource = "arn:aws:states:::lambda:invoke"
              Parameters = {
                "FunctionName" = aws_lambda_function.worker.arn
                "Payload" = {
                  "worker_id.$"    = "States.Format('worker-{}', $)"
                  "num_inserts.$"  = "$.inserts_per_worker"
                  "table_name.$"   = "$.table_name"
                }
              }
              Retry = [
                {
                  ErrorEquals     = ["Lambda.ServiceException", "Lambda.AWSLambdaException", "Lambda.SdkClientException"]
                  IntervalSeconds = 2
                  MaxAttempts     = 3
                  BackoffRate     = 2
                }
              ]
              End = true
            }
          }
        }
        ResultPath = "$.results"
        Next       = "AggregateResults"
      }
      AggregateResults = {
        Type = "Pass"
        Parameters = {
          "message"                  = "Test terminé"
          "num_workers.$"            = "$.config.num_workers"
          "inserts_per_worker.$"     = "$.config.inserts_per_worker"
          "total_expected_inserts.$" = "States.Format('{}', States.MathMultiply($.config.num_workers, $.config.inserts_per_worker))"
          "worker_results.$"         = "$.results[*].Payload.body"
          "execution_arn.$"          = "$.Execution.Id"
          "start_time.$"             = "$.Execution.StartTime"
        }
        End = true
      }
    }
  })
