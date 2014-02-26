package server

import (
	"log"

	"negentropia/webserv/store"
)

func updateAllZones() {
	//
	// Scan zones
	//
	zones := store.QueryKeys("z:*")
	if len(zones) < 1 {
		log.Printf("server.updateAllZones: no zone found")
		return
	}

	for _, zone := range zones {
		instanceList := store.QueryField(zone, "instanceList")
		if instanceList == "" {
			log.Printf("server.updateAllZones: zone=%s: no instanceList", zone)
			continue
		}

		instances := store.QuerySet(instanceList)
		if len(instances) < 1 {
			log.Printf("server.updateAllZones: zone=%s: empty instanceList", zone)
			continue
		}

		for _, inst := range instances {
			coord := store.QueryField(inst, "coord")
			mission := store.QueryField(inst, "mission")

			log.Printf("server.updateAllZones: zone=%s instance=%s coord=%s mission=%s", zone, inst, coord, mission)
		}
	}
}
