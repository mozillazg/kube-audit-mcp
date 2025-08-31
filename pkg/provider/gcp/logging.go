package gcp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider"
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
	"google.golang.org/api/iterator"
	"google.golang.org/genproto/googleapis/cloud/audit"
	"google.golang.org/protobuf/types/known/structpb"
	k8sauth "k8s.io/api/authentication/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stypes "k8s.io/apimachinery/pkg/types"
	k8saudit "k8s.io/apiserver/pkg/apis/audit"
)

const CloudLoggingProviderName = "gcp-cloud-logging"

type CloudLoggingProvider struct {
	client cloudLoggingProviderClientInterface

	projectId   string
	clusterName string
}

type cloudLoggingProviderClientInterface interface {
	Entries(ctx context.Context, opts ...logadmin.EntriesOption) *logadmin.EntryIterator
}

type CloudLoggingProviderConfig struct {
	ProjectId   string `json:"project_id" yaml:"project_id"`
	ClusterName string `json:"cluster_name" yaml:"cluster_name"`
}

var _ provider.Provider = (*CloudLoggingProvider)(nil)

func NewCloudLoggingProvider(config *CloudLoggingProviderConfig) (*CloudLoggingProvider, error) {
	if err := config.Init(); err != nil {
		return nil, fmt.Errorf("invalid %s provider config: %w", CloudLoggingProviderName, err)
	}

	client, err := logadmin.NewClient(context.TODO(),
		fmt.Sprintf("projects/%s", config.ProjectId))
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud logging client: %w", err)
	}

	return &CloudLoggingProvider{
		client:      client,
		projectId:   config.ProjectId,
		clusterName: config.ClusterName,
	}, nil
}

func (c *CloudLoggingProvider) QueryAuditLog(ctx context.Context, params types.QueryAuditLogParams) (types.AuditLogResult, error) {
	var result types.AuditLogResult
	query := c.buildQuery(params)
	query += fmt.Sprintf(` AND timestamp >= %q AND timestamp <= %q`,
		params.StartTime.Format(time.RFC3339), params.EndTime.Format(time.RFC3339))
	log.Printf("query: %s", query)

	result.ProviderQuery = query
	queryResults, err := c.queryLogs(ctx, params, query)
	if err != nil {
		return result, fmt.Errorf("failed to query logs: %w", err)
	}
	log.Printf("got %d results", len(queryResults))

	entries := make([]types.AuditLogEntry, 0, len(queryResults))
	result.Total = len(entries)
	result.ProviderQuery = query
	for _, item := range queryResults {
		entry, err := c.convertLogToK8sAudit(*item)
		if err != nil {
			return result, fmt.Errorf("failed to convert log to k8s audit: %w", err)
		}
		entries = append(entries, types.AuditLogEntry(entry))
	}
	result.Entries = entries

	return result, nil
}

func (c *CloudLoggingProvider) buildQuery(params types.QueryAuditLogParams) string {
	query := fmt.Sprintf(`resource.type="k8s_cluster" AND logName=%q`, c.logName())
	if c.clusterName != "" {
		query += fmt.Sprintf(" AND resource.labels.cluster_name=%q", c.clusterName)
	}

	if params.User != "" && params.User != "*" {
		query += fmt.Sprintf(" AND protoPayload.authenticationInfo.principalEmail: %q", params.User)
	}

	if params.Namespace != "" && params.Namespace != "*" {
		query += fmt.Sprintf(` AND protoPayload.resourceName: "/namespaces/%s/"`, params.Namespace)
	}

	if len(params.Verbs) > 0 {
		verbs := make([]string, len(params.Verbs))
		for i, verb := range params.Verbs {
			verbs[i] = fmt.Sprintf(`".%s"`, verb)
		}
		query += fmt.Sprintf(" AND protoPayload.methodName: (%s)", strings.Join(verbs, " OR "))
	}

	if len(params.ResourceTypes) > 0 {
		resourceTypes := make([]string, len(params.ResourceTypes))
		for i, rt := range params.ResourceTypes {
			resourceTypes[i] = fmt.Sprintf(`"/%s/"`, rt)
		}
		query += fmt.Sprintf(" AND protoPayload.resourceName: (%s)", strings.Join(resourceTypes, " OR "))
	}

	if params.ResourceName != "" && params.ResourceName != "*" {
		query += fmt.Sprintf(" AND protoPayload.resourceName =~ %q",
			fmt.Sprintf("/%s$", params.ResourceName))
	}

	return query
}

