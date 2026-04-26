package vault

import (
	"log"
	"time"
)

// Alert represents a secret that is nearing or past expiration.
type Alert struct {
	Path      string
	ExpiresAt time.Time
	TTL       time.Duration
}

// Monitor checks a list of secret paths and returns alerts for secrets
// expiring within the given warnBefore duration.
func Monitor(client *Client, paths []string, warnBefore time.Duration) ([]Alert, error) {
	var alerts []Alert

	for _, path := range paths {
		meta, err := client.GetSecretMeta(path)
		if err != nil {
			log.Printf("[WARN] skipping path %q: %v", path, err)
			continue
		}

		if isExpiringSoon(meta.ExpiresAt, warnBefore) {
			alerts = append(alerts, Alert{
				Path:      meta.Path,
				ExpiresAt: meta.ExpiresAt,
				TTL:       meta.TTL,
			})
		}
	}

	return alerts, nil
}

// isExpiringSoon returns true if the given expiry time is within the warn window.
func isExpiringSoon(expiresAt time.Time, warnBefore time.Duration) bool {
	return time.Until(expiresAt) <= warnBefore
}
