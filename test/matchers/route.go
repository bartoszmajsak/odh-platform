package matchers

import (
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	openshiftroutev1 "github.com/openshift/api/route/v1"
)

func BeAttachedToService(svcName string) types.GomegaMatcher {
	return &routeSvcMatcher{expectedSvcName: svcName}
}

type routeSvcMatcher struct {
	expectedSvcName string
}

func (r *routeSvcMatcher) Match(actual any) (bool, error) {
	if actual == nil {
		return true, nil
	}

	route, errDeref := deref[openshiftroutev1.Route](actual)
	if errDeref != nil {
		return false, errDeref
	}

	match, err := gomega.Equal("Service").Match(route.Spec.To.Kind)
	if !match || err != nil {
		return match, err
	}

	return gomega.Equal(r.expectedSvcName).Match(route.Spec.To.Name)
}

func (r *routeSvcMatcher) FailureMessage(actual any) string {
	return format.Message(actual, "to be attached to service named", r.expectedSvcName)
}

func (r *routeSvcMatcher) NegatedFailureMessage(actual any) string {
	return format.Message(actual, "not to be attached to service named", r.expectedSvcName)
}

func HaveHostPrefix(name string) types.GomegaMatcher {
	return &routeHostPrefix{expectedHostPrefix: name}
}

type routeHostPrefix struct {
	expectedHostPrefix string
}

func (matcher *routeHostPrefix) Match(actual any) (bool, error) {
	if actual == nil {
		return true, nil
	}

	route, errDeref := deref[openshiftroutev1.Route](actual)
	if errDeref != nil {
		return false, errDeref
	}

	return gomega.HavePrefix(matcher.expectedHostPrefix).Match(route.Spec.Host)
}

func (matcher *routeHostPrefix) FailureMessage(actual any) string {
	return format.Message(actual, "to have host prefix", matcher.expectedHostPrefix)
}

func (matcher *routeHostPrefix) NegatedFailureMessage(actual any) string {
	return format.Message(actual, "to not have host prefix", matcher.expectedHostPrefix)
}
