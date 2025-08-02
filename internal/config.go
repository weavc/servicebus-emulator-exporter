package internal

type Config struct {
	UserConfig UserConfig `json:"UserConfig"`
}

type UserConfig struct {
	Namespaces []Namespace `json:"Namespaces"`
	Logging Logging `json:"Logging"`
}

type Logging struct {
	Type string `json:"Type"`
}

type Namespace struct {
	Name   string  `json:"Name"`
	Queues []Queue `json:"Queues"`
	Topics []Topic `json:"Topics"`
}

type Topic struct {
	Name          string          `json:"Name"`
	Properties    TopicProperties `json:"Properties"`
	Subscriptions []Subscription  `json:"Subscriptions"`
}

type TopicProperties struct {
	DefaultMessageTimeToLive            *string `json:"DefaultMessageTimeToLive"`
	DuplicateDetectionHistoryTimeWindow *string `json:"DuplicateDetectionHistoryTimeWindow"`
	RequiresDuplicateDetection          *bool   `json:"RequiresDuplicateDetection"`
}

type Queue struct {
	Name       string          `json:"Name"`
	Properties QueueProperties `json:"Properties"`
}

type QueueProperties struct {
	DeadLetteringOnMessageExpiration    *bool   `json:"DeadLetteringOnMessageExpiration,omitempty"`
	DefaultMessageTimeToLive            *string `json:"DefaultMessageTimeToLive,omitempty"`
	DuplicateDetectionHistoryTimeWindow *string `json:"DuplicateDetectionHistoryTimeWindow,omitempty"`
	ForwardDeadLetteredMessagesTo       *string `json:"ForwardDeadLetteredMessagesTo,omitempty"`
	ForwardTo                           *string `json:"ForwardTo,omitempty"`
	LockDuration                        *string `json:"LockDuration,omitempty"`
	MaxDeliveryCount                    *int32  `json:"MaxDeliveryCount,omitempty"`
	RequiresDuplicateDetection          *bool   `json:"RequiresDuplicateDetection,omitempty"`
	RequiresSession                     *bool   `json:"RequiresSession,omitempty"`
}

type Subscription struct {
	Name       string                 `json:"Name"`
	Properties SubscriptionProperties `json:"Properties"`
	Rules      []SubscriptionRule     `json:"Rules,omitempty"`
}

type SubscriptionProperties struct {
	DeadLetteringOnMessageExpiration *bool   `json:"DeadLetteringOnMessageExpiration,omitempty"`
	DefaultMessageTimeToLive         *string `json:"DefaultMessageTimeToLive,omitempty"`
	LockDuration                     *string `json:"LockDuration,omitempty"`
	MaxDeliveryCount                 *int32  `json:"MaxDeliveryCount,omitempty"`
	ForwardDeadLetteredMessagesTo    *string `json:"ForwardDeadLetteredMessagesTo,omitempty"`
	ForwardTo                        *string `json:"ForwardTo,omitempty"`
	RequiresSession                  *bool   `json:"RequiresSession,omitempty"`
}

type SubscriptionRule struct {
	Name       string                     `json:"Name"`
	Properties SubscriptionRuleProperties `json:"Properties"`
}

type SubscriptionRuleProperties struct {
	FilterType        string                             `json:"FilterType"`
	CorrelationFilter *SubscriptionRuleCorrelationFilter `json:"CorrelationFilter,omitempty"`
	SqlFilter         *SubscriptionRuleSqlFilter         `json:"SqlFilter,omitempty"`
}

type SubscriptionRuleCorrelationFilter struct {
	ContentType      *string `json:"ContentType,omitempty"`
	CorrelationId    *string `json:"CorrelationId,omitempty"`
	Label            *string `json:"Label"`
	ReplyTo          *string `json:"ReplyTo"`
	ReplyToSessionId *string `json:"ReplyToSessionId,omitempty"`
	SessionId        *string `json:"SessionId,omitempty"`
	To               *string `json:"To,omitempty"`
}

type SubscriptionRuleSqlFilter struct {
	SqlExpression string `json:"SqlExpression"`
}
