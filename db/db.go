package db

import (
	"fmt"
	"strings"
	"sync"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/torusresearch/coordinator-recovery/codec"
	"github.com/torusresearch/coordinator-recovery/service"
	"github.com/torusresearch/coordinator-recovery/types"
)

type Key string

func (k Key) String() string {
	return string(k)
}

const (
	KeyRequest        Key = "request_"
	KeyIgnoredRequest Key = "ignored_request_"
	KeyCluster        Key = "cluster_"
	KeyDKGNodeList    Key = "dkg_node_list"
)

type Cluster struct {
	Region               types.Region                      `json:"region"`
	IsClusterUp          bool                              `json:"is_cluster_up"`
	IsPrimary            bool                              `json:"is_primary"`
	KubernetesRequests   map[uint32]types.RequestStatus    `json:"kubernetes_requests"`
	ServiceGroupRequests map[uint32]types.RequestStatus    `json:"service_group_requests"`
	ServiceGroups        map[string]types.ServiceGroupType `json:"service_groups"`
	Keyspaces            map[string]struct{}               `json:"keyspaces"` // Keeps track of keyspaces created in the region
}

type DB interface {
	SetOrUpdateRequest(requestID uint32, status types.RequestStatus) error
	GetRequest(requestID uint32) (types.RequestStatus, error)
	IgnoreRequest(requestID uint32) error
	IsRequestIgnored(requestID uint32) (bool, error)
	GetAllRequests() (map[string]types.RequestStatus, error)
	StoreCluster(region types.Region, cluster Cluster) error
	UpdateCluster(cluster Cluster) (err error)
	GetCluster(region types.Region) (Cluster, error)
	GetAllClusters() (map[types.Region]Cluster, error)
	CountExistingClusters() (int, error)
	IsPrimary(region types.Region) (bool, error)
	GetPrimary() (types.Region, error)

	GetK8ssandraRegions() (types.Regions, error)

	AddOrUpdateNode(region types.Region, serviceType types.ServiceGroupType, node types.Node) error
	RemoveNode(region types.Region, serviceType types.ServiceGroupType, node types.Node) error
	GetNodeList(region types.Region, serviceType types.ServiceGroupType) (types.NodeList, error)
	GetNode(region types.Region, serviceType types.ServiceGroupType, serviceGroupID string) (types.Node, error)
}

type db struct {
	service         service.DatabaseService
	clusterLock     sync.Mutex
	requestLock     sync.Mutex
	k8ssandraLock   sync.Mutex
	dkgNodeListLock sync.Mutex
	codec           codec.Codec
}

func NewDB(dbpath string) (DB, error) {
	badgerDB, err := badger.Open(badger.DefaultOptions(dbpath))
	if err != nil {
		return nil, err
	}
	databaseService := service.NewDatabaseService(badgerDB)

	c := codec.NewCodec()
	c.Register(Cluster{})

	return &db{
		service: databaseService,
		codec:   c,
	}, nil
}

func (d *db) SetOrUpdateRequest(requestID uint32, status types.RequestStatus) error {
	d.requestLock.Lock()
	defer d.requestLock.Unlock()
	return d.service.Set(getRequestKey(requestID), []byte(status))
}

func (d *db) GetRequest(requestID uint32) (types.RequestStatus, error) {
	value, err := d.service.Get(getRequestKey(requestID))
	if err != nil {
		return "", err
	}

	return types.RequestStatus(value), nil
}

func (d *db) IgnoreRequest(requestID uint32) error {
	d.requestLock.Lock()
	defer d.requestLock.Unlock()
	return d.service.Set(getIgnoredRequestKey(requestID), []byte(types.RequestStatusIgnored))
}

func (d *db) IsRequestIgnored(requestID uint32) (bool, error) {
	value, err := d.service.Get(getIgnoredRequestKey(requestID))
	if err != nil {
		return false, err
	}

	return string(value) == string(types.RequestStatusIgnored), nil
}

func (d *db) GetAllRequests() (map[string]types.RequestStatus, error) {
	requests := make(map[string]types.RequestStatus)
	data, err := d.service.GetFromPrefix(KeyRequest.String())
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		requestID := strings.ReplaceAll(v.Key, KeyRequest.String(), "")
		requests[requestID] = types.RequestStatus(v.Value)
	}
	return requests, nil
}

func (d *db) CountExistingClusters() (int, error) {
	clusters, err := d.GetAllClusters()
	if err != nil {
		return 0, err
	}

	var count int
	for _, c := range clusters {
		if c.IsClusterUp {
			count++
		}
	}

	return count, nil
}

func (d *db) IsPrimary(region types.Region) (bool, error) {
	cluster, err := d.GetCluster(region)
	if err != nil {
		return false, err
	}
	return cluster.IsPrimary, nil
}

