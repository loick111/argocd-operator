package argocd

import (
	"context"
	"testing"

	argoproj "github.com/argoproj-labs/argocd-operator/api/v1beta1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	testclient "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func TestReconcileApplicationControllerStatefulSet_WithLabels(t *testing.T) {
	logf.SetLogger(ZapLogger(true))
	a := makeTestArgoCD(func(a *argoproj.ArgoCD) {
		a.Spec.Controller.Labels = map[string]string{
			"custom-label":  "custom-value",
			"another-label": "another-value",
		}
	})

	resObjs := []client.Object{a}
	subresObjs := []client.Object{a}
	runtimeObjs := []runtime.Object{}
	sch := makeTestReconcilerScheme(argoproj.AddToScheme)
	cl := makeTestReconcilerClient(sch, resObjs, subresObjs, runtimeObjs)
	r := makeTestReconciler(cl, sch, testclient.NewSimpleClientset())

	assert.NoError(t, r.reconcileApplicationControllerStatefulSet(a, false))

	statefulset := &appsv1.StatefulSet{}
	assert.NoError(t, r.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      a.Name + "-application-controller",
			Namespace: a.Namespace,
		},
		statefulset))

	// Verify custom labels are applied to pod template
	assert.Equal(t, "custom-value", statefulset.Spec.Template.Labels["custom-label"])
	assert.Equal(t, "another-value", statefulset.Spec.Template.Labels["another-label"])
}

func TestReconcileApplicationControllerStatefulSet_WithAnnotations(t *testing.T) {
	logf.SetLogger(ZapLogger(true))
	a := makeTestArgoCD(func(a *argoproj.ArgoCD) {
		a.Spec.Controller.Annotations = map[string]string{
			"custom-annotation":  "custom-value",
			"another-annotation": "another-value",
		}
	})

	resObjs := []client.Object{a}
	subresObjs := []client.Object{a}
	runtimeObjs := []runtime.Object{}
	sch := makeTestReconcilerScheme(argoproj.AddToScheme)
	cl := makeTestReconcilerClient(sch, resObjs, subresObjs, runtimeObjs)
	r := makeTestReconciler(cl, sch, testclient.NewSimpleClientset())

	assert.NoError(t, r.reconcileApplicationControllerStatefulSet(a, false))

	statefulset := &appsv1.StatefulSet{}
	assert.NoError(t, r.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      a.Name + "-application-controller",
			Namespace: a.Namespace,
		},
		statefulset))

	// Verify custom annotations are applied to pod template
	assert.Equal(t, "custom-value", statefulset.Spec.Template.Annotations["custom-annotation"])
	assert.Equal(t, "another-value", statefulset.Spec.Template.Annotations["another-annotation"])
}

func TestReconcileApplicationControllerStatefulSet_WithLabelsAndAnnotations(t *testing.T) {
	logf.SetLogger(ZapLogger(true))
	a := makeTestArgoCD(func(a *argoproj.ArgoCD) {
		a.Spec.Controller.Labels = map[string]string{
			"custom-label": "custom-value",
		}
		a.Spec.Controller.Annotations = map[string]string{
			"custom-annotation": "custom-value",
		}
	})

	resObjs := []client.Object{a}
	subresObjs := []client.Object{a}
	runtimeObjs := []runtime.Object{}
	sch := makeTestReconcilerScheme(argoproj.AddToScheme)
	cl := makeTestReconcilerClient(sch, resObjs, subresObjs, runtimeObjs)
	r := makeTestReconciler(cl, sch, testclient.NewSimpleClientset())

	assert.NoError(t, r.reconcileApplicationControllerStatefulSet(a, false))

	statefulset := &appsv1.StatefulSet{}
	assert.NoError(t, r.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      a.Name + "-application-controller",
			Namespace: a.Namespace,
		},
		statefulset))

	// Verify both labels and annotations are applied to pod template
	assert.Equal(t, "custom-value", statefulset.Spec.Template.Labels["custom-label"])
	assert.Equal(t, "custom-value", statefulset.Spec.Template.Annotations["custom-annotation"])
}
