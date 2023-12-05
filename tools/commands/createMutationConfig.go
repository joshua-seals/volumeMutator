package commands

import (
	"context"
	"os"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func CreateMutationConfig(ctx context.Context, caCertPath string) {

	var (
		webhookNamespace = os.Getenv("WEBHOOK_NAMESPACE")
		mutationCfgName  = os.Getenv("MUTATE_CONFIG")
		webhookService   = os.Getenv("WEBHOOK_SERVICE")
	)
	config := ctrl.GetConfigOrDie()
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("failed to set go -client")
	}

	path := "/mutate"
	fail := admissionregistrationv1.Fail
	// Read in caCert created in GenerateTLSCerts
	caCert, err := os.ReadFile(caCertPath)

	mutateconfig := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: mutationCfgName,
		},
		Webhooks: []admissionregistrationv1.MutatingWebhook{{
			Name: "volume-mutator.renci.org",
			ClientConfig: admissionregistrationv1.WebhookClientConfig{
				CABundle: caCert, // CA bundle created in generateTLSCerts command
				Service: &admissionregistrationv1.ServiceReference{
					Name:      webhookService,
					Namespace: webhookNamespace,
					Path:      &path,
				},
			},
			Rules: []admissionregistrationv1.RuleWithOperations{
				{
					Operations: []admissionregistrationv1.OperationType{
						admissionregistrationv1.Create, admissionregistrationv1.Update,
					},
					Rule: admissionregistrationv1.Rule{
						APIGroups:   []string{"apps"},
						APIVersions: []string{"v1"},
						Resources:   []string{"deployments"},
					},
				}},
			AdmissionReviewVersions: []string{"v1"},
			FailurePolicy:           &fail,
			SideEffects: func() *v1.SideEffectClass {
				sideEffect := v1.SideEffectClassNone
				return &sideEffect
			}(),
		}},
	}

	if _, err := kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(ctx, mutateconfig, metav1.CreateOptions{}); err != nil {
		panic(err)
	}
}