func (d *db) StoreCluster(region types.Region, cluster Cluster) error {
	d.clusterLock.Lock()
	defer d.clusterLock.Unlock()

	out, err := d.codec.Encode(cluster)
	if err != nil {
		return err
	}

	key := getClusterKey(region.String())

	return d.service.Set(key, out)
}

func (d *db) UpdateCluster(cluster Cluster) (err error) {
	d.clusterLock.Lock()
	defer d.clusterLock.Unlock()
	out, err := d.codec.Encode(cluster)
	if err != nil {
		return err
	}
	key := getClusterKey(cluster.Region.String())
	return d.service.Set(key, out)
}

func (d *db) GetCluster(region types.Region) (Cluster, error) {
	key := getClusterKey(region.String())
	value, err := d.service.Get(key)
	if err != nil {
		return Cluster{}, err
	}

	var cluster Cluster
	err = d.codec.Decode(value, &cluster)
	if err != nil {
		return Cluster{}, err
	}

	return cluster, nil
}

func (d *db) GetAllClusters() (map[types.Region]Cluster, error) {
	clusters := make(map[types.Region]Cluster)
	data, err := d.service.GetFromPrefix(KeyCluster.String())
	if err != nil {
		return nil, err
	}

	for _, v := range data {
		var cluster Cluster
		err = d.codec.Decode(v.Value, &cluster)
		if err != nil {
			return nil, err
		}
		region := strings.ReplaceAll(v.Key, KeyCluster.String(), "")
		clusters[types.Region(region)] = cluster
	}
	return clusters, nil
}

func (d *db) GetK8ssandraRegions() (types.Regions, error) {
	allClusters, err := d.GetAllClusters()
	if err != nil {
		return types.Regions{}, err
	}
	var k8ssandraRegions types.Regions
	for _, cluster := range allClusters {
		if _, ok := cluster.ServiceGroups["k8ssandra"]; ok {
			k8ssandraRegions = append(k8ssandraRegions, cluster.Region)
		}
	}
	return k8ssandraRegions, err
}

func (d *db) GetPrimary() (types.Region, error) {
	allCluster, err := d.GetAllClusters()
	if err != nil {
		return "", err
	}
	if len(allCluster) == 0 {
		return "", nil
	}
	for _, c := range allCluster {
		if c.IsPrimary {
			return c.Region, nil
		}
	}
	return "", nil
}

func (d *db) AddOrUpdateNode(region types.Region, serviceType types.ServiceGroupType, node types.Node) error {
	d.dkgNodeListLock.Lock()
	defer d.dkgNodeListLock.Unlock()
	nodes, err := d.GetNodeList(region, serviceType)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	if nodes == nil {
		nodes = types.NodeList{}
	}

	if nodes.Contains(node) {
		nodes = nodes.Update(node)
	} else {
		nodes = append(nodes, node)
	}

	out, err := d.codec.Encode(nodes)
	if err != nil {
		return err
	}

	return d.service.Set(getDKGNodeListKey(region.String(), serviceType), out)
}

func (d *db) RemoveNode(region types.Region, serviceType types.ServiceGroupType, node types.Node) error {
	d.dkgNodeListLock.Lock()
	defer d.dkgNodeListLock.Unlock()
	nodes, err := d.GetNodeList(region, serviceType)
	if err != nil {
		return err
	}

	if nodes == nil || !nodes.Contains(node) {
		return fmt.Errorf("node %s does not exist", node.ServiceGroupID)
	}

	nodes = nodes.Remove(node)
	out, err := d.codec.Encode(nodes)
	if err != nil {
		return err
	}

	return d.service.Set(getDKGNodeListKey(region.String(), serviceType), out)
}

func (d *db) GetNodeList(region types.Region, serviceType types.ServiceGroupType) (types.NodeList, error) {
	value, err := d.service.Get(getDKGNodeListKey(region.String(), serviceType))
	if err != nil {
		return nil, err
	}

	var nodes types.NodeList
	err = d.codec.Decode(value, &nodes)
	return nodes, err
}

func (d *db) GetNode(region types.Region, serviceType types.ServiceGroupType, serviceGroupID string) (types.Node, error) {
	nodes, err := d.GetNodeList(region, serviceType)
	if err != nil {
		return types.Node{}, err
	}

	for _, node := range nodes {
		if node.ServiceGroupID == serviceGroupID {
			return node, nil
		}
	}

	return types.Node{}, fmt.Errorf("node %s does not exist", serviceGroupID)
}

func getRequestKey(requestID uint32) string {
	return fmt.Sprintf("%s-%d", KeyRequest.String(), requestID)
}

func getIgnoredRequestKey(requestID uint32) string {
	return fmt.Sprintf("%s-%d", KeyIgnoredRequest.String(), requestID)
}

func getClusterKey(region string) string {
	return KeyCluster.String() + region
}

func getDKGNodeListKey(region string, serviceType types.ServiceGroupType) string {
	return fmt.Sprintf("%s-%s-%s", KeyDKGNodeList.String(), region, serviceType)
}
