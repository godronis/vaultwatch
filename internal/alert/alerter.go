package alert

import "time"

// Alerter is the interface implemented by all alert backends.
type Alerter interface {
	// Send delivers an alert for a secret at the given path that expires at expiresAt.
	Send(path string, expiresAt time.Time) error
}

// MultiAlerter fans out an alert to multiple Alerter implementations.
type MultiAlerter struct {
	Alerters []Alerter
}

// NewMultiAlerter creates a MultiAlerter from the provided list of alerters.
func NewMultiAlerter(alerters ...Alerter) *MultiAlerter {
	return &MultiAlerter{Alerters: alerters}
}

// Send calls Send on every registered alerter, collecting any errors.
// All alerters are attempted even if one fails.
func (m *MultiAlerter) Send(path string, expiresAt time.Time) error {
	var errs []error
	for _, a := range m.Alerters {
		if err := a.Send(path, expiresAt); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return joinErrors(errs)
	}
	return nil
}
