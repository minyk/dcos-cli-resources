package commands

import (
	"github.com/minyk/dcos-resources/queries"
	"gopkg.in/alecthomas/kingpin.v3-unstable"
)

type unreserveResourceHandler struct {
	q             *queries.UnreserveResources
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

func (cmd *unreserveResourceHandler) handleUnreserve(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.UnreserveResource(cmd.agentID, cmd.role, cmd.principal, cmd.cpus, cmd.cpuLabel, cmd.mem, cmd.memLabel, cmd.disk, cmd.diskLabel, cmd.frameworkID)
}

func (cmd *unreserveResourceHandler) handleUnreserveAll(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.UnreserveResourceAll(cmd.agentID, cmd.role, cmd.principal)
}

func (cmd *unreserveResourceHandler) handleDestroyPersistVolume(a *kingpin.Application, e *kingpin.ParseElement, c *kingpin.ParseContext) error {
	return cmd.q.DestroyVolume(cmd.agentID, cmd.role, cmd.principal, cmd.disk, cmd.diskLabel, cmd.frameworkID, cmd.persistid, cmd.containerpath, cmd.hostpath)
}

// HandleScheduleSection
func HandleUnreserveResourcesSection(app *kingpin.Application, q *queries.UnreserveResources) {
	HandleUnreserveResourcesCommands(app.Command("unreserve", "Unreserve resources").Alias("unreserves"), q)
	HandleUnreserveResourcesAllCommands(app.Command("unreserve-all", "Unreserve all resources").Alias("unreservesall"), q)
	HandleDestroyPersistVolume(app.Command("destroy-persist-volume", "Destroy persistence volume").Alias("destroyvolume"), q)
}

// HandleScheduleCommand
func HandleUnreserveResourcesCommands(resources *kingpin.CmdClause, q *queries.UnreserveResources) {
	cmd := &unreserveResourceHandler{q: q}
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
func HandleUnreserveResourcesAllCommands(resources *kingpin.CmdClause, q *queries.UnreserveResources) {
	cmd := &unreserveResourceHandler{q: q}
	unReserve := resources.Action(cmd.handleUnreserveAll)
	unReserve.Flag("agent-id", "Agent ID to unreserve").Required().StringVar(&cmd.agentID)
	unReserve.Flag("role", "Role for unreserve").Required().StringVar(&cmd.role)
	unReserve.Flag("principal", "Principal for unreservce").Required().StringVar(&cmd.principal)
}

func HandleDestroyPersistVolume(resources *kingpin.CmdClause, q *queries.UnreserveResources) {
	cmd := &unreserveResourceHandler{q: q}
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
