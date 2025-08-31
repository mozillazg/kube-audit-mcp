package gcp

import (
	"encoding/json"
	"testing"
	"time"

	"cloud.google.com/go/logging"
	"github.com/stretchr/testify/assert"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
	"google.golang.org/genproto/googleapis/cloud/audit"
	"google.golang.org/protobuf/types/known/structpb"
	k8sauth "k8s.io/api/authentication/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	k8saudit "k8s.io/apiserver/pkg/apis/audit"
)

func TestConvertLogToK8sAudit(t *testing.T) {
	provider := &CloudLoggingProvider{
		projectId:   "test-project",
		clusterName: "test-cluster",
	}

	testTime := time.Now()
	insertID := "test-insert-id"

	tests := []struct {
		name        string
		logEntry    logging.Entry
		expected    k8saudit.Event
		expectedErr string
	}{
		{
			name: "valid audit log with create verb",
			logEntry: logging.Entry{
				Timestamp: testTime,
				InsertID:  insertID,
				Payload: &audit.AuditLog{
					MethodName:   "io.k8s.core.v1.pods.create",
					ResourceName: "core/v1/namespaces/default/pods/test-pod",
					AuthenticationInfo: &audit.AuthenticationInfo{
						PrincipalEmail: "test@example.com",
					},
					RequestMetadata: &audit.RequestMetadata{
						CallerIp:                "192.168.1.1",
						CallerSuppliedUserAgent: "kubectl/1.21.0",
					},
				},
				Labels: map[string]string{
					"label1": "value1",
				},
				Resource: &mrpb.MonitoredResource{
					Labels: map[string]string{
						"resource_label": "resource_value",
					},
				},
			},
			expected: k8saudit.Event{
				TypeMeta: v1.TypeMeta{
					Kind:       "Event",
					APIVersion: "audit.k8s.io/v1",
				},
				Level:      k8saudit.LevelMetadata,
				AuditID:    k8stypes.UID(insertID),
				Stage:      k8saudit.StageResponseComplete,
				RequestURI: "/core/v1/namespaces/default/pods/test-pod",
				Verb:       "create",
				User: k8sauth.UserInfo{
					Username: "test@example.com",
				},
				SourceIPs: []string{"192.168.1.1"},
				UserAgent: "kubectl/1.21.0",
				ObjectRef: &k8saudit.ObjectReference{
					APIGroup:   "core",
					APIVersion: "v1",
					Namespace:  "default",
					Resource:   "pods",
					Name:       "test-pod",
				},
				ResponseStatus: &v1.Status{
					Status:  "Created (inferred)",
					Code:    201,
					Message: "Created (inferred)",
				},
				RequestReceivedTimestamp: v1.NewMicroTime(testTime),
				StageTimestamp:           v1.NewMicroTime(testTime),
				Annotations: map[string]string{
					"label1":         "value1",
					"resource_label": "resource_value",
				},
			},
		},
		{
			name: "valid audit log with get verb",
			logEntry: logging.Entry{
				Timestamp: testTime,
				InsertID:  insertID,
				Payload: &audit.AuditLog{
					MethodName:   "io.k8s.core.v1.pods.get",
					ResourceName: "core/v1/namespaces/default/pods/test-pod",
					AuthenticationInfo: &audit.AuthenticationInfo{
						PrincipalEmail: "user@example.com",
					},
					RequestMetadata: &audit.RequestMetadata{
						CallerIp:                "10.0.0.1",
						CallerSuppliedUserAgent: "curl/7.68.0",
					},
				},
				Labels: map[string]string{},
				Resource: &mrpb.MonitoredResource{
					Labels: map[string]string{},
				},
			},
			expected: k8saudit.Event{
				TypeMeta: v1.TypeMeta{
					Kind:       "Event",
					APIVersion: "audit.k8s.io/v1",
				},
				Level:      k8saudit.LevelMetadata,
				AuditID:    k8stypes.UID(insertID),
				Stage:      k8saudit.StageResponseComplete,
				RequestURI: "/core/v1/namespaces/default/pods/test-pod",
				Verb:       "get",
				User: k8sauth.UserInfo{
					Username: "user@example.com",
				},
				SourceIPs: []string{"10.0.0.1"},
				UserAgent: "curl/7.68.0",
				ObjectRef: &k8saudit.ObjectReference{
					APIGroup:   "core",
					APIVersion: "v1",
					Namespace:  "default",
					Resource:   "pods",
					Name:       "test-pod",
				},
				ResponseStatus: &v1.Status{
					Status:  "OK (inferred)",
					Code:    200,
					Message: "OK (inferred)",
				},
				RequestReceivedTimestamp: v1.NewMicroTime(testTime),
				StageTimestamp:           v1.NewMicroTime(testTime),
				Annotations:              map[string]string{},
			},
		},
		{
			name: "audit log with exec subresource",
			logEntry: logging.Entry{
				Timestamp: testTime,
				InsertID:  insertID,
				Payload: &audit.AuditLog{
					MethodName:   "io.k8s.core.v1.pods.exec.create",
					ResourceName: "core/v1/namespaces/default/pods/test-pod/exec",
					AuthenticationInfo: &audit.AuthenticationInfo{
						PrincipalEmail: "admin@example.com",
					},
					RequestMetadata: &audit.RequestMetadata{
						CallerIp:                "172.16.0.1",
						CallerSuppliedUserAgent: "kubectl/1.22.0",
					},
				},
				Labels: map[string]string{},
				Resource: &mrpb.MonitoredResource{
					Labels: map[string]string{},
				},
			},
			expected: k8saudit.Event{
				TypeMeta: v1.TypeMeta{
					Kind:       "Event",
					APIVersion: "audit.k8s.io/v1",
				},
				Level:      k8saudit.LevelRequest,
				AuditID:    k8stypes.UID(insertID),
				Stage:      k8saudit.StageResponseStarted,
				RequestURI: "/core/v1/namespaces/default/pods/test-pod/exec",
				Verb:       "create",
				User: k8sauth.UserInfo{
					Username: "admin@example.com",
				},
				SourceIPs: []string{"172.16.0.1"},
				UserAgent: "kubectl/1.22.0",
				ObjectRef: &k8saudit.ObjectReference{
					APIGroup:    "core",
					APIVersion:  "v1",
					Namespace:   "default",
					Resource:    "pods",
					Name:        "test-pod",
					Subresource: "exec",
				},
				ResponseStatus: &v1.Status{
					Status:  "Switching Protocols (inferred)",
					Code:    101,
					Message: "Switching Protocols (inferred)",
				},
				RequestReceivedTimestamp: v1.NewMicroTime(testTime),
				StageTimestamp:           v1.NewMicroTime(testTime),
				Annotations:              map[string]string{},
			},
		},
		{
			name: "audit log with request object",
			logEntry: logging.Entry{
				Timestamp: testTime,
				InsertID:  insertID,
				Payload: &audit.AuditLog{
					MethodName:   "io.k8s.core.v1.pods.create",
					ResourceName: "core/v1/namespaces/default/pods/test-pod",
					AuthenticationInfo: &audit.AuthenticationInfo{
						PrincipalEmail: "test@example.com",
					},
					RequestMetadata: &audit.RequestMetadata{
						CallerIp:                "192.168.1.1",
						CallerSuppliedUserAgent: "kubectl/1.21.0",
					},
					Request: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"apiVersion": {
								Kind: &structpb.Value_StringValue{StringValue: "v1"},
							},
							"kind": {
								Kind: &structpb.Value_StringValue{StringValue: "Pod"},
							},
						},
					},
				},
				Labels: map[string]string{},
				Resource: &mrpb.MonitoredResource{
					Labels: map[string]string{},
				},
			},
			expected: k8saudit.Event{
				TypeMeta: v1.TypeMeta{
					Kind:       "Event",
					APIVersion: "audit.k8s.io/v1",
				},
				Level:      k8saudit.LevelRequest,
				AuditID:    k8stypes.UID(insertID),
				Stage:      k8saudit.StageResponseComplete,
				RequestURI: "/core/v1/namespaces/default/pods/test-pod",
				Verb:       "create",
				User: k8sauth.UserInfo{
					Username: "test@example.com",
				},
				SourceIPs: []string{"192.168.1.1"},
				UserAgent: "kubectl/1.21.0",
				ObjectRef: &k8saudit.ObjectReference{
					APIGroup:   "core",
					APIVersion: "v1",
					Namespace:  "default",
					Resource:   "pods",
					Name:       "test-pod",
				},
				ResponseStatus: &v1.Status{
					Status:  "Created (inferred)",
					Code:    201,
					Message: "Created (inferred)",
				},
				RequestObject: &runtime.Unknown{
					Raw: []byte(`{"apiVersion":"v1","kind":"Pod"}`),
				},
				RequestReceivedTimestamp: v1.NewMicroTime(testTime),
				StageTimestamp:           v1.NewMicroTime(testTime),
				Annotations:              map[string]string{},
			},
		},
		{
			name: "audit log with both request and response objects",
			logEntry: logging.Entry{
				Timestamp: testTime,
				InsertID:  insertID,
				Payload: &audit.AuditLog{
					MethodName:   "io.k8s.core.v1.pods.update",
					ResourceName: "core/v1/namespaces/default/pods/test-pod",
					AuthenticationInfo: &audit.AuthenticationInfo{
						PrincipalEmail: "test@example.com",
					},
					RequestMetadata: &audit.RequestMetadata{
						CallerIp:                "192.168.1.1",
						CallerSuppliedUserAgent: "kubectl/1.21.0",
					},
					Request: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"spec": {
								Kind: &structpb.Value_StructValue{
									StructValue: &structpb.Struct{
										Fields: map[string]*structpb.Value{
											"replicas": {
												Kind: &structpb.Value_NumberValue{NumberValue: 3},
											},
										},
									},
								},
							},
						},
					},
					Response: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"status": {
								Kind: &structpb.Value_StringValue{StringValue: "updated"},
							},
						},
					},
				},
				Labels: map[string]string{},
				Resource: &mrpb.MonitoredResource{
					Labels: map[string]string{},
				},
			},
			expected: k8saudit.Event{
				TypeMeta: v1.TypeMeta{
					Kind:       "Event",
					APIVersion: "audit.k8s.io/v1",
				},
				Level:      k8saudit.LevelRequestResponse,
				AuditID:    k8stypes.UID(insertID),
				Stage:      k8saudit.StageResponseComplete,
				RequestURI: "/core/v1/namespaces/default/pods/test-pod",
				Verb:       "update",
				User: k8sauth.UserInfo{
					Username: "test@example.com",
				},
				SourceIPs: []string{"192.168.1.1"},
				UserAgent: "kubectl/1.21.0",
				ObjectRef: &k8saudit.ObjectReference{
					APIGroup:   "core",
					APIVersion: "v1",
					Namespace:  "default",
					Resource:   "pods",
					Name:       "test-pod",
				},
				ResponseStatus: &v1.Status{
					Status:  "OK (inferred)",
					Code:    200,
					Message: "OK (inferred)",
				},
				RequestObject: &runtime.Unknown{
					Raw: []byte(`{"spec":{"replicas":3}}`),
				},
				ResponseObject: &runtime.Unknown{
					Raw: []byte(`{"status":"updated"}`),
				},
				RequestReceivedTimestamp: v1.NewMicroTime(testTime),
				StageTimestamp:           v1.NewMicroTime(testTime),
				Annotations:              map[string]string{},
			},
		},
		{
			name: "invalid payload type",
			logEntry: logging.Entry{
				Timestamp: testTime,
				InsertID:  insertID,
				Payload:   "invalid payload",
				Labels:    map[string]string{},
				Resource: &mrpb.MonitoredResource{
					Labels: map[string]string{},
				},
			},
			expectedErr: "failed to convert log to k8s audit event, type of payload is not *audit.AuditLog but string",
		},
		{
			name: "cluster-scoped resource",
			logEntry: logging.Entry{
				Timestamp: testTime,
				InsertID:  insertID,
				Payload: &audit.AuditLog{
					MethodName:   "io.k8s.core.v1.nodes.get",
					ResourceName: "core/v1/nodes/worker-node-1",
					AuthenticationInfo: &audit.AuthenticationInfo{
						PrincipalEmail: "admin@example.com",
					},
					RequestMetadata: &audit.RequestMetadata{
						CallerIp:                "10.0.0.1",
						CallerSuppliedUserAgent: "kubectl/1.21.0",
					},
				},
				Labels: map[string]string{},
				Resource: &mrpb.MonitoredResource{
					Labels: map[string]string{},
				},
			},
			expected: k8saudit.Event{
				TypeMeta: v1.TypeMeta{
					Kind:       "Event",
					APIVersion: "audit.k8s.io/v1",
				},
				Level:      k8saudit.LevelMetadata,
				AuditID:    k8stypes.UID(insertID),
				Stage:      k8saudit.StageResponseComplete,
				RequestURI: "/core/v1/nodes/worker-node-1",
				Verb:       "get",
				User: k8sauth.UserInfo{
					Username: "admin@example.com",
				},
				SourceIPs: []string{"10.0.0.1"},
				UserAgent: "kubectl/1.21.0",
				ObjectRef: &k8saudit.ObjectReference{
					APIGroup:   "core",
					APIVersion: "v1",
					Resource:   "nodes",
					Name:       "worker-node-1",
				},
				ResponseStatus: &v1.Status{
					Status:  "OK (inferred)",
					Code:    200,
					Message: "OK (inferred)",
				},
				RequestReceivedTimestamp: v1.NewMicroTime(testTime),
				StageTimestamp:           v1.NewMicroTime(testTime),
				Annotations:              map[string]string{},
			},
		},
		{
			name: "empty resource name",
			logEntry: logging.Entry{
				Timestamp: testTime,
				InsertID:  insertID,
				Payload: &audit.AuditLog{
					MethodName:   "io.k8s.core.v1.list",
					ResourceName: "",
					AuthenticationInfo: &audit.AuthenticationInfo{
						PrincipalEmail: "user@example.com",
					},
					RequestMetadata: &audit.RequestMetadata{
						CallerIp:                "192.168.1.1",
						CallerSuppliedUserAgent: "kubectl/1.21.0",
					},
				},
				Labels: map[string]string{},
				Resource: &mrpb.MonitoredResource{
					Labels: map[string]string{},
				},
			},
			expected: k8saudit.Event{
				TypeMeta: v1.TypeMeta{
					Kind:       "Event",
					APIVersion: "audit.k8s.io/v1",
				},
				Level:      k8saudit.LevelMetadata,
				AuditID:    k8stypes.UID(insertID),
				Stage:      k8saudit.StageResponseComplete,
				RequestURI: "/",
				Verb:       "list",
				User: k8sauth.UserInfo{
					Username: "user@example.com",
				},
				SourceIPs: []string{"192.168.1.1"},
				UserAgent: "kubectl/1.21.0",
				ObjectRef: nil,
				ResponseStatus: &v1.Status{
					Status:  "OK (inferred)",
					Code:    200,
					Message: "OK (inferred)",
				},
				RequestReceivedTimestamp: v1.NewMicroTime(testTime),
				StageTimestamp:           v1.NewMicroTime(testTime),
				Annotations:              map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := provider.convertLogToK8sAudit(tt.logEntry)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.TypeMeta, result.TypeMeta)
			assert.Equal(t, tt.expected.Level, result.Level)
			assert.Equal(t, tt.expected.AuditID, result.AuditID)
			assert.Equal(t, tt.expected.Stage, result.Stage)
			assert.Equal(t, tt.expected.RequestURI, result.RequestURI)
			assert.Equal(t, tt.expected.Verb, result.Verb)
			assert.Equal(t, tt.expected.User, result.User)
			assert.Equal(t, tt.expected.SourceIPs, result.SourceIPs)
			assert.Equal(t, tt.expected.UserAgent, result.UserAgent)
			assert.Equal(t, tt.expected.ObjectRef, result.ObjectRef)
			assert.Equal(t, tt.expected.ResponseStatus, result.ResponseStatus)
			assert.Equal(t, tt.expected.RequestReceivedTimestamp, result.RequestReceivedTimestamp)
			assert.Equal(t, tt.expected.StageTimestamp, result.StageTimestamp)
			assert.Equal(t, tt.expected.Annotations, result.Annotations)

			if tt.expected.RequestObject != nil {
				assert.NotNil(t, result.RequestObject)
				ajson, err := json.Marshal(tt.expected.RequestObject)
				assert.NoError(t, err)
				bjson, err := json.Marshal(result.RequestObject)
				assert.NoError(t, err)
				assert.JSONEq(t, string(ajson), string(bjson))
			} else {
				assert.Nil(t, result.RequestObject)
			}

			if tt.expected.ResponseObject != nil {
				assert.NotNil(t, result.ResponseObject)
				assert.Equal(t, tt.expected.ResponseObject.Raw, result.ResponseObject.Raw)
			} else {
				assert.Nil(t, result.ResponseObject)
			}
		})
	}
}

