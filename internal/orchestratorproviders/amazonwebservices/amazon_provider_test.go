package amazonwebservices_test

import (
	"errors"
	"github.com/hexa-org/policy-orchestrator/internal/orchestrator"
	"github.com/hexa-org/policy-orchestrator/internal/orchestratorproviders/amazonwebservices"
	"github.com/hexa-org/policy-orchestrator/internal/orchestratorproviders/amazonwebservices/test"
	"github.com/hexa-org/policy-orchestrator/internal/policysupport"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"
)

func TestAmazonProvider_Credentials(t *testing.T) {
	key := []byte(`
{
  "accessKeyID": "anAccessKeyID",
  "secretAccessKey": "aSecretAccessKey",
  "region": "aRegion"
}
`)
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: &cognitoidentityprovider.Client{}}
	c := p.Credentials(key)
	assert.Equal(t, "anAccessKeyID", c.AccessKeyID)
}

func TestAmazonProvider_DiscoverApplications(t *testing.T) {
	key := []byte(`
{
  "accessKeyID": "anAccessKeyID",
  "secretAccessKey": "aSecretAccessKey",
  "region": "aRegion"
}
`)
	info := orchestrator.IntegrationInfo{Name: "amazon", Key: key}
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: &cognitoidentityprovider.Client{}}
	_, err := p.DiscoverApplications(info)
	assert.Error(t, err, "operation error Cognito Identity Provider: ListUserPools, expected endpoint resolver to not be nil")
}

func TestAmazonProvider_DiscoverApplications_withOtherProvider(t *testing.T) {
	p := &amazonwebservices.AmazonProvider{}
	info := orchestrator.IntegrationInfo{Name: "not_amazon", Key: []byte("aKey")}
	_, err := p.DiscoverApplications(info)
	assert.NoError(t, err)
	assert.Nil(t, p.CognitoClientOverride)
}

func TestAmazonProvider_ListUserPools(t *testing.T) {
	key := []byte(`
{
  "accessKeyID": "anAccessKeyID",
  "secretAccessKey": "aSecretAccessKey",
  "region": "aRegion"
}
`)
	mockClient := &amazonwebservices_test.MockClient{}
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}
	info := orchestrator.IntegrationInfo{Name: "amazon", Key: key}
	pools, err := p.ListUserPools(info)
	assert.NoError(t, err)
	assert.Len(t, pools, 1)
	assert.Equal(t, "anId", pools[0].ObjectID)
	assert.Equal(t, "aName", pools[0].Name)
	assert.Equal(t, "Cognito", pools[0].Service)
}

func TestAmazonProvider_ListUserPools_withError(t *testing.T) {
	key := []byte(`
{
  "accessKeyID": "anAccessKeyID",
  "secretAccessKey": "aSecretAccessKey",
  "region": "aRegion"
}
`)
	mockClient := &amazonwebservices_test.MockClient{Errs: map[string]error{}}
	mockClient.Errs["ListUserPools"] = errors.New("oops")
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}
	info := orchestrator.IntegrationInfo{Name: "amazon", Key: key}
	_, err := p.ListUserPools(info)
	assert.Error(t, err)
}

func TestAmazonProvider_GetPolicyInfo(t *testing.T) {
	mockClient := &amazonwebservices_test.MockClient{}
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}
	info, _ := p.GetPolicyInfo(orchestrator.IntegrationInfo{}, orchestrator.ApplicationInfo{ObjectID: "anObjectId"})
	assert.Equal(t, 1, len(info))
	assert.Equal(t, "aws:amazon.cognito/access", info[0].Actions[0].ActionUri)
	assert.Equal(t, "aUser:aUser@amazon.com", info[0].Subject.Members[0])
	assert.Equal(t, "anObjectId", info[0].Object.ResourceID)
}

func TestAmazonProvider_GetPolicyInfo_withError(t *testing.T) {
	mockClient := &amazonwebservices_test.MockClient{Errs: map[string]error{}}
	mockClient.Errs["ListUsers"] = errors.New("oops")
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}
	_, err := p.GetPolicyInfo(orchestrator.IntegrationInfo{}, orchestrator.ApplicationInfo{ObjectID: "anObjectId"})
	assert.Error(t, err)
}

func TestAmazonProvider_ShouldEnable(t *testing.T) {
	mockClient := &amazonwebservices_test.MockClient{}
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}
	shouldAdd := p.ShouldEnable([]string{"aUser@amazon.com", "yetAnotherUser@amazon.com"}, []string{"anotherUser@amazon.com"})
	assert.Equal(t, []string{"anotherUser@amazon.com"}, shouldAdd)
}

