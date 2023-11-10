package types

import (
	"fmt"
)

// RedisEndpoint has format redis-master.<NAMESPACE>.svc.<KUBE_CLUSTER_NAME>-<REGION>
// eg: redis-master.redis.svc.sapphire-asia-southeast1
func RedisEndpoint(clusterName, namespace string, region Region) string {
	clusterDomain := GetClusterDomain(clusterName, region)
	return fmt.Sprintf("redis-master.%s.svc.%s:6379", namespace, clusterDomain)
}
