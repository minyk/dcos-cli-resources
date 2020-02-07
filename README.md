DC/OS CLI Subcommand for mesos resource reserve/unreserve
==========================================

**CAUTION: This is not battle-hardened yet. USE AT YOUR OWN RISK.** 

# Subcommands

## resources

```
$ dcos resources

usage: dcos resources [<flags>] <command>


Flags:
  -h, --help              Show context-sensitive help.
  -v, --verbose           Enable extra logging of requests/responses
      --name="resources"  Name of the service instance to query

Commands:
  help [<command> ...]
    Show help.


  reserve --agent-id=AGENT-ID --role=ROLE [<flags>]
    Reserve resources

    --agent-id=AGENT-ID         Agent ID to reserve
    --role=ROLE                 Role for reserve
    --principal="my-principal"  Principal for reserve
    --cpus=0                    Amount of cpus to reserve
    --mem=0                     Amount of memory to reserve. The unit is MB.


  unreserve --agent-id=AGENT-ID --role=ROLE [<flags>]
    Unreserve resources

    --agent-id=AGENT-ID         Agent ID to unreserve
    --role=ROLE                 Role for unreserve
    --principal="my-principal"  Principal for unreserve.
    --cpus=0                    Amount of cpus to unreserve
    --cpus-resource-id=""       Resource id for unreserve action.
    --mem=0                     Amount of memory to unreserve. The unit is MB.
    --mem-resource-id=""        Resource id for unreserve action.

```

### Examples

* reserve

```sh
$ dcos resources reserve --agent-id="AAA-BBB-CCCC" --role="role1" --cpus=1 --mem=1024
```


* unreserve

```sh
$ dcos resources unreserve --agent-id="AAA-BBB-CCCC" --role="role1" --cpus=1 --cpus-resource-id="xxxx" --mem=1024 --mem-resource-id="yyyy"
```

* list resources

```sh
$ dcos resources list-resources --agent-id="ef71ac72-3f3e-4bd8-904a-4db098706e06-S0" --role="ccdb-role"
Role		Principal		Name		Value		ID		PersistentID		ContainerPath
ccdb-role		/ccdb-principal		disk		{5000.000000}		bf6c0a6f-32d5-4ce8-af67-b797c2b2437a		d70914c6-3714-41f0-9532-cd54fd1441d2		cockroach-data
ccdb-role		/ccdb-principal		cpus		{0.100000}		6eccfee1-44bc-42de-92e4-290ee5686394
ccdb-role		/ccdb-principal		mem		{32.000000}		864f5c23-c54f-4608-bdb9-7aec0df3b63f
ccdb-role		/ccdb-principal		disk		{256.000000}		dbb5a068-8bcf-4c15-a1bf-c0795f79d74c
```

# How to

## Build

* Install docker >= 1.13
* Run `make build`(or manually run command in `Makefile`).

## Install

* First of all, see this document for dc/os cli subcommand plugin: https://github.com/dcos/dcos-cli/blob/master/design/plugin.md

* Make direcdories like:
```
$ mkdir -p ~/.dcos/clusters/<cluster-hash>/subcommands/resources
$ mkdir -p ~/.dcos/clusters/<cluster-hash>/subcommands/resources/env/bin
```

* Copy binary into `env/bin`

```
$ cp build/dcos-resources-<OS> ~/.dcos/clusters/<cluster-hash>/subcommands/resources/env/bin/dcos-resources
```

* Copy `plugin.toml`

```
$ cp plugin.toml ~/.dcos/clusters/<cluster-hash>/subcommands/resources/env/
```

* Touch `package.json`

```
$ touch ~/.dcos/clusters/<cluster-hash>/subcommands/resources/package.json
```

From this moment, your dcos-cli should recognize resources subcommand plugin.

```
$ dcos
Command line utility for the Mesosphere Datacenter Operating
System (DC/OS). The Mesosphere DC/OS is a distributed operating
system built around Apache Mesos. This utility provides tools
for easy management of a DC/OS installation.

Available DC/OS commands:

	auth           	Authenticate to DC/OS cluster
	cluster        	Manage your DC/OS clusters
	config         	Manage the DC/OS configuration file
	experimental   	Manage commands that are under development
	help           	Display help information about DC/OS
	job            	Deploy and manage jobs in DC/OS
	marathon       	Deploy and manage applications to DC/OS
	node           	View DC/OS node information
	package        	Install and manage DC/OS software packages
	resources      	Reserve/Unreserve Resources
	service        	Manage DC/OS services
	task           	Manage DC/OS tasks

Get detailed command description with 'dcos <command> --help'.
```

# Acknowledgement

* The client code is heavily adopted from https://github.com/mesosphere/dcos-commons/tree/master/cli
