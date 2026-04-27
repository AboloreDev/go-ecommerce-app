#!/bin/bash
# Create SQS QUEUE
awslocal sqs create-queue --queue-name armory-events

# Create bucket
awslocal s3 mb s3://armory-db-bucket

echo "Localstack initialisation complete"