func TestConvertLogToK8sAudit_EdgeCases(t *testing.T) {
	provider := &CloudLoggingProvider{
		projectId:   "test-project",
		clusterName: "test-cluster",
	}

	testTime := time.Now()
	insertID := "test-insert-id"

	t.Run("nil authentication info", func(t *testing.T) {
		logEntry := logging.Entry{
			Timestamp: testTime,
			InsertID:  insertID,
			Payload: &audit.AuditLog{
				MethodName:         "io.k8s.core.v1.pods.get",
				ResourceName:       "core/v1/namespaces/default/pods/test-pod",
				AuthenticationInfo: nil,
				RequestMetadata: &audit.RequestMetadata{
					CallerIp:                "192.168.1.1",
					CallerSuppliedUserAgent: "kubectl/1.21.0",
				},
			},
			Labels: map[string]string{},
			Resource: &mrpb.MonitoredResource{
				Labels: map[string]string{},
			},
		}

		result, err := provider.convertLogToK8sAudit(logEntry)
		assert.NoError(t, err)
		assert.Equal(t, "", result.User.Username)
	})

	t.Run("nil request metadata", func(t *testing.T) {
		logEntry := logging.Entry{
			Timestamp: testTime,
			InsertID:  insertID,
			Payload: &audit.AuditLog{
				MethodName:   "io.k8s.core.v1.pods.get",
				ResourceName: "core/v1/namespaces/default/pods/test-pod",
				AuthenticationInfo: &audit.AuthenticationInfo{
					PrincipalEmail: "test@example.com",
				},
				RequestMetadata: nil,
			},
			Labels: map[string]string{},
			Resource: &mrpb.MonitoredResource{
				Labels: map[string]string{},
			},
		}

		result, err := provider.convertLogToK8sAudit(logEntry)
		assert.NoError(t, err)
		assert.Equal(t, []string(nil), result.SourceIPs)
		assert.Equal(t, "", result.UserAgent)
	})
}
