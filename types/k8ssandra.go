package types

import (
	"fmt"
	"strings"
)

func FormatSeeds(seeds []string) string {
	if len(seeds) == 0 {
		return ""
	}
	additionalSeeds := ""
	for _, seed := range seeds {
		additionalSeeds += `"` + seed + `",`
	}
	return "[ " + strings.TrimSuffix(additionalSeeds, ",") + " ]"
}

// NewSeed Seeds are of the format <K8SSANDRA_CLUSTER_NAME>-<DATACENTER_NAME>-service.<NAMESPACE>.svc.<KUBE_CLUSTER_NAME>-<REGION>
// eg: sapphire-datacenter-service.k8ssandra.svc.sapphire-asia-southeast1
func NewSeed(clusterName, namespace string, region Region) string {
	clusterDomain := GetClusterDomain(clusterName, region)

	return fmt.Sprintf("%s-%s-service.%s.svc.%s", clusterName, region.Human(), namespace, clusterDomain)
}

func GetDatacenterPodName(clusterName string, r Region) string {
	return clusterName + "-" + r.Human() + "-default-sts-0"
}
