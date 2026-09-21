package flows

import (
	"fmt"
	"server/messages"
)

var entries map[string]*messages.FlowEntry

func addEntry(entry messages.FlowEntry) {
	if entries == nil {
		entries = make(map[string]*messages.FlowEntry)
	}

	key := fmt.Sprintf("%s|%s|%s|%s|%s", entry.SrcMac, entry.DstMac, entry.HwProto, entry.SrcIp, entry.DstIp)
	if _, ok := entries[key]; !ok {
		entries[key] = &entry
	} else {
		entries[key].Count++
	}
}

func GetAllEntriesAsList() []*messages.FlowEntry {
	list := make([]*messages.FlowEntry, 0, len(entries))
	for _, v := range entries {
		list = append(list, v)
	}
	return list
}