func TestAmazonProvider_ShouldDisable(t *testing.T) {
	mockClient := &amazonwebservices_test.MockClient{}
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}
	shouldAdd := p.ShouldDisable([]string{"aUser@amazon.com", "yetAnotherUser@amazon.com"}, []string{"yetAnotherUser@amazon.com"})
	assert.Equal(t, []string{"aUser@amazon.com"}, shouldAdd)
}

func TestAmazonProvider_SetPolicyInfo(t *testing.T) {
	mockClient := &amazonwebservices_test.MockClient{}
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}
	status, err := p.SetPolicyInfo(orchestrator.IntegrationInfo{}, orchestrator.ApplicationInfo{ObjectID: "anObjectId"}, []policysupport.PolicyInfo{{
		Meta:    policysupport.MetaInfo{Version: "0"},
		Actions: []policysupport.ActionInfo{{"aws:amazon.cognito/access"}},
		Subject: policysupport.SubjectInfo{Members: []string{"aUser:aUser@amazon.com", "anotherUser:anotherUser@amazon.com"}},
		Object:  policysupport.ObjectInfo{ResourceID: "aResourceId"},
	}})
	assert.Equal(t, http.StatusCreated, status)
	assert.NoError(t, err)
}

func TestAmazonProvider_SetPolicyInfo_withInvalidArguments(t *testing.T) {
	mockClient := &amazonwebservices_test.MockClient{}
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}

	status, _ := p.SetPolicyInfo(orchestrator.IntegrationInfo{}, orchestrator.ApplicationInfo{}, []policysupport.PolicyInfo{})
	assert.Equal(t, http.StatusInternalServerError, status)

	status, _ = p.SetPolicyInfo(orchestrator.IntegrationInfo{}, orchestrator.ApplicationInfo{ObjectID: "anObjectId"}, []policysupport.PolicyInfo{{
		Actions: []policysupport.ActionInfo{},
		Subject: policysupport.SubjectInfo{Members: []string{}},
		Object:  policysupport.ObjectInfo{ResourceID: "aResourceId"},
	}})
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestAmazonProvider_SetPolicyInfo_withListErr(t *testing.T) {
	mockClient := &amazonwebservices_test.MockClient{Errs: map[string]error{}}
	mockClient.Errs["ListUsers"] = errors.New("oops")
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}
	status, err := p.SetPolicyInfo(orchestrator.IntegrationInfo{}, orchestrator.ApplicationInfo{ObjectID: "anObjectId"}, []policysupport.PolicyInfo{{
		Meta:    policysupport.MetaInfo{Version: "0"},
		Actions: []policysupport.ActionInfo{},
		Subject: policysupport.SubjectInfo{Members: []string{}},
		Object:  policysupport.ObjectInfo{ResourceID: "aResourceId"},
	}})
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Error(t, err)
}

func TestAmazonProvider_SetPolicyInfo_withEnableErr(t *testing.T) {
	mockClient := &amazonwebservices_test.MockClient{Errs: map[string]error{}}
	mockClient.Errs["AdminEnableUser"] = errors.New("oops")
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}
	status, err := p.SetPolicyInfo(orchestrator.IntegrationInfo{}, orchestrator.ApplicationInfo{ObjectID: "anObjectId"}, []policysupport.PolicyInfo{{
		Meta:    policysupport.MetaInfo{Version: "0"},
		Actions: []policysupport.ActionInfo{{"aws:amazon.cognito/access"}},
		Subject: policysupport.SubjectInfo{Members: []string{"aUser:aUser@amazon.com", "anotherUser:anotherUser@amazon.com"}},
		Object:  policysupport.ObjectInfo{ResourceID: "aResourceId"},
	}})
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Error(t, err)
}

func TestAmazonProvider_SetPolicyInfo_withDisableErr(t *testing.T) {
	mockClient := &amazonwebservices_test.MockClient{Errs: map[string]error{}}
	mockClient.Errs["AdminDisableUser"] = errors.New("oops")
	p := &amazonwebservices.AmazonProvider{CognitoClientOverride: mockClient}
	status, err := p.SetPolicyInfo(orchestrator.IntegrationInfo{}, orchestrator.ApplicationInfo{ObjectID: "anObjectId"}, []policysupport.PolicyInfo{{
		Meta:    policysupport.MetaInfo{Version: "0"},
		Actions: []policysupport.ActionInfo{},
		Subject: policysupport.SubjectInfo{Members: []string{}},
		Object:  policysupport.ObjectInfo{ResourceID: "aResourceId"},
	}})
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Error(t, err)
}
