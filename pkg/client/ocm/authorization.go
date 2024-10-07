package ocm

import (
	"context"
	"fmt"

	azv1 "github.com/openshift-online/ocm-sdk-go/authorizations/v1"
)

//go:generate mockgen -source=authorization.go -package=ocm -destination=mock_authorization.go
type OCMAuthorization interface {
	SelfAccessReview(ctx context.Context, action, resourceType, organizationID, subscriptionID, clusterID string) (allowed bool, err error)
	AccessReview(ctx context.Context, username, action, resourceType, organizationID, subscriptionID, clusterID string) (allowed bool, err error)
	ResourceReview(ctx context.Context, username string, action string, resource string) (*azv1.ResourceReview, error)
}

type authorization service

var _ OCMAuthorization = &authorization{}

func (a authorization) SelfAccessReview(ctx context.Context, action, resourceType, organizationID, subscriptionID, clusterID string) (allowed bool, err error) {
	con := a.client.connection
	selfAccessReview := con.Authorizations().V1().SelfAccessReview()

	request, err := azv1.NewSelfAccessReviewRequest().
		Action(action).
		ResourceType(resourceType).
		OrganizationID(organizationID).
		ClusterID(clusterID).
		SubscriptionID(subscriptionID).
		Build()
	if err != nil {
		return false, err
	}

	postResp, err := selfAccessReview.Post().
		Request(request).
		SendContext(ctx)
	if err != nil {
		return false, err
	}
	response, ok := postResp.GetResponse()
	if !ok {
		return false, fmt.Errorf("Empty response from authorization post request")
	}

	return response.Allowed(), nil
}

func (a authorization) AccessReview(ctx context.Context, username, action, resourceType, organizationID, subscriptionID, clusterID string) (allowed bool, err error) {
	con := a.client.connection
	accessReview := con.Authorizations().V1().AccessReview()

	request, err := azv1.NewAccessReviewRequest().
		AccountUsername(username).
		Action(action).
		ResourceType(resourceType).
		OrganizationID(organizationID).
		ClusterID(clusterID).
		SubscriptionID(subscriptionID).
		Build()
	if err != nil {
		return false, err
	}

	postResp, err := accessReview.Post().
		Request(request).
		SendContext(ctx)
	if err != nil {
		return false, err
	}

	response, ok := postResp.GetResponse()
	if !ok {
		return false, fmt.Errorf("Empty response from authorization post request")
	}

	return response.Allowed(), nil
}

func (a authorization) ResourceReview(ctx context.Context, username string, action string, resource string) (*azv1.ResourceReview, error) {
	con := a.client.connection
	resourceReviewClient := con.Authorizations().V1().ResourceReview()

	request, err := azv1.NewResourceReviewRequest().AccountUsername(username).Action(action).ResourceType(resource).Build()
	if err != nil {
		return nil, err
	}
	response, err := resourceReviewClient.Post().Request(request).SendContext(ctx)
	if err != nil {
		return nil, err
	}
	return response.Review(), nil
}
