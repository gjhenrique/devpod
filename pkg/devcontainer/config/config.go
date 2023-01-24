package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
)

var (
	// ErrUnsupportedType is returned if the type is not implemented
	ErrUnsupportedType = errors.New("unsupported type")
)

type DevContainerConfig struct {
	// A name for the dev container which can be displayed to the user.
	Name string `json:"name,omitempty"`

	// Features to add to the dev container.
	Features map[string]interface{} `json:"features,omitempty"`

	// Array consisting of the Feature id (without the semantic version) of Features in the order the user wants them to be installed.
	OverrideFeatureInstallOrder []string `json:"overrideFeatureInstallOrder,omitempty"`

	// Ports that are forwarded from the container to the local machine. Can be an integer port number, or a string of the format "host:port_number".
	ForwardPorts []json.Number `json:"forwardPorts,omitempty"`

	// Set default properties that are applied when a specific port number is forwarded.
	PortsAttributes map[string]PortAttribute `json:"portAttributes,omitempty"`

	// Set default properties that are applied to all ports that don't get properties from the setting `remote.portsAttributes`.
	OtherPortsAttributes map[string]PortAttribute `json:"otherPortsAttributes,omitempty"`

	// Controls whether on Linux the container's user should be updated with the local user's UID and GID. On by default when opening from a local folder.
	UpdateRemoteUserUID bool `json:"updateRemoteUserUID,omitempty"`

	// Remote environment variables to set for processes spawned in the container including lifecycle scripts and any remote editor/IDE server process.
	RemoteEnv map[string]string `json:"remoteEnv,omitempty"`

	// The username to use for spawning processes in the container including lifecycle scripts and any remote editor/IDE server process. The default is the same user as the container.
	RemoteUser string `json:"remoteUser,omitempty"`

	// A command to run locally before anything else. This command is run before "onCreateCommand". If this is a single string, it will be run in a shell. If this is an array of strings, it will be run as a single command without shell.
	InitializeCommand StrArray `json:"initializeCommand,omitempty"`

	// A command to run when creating the container. This command is run after "initializeCommand" and before "updateContentCommand". If this is a single string, it will be run in a shell. If this is an array of strings, it will be run as a single command without shell.
	OnCreateCommand StrArray `json:"onCreateCommand,omitempty"`

	// A command to run when creating the container and rerun when the workspace content was updated while creating the container.
	// This command is run after "onCreateCommand" and before "postCreateCommand". If this is a single string, it will be run in a shell.
	// If this is an array of strings, it will be run as a single command without shell.
	UpdateContentCommand StrArray `json:"updateContentCommand,omitempty"`

	// A command to run after creating the container. This command is run after "updateContentCommand" and before "postStartCommand".
	// If this is a single string, it will be run in a shell. If this is an array of strings, it will be run as a single command without shell.
	PostCreateCommand StrArray `json:"postCreateCommand,omitempty"`

	// A command to run after starting the container. This command is run after "postCreateCommand" and before "postAttachCommand".
	// If this is a single string, it will be run in a shell. If this is an array of strings, it will be run as a single command without shell.
	PostStartCommand StrArray `json:"postStartCommand,omitempty"`

	// A command to run when attaching to the container. This command is run after "postStartCommand".
	// If this is a single string, it will be run in a shell. If this is an array of strings, it will be run as a single command without shell.
	PostAttachCommand StrArray `json:"postAttachCommand,omitempty"`

	// The user command to wait for before continuing execution in the background while the UI is starting up. The default is "updateContentCommand".
	WaitFor string `json:"waitFor,omitempty"`

	// User environment probe to run. The default is "loginInteractiveShell".
	UserEnvProbe string `json:"userEnvProbe,omitempty"`

	// Host hardware requirements.
	HostRequirements HostRequirements `json:"hostRequirements,omitempty"`

	// Tool-specific configuration. Each tool should use a JSON object subproperty with a unique name to group its customizations.
	Customizations map[string]interface{} `json:"customizations,omitempty"`

	// Action to take when the user disconnects from the container in their editor. The default is to stop the container.
	ShutdownAction string `json:"shutdownAction,omitempty"`

	// Whether to overwrite the command specified in the image. The default is true.
	OverrideCommand bool `json:"overrideCommand,omitempty"`

	// The path of the workspace folder inside the container.
	WorkspaceFolder string `json:"workspaceFolder,omitempty"`

	// DEPRECATED: Use 'customizations/vscode/settings' instead
	// Machine specific settings that should be copied into the container. These are only copied when connecting to the container for the first time, rebuilding the container then triggers it again.
	Settings map[string]interface{} `json:"settings,omitempty"`

	// DEPRECATED: Use 'customizations/vscode/extensions' instead
	// An array of extensions that should be installed into the container.
	Extensions []string `json:"extensions,omitempty"`

	// DEPRECATED: Use 'customizations/vscode/devPort' instead
	// The port VS Code can use to connect to its backend.
	DevPort int `json:"devPort,omitempty"`

	ImageContainer      `json:",inline"`
	NonComposeBase      `json:",inline"`
	ComposeContainer    `json:",inline"`
	DockerfileContainer `json:",inline"`

	// Origin is the origin from where this config was loaded
	Origin string `json:"-"`
}

