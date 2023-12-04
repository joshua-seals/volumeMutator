package main

import (
	"log"

	"github.com/helxplatform/volumeMutator/tools/commands"
)

// If run locally via go run main.go certs will be local
// If run via makefile worfkow
// path is CERT_PATH /etc/webhook/certs/
var certPath = "./certs/"

func main() {
	err := commands.GenerateTLSCerts(certPath)
	if err != nil {
		log.Panic(err)
	}
	// ctx := context.Background()
	// commands.CreateMutationConfig(ctx, certPath)
}

// Reference https://www.velotio.com/engineering-blog/managing-tls-certificate-for-kubernetes-admission-webhook
