package queries

import (
	"encoding/json"
	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/agent"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/minyk/dcos-resources/client"
)

type Resources struct {
	PrefixCb func() string
}

func NewResources() *Resources {
	return &Resources{}
}

func (q *Resources) ReserveResource(agentid string, role string, principal string, cpus float64, mem float64) error {

	resources := []mesos.Resource{}
	resources = append(resources, resourceCPU(role, principal, cpus))
	resources = append(resources, resourceMEM(role, principal, mem))

	body := master.Call{
		Type: master.Call_RESERVE_RESOURCES,
		ReserveResources: &master.Call_ReserveResources{
			AgentID:   mesos.AgentID{Value: agentid},
			Resources: resources,
		},
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON("/api/v1", requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Reservation is successful.")
	}

	return nil
}

func (q *Resources) UnreserveResource(agentid string, role string, principal string, cpus float64, cpusLabel string, mem float64, memLabel string, disk float64, diskLabel string) error {

	resources := []mesos.Resource{}
	if cpus > 0 {
		resources = append(resources, resourceWithLabel("cpus", role, principal, cpus, cpusLabel))
	}
	if mem > 0 {
		resources = append(resources, resourceWithLabel("mem", role, principal, mem, memLabel))
	}
	if disk > 0 {
		resources = append(resources, resourceWithLabel("disk", role, principal, disk, diskLabel))
	}

	body := master.Call{
		Type: master.Call_UNRESERVE_RESOURCES,
		UnreserveResources: &master.Call_UnreserveResources{
			AgentID:   mesos.AgentID{Value: agentid},
			Resources: resources,
		},
	}

	requestContent, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON("/api/v1", requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Unreservation is successful.")
	}

	return nil

}

func (q *Resources) DestroyVolume(agentid string, role string, principal string, disk float64, resourceid string, persistid string, containerpath string, hostpath string) error {

	resources := []mesos.Resource{}

	resources = append(resources, resourceDiskWithLabel(role, principal, disk, resourceid, persistid, containerpath, ""))

	request_body := master.Call{
		Type: master.Call_DESTROY_VOLUMES,
		DestroyVolumes: &master.Call_DestroyVolumes{
			AgentID: mesos.AgentID{Value: agentid},
			Volumes: resources,
		},
	}

	requestContent, err := json.Marshal(request_body)
	_, err = client.HTTPServicePostJSON("/api/v1", requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Persistence Volume is successfully removed.")
	}

	return nil
}

func (q *Resources) UnreserveResourceAll(agentid string, role string, principal string) error {

	// query current state from agent
	state_body := agent.Call{
		Type: agent.Call_GET_STATE,
	}
	requestContent, err := json.Marshal(state_body)
	response, err := client.HTTPServicePostJSON("/api/v1", requestContent)
	if err != nil {
		return err
	}

	agentStateReponse := master.Response{}
	json.Unmarshal(response, &agentStateReponse)

	// TODO find all resources with role and principal
	cpus := float64(4)
	cpus_resource_id := ""
	mem := float64(4)
	mem_resource_id := ""

	// unreserve cpus, mem resources with resource_id
	resources := []mesos.Resource{}
	resources = append(resources, resourceWithLabel("cpus", role, principal, cpus, cpus_resource_id))
	resources = append(resources, resourceWithLabel("mem", role, principal, mem, mem_resource_id))

	body := master.Call{
		Type: master.Call_UNRESERVE_RESOURCES,
		UnreserveResources: &master.Call_UnreserveResources{
			AgentID:   mesos.AgentID{Value: agentid},
			Resources: resources,
		},
	}

	requestContent, err = json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = client.HTTPServicePostJSON("/api/v1", requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Unreservation is successful.")
	}

	return nil

}

func resourceCPU(role string, principal string, cpus float64) mesos.Resource {
	return resource("cpus", role, principal, cpus)
}

func resourceMEM(role string, principal string, mem float64) mesos.Resource {
	return resource("mem", role, principal, mem)
}

func resource(resourceType string, role string, principal string, cpus float64) mesos.Resource {

	scala := mesos.Value_Type(mesos.SCALAR)
	reservation := mesos.Resource_ReservationInfo{
		Principal: &principal,
	}

	return mesos.Resource{
		Type:        &scala,
		Name:        resourceType,
		Role:        &role,
		Reservation: &reservation,
		Scalar:      &mesos.Value_Scalar{Value: cpus},
	}
}

func resourceWithLabel(resourceType string, role string, principal string, cpus float64, resourceid string) mesos.Resource {

	label := mesos.Label{
		Key:   "resource_id",
		Value: &resourceid,
	}

	labels := []mesos.Label{}
	labels = append(labels, label)
	mesosLabels := mesos.Labels{Labels: labels}

	scala := mesos.Value_Type(mesos.SCALAR)
	reservation := mesos.Resource_ReservationInfo{
		Principal: &principal,
		Labels:    &mesosLabels,
	}

	return mesos.Resource{
		Type:        &scala,
		Name:        resourceType,
		Role:        &role,
		Reservation: &reservation,
		Scalar:      &mesos.Value_Scalar{Value: cpus},
	}
}

func resourceDiskWithLabel(role string, principal string, disk float64, resourceid string, persistid string, containerPath string, hostPath string) mesos.Resource {

	label := mesos.Label{
		Key:   "resource_id",
		Value: &resourceid,
	}

	labels := []mesos.Label{}
	labels = append(labels, label)
	mesosLabels := mesos.Labels{Labels: labels}

	scala := mesos.Value_Type(mesos.SCALAR)
	reservation := mesos.Resource_ReservationInfo{
		Principal: &principal,
		Labels:    &mesosLabels,
	}

	persist := mesos.Resource_DiskInfo_Persistence{
		ID:        persistid,
		Principal: &principal,
	}

	rw := mesos.RW

	volume := mesos.Volume{
		Mode:          &rw,
		ContainerPath: containerPath,
	}

	diskinfo := mesos.Resource_DiskInfo{
		Persistence: &persist,
		Volume:      &volume,
		Source:      nil,
	}

	return mesos.Resource{
		Type:        &scala,
		Name:        "disk",
		Role:        &role,
		Reservation: &reservation,
		Scalar:      &mesos.Value_Scalar{Value: disk},
		Disk:        &diskinfo,
	}
}
