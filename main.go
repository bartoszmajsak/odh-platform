package main

import (
	"flag"
	"os"

	"github.com/opendatahub-io/odh-platform/controllers/authorization"
	"github.com/opendatahub-io/odh-platform/controllers/routing"
	"github.com/opendatahub-io/odh-platform/pkg/config"
	pschema "github.com/opendatahub-io/odh-platform/pkg/schema"
	"github.com/opendatahub-io/odh-platform/pkg/spi"
	"github.com/opendatahub-io/odh-platform/version"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.) to ensure that exec-entrypoint and run can make use of them.
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

//nolint:gochecknoglobals //reason: used only here
var (
	scheme               = runtime.NewScheme()
	setupLog             = ctrl.Log.WithName("setup")
	metricsAddr          string
	enableLeaderElection bool
	probeAddr            string
)

func init() { //nolint:gochecknoinits //reason this way we ensure schemes are always registered before we start anything
	pschema.RegisterSchemes(scheme)
}

func main() {
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "odh-platform",
		Metrics: metricsserver.Options{
			BindAddress: metricsAddr,
		},
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	ctrlLog := ctrl.Log.WithName("controllers").
		WithName("odh-platform")
	ctrlLog.Info("creating controller instance", "version", version.Version, "commit", version.Commit, "build-time", version.BuildTime)

	authorizationComponents, errLoad := config.Load(spi.AuthorizationComponent{}, config.GetConfigFile())
	if errLoad != nil {
		setupLog.Error(errLoad, "unable to load config from "+config.GetConfigFile())
		os.Exit(1)
	}

	authorizationConfig := authorization.PlatformAuthorizationConfig{
		Label:        config.GetAuthorinoLabel(),
		Audiences:    config.GetAuthAudience(),
		ProviderName: config.GetAuthProvider(),
	}
	for _, component := range authorizationComponents {
		if err = authorization.NewPlatformAuthorizationReconciler(mgr.GetClient(), ctrlLog, component, authorizationConfig).
			SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "odh-platform-"+component.CustomResourceType.Kind)
			os.Exit(1)
		}
	}

	routingComponents, errLoad := config.Load(spi.RoutingComponent{}, config.GetConfigFile())
	if errLoad != nil {
		setupLog.Error(errLoad, "unable to load config from "+config.GetConfigFile())
		os.Exit(1)
	}

	routingConfig := routing.PlatformRoutingConfiguration{
		IngressSelectorLabel: config.GetIngressSelectorKey(),
		IngressSelectorValue: config.GetIngressSelectorValue(),
		IngressService:       config.GetGatewayService(),
		GatewayNamespace:     config.GetGatewayNamespace(),
	}

	for _, component := range routingComponents {
		if err = routing.NewPlatformRoutingReconciler(
			mgr.GetClient(),
			ctrlLog,
			component,
			routingConfig,
		).
			SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "odh-platform-"+component.CustomResourceType.Kind)
			os.Exit(1)
		}
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("Starting manager")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
