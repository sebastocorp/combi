# COMBI (Config Combinator)

![GitHub Release](https://img.shields.io/github/v/release/freepik-company/combi)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/freepik-company/combi)
![GitHub License](https://img.shields.io/github/license/freepik-company/combi)

![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/freepik-company/combi/total)

![GitHub User's stars](https://img.shields.io/github/stars/sebastocorp)
![GitHub followers](https://img.shields.io/github/followers/sebastocorp)

## Description

Combi is a simple tool that consumes, update and merge diferents configurations in different formats to generate a final usable configuration and performs some defined actions at the end of the process. This tool consumes its own configuration from a source (like a file in a git repository), and perform a merge with the patches defined in that configuration in a local configuration file with the same format as the patches (libconfig, yaml, json, etc).

## Motivation

Many services used daily in the industry receive their configuration mostly from a file with a specific format (yaml, json, libconfig, etc.), and many of these services play a critical role, so downtime has to be minimal. or, if possible, it should be none.

For this reason, many of these services have functionalities to collect the updated configuration again without stopping its execution. Thanks to these functionalities, the configuration can be modified and applied without practically any downtime, but the problem appears, as always, when you have a multitude of instances of that service, with different configurations, in different environments.

The idea is to avoid restarting the service and easily update the configuration, so some problems arise here:

- In a container environment you do't want to rotate them so as not to generate downtime and lost requests, even if it is minimal and you can control it more or less with a rotation strategy.
- And whether with containers or vms, if you want to avoid the restart, you will have to modify the specific configuration they have, enter each of the instances to update it, and execute the functionality that refreshes the configuration.
- If you have different configurations of the same service in different instances, those configurations will often be splitted in different repositories, or if it is mono-repo, they will be splitted in different parts of the repository, or in different files.

It may also be the case that you want 2 parts of the same config separated from each other, and different teams modify one of these parts separately, one of the two configurations will be the base config, with optional fields, and the other it will be a chore config, with mandatory fields, and you want performing a merge of the configurations with precedence of the chore configuration.

Thinking about these problems and possible solutions, we have decided to create this tool, which is not only capable of centralizing the different configurations of the same format, but is also capable of performing checks on the final configuration and executing the commands necessary to refresh the configuration.

## How to use

This project provides the binary files and container image in differents architectures to make it easy to use wherever wanted.

```sh
combi run \
    --config=/config/combi.yaml
```

### Flags

| Name     | Command | Default      | Description |
|:---      |:---     |:---          |:---         |
| `config` | `run`   | `combi.yaml` | Filepath where configuration in located. |

### Configuration

Current configuration api: `v1alpha4`

| Field | Description |
|:--- |:--- |
| `kind`                                 | Type of the main configuration, specified as a string (values: "LIBCONFIG", "JSON", "YAML", "NGINX"). |
| `settings.logger.level`                | Log level (e.g., "info", "debug", "error"). |
| `settings.syncTime`                    | Sync time for the configuration, specified as `time.Duration`. |
| `settings.tmpObjs.path`                | Path where temporary objects will be stored. |
| `settings.tmpObjs.mode`                | Access mode for the temporary objects, specified as a `uint32`. |
| `settings.target.path`                 | Path where final configuration file will be stored. |
| `settings.target.file`                 | Name of the file at the target location. |
| `settings.target.mode`                 | Access mode for the final configuration file, specified as a `uint32`. |
| `sources`                              | List of data source configurations, which can be of type RAW, FILE, GIT, or K8S. |
| `sources[].name`                       | Name of the data source or Kubernetes resource. |
| `sources[].type`                       | Type of the source (values: "RAW", "FILE", "GIT" "K8S"). |
| `sources[].raw`                        | Raw data for sources of type "RAW". |
| `sources[].file`                       | Filepath for sources of type "FILE". |
| `sources[].git.sshUrl`                 | SSH URL for the Git repository to clone. |
| `sources[].git.sshKeyFilepath`         | Path to the SSH key file for authenticating with the Git repository. |
| `sources[].git.branch`                 | Git branch to use. |
| `sources[].git.filepath`               | Path to the file within the Git repository. |
| `sources[].k8s.context.inCluster`      | Boolean indicating whether the Kubernetes configuration is within the cluster. |
| `sources[].k8s.context.configFilepath` | Path to the Kubernetes configuration file. |
| `sources[].k8s.context.masterUrl`      | URL of the Kubernetes master server. |
| `sources[].k8s.kind`                   | Kind of the Kubernetes resource (values: "Secret", "ConfigMap"). |
| `sources[].k8s.namespace`              | Kubernetes namespace where the resource is located. |
| `sources[].k8s.name`                   | Kubernetes resource name. |
| `sources[].k8s.key`                    | Key to identify the resource within Kubernetes. |
| `behavior.conditions`                  | List of conditions to evaluate before taking actions. |
| `behavior.conditions[].name`           | Name of the condition. |
| `behavior.conditions[].mandatory`      | Indicates whether the condition is mandatory or optional. |
| `behavior.conditions[].template`       | Template that defines how the condition should be evaluated. |
| `behavior.conditions[].expect`         | Expected value for the condition to be considered true. |
| `behavior.actions`                     | List of actions to perform when conditions are met. |
| `behavior.actions[].name`              | Name of the action. |
| `behavior.actions[].on`                | Event that triggers the action whn the conditions ends with fail or success  (values: "SUCCESS", "FAILURE"). |
| `behavior.actions[].command`           | Command to execute as part of the action. |

> **WARNING**
>
> The list of actions are commands that are executed on the machine after checking the conditions. Please be careful.

## How does it work?

Synchronization flow diagram:

```sh
               ┌─────────────┐                                   
               │             │                                   
 ┌────┬────┬───►  sync time  │                                   
 │    │    │   │             │                                   
 │    │    │   └──────┬──────┘                                   
 │    │    │          │                                          
 │    │    │  ┌───────▼───────┐    ┌──────────┐                  
 │    │    │  │               │    │          │  │local file     
 │    │    no │  get  config  ◄────┤  source  ├─►│git repo       
 │    │    │  │               │    │          │  │...            
 │    │    │  └───────┬───────┘    └──────────┘                  
 │    │    │          │                                          
 │    │    │    ┌─────▼─────┐                                    
 │    │    │    │           │                                    
 │    │    └────┤  update?  │                                    
 │    │         │           │                                    
 │    │         └─────┬─────┘                                    
 │    │               │                                          
 │    │              yes                                         
 │    │               │                                          
 │    │      ┌────────▼────────┐                                 
 │    │      │                 │                                 
 │    │      │  decode config  ◄─────────────┐                   
 │    │      │                 │             │                   
 │    │      └────────┬────────┘             │                   
 │    │               │                      │                   
 │    │      ┌────────▼────────┐             │                   
 │    │      │                 │             │                   
 │    │      │  merge configs  │             │                   
 │    │      │                 │       ┌─────┴─────┐  │libconfig 
 │    │      └────────┬────────┘       │           │  │nginx conf
 │    │               │                │  encoder  ├─►│json      
 │    │    ┌──────────▼──────────┐     │           │  │yaml      
 │    │    │                     │     └─────┬─────┘  │...       
 │    │    │  check  conditions  │           │                   
 │    │    │                     │           │                   
 │    │    └──────────┬──────────┘           │                   
 │    │               │                      │                   
 │    │     ┌─────────┴─────────┐            │                   
 │    │     │                   │            │                   
 │ ┌──┴─────▼────────┐ ┌────────▼─────────┐  │                   
 │ │                 │ │                  │  │                   
 │ │  fail acttions  │ │  encode  config  ◄──┘                   
 │ │                 │ │                  │                      
 │ └─────────────────┘ └────────┬─────────┘                      
 │                              │                                
┌┴───────────────────┐ ┌────────▼─────────┐                      
│                    │ │                  │                      
│  success acttions  ◄─┤  update  config  │                      
│                    │ │                  │                      
└────────────────────┘ └──────────────────┘                      
```

## Supported configuration formats

| Format      | Status |
|:---         |:---    |
| `yaml`      | ✅     |
| `json`      | ✅     |
| `libconfig` | ✅     |
| `nginx`     | ❌     |
| `hcl`       | ❌     |

## How to collaborate

We are open to external collaborations for this project: improvements, bugfixes, whatever.

For doing it, open an issue to discuss the need of the changes, then:

- Fork the repository
- Make your changes to the code
- Open a PR and wait for review
