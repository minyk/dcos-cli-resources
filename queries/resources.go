package queries

import (
	"encoding/json"
	"errors"
	"github.com/mesos/mesos-go/api/v1/lib"
	"github.com/mesos/mesos-go/api/v1/lib/agent"
	agentcalls "github.com/mesos/mesos-go/api/v1/lib/agent/calls"
	"github.com/mesos/mesos-go/api/v1/lib/master"
	"github.com/minyk/dcos-resources/client"
)

type Resources struct {
	PrefixMesosMasterApiV1 func() string
	PrefixMesosSlaveApiV0  func(string) string
	PrefixMesosSlaveApiV1  func(string) string
}

func NewResources() *Resources {
	return &Resources{
		PrefixMesosMasterApiV1: func() string { return "/mesos/api/v1/" },
		PrefixMesosSlaveApiV0:  func(agentid string) string { return "/agent/" + agentid },
		PrefixMesosSlaveApiV1:  func(agentid string) string { return "/agent/" + agentid + "/api/v1" },
	}
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

	_, err = client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Reservation is successful.")
	}

	return nil
}

func (q *Resources) UnreserveResource(agentid string, role string, principal string, cpus float64, cpusLabel string, mem float64, memLabel string, disk float64, diskLabel string, frameworkLabel string) error {

	resources := []mesos.Resource{}
	if cpus > 0 {
		resources = append(resources, resourceWithLabel("cpus", role, principal, cpus, cpusLabel, frameworkLabel))
	}
	if mem > 0 {
		resources = append(resources, resourceWithLabel("mem", role, principal, mem, memLabel, frameworkLabel))
	}
	if disk > 0 {
		resources = append(resources, resourceWithLabel("disk", role, principal, disk, diskLabel, frameworkLabel))
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

	_, err = client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Unreservation is successful.")
	}

	return nil

}

func (q *Resources) DestroyVolume(agentid string, role string, principal string, disk float64, resourceid string, frameworkid string, persistid string, containerpath string, hostpath string) error {

	resources := []mesos.Resource{}

	resources = append(resources, resourceDiskWithLabel(role, principal, disk, resourceid, frameworkid, persistid, containerpath, ""))

	request_body := master.Call{
		Type: master.Call_DESTROY_VOLUMES,
		DestroyVolumes: &master.Call_DestroyVolumes{
			AgentID: mesos.AgentID{Value: agentid},
			Volumes: resources,
		},
	}

	requestContent, err := json.Marshal(request_body)
	_, err = client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Persistence Volume is successfully removed.")
	}

	return nil
}

func (q *Resources) UnreserveOneResource(agentid string, role string, principal string, resourceType string, resourceValue float64, resourceLabel string, frameworkLabel string) error {

	resources := []mesos.Resource{}
	resources = append(resources, resourceWithLabel(resourceType, role, principal, resourceValue, resourceLabel, frameworkLabel))

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

	_, err = client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
	if err != nil {
		return err
	} else {
		client.PrintMessage("Unreservation is successful.")
	}

	return nil

}

