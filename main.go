package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/styleila/analyser/app/analyser"
)

func main() {
	lambda.Start(analyser.HandleEvent)
}
