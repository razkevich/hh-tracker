// Package health checks the health of the system for K8s
package health

// IHealther is an interface that packages may satisfy to indicate health
type IHealther interface {
	Health() bool
}

// IHealth is a collection of health checkers
type IHealth interface {
	// Readiness returns true if all registered checkers indicate health
	Readiness() bool

	// Register is used to register a checker for future health checks
	Register(c IHealther)
}

// Health contains controller context for health
type Health struct {
	checkers []IHealther
}

// ProvideHealth is a wire provider for the health checks
func ProvideHealth() IHealth {
	return &Health{}
}

// Readiness of the servcie
func (h *Health) Readiness() bool {
	for _, c := range h.checkers {
		if !c.Health() {
			return false
		}
	}
	return true
}

// Register a checker
func (h *Health) Register(c IHealther) {
	h.checkers = append(h.checkers, c)
}
