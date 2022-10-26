package audit

import (
	"context"
	"testing"

	"github.com/open-policy-agent/gatekeeper/pkg/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	k8schema "k8s.io/apimachinery/pkg/runtime/schema"
)

func TestPublishViolationToDapr(t *testing.T) {
	// Build test input.
	tc := struct {
		Context           context.Context
		Constraint        *unstructured.Unstructured
		ResourceGroupGVK  k8schema.GroupVersionKind
		ResourceGroupName string
		ResourceGroupNS   string
		Message           string
		Action            util.EnforcementAction
		Details           interface{}
	}{
		Context:    context.Background(),
		Constraint: newConstraint("test_ckind", "test_cname", "dryrun", t),
		ResourceGroupGVK: k8schema.GroupVersionKind{
			Group:   "test_rg",
			Version: "test_rg_v",
			Kind:    "test_rg_kind",
		},
		Action:  "dryrun",
		Details: "",
	}

	// Test publishViolationToDapr.
	// TODO: Create mock manager.
	var testManager Manager

	err := (&testManager).publishViolationToDapr(
		tc.Context, tc.Constraint, tc.Action,
		tc.ResourceGroupGVK, tc.ResourceGroupName, tc.ResourceGroupNS, tc.Message, tc.Details)

	if err != nil {
		t.Errorf("err = %s; want nil", err)
	}
}

func newConstraint(kind, name string, enforcementAction string, t *testing.T) *unstructured.Unstructured {
	c := &unstructured.Unstructured{}
	c.SetGroupVersionKind(k8schema.GroupVersionKind{
		Group:   "constraints.gatekeeper.sh",
		Version: "v1alpha1",
		Kind:    kind,
	})
	c.SetName(name)
	if err := unstructured.SetNestedField(c.Object, enforcementAction, "spec", "enforcementAction"); err != nil {
		t.Errorf("unable to set enforcementAction for constraint resources: %s", err)
	}
	return c
}