type ComposeContainer struct {
	// The name of the docker-compose file(s) used to start the services.
	DockerComposeFile StrArray `json:"dockerComposeFile,omitempty"`

	// The service you want to work on. This is considered the primary container for your dev environment which your editor will connect to.
	Service string `json:"string,omitempty"`

	// An array of services that should be started and stopped.
	RunServices []string `json:"runServices,omitempty"`
}

type ImageContainer struct {
	// The docker image that will be used to create the container.
	Image string `json:"image,omitempty"`
}

type NonComposeBase struct {
	// Application ports that are exposed by the container. This can be a single port or an array of ports. Each port can be a number or a string.
	// A number is mapped to the same port on the host. A string is passed to Docker unchanged and can be used to map ports differently,
	// e.g. "8000:8010".
	AppPorts StrIntArray `json:"appPorts,omitempty"`

	// Container environment variables.
	ContainerEnv map[string]string `json:"containerEnv,omitempty"`

	// The user the container will be started with. The default is the user on the Docker image.
	ContainerUser string `json:"containerUser,omitempty"`

	// Mount points to set up when creating the container. See Docker's documentation for the --mount option for the supported syntax.
	Mounts []string `json:"mounts,omitempty"`

	// The arguments required when starting in the container.
	RunArgs []string `json:"runArgs,omitempty"`

	// The --mount parameter for docker run. The default is to mount the project folder at /workspaces/$project.
	WorkspaceMount string `json:"workspaceMount,omitempty"`
}

type DockerfileContainer struct {
	// The location of the Dockerfile that defines the contents of the container. The path is relative to the folder containing the `devcontainer.json` file.
	Dockerfile string `json:"dockerFile,omitempty"`

	// The location of the context folder for building the Docker image. The path is relative to the folder containing the `devcontainer.json` file.
	Context string `json:"context,omitempty"`

	// Docker build-related options.
	Build BuildOptions `json:"build,omitempty"`
}

type BuildOptions struct {
	// The location of the Dockerfile that defines the contents of the container. The path is relative to the folder containing the `devcontainer.json` file.
	Dockerfile string `json:"dockerfile,omitempty"`

	// The location of the context folder for building the Docker image. The path is relative to the folder containing the `devcontainer.json` file.
	Context string `json:"context,omitempty"`

	// Target stage in a multi-stage build.
	Target string `json:"target,omitempty"`

	// Build arguments.
	Args map[string]string `json:"args,omitempty"`

	// The image to consider as a cache. Use an array to specify multiple images.
	CacheFrom StrArray `json:"cacheFrom,omitempty"`
}

type HostRequirements struct {
	// Number of required CPUs.
	CPUs int `json:"cpus,omitempty"`

	// Amount of required RAM in bytes. Supports units tb, gb, mb and kb.
	Memory string `json:"memory,omitempty"`

	// Amount of required disk space in bytes. Supports units tb, gb, mb and kb.
	Storage string `json:"storage,omitempty"`
}

type PortAttribute struct {
	// Defines the action that occurs when the port is discovered for automatic forwarding
	// default=notify
	OnAutoForward string `json:"onAutoForward,omitempty"`

	// Automatically prompt for elevation (if needed) when this port is forwarded. Elevate is required if the local port is a privileged port.
	ElevateIfNeeded bool `json:"elevateIfNeeded,omitempty"`

	// Label that will be shown in the UI for this port.
	// default=Application
	Label string `json:"label,omitempty"`

	// When true, a modal dialog will show if the chosen local port isn't used for forwarding.
	RequireLocalPort bool `json:"requireLocalPort,omitempty"`

	// The protocol to use when forwarding this port.
	Protocol string `json:"protocol,omitempty"`
}

// StrIntArray string array to be used on JSON UnmarshalJSON
type StrIntArray []string

// UnmarshalJSON convert JSON object array of string or
// a string format strings to a golang string array
func (sa *StrIntArray) UnmarshalJSON(data []byte) error {
	var jsonObj interface{}
	err := json.Unmarshal(data, &jsonObj)
	if err != nil {
		return err
	}
	switch obj := jsonObj.(type) {
	case string:
		*sa = StrIntArray([]string{obj})
		return nil
	case int:
		*sa = StrIntArray([]string{strconv.Itoa(obj)})
		return nil
	case []interface{}:
		s := make([]string, 0, len(obj))
		for _, v := range obj {
			value, ok := v.(string)
			if !ok {
				return ErrUnsupportedType
			}
			s = append(s, value)
		}
		*sa = StrIntArray(s)
		return nil
	}
	return ErrUnsupportedType
}

// StrArray string array to be used on JSON UnmarshalJSON
type StrArray []string

// UnmarshalJSON convert JSON object array of string or
// a string format strings to a golang string array
func (sa *StrArray) UnmarshalJSON(data []byte) error {
	var jsonObj interface{}
	err := json.Unmarshal(data, &jsonObj)
	if err != nil {
		return err
	}
	switch obj := jsonObj.(type) {
	case string:
		*sa = StrArray([]string{obj})
		return nil
	case []interface{}:
		s := make([]string, 0, len(obj))
		for _, v := range obj {
			value, ok := v.(string)
			if !ok {
				return ErrUnsupportedType
			}
			s = append(s, value)
		}
		*sa = StrArray(s)
		return nil
	}
	return ErrUnsupportedType
}
