package types

import "fmt"

func GetClusterDomain(clusterName string, region Region) string {
	// eg: sapphire-asia-southeast1
	// DONOT MODIFY THIS FORMAT
	// GKE Terraform scripts have been configured to work this way
	return fmt.Sprintf("%s-%s", clusterName, region)
}

func GetFQDN(service, namespace, clusterDomain string) string {
	// eg: http://dkgnode.ds1.svc.merge-test-asia-southeast1
	// DONOT MODIFY THIS FORMAT
	return fmt.Sprintf("http://%s.%s.svc.%s", service, namespace, clusterDomain)
}
