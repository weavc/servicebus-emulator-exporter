package main

type Config struct {
	UserConfig struct {
		Namespaces []struct {
			Name   string `json:"Name"`
			Queues []struct {
				Name       string `json:"Name"`
				Properties struct {
					DeadLetteringOnMessageExpiration    bool   `json:"DeadLetteringOnMessageExpiration"`
					DefaultMessageTimeToLive            string `json:"DefaultMessageTimeToLive"`
					DuplicateDetectionHistoryTimeWindow string `json:"DuplicateDetectionHistoryTimeWindow"`
					ForwardDeadLetteredMessagesTo       string `json:"ForwardDeadLetteredMessagesTo"`
					ForwardTo                           string `json:"ForwardTo"`
					LockDuration                        string `json:"LockDuration"`
					MaxDeliveryCount                    int    `json:"MaxDeliveryCount"`
					RequiresDuplicateDetection          bool   `json:"RequiresDuplicateDetection"`
					RequiresSession                     bool   `json:"RequiresSession"`
				} `json:"Properties"`
			} `json:"Queues"`
			Topics []struct {
				Name       string `json:"Name"`
				Properties struct {
					DefaultMessageTimeToLive            string `json:"DefaultMessageTimeToLive"`
					DuplicateDetectionHistoryTimeWindow string `json:"DuplicateDetectionHistoryTimeWindow"`
					RequiresDuplicateDetection          bool   `json:"RequiresDuplicateDetection"`
				} `json:"Properties"`
				Subscriptions []struct {
					Name       string `json:"Name"`
					Properties struct {
						DeadLetteringOnMessageExpiration bool   `json:"DeadLetteringOnMessageExpiration"`
						DefaultMessageTimeToLive         string `json:"DefaultMessageTimeToLive"`
						LockDuration                     string `json:"LockDuration"`
						MaxDeliveryCount                 int    `json:"MaxDeliveryCount"`
						ForwardDeadLetteredMessagesTo    string `json:"ForwardDeadLetteredMessagesTo"`
						ForwardTo                        string `json:"ForwardTo"`
						RequiresSession                  bool   `json:"RequiresSession"`
					} `json:"Properties"`
					Rules []struct {
						Name       string `json:"Name"`
						Properties struct {
							FilterType        string `json:"FilterType"`
							CorrelationFilter struct {
								ContentType string `json:"ContentType"`
							} `json:"CorrelationFilter"`
						} `json:"Properties"`
					} `json:"Rules,omitempty"`
				} `json:"Subscriptions"`
			} `json:"Topics"`
		} `json:"Namespaces"`
		Logging struct {
			Type string `json:"Type"`
		} `json:"Logging"`
	} `json:"UserConfig"`
}
