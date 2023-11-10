package types

type RequestStatus string

const (
	RequestStatusQueued     RequestStatus = "queued"
	RequestStatusProcessing RequestStatus = "processing"
	RequestStatusSuccess    RequestStatus = "success"
	RequestStatusInvalid    RequestStatus = "invalid"
	RequestStatusError      RequestStatus = "error"
	RequestStatusFailed     RequestStatus = "failed"
	RequestStatusIgnored    RequestStatus = "ignored"
)

func (s RequestStatus) String() string {
	return string(s)
}

type Action string

const (
	ActionDeployCluster  Action = "deploy-cluster"
	ActionDestroyCluster Action = "destroy-cluster"
	ActionInstall        Action = "install"
	ActionUnInstall      Action = "uninstall"
	ActionUpgrade        Action = "upgrade"
	ActionPause        	 Action = "pause"
	ActionResume         Action = "resume"
)

func (a Action) String() string {
	return string(a)
}

type ServiceGroupType string

const (
	ServiceGroupTypeDKGNodeSecp256k1 ServiceGroupType = "dkgsecp256k1"
	ServiceGroupTypeDKGNodeEd25519   ServiceGroupType = "dkged25519"
	ServiceGroupTypeK8ssandra        ServiceGroupType = "k8ssandra"
	ServiceGroupTypeSSS              ServiceGroupType = "sss"
	ServiceGroupTypeTSS              ServiceGroupType = "tss"
	ServiceGroupTypeRSS              ServiceGroupType = "rss"
	ServiceGroupTypeMetadata         ServiceGroupType = "metadata"
	ServiceGroupTypeHAProxy          ServiceGroupType = "haproxy"
	ServiceGroupTypePDB              ServiceGroupType = "pdb"
	ServiceGroupTypeRedis            ServiceGroupType = "redis"
	ServiceGroupTypeMCI              ServiceGroupType = "mci"
	ServiceGroupTypeMonitoring       ServiceGroupType = "monitoring"
)

func (s ServiceGroupType) String() string {
	return string(s)
}
