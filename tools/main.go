package main

import (
	"context"
	"flag"
	"log"

	"github.com/helxplatform/volumeMutator/tools/commands"
)

// If run locally via go run main.go certs will be ./webhook/certs/
// If run via makefile worfkow in Docker
// path is CERT_PATH /helx/etc/webhook/certs/
var certPath = "./helx/webhook/certs/"

func main() {
	// Flag to enable Creation of Webhook Mutator Config
	mutationConfig := flag.Bool("M", false, "Create Webhook Mutator Configuration.")
	flag.Parse()

	caPEM, err := commands.GenerateTLSCerts(certPath)
	if err != nil {
		log.Panic(err)
	}
	if *mutationConfig {
		ctx := context.Background()
		// Use CABundle to Register new MutationWebhook
		commands.CreateMutationConfig(ctx, caPEM)
	}
}

// Reference https://www.velotio.com/engineering-blog/managing-tls-certificate-for-kubernetes-admission-webhook
