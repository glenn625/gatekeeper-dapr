package audit

import (
	"context"

	"github.com/open-policy-agent/gatekeeper/pkg/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	pubsubName = "redis-pubsub"
	topicName  = "constraint-violations"
)

// The Dapr message contents.
type DetailedStatusViolation struct {
	Message              string
	Details              interface{}
	ConstraintGroup      string
	ConstraintAPIVersion string
	ConstraintKind       string
	ConstraintName       string
	ConstraintNamespace  string
	ConstraintAction     string
	ResourceGroup        string
	ResourceAPIVersion   string
	ResourceKind         string
	ResourceNamespace    string
	ResourceName         string
}

// Use Dapr to publish constraint violation to Redis broker.
func (am *Manager) publishViolationToDapr(
	ctx context.Context,
	constraint *unstructured.Unstructured,
	enforcementAction util.EnforcementAction,
	resourceGroupVersionKind schema.GroupVersionKind,
	rnamespace, rname, message string,
	details interface{}) {

	client, err := dapr.NewClient()
	if err != nil {
		am.log.Error(err, "Error creating Dapr client.")
	}

	defer client.Close()

	if err := client.PublishEvent(ctx, pubsubName, topicName, DetailedStatusViolation{
		Message:              message,
		Details:              details,
		ConstraintGroup:      constraint.GroupVersionKind().Group,
		ConstraintAPIVersion: constraint.GroupVersionKind().Version,
		ConstraintKind:       constraint.GetKind(),
		ConstraintName:       constraint.GetName(),
		ConstraintNamespace:  constraint.GetNamespace(),
		ConstraintAction:     string(enforcementAction),
		ResourceGroup:        resourceGroupVersionKind.Group,
		ResourceAPIVersion:   resourceGroupVersionKind.Version,
		ResourceKind:         resourceGroupVersionKind.Kind,
		ResourceNamespace:    rnamespace,
		ResourceName:         rname,
	}); err != nil {
		am.log.Error(err, "Could not publish violation to message broker.")
	}
}
