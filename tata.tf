{
  "Comment": "Test PostgreSQL avec Lambdas parallèles",
  "StartAt": "GenerateWorkers",
  "States": {
    "GenerateWorkers": {
      "Type": "Pass",
      "Parameters": {
        "workers.$": "States.ArrayRange(1, $.num_workers, 1)"
      },
      "Next": "MapWorkers"
    },
    "MapWorkers": {
      "Type": "Map",
      "MaxConcurrency": 0,
      "ItemsPath": "$.workers",
      "Iterator": {
        "StartAt": "InvokeWorker",
        "States": {
          "InvokeWorker": {
            "Type": "Task",
            "Resource": "arn:aws:lambda:eu-west-1:VOTRE-ACCOUNT-ID:function:postgres-test-worker",
            "Parameters": {
              "num_inserts": 2000,
              "table_name": "test_data",
              "insert_mode": "batch"
            },
            "End": true
          }
        }
      },
      "End": true
    }
  }
}
