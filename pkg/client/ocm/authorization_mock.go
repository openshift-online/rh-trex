package ocm

import (
	"context"

	azv1 "github.com/openshift-online/ocm-sdk-go/authorizations/v1"
)

// authorizationMock returns allowed=true for every request
type authorizationMock service

var _ OCMAuthorization = &authorizationMock{}

func (a authorizationMock) SelfAccessReview(ctx context.Context, action, resourceType, organizationID, subscriptionID, clusterID string) (allowed bool, err error) {
	return true, nil
}

func (a authorizationMock) AccessReview(ctx context.Context, username, action, resourceType, organizationID, subscriptionID, clusterID string) (allowed bool, err error) {
	return true, nil
}

func (a authorizationMock) ResourceReview(ctx context.Context, username string, action string, resource string) (*azv1.ResourceReview, error) {
	response, err := azv1.NewResourceReview().AccountUsername(username).Action(action).ResourceType(resource).OrganizationIDs("*").Build()
	if err != nil {
		return nil, err
	}
	return response, nil
}
