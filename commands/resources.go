package commands

import (
	"github.com/minyk/dcos-resources/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type resourcesHandler struct {
	q             *queries.Resources
	agentID       string
	role          string
	principal     string
	frameworkID   string
	cpus          float64
	cpuLabel      string
	mem           float64
	memLabel      string
	disk          float64
	diskLabel     string
	persistid     string
	containerpath string
	hostpath      string
}

func (cmd *resourcesHandler) handleReserve(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.ReserveResource(cmd.agentID, cmd.role, cmd.principal, cmd.cpus, cmd.mem)
}

func (cmd *resourcesHandler) handleUnreserve(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.UnreserveResource(cmd.agentID, cmd.role, cmd.principal, cmd.cpus, cmd.cpuLabel, cmd.mem, cmd.memLabel, cmd.disk, cmd.diskLabel, cmd.frameworkID)
}

func (cmd *resourcesHandler) handleUnreserveAll(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.UnreserveResourceAll(cmd.agentID, cmd.role, cmd.principal)
}

func (cmd *resourcesHandler) handleDestroyPersistVolume(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.DestroyVolume(cmd.agentID, cmd.role, cmd.principal, cmd.disk, cmd.diskLabel, cmd.frameworkID, cmd.persistid, cmd.containerpath, cmd.hostpath)
}

func (cmd *resourcesHandler) handleListResources(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.ListResourcesFromNode(cmd.agentID, cmd.role)
}

// HandleScheduleSection
func HandleResourcesSection(app *kingpin.Application, q *queries.Resources) {
	HandleReserveResourcesCommands(app.Command("reserve", "Reserve resources").Alias("reserves"), q)
	HandleUnreserveResourcesCommands(app.Command("unreserve", "Unreserve resources").Alias("unreserves"), q)
	HandleUnreserveResourcesAllCommands(app.Command("unreserve-all", "Unreserve all resources").Alias("unreservesall"), q)
	HandleDestroyPersistVolume(app.Command("destroy-persist-volume", "Destroy persistence volume").Alias("destroyvolume"), q)
	HandleListResourcesCommands(app.Command("list-resources", "List reserved resources from agent").Alias("listresources"), q)
}

// HandleScheduleCommand
func HandleReserveResourcesCommands(resources *kingpin.CmdClause, q *queries.Resources) {
	cmd := &resourcesHandler{q: q}
	reserve := resources.Action(cmd.handleReserve)
	reserve.Flag("agent-id", "Agent ID to reserve").Required().StringVar(&cmd.agentID)
	reserve.Flag("role", "Role for reserve").Required().StringVar(&cmd.role)
	reserve.Flag("principal", "Principal for reserve").Default("my-principal").StringVar(&cmd.principal)
	reserve.Flag("cpus", "Amount of cpus to reserve").Default("0").Float64Var(&cmd.cpus)
	reserve.Flag("mem", "Amount of memory to reserve. The unit is MB.").Default("0").Float64Var(&cmd.mem)
}

// HandleScheduleCommand
func HandleUnreserveResourcesCommands(resources *kingpin.CmdClause, q *queries.Resources) {
	cmd := &resourcesHandler{q: q}
	unReserve := resources.Action(cmd.handleUnreserve)
	unReserve.Flag("agent-id", "Agent ID to unreserve").Required().StringVar(&cmd.agentID)
	unReserve.Flag("role", "Role for unreserve").Required().StringVar(&cmd.role)
	unReserve.Flag("principal", "Principal for unreserve.").Default("my-principal").StringVar(&cmd.principal)
	unReserve.Flag("framework-id", "Framework ID").Default("").StringVar(&cmd.frameworkID)
	unReserve.Flag("cpus", "Amount of cpus to unreserve").Default("0").Float64Var(&cmd.cpus)
	unReserve.Flag("cpus-resource-id", "Resource id for unreserve action.").Default("").StringVar(&cmd.cpuLabel)
	unReserve.Flag("mem", "Amount of memory to unreserve. The unit is MB.").Default("0").Float64Var(&cmd.mem)
	unReserve.Flag("mem-resource-id", "Resource id for unreserve action.").Default("").StringVar(&cmd.memLabel)
	unReserve.Flag("disk", "Amount of disk to unreserve").Default("0").Float64Var(&cmd.disk)
	unReserve.Flag("disk-resource-id", "Resource id for unreserve action.").Default("").StringVar(&cmd.diskLabel)
}

// Unreserve all resources with role and principal
func HandleUnreserveResourcesAllCommands(resources *kingpin.CmdClause, q *queries.Resources) {
	cmd := &resourcesHandler{q: q}
	unReserve := resources.Action(cmd.handleUnreserveAll)
	unReserve.Flag("agent-id", "Agent ID to unreserve").Required().StringVar(&cmd.agentID)
	unReserve.Flag("role", "Role for unreserve").Required().StringVar(&cmd.role)
	unReserve.Flag("principal", "Principal for unreservce").Required().StringVar(&cmd.principal)
}

func HandleDestroyPersistVolume(resources *kingpin.CmdClause, q *queries.Resources) {
	cmd := &resourcesHandler{q: q}
	destroyPersistVolume := resources.Action(cmd.handleDestroyPersistVolume)
	destroyPersistVolume.Flag("agent-id", "Agent ID to unreserve").Required().StringVar(&cmd.agentID)
	destroyPersistVolume.Flag("role", "Role for unreserve").Required().StringVar(&cmd.role)
	destroyPersistVolume.Flag("principal", "Principal for unreserve.").Default("my-principal").StringVar(&cmd.principal)
	destroyPersistVolume.Flag("disk", "Amount of disk to unreserve").Default("0").Float64Var(&cmd.disk)
	destroyPersistVolume.Flag("disk-resource-id", "Resource id for unreserve action.").Default("").StringVar(&cmd.diskLabel)
	destroyPersistVolume.Flag("disk-persist-id", "Persistence id for unreserve action.").Default("").StringVar(&cmd.persistid)
	destroyPersistVolume.Flag("container-path", "Container path of disk.").Default("").StringVar(&cmd.containerpath)
	destroyPersistVolume.Flag("host-path", "host path of disk.").Default("").StringVar(&cmd.hostpath)
}

func HandleListResourcesCommands(resources *kingpin.CmdClause, q *queries.Resources) {
	cmd := &resourcesHandler{q: q}
	listResources := resources.Action(cmd.handleListResources)
	listResources.Flag("agent-id", "Agent ID to unreserve").Required().StringVar(&cmd.agentID)
	listResources.Flag("role", "Role for unreserve").Required().StringVar(&cmd.role)
}
