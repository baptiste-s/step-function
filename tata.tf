#!/bin/bash

# --- Configuration ---
STREAM_NAME="mon-stream-kinesis"
DATA="Hello Kinesis!"  # contenu du record
COUNT=50000            # nombre d'envois
BATCH_SIZE=500         # max 500 par appel put-records

# Encode la donnée une seule fois (base64 obligatoire)
ENCODED_DATA=$(echo -n "$DATA" | base64)

echo "Envoi de $COUNT records dans le stream '$STREAM_NAME'..."
for ((i=0; i<COUNT; i+=BATCH_SIZE)); do
  RECORDS="["

  for ((j=0; j<BATCH_SIZE && i+j<COUNT; j++)); do
    PARTITION_KEY="key-$((i+j))"
    RECORDS+="{\"Data\":\"$ENCODED_DATA\", \"PartitionKey\":\"$PARTITION_KEY\"},"
  done

  # Supprime la dernière virgule et ferme le tableau JSON
  RECORDS="${RECORDS%,}]"

  aws kinesis put-records \
    --stream-name "$STREAM_NAME" \
    --records "$RECORDS" \
    >/dev/null

  echo "$((i+BATCH_SIZE)) records envoyés..."
done

echo "✅ Envoi terminé."
