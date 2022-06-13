package health_test

import (
	"testing"

	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/health"
)

type testCheck bool

func (c testCheck) Health() bool {
	return bool(c)
}

func TestHealth(t *testing.T) {
	for name, test := range map[string]struct {
		checks         []testCheck
		expectedResult bool
	}{
		"no checks": {
			expectedResult: true,
		},
		"only healthy": {
			checks:         []testCheck{true},
			expectedResult: true,
		},
		"only unhealthy": {
			checks:         []testCheck{false},
			expectedResult: false,
		},
		"some unhealthy": {
			checks:         []testCheck{true, false},
			expectedResult: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			h := health.ProvideHealth()
			for _, c := range test.checks {
				h.Register(c)
			}
			result := h.Readiness()
			if result != test.expectedResult {
				t.Errorf("unexpected result %v", result)
			}
		})
	}
}
