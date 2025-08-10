package alibaba

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/util"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider"
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
	k8stypes "k8s.io/apimachinery/pkg/types"
	k8saudit "k8s.io/apiserver/pkg/apis/audit"
)

const SLSProviderName = "alibaba-sls"

type SLSProvider struct {
	client SLSClientInterface

	project  string
	logstore string
}

type SLSProviderConfig struct {
	Endpoint    string `yaml:"endpoint" json:"endpoint"`
	Region      string `yaml:"region,omitempty" json:"region,omitempty"`
	AuthVersion string `yaml:"auth_version,omitempty" json:"auth_version,omitempty"`

	Project  string `yaml:"project" json:"project"`
	LogStore string `yaml:"logstore" json:"logstore"`
}

type SLSAuthProvider struct {
	cred credentials.Credential
}

type SLSClientInterface interface {
	GetLogs(project, logstore, topic string, from, to int64, query string,
		lines, offset int64, reverse bool) (*sls.GetLogsResponse, error)
}

var _ provider.Provider = (*SLSProvider)(nil)

func NewSLSProvider(config *SLSProviderConfig) (*SLSProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid SLS provider config: %w", err)
	}
	cred, err := credentials.NewCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("create credential error: %w", err)
	}
	if _, err := cred.GetCredential(); err != nil {
		return nil, fmt.Errorf("get credential error: %w", err)
	}

	client := sls.CreateNormalInterfaceV2(config.Endpoint, &SLSAuthProvider{cred: cred})
	if config.V4Auth() {
		client.SetRegion(config.Region)
		client.SetAuthVersion(sls.AuthV4)
	}

	return &SLSProvider{
		client:   client,
		project:  config.Project,
		logstore: config.LogStore,
	}, nil
}

func (s *SLSProvider) QueryAuditLog(ctx context.Context, params types.QueryAuditLogParams) (types.AuditLogResult, error) {
	var result types.AuditLogResult
	query := s.buildQuery(params)
	req := &sls.GetLogRequest{
		From:    params.StartTime.Unix(),
		To:      params.EndTime.Unix(),
		Topic:   "",
		Lines:   int64(params.Limit),
		Offset:  0,
		Reverse: true,
		Query:   query,
	}

	resp, err := s.client.GetLogs(s.project, s.logstore, req.Topic,
		req.From, req.To, req.Query, req.Lines, req.Offset, req.Reverse)
	if err != nil {
		return result, fmt.Errorf("get logs error: %w", err)
	}

	entries := make([]types.AuditLogEntry, 0, len(resp.Logs))
	for _, item := range resp.Logs {
		entry := s.convertLogToK8sAudit(item)
		entries = append(entries, types.AuditLogEntry(entry))
	}
	result.ProviderQuery = query
	result.Entries = entries
	result.Total = len(entries)

	return result, nil
}

func (s *SLSProvider) buildQuery(params types.QueryAuditLogParams) string {
	query := "*"

	if params.User != "" && params.User != "*" {
		query += fmt.Sprintf(" and user.username: %q", params.User)
	}

	if params.Namespace != "" && params.Namespace != "*" {
		query += fmt.Sprintf(" and objectRef.namespace: %q", params.Namespace)
	}

	if len(params.Verbs) > 0 {
		verbs := make([]string, len(params.Verbs))
		for i, verb := range params.Verbs {
			verbs[i] = fmt.Sprintf("verb: %q", verb)
		}
		query += fmt.Sprintf(" and (%s)", strings.Join(verbs, " or "))
	}

	if len(params.ResourceTypes) > 0 {
		resourceTypes := make([]string, len(params.ResourceTypes))
		for i, rt := range params.ResourceTypes {
			resourceTypes[i] = fmt.Sprintf("objectRef.resource: %q", rt)
		}
		query += fmt.Sprintf(" and (%s)", strings.Join(resourceTypes, " or "))
	}

	if params.ResourceName != "" && params.ResourceName != "*" {
		query += fmt.Sprintf(" and objectRef.name: %q", params.ResourceName)
	}

	return query
}

func (s *SLSProvider) convertLogToK8sAudit(rawLog map[string]string) k8saudit.Event {
	var event k8saudit.Event

	json.Unmarshal([]byte(rawLog["annotations"]), &event.Annotations)
	event.APIVersion = rawLog["apiVersion"]
	event.AuditID = k8stypes.UID(rawLog["auditID"])
	json.Unmarshal([]byte(rawLog["impersonatedUser"]), &event.ImpersonatedUser)
	event.Kind = rawLog["kind"]
	event.Level = k8saudit.Level(rawLog["level"])
	json.Unmarshal([]byte(rawLog["objectRef"]), &event.ObjectRef)
	json.Unmarshal([]byte(rawLog["requestReceivedTimestamp"]), &event.RequestReceivedTimestamp)
	json.Unmarshal([]byte(rawLog["requestObject"]), &event.RequestObject)
	event.RequestURI = rawLog["requestURI"]
	json.Unmarshal([]byte(rawLog["responseStatus"]), &event.ResponseStatus)
	json.Unmarshal([]byte(rawLog["responseObject"]), &event.ResponseObject)
	json.Unmarshal([]byte(rawLog["sourceIPs"]), &event.SourceIPs)
	event.Stage = k8saudit.Stage(rawLog["stage"])
	json.Unmarshal([]byte(rawLog["stageTimestamp"]), &event.StageTimestamp)
	json.Unmarshal([]byte(rawLog["user"]), &event.User)
	event.UserAgent = rawLog["userAgent"]
	event.Verb = rawLog["verb"]

	return event
}

func (c *SLSProviderConfig) Validate() error {
	if c.Endpoint == "" {
		if c.Region != "" {
			c.Endpoint = fmt.Sprintf("%s.log.aliyuncs.com", c.Region)
		}
	}
	if c.Endpoint == "" && c.Region == "" {
		return errors.New("either endpoint or region must be provided")
	}

	if c.V4Auth() && c.Region == "" {
		region, err := util.ParseRegion(c.Endpoint)
		if err == nil && region != "" {
			c.Region = region
		} else {
			return errors.New("region is required when auth_version is v4")
		}
	}

	if c.Project == "" {
		return errors.New("project is required")
	}
	if c.LogStore == "" {
		return errors.New("logstore is required")
	}
	return nil
}

func (c *SLSProviderConfig) V4Auth() bool {
	return c.AuthVersion == "v4" || c.AuthVersion == ""
}

func (a *SLSAuthProvider) GetCredentials() (sls.Credentials, error) {
	cred, err := a.cred.GetCredential()
	if err != nil {
		return sls.Credentials{}, fmt.Errorf("get credential error: %w", err)
	}

	return sls.Credentials{
		AccessKeyID:     tea.StringValue(cred.AccessKeyId),
		AccessKeySecret: tea.StringValue(cred.AccessKeySecret),
		SecurityToken:   tea.StringValue(cred.SecurityToken),
	}, nil
}
