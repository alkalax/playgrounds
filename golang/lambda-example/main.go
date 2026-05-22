package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func main() {
	functionName := "example_lambda_function"
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config, %v", err)
	}

	client := lambda.NewFromConfig(cfg)

	input := &lambda.InvokeInput{
		FunctionName:   &functionName,
		InvocationType: types.InvocationTypeRequestResponse,
		Payload:        []byte("{}"),
	}

	log.Printf("Invoking lambda function: %s...", functionName)
	result, err := client.Invoke(ctx, input)
	if err != nil {
		log.Fatalf("failed to invoke lambda function, %v", err)
	}

	log.Printf("Response status code: %d", result.StatusCode)
	if result.FunctionError != nil {
		log.Printf("function error returned: %s\n", *result.FunctionError)
	}

	log.Printf("Response payload: %s", string(result.Payload))
}
