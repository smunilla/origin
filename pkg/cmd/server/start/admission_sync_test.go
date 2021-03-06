package start

import (
	"testing"

	// this package has imports for all the admission controllers used in the kube api server
	// it causes all the admission plugins to be registered, giving us a full listing.
	_ "k8s.io/kubernetes/cmd/kube-apiserver/app"

	"k8s.io/kubernetes/pkg/admission"
	"k8s.io/kubernetes/pkg/util"

	"github.com/openshift/origin/pkg/cmd/server/kubernetes"
)

var admissionPluginsNotUsedByKube = util.NewStringSet(
	"AlwaysAdmit",            // from kube, no need for this by default
	"AlwaysDeny",             // from kube, no need for this by default
	"NamespaceAutoProvision", // from kube, it creates a namespace if a resource is created in a non-existent namespace.  We don't want this behavior
	"SecurityContextDeny",    // from kube, it denies pods that want to use SCC capabilities.  We have different rules to allow this in openshift.

	"BuildByStrategy",          // from origin, only needed for managing builds, not kubernetes resources
	"OriginNamespaceLifecycle", // from origin, only needed for rejecting openshift resources, so not needed by kube
)

func TestKubeAdmissionControllerUsage(t *testing.T) {
	registeredKubePlugins := util.NewStringSet(admission.GetPlugins()...)

	usedAdmissionPlugins := util.NewStringSet(kubernetes.AdmissionPlugins...)

	if missingPlugins := usedAdmissionPlugins.Difference(registeredKubePlugins); len(missingPlugins) != 0 {
		t.Errorf("%v not found", missingPlugins.List())
	}

	if notUsed := registeredKubePlugins.Difference(usedAdmissionPlugins); len(notUsed) != 0 {
		for pluginName := range notUsed {
			if !admissionPluginsNotUsedByKube.Has(pluginName) {
				t.Errorf("%v not used", pluginName)
			}
		}
	}
}
