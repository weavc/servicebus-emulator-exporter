package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/spf13/cobra"
)

func errHandler(err error) {
	if err != nil {
		v := "[Error] "
		fmt.Printf("%s Encountered unexpected error: %s", v, err.Error())
		os.Exit(-1)
	}
}

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

func generateConfig(ctx context.Context, namespace []string) Config {

	namespaces := []Namespace{}

	for _, ns := range namespace {

		client, err := admin.NewClientFromConnectionString(ns, nil)
		errHandler(err)

		nsProperties, err := client.GetNamespaceProperties(ctx, nil)
		errHandler(err)

		ns := Namespace{
			Name:   nsProperties.Name,
			Queues: getQueues(ctx, client),
			Topics: getTopics(ctx, client),
		}
		namespaces = append(namespaces, ns)
	}

	return Config{UserConfig{Namespaces: namespaces}}
}

func getSubscriptions(ctx context.Context, client *admin.Client, topicName string) []Subscription {

	var subs []Subscription = []Subscription{}

	pager := client.NewListSubscriptionsPager(topicName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Subscriptions {

			mapped := Subscription{
				Name: q.SubscriptionName,
				Properties: SubscriptionProperties{
					DeadLetteringOnMessageExpiration: q.DeadLetteringOnMessageExpiration,
					DefaultMessageTimeToLive:         q.DefaultMessageTimeToLive,
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

func getSubscriptionRules(ctx context.Context, client *admin.Client, topicName string, subscriptionName string) []SubscriptionRule {

	var subRules []SubscriptionRule = []SubscriptionRule{}

	pager := client.NewListRulesPager(topicName, subscriptionName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Rules {

			mapped := SubscriptionRule{
				Name: q.Name,
			}

			sf, ok := q.Filter.(*admin.SQLFilter)
			if ok {
				mapped.Properties = SubscriptionRuleProperties{
					FilterType: "Sql",
					SqlFilter: &SubscriptionRuleSqlFilter{
						SqlExpression: sf.Expression,
					},
				}
			}

			cf, ok := q.Filter.(*admin.CorrelationFilter)
			if ok {
				mapped.Properties = SubscriptionRuleProperties{
					FilterType: "Correlation",
					CorrelationFilter: &SubscriptionRuleCorrelationFilter{
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

func getTopics(ctx context.Context, client *admin.Client) []Topic {
	var topics []Topic = []Topic{}

	pager := client.NewListTopicsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Topics {

			mapped := Topic{
				Name: q.TopicName,
				Properties: TopicProperties{
					DefaultMessageTimeToLive:            q.DefaultMessageTimeToLive,
					DuplicateDetectionHistoryTimeWindow: q.DuplicateDetectionHistoryTimeWindow,
					RequiresDuplicateDetection:          q.RequiresDuplicateDetection,
				},
				Subscriptions: getSubscriptions(ctx, client, q.TopicName),
			}

			topics = append(topics, mapped)
		}
	}

	return topics
}

func getQueues(ctx context.Context, client *admin.Client) []Queue {

	var queues []Queue = []Queue{}

	pager := client.NewListQueuesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Queues {
			mapped := Queue{
				Name: q.QueueName,
				Properties: QueueProperties{
					DeadLetteringOnMessageExpiration:    q.DeadLetteringOnMessageExpiration,
					DefaultMessageTimeToLive:            q.DefaultMessageTimeToLive,
					DuplicateDetectionHistoryTimeWindow: q.DuplicateDetectionHistoryTimeWindow,
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
