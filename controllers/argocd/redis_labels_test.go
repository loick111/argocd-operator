package argocd

import (
	"context"
	"testing"

	argoproj "github.com/argoproj-labs/argocd-operator/api/v1beta1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	testclient "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func TestReconcileRedisDeployment_WithLabels(t *testing.T) {
	logf.SetLogger(ZapLogger(true))
	a := makeTestArgoCD(func(a *argoproj.ArgoCD) {
		a.Spec.Redis.Labels = map[string]string{
			"custom-label":  "custom-value",
			"another-label": "another-value",
		}
	})

	resObjs := []client.Object{a}
	subresObjs := []client.Object{a}
	sch := makeTestReconcilerScheme(argoproj.AddToScheme)
	cl := makeTestReconcilerClient(sch, resObjs, subresObjs, nil)
	r := makeTestReconciler(cl, sch, testclient.NewSimpleClientset())

	assert.NoError(t, r.reconcileRedisDeployment(a, false))

	deployment := &appsv1.Deployment{}
	assert.NoError(t, r.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      a.Name + "-redis",
			Namespace: a.Namespace,
		},
		deployment))

	// Verify custom labels are applied to pod template
	assert.Equal(t, "custom-value", deployment.Spec.Template.Labels["custom-label"])
	assert.Equal(t, "another-value", deployment.Spec.Template.Labels["another-label"])
}

func TestReconcileRedisDeployment_WithAnnotations(t *testing.T) {
	logf.SetLogger(ZapLogger(true))
	a := makeTestArgoCD(func(a *argoproj.ArgoCD) {
		a.Spec.Redis.Annotations = map[string]string{
			"custom-annotation":  "custom-value",
			"another-annotation": "another-value",
		}
	})

	resObjs := []client.Object{a}
	subresObjs := []client.Object{a}
	sch := makeTestReconcilerScheme(argoproj.AddToScheme)
	cl := makeTestReconcilerClient(sch, resObjs, subresObjs, nil)
	r := makeTestReconciler(cl, sch, testclient.NewSimpleClientset())

	assert.NoError(t, r.reconcileRedisDeployment(a, false))

	deployment := &appsv1.Deployment{}
	assert.NoError(t, r.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      a.Name + "-redis",
			Namespace: a.Namespace,
		},
		deployment))

	// Verify custom annotations are applied to pod template
	assert.Equal(t, "custom-value", deployment.Spec.Template.Annotations["custom-annotation"])
	assert.Equal(t, "another-value", deployment.Spec.Template.Annotations["another-annotation"])
}

func TestReconcileRedisDeployment_WithLabelsAndAnnotations(t *testing.T) {
	logf.SetLogger(ZapLogger(true))
	a := makeTestArgoCD(func(a *argoproj.ArgoCD) {
		a.Spec.Redis.Labels = map[string]string{
			"custom-label": "custom-value",
		}
		a.Spec.Redis.Annotations = map[string]string{
			"custom-annotation": "custom-value",
		}
	})

	resObjs := []client.Object{a}
	subresObjs := []client.Object{a}
	sch := makeTestReconcilerScheme(argoproj.AddToScheme)
	cl := makeTestReconcilerClient(sch, resObjs, subresObjs, nil)
	r := makeTestReconciler(cl, sch, testclient.NewSimpleClientset())

	assert.NoError(t, r.reconcileRedisDeployment(a, false))

	deployment := &appsv1.Deployment{}
	assert.NoError(t, r.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      a.Name + "-redis",
			Namespace: a.Namespace,
		},
		deployment))

	// Verify both labels and annotations are applied to pod template
	assert.Equal(t, "custom-value", deployment.Spec.Template.Labels["custom-label"])
	assert.Equal(t, "custom-value", deployment.Spec.Template.Annotations["custom-annotation"])
}
