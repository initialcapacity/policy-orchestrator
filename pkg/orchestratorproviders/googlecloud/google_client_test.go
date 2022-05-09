package googlecloud_test

import (
	"errors"
	"github.com/hexa-org/policy-orchestrator/pkg/orchestratorproviders/googlecloud"
	"github.com/hexa-org/policy-orchestrator/pkg/orchestratorproviders/googlecloud/test"
	"github.com/hexa-org/policy-orchestrator/pkg/policysupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoogleClient_GetAppEngineApplications(t *testing.T) {
	m := google_cloud_test.NewMockClient()
	m.ResponseBody["appengine"] = google_cloud_test.Resource("appengine.json")
	client := googlecloud.GoogleClient{HttpClient: m}
	applications, _ := client.GetAppEngineApplications()

	assert.Equal(t, 1, len(applications))
	assert.Equal(t, "hexa-demo", applications[0].ObjectID)
	assert.Equal(t, "apps/hexa-demo", applications[0].Name)
	assert.Equal(t, "hexa-demo.uc.r.appspot.com", applications[0].Description)
}

func TestClient_GetAppEngineApplications_withRequestError(t *testing.T) {
	m := google_cloud_test.NewMockClient()
	m.Err = errors.New("oops")

	client := googlecloud.GoogleClient{HttpClient: m}
	_, err := client.GetAppEngineApplications()
	assert.Error(t, err)
}

func TestClient_GetAppEngineApplications_withBadJson(t *testing.T) {
	m := google_cloud_test.NewMockClient()
	m.ResponseBody["compute"] = []byte("-")
	client := googlecloud.GoogleClient{HttpClient: m}
	_, err := client.GetAppEngineApplications()
	assert.Error(t, err)
}

func TestClient_GetBackendApplications(t *testing.T) {
	m := google_cloud_test.NewMockClient()
	m.ResponseBody["compute"] = google_cloud_test.Resource("backends.json")
	client := googlecloud.GoogleClient{HttpClient: m}
	applications, _ := client.GetBackendApplications()

	assert.Equal(t, 2, len(applications))
	assert.Equal(t, "k8s1-aName", applications[0].Name)
	assert.Equal(t, "k8s1-anotherName", applications[1].Name)
}

func TestClient_GetBackendApplications_withRequestError(t *testing.T) {
	m := google_cloud_test.NewMockClient()
	m.Err = errors.New("oops")

	client := googlecloud.GoogleClient{HttpClient: m}
	_, err := client.GetBackendApplications()
	assert.Error(t, err)
}

func TestClient_GetBackendApplications_withBadJson(t *testing.T) {
	m := google_cloud_test.NewMockClient()
	m.ResponseBody["compute"] = []byte("-")
	client := googlecloud.GoogleClient{HttpClient: m}
	_, err := client.GetBackendApplications()
	assert.Error(t, err)
}

func TestGoogleClient_GetBackendPolicies(t *testing.T) {
	m := google_cloud_test.NewMockClient()
	m.ResponseBody["compute"] = google_cloud_test.Resource("policy.json")
	client := googlecloud.GoogleClient{HttpClient: m}
	infos, _ := client.GetBackendPolicy("anObjectId")

	expectedUsers := []string{
		"user:phil@example.com",
		"group:admins@example.com",
		"domain:google.com",
		"serviceAccount:my-project-id@appspot.gserviceaccount.com",
	}
	assert.Equal(t, 2, len(infos))
	assert.Equal(t, expectedUsers, infos[0].Subject.AuthenticatedUsers)
	assert.Equal(t, []string{"/"}, infos[1].Object.Resources)
}

func TestGoogleClient_GetBackendPolicies_withRequestError(t *testing.T) {
	m := google_cloud_test.NewMockClient()
	m.Err = errors.New("oops")

	client := googlecloud.GoogleClient{HttpClient: m}
	_, err := client.GetBackendPolicy("anObjectId")
	assert.Error(t, err)
}

func TestGoogleClient_GetBackendPolicies_withBadJson(t *testing.T) {
	m := google_cloud_test.NewMockClient()
	m.ResponseBody["compute"] = []byte("-")
	client := googlecloud.GoogleClient{HttpClient: m}
	_, err := client.GetBackendPolicy("anObjectId")
	assert.Error(t, err)
}

func TestGoogleClient_SetBackendPolicies(t *testing.T) {
	policy := policysupport.PolicyInfo{
		Version: "aVersion", Action: "anAction", Subject: policysupport.SubjectInfo{AuthenticatedUsers: []string{"aUser"}}, Object: policysupport.ObjectInfo{Resources: []string{"/"}},
	}
	m := google_cloud_test.NewMockClient()
	client := googlecloud.GoogleClient{HttpClient: m}
	err := client.SetBackendPolicy("anObjectId", policy)
	assert.NoError(t, err)
	assert.Equal(t, "{\"policy\":{\"bindings\":[{\"role\":\"anAction\",\"members\":[\"aUser\"]}]}}\n", string(m.RequestBody))
}

func TestGoogleClient_SetBackendPolicies_withRequestError(t *testing.T) {
	policy := policysupport.PolicyInfo{
		Version: "aVersion", Action: "anAction", Subject: policysupport.SubjectInfo{AuthenticatedUsers: []string{"aUser"}}, Object: policysupport.ObjectInfo{Resources: []string{"/"}},
	}
	m := google_cloud_test.NewMockClient()
	m.Err = errors.New("oops")
	client := googlecloud.GoogleClient{HttpClient: m}
	err := client.SetBackendPolicy("anObjectId", policy)
	assert.Error(t, err)
}