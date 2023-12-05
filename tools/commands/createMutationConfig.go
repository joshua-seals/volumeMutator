package commands

import (
	"bytes"
	"context"
	"os"

	v1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func CreateMutationConfig(ctx context.Context, caPEM *bytes.Buffer) {

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

	path := "/volume-mutator"
	fail := v1.Fail
	port := int32(8443)

	if err != nil {
		panic("failed to read certPath")
	}

	mutateconfig := &v1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: mutationCfgName,
		},
		Webhooks: []v1.MutatingWebhook{{
			Name: "volume-mutator.renci.org",
			ClientConfig: v1.WebhookClientConfig{
				CABundle: caPEM.Bytes(), // CA bundle created in generateTLSCerts command
				Service: &v1.ServiceReference{
					Name:      webhookService,
					Namespace: webhookNamespace,
					Path:      &path,
					Port:      &port,
				},
			},
			Rules: []v1.RuleWithOperations{
				{
					Operations: []v1.OperationType{
						v1.Create, v1.Update,
					},
					Rule: v1.Rule{
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

	if _, err := kubeClient.AdmissionregistrationV1().MutatingWebhookConfigurations().Create(
		ctx, mutateconfig, metav1.CreateOptions{},
	); err != nil {
		panic(err)
	}
}
