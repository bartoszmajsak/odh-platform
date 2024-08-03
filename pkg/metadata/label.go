package metadata

import "github.com/opendatahub-io/odh-platform/version"

var Labels = struct { //nolint:gochecknoglobals //reason: anonymous struct is used for grouping labels together instead of consts
	AppPartOf       string
	AppComponent    string
	AppName         string
	AppVersion      string
	AppManagedBy    string
	OwnerName       string
	RoutingExported string
}{
	AppPartOf:       "app.kubernetes.io/part-of",
	AppComponent:    "app.kubernetes.io/component",
	AppName:         "app.kubernetes.io/name",
	AppVersion:      "app.kubernetes.io/version",
	AppManagedBy:    "app.kubernetes.io/managed-by",
	OwnerName:       "platform.opendatahub.io/owner-name",
	RoutingExported: "routing.opendatahub.io/exported",
}

func ApplyStandard(source map[string]string) map[string]string {
	target := map[string]string{}

	target[Labels.AppPartOf] = source[Labels.AppName]
	target[Labels.AppComponent] = source[Labels.AppComponent]

	target[Labels.AppVersion] = version.Version
	target[Labels.AppManagedBy] = "odh-platform"

	return target
}