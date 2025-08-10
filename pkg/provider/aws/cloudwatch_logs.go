package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/mozillazg/kube-audit-mcp/pkg/provider"
	"github.com/mozillazg/kube-audit-mcp/pkg/types"
	k8saudit "k8s.io/apiserver/pkg/apis/audit"
)

const CloudWatchProviderName = "aws-cloudwatch-logs"

type CloudWatchLogsProvider struct {
	client *cloudwatchlogs.Client

	logGroupName       string
	logGroupIdentifier string
}

type CloudWatchLogsProviderConfig struct {
	Region             string `yaml:"region,omitempty" json:"region,omitempty"`
	LogGroupName       string `yaml:"log_group_name" json:"log_group_name"`
	LogGroupIdentifier string `yaml:"log_group_identifier,omitempty" json:"log_group_identifier,omitempty"`
}

var _ provider.Provider = (*CloudWatchLogsProvider)(nil)

func NewCloudWatchLogsProvider(config *CloudWatchLogsProviderConfig) (*CloudWatchLogsProvider, error) {
	if err := config.Init(); err != nil {
		return nil, err
	}

	var opts []func(*awsconfig.LoadOptions) error
	if config.Region != "" {
		opts = append(opts, awsconfig.WithRegion(config.Region))
	}
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}
	client := cloudwatchlogs.NewFromConfig(cfg)

	return &CloudWatchLogsProvider{
		client:             client,
		logGroupName:       config.LogGroupName,
		logGroupIdentifier: config.LogGroupIdentifier,
	}, nil
}

func (c *CloudWatchLogsProvider) QueryAuditLog(ctx context.Context, params types.QueryAuditLogParams) (types.AuditLogResult, error) {
	var result types.AuditLogResult
	query := c.buildQuery(params)
	log.Printf("query: %s", query)

	queryResults, err := c.queryLogs(ctx, params, query)
	if err != nil {
		return result, fmt.Errorf("failed to query logs: %w", err)
	}

	entries := make([]types.AuditLogEntry, 0, len(queryResults))
	result.Total = len(entries)
	result.ProviderQuery = query
	for _, item := range queryResults {
		entry, err := c.convertLogToK8sAudit(item)
		if err != nil {
			return result, fmt.Errorf("failed to convert log to k8s audit: %w", err)
		}
		entries = append(entries, types.AuditLogEntry(entry))
	}
	result.Entries = entries

	return result, nil
}

func (c *CloudWatchLogsProvider) queryLogs(ctx context.Context, params types.QueryAuditLogParams, query string) ([]string, error) {
	var logGroupIdentifiers []string
	var logGroupName *string
	if c.logGroupName != "" {
		logGroupName = aws.String(c.logGroupName)
	}
	if c.logGroupIdentifier != "" {
		logGroupIdentifiers = append(logGroupIdentifiers, c.logGroupIdentifier)
	}

	req := cloudwatchlogs.StartQueryInput{
		StartTime:           aws.Int64(params.StartTime.Unix()),
		EndTime:             aws.Int64(params.EndTime.Unix()),
		Limit:               aws.Int32(int32(params.Limit)),
		LogGroupName:        logGroupName,
		LogGroupIdentifiers: logGroupIdentifiers,
		QueryString:         aws.String(query),
	}

	resp, err := c.client.StartQuery(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to start query: %w", err)
	}

	var queryResults []string
getResults:
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("query was cancled: %w", ctx.Err())
		default:
		}

		output, err := c.client.GetQueryResults(context.TODO(), &cloudwatchlogs.GetQueryResultsInput{
			QueryId: resp.QueryId,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get query results: %w", err)
		}

		log.Printf("query status: %s", output.Status)
		switch output.Status {
		case cloudwatchlogstypes.QueryStatusComplete:
			for _, kvals := range output.Results {
				for _, item := range kvals {
					if item.Field == nil || *item.Field != "@message" {
						log.Printf("skipping field %s, only @message is processed", *item.Field)
						continue
					}
					queryResults = append(queryResults, aws.ToString(item.Value))
				}
			}
			break getResults
		case cloudwatchlogstypes.QueryStatusFailed, cloudwatchlogstypes.QueryStatusCancelled, cloudwatchlogstypes.QueryStatusTimeout:
			return nil, fmt.Errorf("query failed with status: %s", output.Status)
		default:
			break
		}

		time.Sleep(1 * time.Second)
	}

	return queryResults, nil
}

func (c *CloudWatchLogsProvider) buildQuery(params types.QueryAuditLogParams) string {
	var filters []string
	query := `fields @message | filter @logStream like "kube-apiserver-audit"`

	if params.User != "" && params.User != "*" {
		exp, val := getFilterExp(params.User)
		filters = append(filters, fmt.Sprintf("user.username %s %q", exp, val))
	}

	if params.Namespace != "" && params.Namespace != "*" {
		exp, val := getFilterExp(params.Namespace)
		filters = append(filters, fmt.Sprintf("objectRef.namespace %s %q", exp, val))
	}

	if len(params.Verbs) > 0 {
		verbs := make([]string, len(params.Verbs))
		for i, verb := range params.Verbs {
			verbs[i] = fmt.Sprintf("%q", verb)
		}
		filters = append(filters, fmt.Sprintf("verb in [%s]", strings.Join(verbs, ", ")))
	}

	if len(params.ResourceTypes) > 0 {
		resourceTypes := make([]string, len(params.ResourceTypes))
		for i, rt := range params.ResourceTypes {
			resourceTypes[i] = fmt.Sprintf("%q", rt)
		}
		filters = append(filters, fmt.Sprintf("objectRef.resource in [%s]",
			strings.Join(resourceTypes, ", ")))
	}

	if params.ResourceName != "" && params.ResourceName != "*" {
		exp, val := getFilterExp(params.ResourceName)
		filters = append(filters, fmt.Sprintf("objectRef.name %s %q", exp, val))
	}

	if len(filters) > 0 {
		query = fmt.Sprintf("%s | filter %s", query, strings.Join(filters, " and "))
	}
	query = fmt.Sprintf("%s | sort @timestamp desc | limit %d", query, params.Limit)

	return query
}

func (c *CloudWatchLogsProvider) convertLogToK8sAudit(rawLog string) (k8saudit.Event, error) {
	var event k8saudit.Event

	err := json.Unmarshal([]byte(rawLog), &event)

	return event, err
}

func (c *CloudWatchLogsProviderConfig) Init() error {
	if c.LogGroupName == "" && c.LogGroupIdentifier == "" {
		return errors.New("either log_group_name or log_group_identifier must be provided")
	}
	if c.LogGroupName != "" && c.LogGroupIdentifier != "" {
		return errors.New("only one of log_group_name or log_group_identifier can be provided")
	}
	return nil
}

func getFilterExp(keyword string) (exp, val string) {
	switch {
	case strings.HasSuffix(keyword, "*"):
		exp = "like"
		val = strings.TrimSuffix(keyword, "*") + "."
		break
	default:
		exp = "="
		val = keyword
	}
	return
}
