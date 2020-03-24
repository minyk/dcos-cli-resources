package queries

import mesos "github.com/mesos/mesos-go/api/v1/lib"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Struct for Mesos API V0
type AgentState struct {
	AgentReservedResourcesFull ReservedResourcesFull `json:"reserved_resources_full,omitempty"`
}

type ReservedResourcesFull map[string]ResourceRole

type ResourceRole []mesos.Resource
