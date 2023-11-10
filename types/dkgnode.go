package types

import "fmt"

const (
	DefaultDKGNodeServiceName = "dkgnode"
)

type NodeList []Node

func (n NodeList) Remove(node Node) NodeList {
	for i, v := range n {
		if v == node {
			return append(n[:i], n[i+1:]...)
		}
	}
	return n
}

func (n NodeList) Contains(node Node) bool {
	for _, v := range n {
		if v.ServiceGroupID == node.ServiceGroupID {
			return true
		}
	}
	return false
}

func (n NodeList) Update(node Node) NodeList {
	for i, v := range n {
		if v.ServiceGroupID == node.ServiceGroupID {
			n[i] = node
			return n
		}
	}
	return n
}

type Node struct {
	// Service Group ID
	ServiceGroupID string `json:"service_group_id"`
	// FQDN for connecting to the node in a shared network
	FQDN string `json:"fqdn"`
	// AddAtBlock is the block height at which the node is meant to be added
	AddAtBlock uint32 `json:"add_at_block"`
	// RemoveAtBlock is the block height at which the node is meant to be removed
	RemoveAtBlock uint32 `json:"remove_at_block"`
}

func GetDKGNodePublicUrl(serviceGroupID, region, domain string) string {
	return fmt.Sprintf("%s-%s.%s", serviceGroupID, region, domain)
}