func (c *CloudLoggingProvider) queryLogs(ctx context.Context, params types.QueryAuditLogParams, query string) ([]*logging.Entry, error) {
	var entries = make([]*logging.Entry, 0, params.Limit)
	iter := c.client.Entries(ctx, logadmin.Filter(query), logadmin.NewestFirst())

	for len(entries) < params.Limit {
		entry, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			return entries, nil
		}
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (c *CloudLoggingProvider) convertLogToK8sAudit(logEntry logging.Entry) (k8saudit.Event, error) {
	var event k8saudit.Event
	var err error

	auditLog, ok := logEntry.Payload.(*audit.AuditLog)
	if !ok {
		return event, fmt.Errorf("failed to convert log to k8s audit event, type of payload is not *audit.AuditLog but %T", logEntry.Payload)
	}

	timestampMicro := v1.NewMicroTime(logEntry.Timestamp)

	verb := c.getVerb(auditLog.MethodName)
	status := c.getStatus(verb)
	objRef := c.getObjectReference(auditLog.GetResourceName())
	level, stage := c.getLevelAndStage(objRef, status, auditLog.Request, auditLog.Response)

	annotations := make(map[string]string, len(logEntry.Labels)+len(logEntry.Resource.Labels))
	for l, v := range logEntry.Labels {
		annotations[l] = v
	}
	for l, v := range logEntry.Resource.Labels {
		annotations[l] = v
	}

	requestObj := c.unmarshalResourceObject(auditLog.Request)
	responseObj := c.unmarshalResourceObject(auditLog.Response)
	userName := ""
	if auditLog.AuthenticationInfo != nil {
		userName = auditLog.AuthenticationInfo.PrincipalEmail
	}
	var sourceIps []string
	if auditLog.RequestMetadata != nil && auditLog.RequestMetadata.CallerIp != "" {
		sourceIps = []string{auditLog.RequestMetadata.CallerIp}
	}
	userAgent := ""
	if auditLog.RequestMetadata != nil {
		userAgent = auditLog.RequestMetadata.CallerSuppliedUserAgent
	}

	event = k8saudit.Event{
		TypeMeta: v1.TypeMeta{
			Kind:       "Event",
			APIVersion: fmt.Sprintf("%s/v1", k8saudit.GroupName),
		},
		Level:      level,
		AuditID:    k8stypes.UID(logEntry.InsertID),
		Stage:      stage,
		RequestURI: fmt.Sprintf("/%s", auditLog.ResourceName),
		Verb:       verb,
		User: k8sauth.UserInfo{
			Username: userName,
		},
		ImpersonatedUser:         nil,
		SourceIPs:                sourceIps,
		UserAgent:                userAgent,
		ObjectRef:                objRef,
		ResponseStatus:           status,
		RequestObject:            requestObj,
		ResponseObject:           responseObj,
		RequestReceivedTimestamp: timestampMicro,
		StageTimestamp:           timestampMicro,
		Annotations:              annotations,
	}

	return event, err
}

func (c *CloudLoggingProvider) logName() string {
	return fmt.Sprintf("projects/%s/logs/cloudaudit.googleapis.com%%2Factivity", c.projectId)
}

func (c *CloudLoggingProviderConfig) Init() error {
	if c.ProjectId == "" {
		return errors.New("project_id is required")
	}

	return nil
}

func (c *CloudLoggingProvider) getObjectReference(resourceName string) *k8saudit.ObjectReference {
	if resourceName == "" {
		return nil
	}

	resourceNameParts := strings.Split(string(resourceName), "/")

	var objRef *k8saudit.ObjectReference
	if len(resourceNameParts) == 6 &&
		resourceNameParts[2] == "namespaces" {
		// The object reference includes a namespace and object name
		objRef = &k8saudit.ObjectReference{
			APIGroup:   resourceNameParts[0],
			APIVersion: resourceNameParts[1],
			Namespace:  resourceNameParts[3],
			Resource:   resourceNameParts[4],
			Name:       resourceNameParts[5],
		}
	} else if len(resourceNameParts) == 5 &&
		resourceNameParts[2] == "namespaces" {
		// The object reference does include a namespace but does not have an
		// object name
		objRef = &k8saudit.ObjectReference{
			APIGroup:   resourceNameParts[0],
			APIVersion: resourceNameParts[1],
			Namespace:  resourceNameParts[3],
			Resource:   resourceNameParts[4],
		}
	} else if len(resourceNameParts) == 4 {
		// The object reference does not include a namespace
		objRef = &k8saudit.ObjectReference{
			APIGroup:   resourceNameParts[0],
			APIVersion: resourceNameParts[1],
			Resource:   resourceNameParts[2],
			Name:       resourceNameParts[3],
		}
	} else if len(resourceNameParts) >= 7 &&
		resourceNameParts[2] == "namespaces" {
		// subresource with namespace
		objRef = &k8saudit.ObjectReference{
			APIGroup:    resourceNameParts[0],
			APIVersion:  resourceNameParts[1],
			Namespace:   resourceNameParts[3],
			Resource:    resourceNameParts[4],
			Name:        resourceNameParts[5],
			Subresource: resourceNameParts[6],
		}
	} else if len(resourceNameParts) >= 5 &&
		resourceNameParts[2] != "namespaces" {
		// subresource without namespace
		objRef = &k8saudit.ObjectReference{
			APIGroup:    resourceNameParts[0],
			APIVersion:  resourceNameParts[1],
			Resource:    resourceNameParts[2],
			Name:        resourceNameParts[3],
			Subresource: resourceNameParts[4],
		}
	} else {
		log.Printf("unable to parse resourcename: %q", resourceName)
		return nil
	}

	return objRef
}

func (c *CloudLoggingProvider) getVerb(methodName string) string {
	methodNameParts := strings.Split(methodName, ".")
	return methodNameParts[len(methodNameParts)-1]
}

func (c *CloudLoggingProvider) getLevelAndStage(objRef *k8saudit.ObjectReference,
	status *v1.Status, req, resp *structpb.Struct) (k8saudit.Level, k8saudit.Stage) {
	var level k8saudit.Level
	var stage k8saudit.Stage

	if objRef != nil && (objRef.Subresource == "attach" ||
		objRef.Subresource == "exec") {
		level = k8saudit.LevelRequest
		stage = k8saudit.StageResponseStarted
		status.Code = 101
		status.Status = "Switching Protocols (inferred)"
		status.Message = "Switching Protocols (inferred)"
	} else {
		switch {
		case req == nil && resp == nil:
			level = k8saudit.LevelMetadata
			break
		case req != nil && resp == nil:
			level = k8saudit.LevelRequest
			break
		default:
			level = k8saudit.LevelRequestResponse
		}
		stage = k8saudit.StageResponseComplete
	}

	return level, stage
}

func (c *CloudLoggingProvider) getStatus(verb string) *v1.Status {
	if verb == "create" {
		return &v1.Status{
			Status:  "Created (inferred)",
			Code:    201,
			Message: "Created (inferred)",
		}
	}

	return &v1.Status{
		Status:  "OK (inferred)",
		Code:    200,
		Message: "OK (inferred)",
	}
}

func (c *CloudLoggingProvider) unmarshalResourceObject(obj *structpb.Struct) *runtime.Unknown {
	if obj == nil {
		return nil
	}
	objStruct, err := obj.MarshalJSON()
	if err != nil {
		log.Printf("failed to marshal to json: %v", err)
		return nil
	}
	return &runtime.Unknown{Raw: objStruct}
}
