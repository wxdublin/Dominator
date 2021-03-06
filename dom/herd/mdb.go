package herd

import (
	"github.com/Symantec/Dominator/dom/mdb"
)

func (herd *Herd) mdbUpdate(mdb *mdb.Mdb) {
	herd.waitForCompletion()
	herd.Lock()
	defer herd.Unlock()
	herd.subsByIndex = make([]*Sub, 0, len(mdb.Machines))
	// Mark for delete all current subs, then later unmark ones in the new MDB.
	subsToDelete := make(map[string]bool)
	for _, sub := range herd.subsByName {
		subsToDelete[sub.hostname] = true
	}
	for _, machine := range mdb.Machines { // Sorted by Hostname.
		sub := herd.subsByName[machine.Hostname]
		if sub == nil {
			sub = new(Sub)
			sub.herd = herd
			sub.hostname = machine.Hostname
			herd.subsByName[machine.Hostname] = sub
		}
		subsToDelete[sub.hostname] = false
		herd.subsByIndex = append(herd.subsByIndex, sub)
		sub.requiredImage = machine.RequiredImage
		sub.plannedImage = machine.PlannedImage
		herd.getImageHaveLock(sub.requiredImage) // Preload.
		herd.getImageHaveLock(sub.plannedImage)
	}
	// Delete flagged subs (those not in the new MDB).
	for subHostname, toDelete := range subsToDelete {
		if toDelete {
			delete(herd.subsByName, subHostname)
		}
	}
	imageUseMap := make(map[string]bool) // Unreferenced by default.
	for _, sub := range herd.subsByName {
		imageUseMap[sub.requiredImage] = true
		imageUseMap[sub.plannedImage] = true
	}
	// Clean up unreferenced images.
	for name, used := range imageUseMap {
		if !used {
			delete(herd.imagesByName, name)
		}
	}
}
