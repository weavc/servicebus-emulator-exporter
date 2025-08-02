package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/spf13/cobra"

	"github.com/weavc/servicebus-emulator-exporter/config"
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

func generateConfig(ctx context.Context, namespace []string) config.Config {

	namespaces := []config.Namespace{}

	for _, ns := range namespace {

		client, err := admin.NewClientFromConnectionString(ns, nil)
		errHandler(err)

		nsProperties, err := client.GetNamespaceProperties(ctx, nil)
		errHandler(err)

		ns := config.Namespace{
			Name:   nsProperties.Name,
			Queues: getQueues(ctx, client),
			Topics: getTopics(ctx, client),
		}
		namespaces = append(namespaces, ns)
	}

	return config.Config{UserConfig: config.UserConfig{Namespaces: namespaces}}
}

func getSubscriptions(ctx context.Context, client *admin.Client, topicName string) []config.Subscription {

	var subs []config.Subscription = []config.Subscription{}

	pager := client.NewListSubscriptionsPager(topicName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Subscriptions {

			mapped := config.Subscription{
				Name: q.SubscriptionName,
				Properties: config.SubscriptionProperties{
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

func getSubscriptionRules(ctx context.Context, client *admin.Client, topicName string, subscriptionName string) []config.SubscriptionRule {

	var subRules []config.SubscriptionRule = []config.SubscriptionRule{}

	pager := client.NewListRulesPager(topicName, subscriptionName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Rules {

			mapped := config.SubscriptionRule{
				Name: q.Name,
			}

			sf, ok := q.Filter.(*admin.SQLFilter)
			if ok {
				mapped.Properties = config.SubscriptionRuleProperties{
					FilterType: "Sql",
					SqlFilter: &config.SubscriptionRuleSqlFilter{
						SqlExpression: sf.Expression,
					},
				}
			}

			cf, ok := q.Filter.(*admin.CorrelationFilter)
			if ok {
				mapped.Properties = config.SubscriptionRuleProperties{
					FilterType: "Correlation",
					CorrelationFilter: &config.SubscriptionRuleCorrelationFilter{
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

func getTopics(ctx context.Context, client *admin.Client) []config.Topic {
	var topics []config.Topic = []config.Topic{}

	pager := client.NewListTopicsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Topics {

			mapped := config.Topic{
				Name: q.TopicName,
				Properties: config.TopicProperties{
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

func getQueues(ctx context.Context, client *admin.Client) []config.Queue {

	var queues []config.Queue = []config.Queue{}

	pager := client.NewListQueuesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		errHandler(err)

		for _, q := range page.Queues {
			mapped := config.Queue{
				Name: q.QueueName,
				Properties: config.QueueProperties{
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
