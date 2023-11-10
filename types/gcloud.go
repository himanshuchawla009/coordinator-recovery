package types

import "fmt"

type StaticIPMap struct {
	Region string `json:"region"`
	IP     string `json:"address"`
	Name   string `json:"name"`
}

type Region string

func (r Region) String() string {
	return string(r)
}

func (r Region) Metrics() string {
	switch r {
	case RegionSingapore:
		return "asia_southeast1"
	case RegionUSCentral:
		return "us_central1"
	case RegionBrazil:
		return "southamerica_east1"
	case RegionEUCentral:
		return "europe_central2"
	default:
		return "unknown"
	}
}

func (r Region) Human() string {
	switch r {
	case RegionSingapore:
		return "singapore"
	case RegionUSCentral:
		return "us-central"
	case RegionBrazil:
		return "brazil"
	case RegionEUCentral:
		return "eu-central"
	default:
		return "unknown"
	}
}

func (r Region) Short() string {
	switch r {
	case RegionSingapore:
		return "sg"
	case RegionUSCentral:
		return "usc"
	case RegionBrazil:
		return "bz"
	case RegionEUCentral:
		return "euc"
	default:
		return "unknown"
	}
}

const (
	RegionSingapore Region = "asia-southeast1"
	RegionUSCentral Region = "us-central1"
	RegionBrazil    Region = "southamerica-east1"
	RegionEUCentral Region = "europe-central2"
)

type Regions []Region

var (
	ValidGCloudRegions Regions = []Region{RegionSingapore, RegionUSCentral, RegionBrazil, RegionEUCentral}
)

func (r Regions) String() string {
	if len(r) == 0 {
		return ""
	}
	regions := ""
	for _, region := range r {
		regions += region.String() + ","
	}
	regions = regions[:len(regions)-1]
	return regions
}

// Remove a region from the list of regions
func (r Regions) Remove(region Region) Regions {
	var regions Regions
	for _, r := range r {
		if r != region {
			regions = append(regions, r)
		}
	}
	return regions
}

// Contains returns true if the region is in the list of regions
func (r Regions) Contains(region Region) bool {
	for _, r := range r {
		if r == region {
			return true
		}
	}
	return false
}

func ConstructKubeConfigClusterName(clusterName, projectName string, region Region) string {
	return fmt.Sprintf("gke_%s_%s_%s", projectName, region, clusterName)
}

func IsValidGCloudRegion(region string) bool {
	for _, gcloudRegion := range ValidGCloudRegions {
		if string(gcloudRegion) == region {
			return true
		}
	}

	return false
}

func NewStaticIPName(clusterName, serviceGroupID string) string {
	return fmt.Sprintf("%s-dkgnode-%s", clusterName, serviceGroupID)
}
