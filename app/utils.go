package app

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aplr/lacuna/docker"
	"github.com/aplr/lacuna/pubsub"
	log "github.com/sirupsen/logrus"
)

func extractSubscriptions(container docker.Container) []pubsub.Subscription {
	subscriptions := make([]pubsub.Subscription, 0)
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)

	// Intermediate storage to hold subscriptions as we process labels
	subscriptionMap := make(map[string]*pubsub.Subscription)

	// Gather subscriptions by processing a container's labels
	for key, value := range container.Labels {
		keyParts := strings.Split(key, ".")

		// Check that the key starts with pubsub.subscription
		if len(keyParts) < 2 || keyParts[0] != labelPrefix || keyParts[1] != "subscription" {
			continue
		}

		// Check that the key has the correct number of parts
		if len(keyParts) != 4 {
			log.Warnf("invalid subscription key: %s, must be in the format 'lacuna.subscription.<name>.<option>'\n", key)
			continue
		}

		// Check that the subscription name is valid
		if !nameRegex.MatchString(keyParts[2]) {
			log.Warnf("invalid subscription name in key: %s, subscription name should be alphanumeric and may contain dashes\n", key)
			continue
		}

		name := strings.ToLower(keyParts[2])

		// Check if subscription already exists in the map
		if _, ok := subscriptionMap[name]; !ok {
			subscriptionMap[name] = &pubsub.Subscription{
				Service: container.Name(),
				Name:    name,
			}
		}

		// Assign the value to the correct field
		switch keyParts[3] {
		case "topic":
			subscriptionMap[name].Topic = value
		case "endpoint":
			subscriptionMap[name].Endpoint = value
		case "ack-deadline":
			deadline, err := time.ParseDuration(value)
			if err != nil {
				log.Warnf("invalid ack-deadline: %s, must be a valid duration\n", value)
				continue
			}
			subscriptionMap[name].AckDeadline = deadline
		case "retain-acked-messages":
			retain, err := strconv.ParseBool(value)
			if err != nil {
				log.Warnf("invalid retain-acked-messages value: %s, must be a valid boolean\n", value)
				continue
			}
			subscriptionMap[name].RetainAckedMessages = retain
		case "retention-duration":
			duration, err := time.ParseDuration(value)
			if err != nil {
				log.Warnf("invalid retention-duration: %s, must be a valid duration\n", value)
				continue
			}
			subscriptionMap[name].RetentionDuration = duration
		case "enable-ordering":
			enable, err := strconv.ParseBool(value)
			if err != nil {
				log.Warnf("invalid enable-ordering value: %s, must be a valid boolean\n", value)
				continue
			}
			subscriptionMap[name].EnableOrdering = enable
		case "expiration-ttl":
			ttl, err := time.ParseDuration(value)
			if err != nil {
				log.Warnf("invalid expiration-ttl: %s, must be a valid duration\n", value)
				continue
			}
			subscriptionMap[name].ExpirationTTL = ttl
		case "filter":
			subscriptionMap[name].Filter = value
		case "deliver-exactly-once":
			deliver, err := strconv.ParseBool(value)
			if err != nil {
				log.Warnf("invalid deliver-exactly-once value: %s, must be a valid boolean\n", value)
				continue
			}
			subscriptionMap[name].DeliverExactlyOnce = deliver
		case "dead-letter-topic":
			subscriptionMap[name].DeadLetterTopic = value
		case "max-dead-letter-delivery-attempts":
			attempts, err := strconv.Atoi(value)
			if err != nil {
				log.Warnf("invalid max-dead-letter-delivery-attempts value: %s, must be a valid integer\n", value)
				continue
			}
			subscriptionMap[name].MaxDeadLetterDeliveryAttempts = attempts
		case "retry-minimum-backoff":
			backoff, err := time.ParseDuration(value)
			if err != nil {
				log.Warnf("invalid retry-minimum-backoff: %s, must be a valid duration\n", value)
				continue
			}
			subscriptionMap[name].RetryMinimumBackoff = &backoff
		case "retry-maximum-backoff":
			backoff, err := time.ParseDuration(value)
			if err != nil {
				log.Warnf("invalid retry-maximum-backoff: %s, must be a valid duration\n", value)
				continue
			}
			subscriptionMap[name].RetryMaximumBackoff = &backoff
		default:
			log.Warnf("skipping invalid subscription key: %s, must be one of 'topic' or 'endpoint'\n", key)
		}
	}

	for _, subscription := range subscriptionMap {

		subscriptions = append(subscriptions, *subscription)

	}

	return subscriptions
}
