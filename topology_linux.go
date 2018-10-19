// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"io/ioutil"
	"strconv"
	"strings"
)

func topologyFillInfo(info *TopologyInfo) error {
	nodes, err := TopologyNodes()
	if err != nil {
		return err
	}
	info.Nodes = nodes
	if len(info.Nodes) == 1 {
		info.Architecture = SMP
	} else {
		info.Architecture = NUMA
	}
	return nil
}

func TopologyNodes() ([]*TopologyNode, error) {
	nodes := make([]*TopologyNode, 0)

	files, err := ioutil.ReadDir(pathSysDevicesSystemNode())
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, "node") {
			continue
		}
		node := &TopologyNode{}
		nodeId, err := strconv.Atoi(filename[4:])
		if err != nil {
			return nil, err
		}
		node.Id = uint32(nodeId)
		cores, err := coresForNode(int(node.Id))
		if err != nil {
			return nil, err
		}
		node.Cores = cores
		caches, err := cachesForNode(int(node.Id))
		if err != nil {
			return nil, err
		}
		node.Caches = caches
		nodes = append(nodes, node)
	}
	return nodes, nil
}
