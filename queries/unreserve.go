package queries

import (
	"encoding/json"
	"github.com/mesos/mesos-go/api/v1/lib"
	mastercalls "github.com/mesos/mesos-go/api/v1/lib/master/calls"
	"github.com/minyk/dcos-resources/client"
)

type UnreserveResources struct {
	PrefixMesosMasterApiV1 func() string
	PrefixMesosSlaveApiV0  func(string) string
	PrefixMesosSlaveApiV1  func(string) string
}

func NewUnreserveResources() *UnreserveResources {
	return &UnreserveResources{
		PrefixMesosMasterApiV1: func() string { return "/mesos/api/v1/" },
		PrefixMesosSlaveApiV0:  func(agentid string) string { return "/agent/" + agentid },
		PrefixMesosSlaveApiV1:  func(agentid string) string { return "/agent/" + agentid + "/api/v1" },
	}
}

func (q *UnreserveResources) UnreserveResource(agentid string, role string, principal string, cpus float64, cpusLabel string, mem float64, memLabel string, disk float64, diskLabel string, frameworkLabel string) error {

	var resources []mesos.Resource
	if cpus > 0 {
		resources = append(resources, resourceWithLabel("cpus", role, principal, cpus, cpusLabel, frameworkLabel))
	}
	if mem > 0 {
		resources = append(resources, resourceWithLabel("mem", role, principal, mem, memLabel, frameworkLabel))
	}
	if disk > 0 {
		resources = append(resources, resourceWithLabel("disk", role, principal, disk, diskLabel, frameworkLabel))
	}

	q.UnreserveMesosResource(agentid, resources...)

	return nil

}

func (q *UnreserveResources) DestroyVolume(agentid string, role string, principal string, disk float64, resourceid string, frameworkid string, persistid string, containerpath string, hostpath string) error {

	var resources []mesos.Resource

	resources = append(resources, resourceDiskWithLabel(role, principal, disk, resourceid, frameworkid, persistid, containerpath, ""))
	requestBody := mastercalls.DestroyVolumes(mesos.AgentID{Value: agentid}, resources...)
	requestContent, err := json.Marshal(requestBody)
	_, err = client.HTTPServicePostJSON(q.PrefixMesosMasterApiV1(), requestContent)
	if err != nil {
		return err
	}

	return nil
}

func (q *UnreserveResources) UnreserveResourceAll(agentid string, role string, principal string) error {

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
		rid, _ := getIDsFromLabels(r.GetReservation().GetLabels().GetLabels())
		client.PrintMessage("unreserve resouce: %s", rid)
		err = q.UnreserveMesosResource(agentid, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (q *UnreserveResources) UnreserveMesosResource(agentid string, resources ...mesos.Resource) error {

	body := mastercalls.UnreserveResources(mesos.AgentID{Value: agentid}, resources...)
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

	reservation := mesos.Resource_ReservationInfo{
		Principal: &principal,
		Labels:    &mesosLabels,
	}

	dynamicReservation := mesos.Resource_ReservationInfo{
		Type:      mesos.Resource_ReservationInfo_DYNAMIC.Enum(),
		Role:      &role,
		Principal: &principal,
		Labels:    &mesosLabels,
	}

	var reservations []mesos.Resource_ReservationInfo
	reservations = append(reservations, dynamicReservation)

	return mesos.Resource{
		Type:         mesos.SCALAR.Enum(),
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

	reservation := mesos.Resource_ReservationInfo{
		Principal: &principal,
		Labels:    &mesosLabels,
	}

	dynamicReservation := mesos.Resource_ReservationInfo{
		Type:      mesos.Resource_ReservationInfo_DYNAMIC.Enum(),
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

	volume := mesos.Volume{
		Mode:          mesos.RW.Enum(),
		ContainerPath: containerPath,
	}

	diskinfo := mesos.Resource_DiskInfo{
		Persistence: &persist,
		Volume:      &volume,
		Source:      nil,
	}

	return mesos.Resource{
		Type:         mesos.SCALAR.Enum(),
		Name:         "disk",
		Role:         &role,
		Reservation:  &reservation,
		Reservations: reservations,
		Scalar:       &mesos.Value_Scalar{Value: disk},
		Disk:         &diskinfo,
	}
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
