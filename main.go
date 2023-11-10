package main

import (
	"fmt"
	"log"

	db "github.com/torusresearch/coordinator-recovery/db"
	"github.com/torusresearch/coordinator-recovery/types"
)

func main() {
	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	// db, err := badger.Open(badger.DefaultOptions("../badger/badger"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Connected with db")

	// db.View(func(txn *badger.Txn) error {
	// 	it := txn.NewIterator(badger.DefaultIteratorOptions)
	// 	defer it.Close()
	// 	prefix := []byte("dkg_node_list")
	// 	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
	// 		item := it.Item()
	// 		k := item.Key()
	// 		err := item.Value(func(v []byte) error {
	// 			fmt.Println()
	// 			fmt.Printf("key=%s, value=%v", k, v)
	// 			return nil
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// })

	// defer db.Close()
	// Your code here…

	db, err := db.NewDB("../badger/badger")

	if err != nil {
		log.Fatal(err)
	}
	clusters, err := db.GetCluster(types.RegionBrazil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("clusters", clusters.Keyspaces)
	// region := []types.Region{
	// 	types.RegionBrazil,
	// 	types.RegionSingapore,
	// 	types.RegionEUCentral,
	// 	types.RegionUSCentral,
	// }

	// groupTypes := []types.ServiceGroupType{
	// 	types.ServiceGroupTypeDKGNodeSecp256k1,
	// 	types.ServiceGroupTypeSSS,
	// 	types.ServiceGroupTypeTSS,
	// 	types.ServiceGroupTypeRSS,
	// 	types.ServiceGroupTypeMetadata,
	// 	types.ServiceGroupTypePDB,
	// 	types.ServiceGroupTypeRedis,
	// 	types.ServiceGroupTypeMCI,
	// 	types.ServiceGroupTypeHAProxy,
	// 	types.ServiceGroupTypeK8ssandra,
	// 	types.ServiceGroupTypeMonitoring,
	// }

	// nodes, err := db.GetNodeList(types.RegionBrazil, types.ServiceGroupTypeDKGNodeSecp256k1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // fmt.Println("nodes", nodes)

	// if len(nodes) > 0 {
	// 	// Loop over the map using a for range loop
	// 	for _, value := range nodes {
	// 		oldSubstr := "sapphire-dev-2-1"
	// 		newSubstr := "sapphire-dev-2-5"

	// 		// Replace the old substring with the new substring
	// 		replacedURL := strings.Replace(value.FQDN, oldSubstr, newSubstr, -1)
	// 		fmt.Printf("old: %v, new: %v\n", value.FQDN, replacedURL)
	// 		value.FQDN = replacedURL
	// 		newNode := value
	// 		db.AddOrUpdateNode(types.RegionBrazil, types.ServiceGroupTypeDKGNodeSecp256k1, newNode)
	// 	}
	// }

	// 2D for loop to print all combinations of keys
	// for _, r := range region {
	// 	for _, gt := range groupTypes {
	// 		nodes, err := db.GetNodeList(r, gt)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}

	// 		if len(nodes) > 0 {
	// 			// Loop over the map using a for range loop
	// 			for _, value := range nodes {
	// 				oldSubstr := "sapphire-dev-2-1"
	// 				newSubstr := "sapphire-dev-2-5"

	// 				// Replace the old substring with the new substring
	// 				replacedURL := strings.Replace(value.FQDN, oldSubstr, newSubstr, -1)
	// 				fmt.Printf("old: %v, new: %v\n", value.FQDN, replacedURL)
	// 				value.FQDN = replacedURL
	// 				// newNode := value
	// 				// db.AddOrUpdateNode(r, gt, newNode)
	// 			}
	// 		}

	// 	}
	// }

}
