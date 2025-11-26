#!/bin/bash

# Configuration
# Note: Adjusted filename to match the project structure: test_input/transactions.json
INPUT_FILE="transactions.json"
TOPIC="TransactionTopic"
BOOTSTRAP_SERVER="kafka.curryware.org:9092"

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: File $INPUT_FILE not found."
    exit 1
fi

# 1. cat reads the file
# 2. base64 encodes it
# 3. tr -d '\n' removes any newlines added by the base64 tool (ensuring it is treated as one message)
# 4. kafka-console-producer sends it to the topic
cat "$INPUT_FILE" | base64 | tr -d '\n' | /opt/homebrew/opt/kafka/bin/kafka-console-producer --bootstrap-server "$BOOTSTRAP_SERVER" --topic "$TOPIC"

echo ""
echo "Encoded content of $INPUT_FILE sent to topic $TOPIC"