func (q *Resources) UnreserveResourceAll(agentid string, role string, principal string) error {

	client.PrintMessage("Unreserve all resources for %s", role)

	resources, err := getResourcesOnRole(q.PrefixMesosSlaveApiV0(agentid), role, principal)
	if err != nil {
		return err
	}

	// first, trying to destroy persistent volume
	for _, r := range resources {
		if r.GetName() == "disk" && r.GetDisk().GetPersistence().GetID() != "" {
			rid, fid := getIDsFromLabels(r.GetReservation().GetLabels().GetLabels())
			client.PrintMessage("Destroying persistent volumes: %s", rid)
			err = q.DestroyVolume(agentid, role, principal, r.GetScalar().GetValue(), rid, fid, r.GetDisk().GetPersistence().GetID(), r.GetDisk().GetVolume().GetContainerPath(), "")
			if err != nil {
				return err
			}
		}
	}

	for _, r := range resources {
		// TODO we should handle ports resource.
		if r.GetName() != "ports" {
			rid, fid := getIDsFromLabels(r.GetReservation().GetLabels().GetLabels())
			client.PrintMessage("unreserve resouce: %s", rid)
			err = q.UnreserveOneResource(agentid, role, principal, r.GetName(), r.GetScalar().GetValue(), rid, fid)
			if err != nil {
				return err
			}
		}
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

func resourceWithLabel(resourceType string, role string, principal string, cpus float64, resourceid string, frameworkid string) mesos.Resource {
	var labels []mesos.Label
	if frameworkid != "" {
		labelFrameworkID := mesos.Label{
			Key:   "framework_id",
			Value: &frameworkid,
		}
		labels = append(labels, labelFrameworkID)
	}

	labelResourceID := mesos.Label{
		Key:   "resource_id",
		Value: &resourceid,
	}
	labels = append(labels, labelResourceID)

	mesosLabels := mesos.Labels{Labels: labels}

	scala := mesos.Value_Type(mesos.SCALAR)
	reservation := mesos.Resource_ReservationInfo{
		Principal: &principal,
		Labels:    &mesosLabels,
	}

	var rtype mesos.Resource_ReservationInfo_Type
	rtype = 2

	dynamicReservation := mesos.Resource_ReservationInfo{
		Type:      &rtype,
		Role:      &role,
		Principal: &principal,
		Labels:    &mesosLabels,
	}

	var reservations []mesos.Resource_ReservationInfo
	reservations = append(reservations, dynamicReservation)

	return mesos.Resource{
		Type:         &scala,
		Name:         resourceType,
		Role:         &role,
		Reservation:  &reservation,
		Reservations: reservations,
		Scalar:       &mesos.Value_Scalar{Value: cpus},
	}
}

func resourceDiskWithLabel(role string, principal string, disk float64, resourceid string, frameworkid string, persistid string, containerPath string, hostPath string) mesos.Resource {
	var labels []mesos.Label
	if frameworkid != "" {
		labelFrameworkID := mesos.Label{
			Key:   "framework_id",
			Value: &frameworkid,
		}
		labels = append(labels, labelFrameworkID)
	}

	labelResourceID := mesos.Label{
		Key:   "resource_id",
		Value: &resourceid,
	}
	labels = append(labels, labelResourceID)
	mesosLabels := mesos.Labels{Labels: labels}

	scala := mesos.Value_Type(mesos.SCALAR)
	reservation := mesos.Resource_ReservationInfo{
		Principal: &principal,
		Labels:    &mesosLabels,
	}

	var rtype mesos.Resource_ReservationInfo_Type
	rtype = 2

	dynamicReservation := mesos.Resource_ReservationInfo{
		Type:      &rtype,
		Role:      &role,
		Principal: &principal,
		Labels:    &mesosLabels,
	}

	var reservations []mesos.Resource_ReservationInfo
	reservations = append(reservations, dynamicReservation)

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
		Type:         &scala,
		Name:         "disk",
		Role:         &role,
		Reservation:  &reservation,
		Reservations: reservations,
		Scalar:       &mesos.Value_Scalar{Value: disk},
		Disk:         &diskinfo,
	}
}

func (q *Resources) ListResourcesFromNode(agentid string, role string) error {

	resources, err := getResourcesOnRole(q.PrefixMesosSlaveApiV0(agentid), role, "")
	if err != nil {
		return err
	}

	client.PrintMessage("Role\t\tPrincipal\t\tFrameworkID\t\tType\t\tValue\t\tID\t\tPersistentID\t\tContainerPath")
	for i := range resources {
		resource := resources[i]
		rid, fid := getIDsFromLabels(resource.GetReservation().GetLabels().GetLabels())
		if resource.GetName() == "disk" {
			client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s\t\t%s\t\t%s", resource.GetRole(), resource.GetReservation().GetPrincipal(), fid, resource.GetName(), resource.GetScalar().GetValue(), rid, resource.GetDisk().GetPersistence().GetID(), resource.GetDisk().GetVolume().GetContainerPath())
		} else if resource.GetName() == "ports" {
			client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s", resource.GetRole(), resource.GetReservation().GetPrincipal(), fid, resource.GetName(), resource.GetRanges().GoString(), rid)
		} else {
			client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s", resource.GetRole(), resource.GetReservation().GetPrincipal(), fid, resource.GetName(), resource.GetScalar().GetValue(), rid)
		}
	}

	getResourceOnExecutors(q.PrefixMesosSlaveApiV1(agentid), role)

	return nil
}

func getIDsFromLabels(labels []mesos.Label) (string, string) {
	var rid = ""
	var fid = ""
	for _, value := range labels {
		if value.Key == "resource_id" {
			rid = *value.Value
		}
		if value.Key == "framework_id" {
			fid = *value.Value
		}
	}
	return rid, fid
}

//func getResourcesOnRole(urlPath string, role string) (ResourceRole, error) {
//	resourcesFull, err := listResources(urlPath)
//	if err != nil {
//		return nil, err
//	}
//
//	resources := resourcesFull[role]
//	if len(resources) == 0 {
//		return nil, errors.New("no resources are reserved for role")
//	}
//
//	return resources, nil
//}

func getResourcesOnRole(urlPath string, role string, principal string) (ResourceRole, error) {
	resourcesFull, err := listResources(urlPath)
	if err != nil {
		return nil, err
	}

	resources := resourcesFull[role]
	if len(resources) == 0 {
		return nil, errors.New("no resources are reserved for role")
	}

	if principal == "" {
		return resources, nil
	}

	var resourcesOfPrinciapl ResourceRole
	for _, r := range resources {
		if r.GetReservation().GetPrincipal() == principal {
			resourcesOfPrinciapl = append(resourcesOfPrinciapl, r)
		}
	}

	return resourcesOfPrinciapl, nil
}

func listResources(urlPath string) (ReservedResourcesFull, error) {
	response, err := client.HTTPServiceGet(urlPath + "/state")
	if err != nil {
		return nil, err
	}

	agentStateReponse := AgentState{}
	err = json.Unmarshal(response, &agentStateReponse)
	if err != nil {
		return nil, err
	}

	for index := range agentStateReponse.AgentReservedResourcesFull {
		client.PrintVerbose("index: %s", index)
	}

	return agentStateReponse.AgentReservedResourcesFull, nil
}

func getResourceOnExecutors(urlPath string, role string) ([]agent.Response_GetExecutors_Executor, error) {
	allExec, err := getExecutors(urlPath)
	if err != nil {
		return nil, err
	}

	client.PrintMessage("ExecutorID\t\tRole\t\tPrincipal\t\tFrameworkID\t\tType\t\tValue\t\tID\t\tPersistentID\t\tContainerPath")
	for _, exec := range allExec {
		execInfo := exec.GetExecutorInfo()
		for _, r := range execInfo.GetResources() {
			if r.GetAllocationInfo().GetRole() == role && len(r.GetReservations()) > 0 {
				rid, fid := getIDsFromLabels(r.GetReservations()[0].GetLabels().GetLabels())
				if r.GetName() == "disk" {
					client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s\t\t%s\t\t%s", execInfo.GetExecutorID(), r.GetRole(), r.GetReservations()[0].GetPrincipal(), fid, r.GetName(), r.GetScalar().GetValue(), rid, r.GetDisk().GetPersistence().GetID(), r.GetDisk().GetVolume().GetContainerPath())
				} else if r.GetName() == "ports" {
					client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s", execInfo.GetExecutorID(), r.GetRole(), r.GetReservations()[0].GetPrincipal(), fid, r.GetName(), r.GetRanges().GoString(), rid)
				} else {
					client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t%s", execInfo.GetExecutorID(), r.GetRole(), r.GetReservations()[0].GetPrincipal(), fid, r.GetName(), r.GetScalar().GetValue(), rid)
				}
			} else if r.GetAllocationInfo().GetRole() == role && len(r.GetReservations()) <= 0 {
				client.PrintMessage("%s\t\t%s\t\t%s\t\t%s\t\t%s\t\t%f\t\t", execInfo.GetExecutorID(), r.GetRole(), "", "", r.GetName(), r.GetScalar().GetValue())
			}
		}
	}

	return nil, nil
}

func getExecutors(urlPath string) ([]agent.Response_GetExecutors_Executor, error) {
	body := agentcalls.GetExecutors()
	requestContent, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := client.HTTPServicePostJSON(urlPath, requestContent)
	if err != nil {
		return nil, err
	}

	executors := agent.Response{}
	err = json.Unmarshal(resp, &executors)
	if err != nil {
		return nil, err
	}

	return executors.GetExecutors.Executors, nil
}
