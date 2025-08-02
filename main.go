package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/spf13/cobra"

	"github.com/weavc/servicebus-emulator-exporter/internal"
)

const (
	MaxTimeToLive         time.Duration = time.Hour
	MaxDuplicateDetection time.Duration = time.Minute * 5
)

func main() {

	root := &cobra.Command{
		Use:   "servicebus-emulator-exporter --cs=\"<Connection String>\"",
		Short: "Run serivce bus exporter tool for provided namespaces",
		Run: func(cmd *cobra.Command, args []string) {
			namespaces, err := cmd.Flags().GetStringArray("cs")
			errHandler(err)
			config := generateConfig(cmd.Context(), namespaces)

			b, err := json.MarshalIndent(config, "", "  ")
			errHandler(err)

			fmt.Printf(string(b))
		},
	}

	root.Flags().StringArray("cs", []string{}, "Run exporter for this connection string. Multiple can be provided.")

	root.Execute()
}

func generateConfig(ctx context.Context, namespace []string) internal.Config {

	namespaces := []internal.Namespace{}

	for _, ns := range namespace {

		client, err := admin.NewClientFromConnectionString(ns, nil)
		errHandler(err)

		nsProperties, err := client.GetNamespaceProperties(ctx, nil)
		errHandler(err)

		ns := internal.Namespace{
			Name:   nsProperties.Name,
			Queues: getQueues(ctx, client),
			Topics: getTopics(ctx, client),
		}
		namespaces = append(namespaces, ns)
	}

	return internal.Config{UserConfig: internal.UserConfig{Namespaces: namespaces, Logging: internal.Logging{Type:"console"}}}
}

func getSubscriptions(ctx context.Context, client *admin.Client, topicName string) []internal.Subscription {

	var subs []internal.Subscription = []internal.Subscription{}

	pager := client.NewListSubscriptionsPager(topicName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Subscriptions {

			mapped := internal.Subscription{
				Name: q.SubscriptionName,
				Properties: internal.SubscriptionProperties{
					DefaultMessageTimeToLive:         capDuration(q.DefaultMessageTimeToLive, MaxTimeToLive),
					DeadLetteringOnMessageExpiration: q.DeadLetteringOnMessageExpiration,
					ForwardDeadLetteredMessagesTo:    q.ForwardDeadLetteredMessagesTo,
					ForwardTo:                        q.ForwardTo,
					MaxDeliveryCount:                 q.MaxDeliveryCount,
					LockDuration:                     q.LockDuration,
					RequiresSession:                  q.RequiresSession,
				},
				Rules: getSubscriptionRules(ctx, client, topicName, q.SubscriptionName),
			}

			subs = append(subs, mapped)
		}
	}

	return subs
}

func getSubscriptionRules(ctx context.Context, client *admin.Client, topicName string, subscriptionName string) []internal.SubscriptionRule {

	var subRules []internal.SubscriptionRule = []internal.SubscriptionRule{}

	pager := client.NewListRulesPager(topicName, subscriptionName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Rules {

			mapped := internal.SubscriptionRule{
				Name: q.Name,
			}

			sf, ok := q.Filter.(*admin.SQLFilter)
			if ok {
				mapped.Properties = internal.SubscriptionRuleProperties{
					FilterType: "Sql",
					SqlFilter: &internal.SubscriptionRuleSqlFilter{
						SqlExpression: sf.Expression,
					},
				}
			}

			cf, ok := q.Filter.(*admin.CorrelationFilter)
			if ok {
				mapped.Properties = internal.SubscriptionRuleProperties{
					FilterType: "Correlation",
					CorrelationFilter: &internal.SubscriptionRuleCorrelationFilter{
						ContentType:      cf.ContentType,
						CorrelationId:    cf.CorrelationID,
						Label:            cf.Subject,
						ReplyTo:          cf.ReplyTo,
						ReplyToSessionId: cf.ReplyToSessionID,
						SessionId:        cf.SessionID,
						To:               cf.To,
					},
				}
			}

			subRules = append(subRules, mapped)
		}
	}

	return subRules

}

func getTopics(ctx context.Context, client *admin.Client) []internal.Topic {
	var topics []internal.Topic = []internal.Topic{}

	pager := client.NewListTopicsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Topics {

			mapped := internal.Topic{
				Name: q.TopicName,
				Properties: internal.TopicProperties{
					DuplicateDetectionHistoryTimeWindow: capDuration(q.DuplicateDetectionHistoryTimeWindow, MaxDuplicateDetection),
					DefaultMessageTimeToLive:            capDuration(q.DefaultMessageTimeToLive, MaxTimeToLive),
					RequiresDuplicateDetection:          q.RequiresDuplicateDetection,
				},
				Subscriptions: getSubscriptions(ctx, client, q.TopicName),
			}

			topics = append(topics, mapped)
		}
	}

	return topics
}

func getQueues(ctx context.Context, client *admin.Client) []internal.Queue {

	var queues []internal.Queue = []internal.Queue{}

	pager := client.NewListQueuesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Queues {
			mapped := internal.Queue{
				Name: q.QueueName,
				Properties: internal.QueueProperties{
					DeadLetteringOnMessageExpiration:    q.DeadLetteringOnMessageExpiration,
					DuplicateDetectionHistoryTimeWindow: capDuration(q.DuplicateDetectionHistoryTimeWindow, MaxDuplicateDetection),
					DefaultMessageTimeToLive:            capDuration(q.DefaultMessageTimeToLive, MaxTimeToLive),
					ForwardDeadLetteredMessagesTo:       q.ForwardDeadLetteredMessagesTo,
					ForwardTo:                           q.ForwardTo,
					LockDuration:                        q.LockDuration,
					MaxDeliveryCount:                    q.MaxDeliveryCount,
					RequiresDuplicateDetection:          q.RequiresDuplicateDetection,
					RequiresSession:                     q.RequiresSession,
				},
			}

			queues = append(queues, mapped)
		}
	}

	return queues
}

func errHandler(err error) {
	if err != nil {
		v := "[Error] "
		fmt.Printf("%s Encountered unexpected error: %s", v, err.Error())
		os.Exit(-1)
	}
}

func capDuration(duration *string, maxDuration time.Duration) *string {
	d, err := internal.ISO8601StringToDuration(duration)
	errHandler(err)

	if *d > maxDuration {
		return internal.DurationToStringPtr(&maxDuration)
	}

	return duration
}
