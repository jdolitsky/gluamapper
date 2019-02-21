# goluamapper: Lua to YAML example


The example code in [main.go](./main.go) shows how a global variable defined in a Lua script can be converted into YAML.

Contents of [porter.lua](./porter.lua):
```lua
local name = "cloud-communicator"
local version = "0.1.0"

local registryHost = "r.mysite.io"
local registryRepo = "trogdor-burninator/" .. name

local command = "/usr/local/bin/talk-to-the-cloud"
local provider = "aws"
local installArgs = {"--setup", "--provider", provider, "-f", provider .. "-setup.conf"}
local uninstallArgs = {"--teardown", "--provider", provider, "-f", provider .. "-teardown.conf"}

bundle = {
    name =  name,
    version = version,
    description = "this thing talks to the cloud. no joke.",
    invocationImage = registryHost .. "/" .. registryRepo .. ":" .. version,
    mixins = {"exec"},
    install = {
        {
            description = "Install " .. name,
            exec = {command = command, arguments = installArgs}
        }
    },
    uninstall = {
        {
            description = "Uninstall " .. name,
            exec = {command = command, arguments = uninstallArgs}
        }
    }
}
```

To run, clone the repo, `cd` to this directory and run:

```
go run . porter.lua
```

Resulting YAML ([porter.yaml](./porter.yaml)):
```
name: cloud-communicator
version: 0.1.0
description: this thing talks to the cloud. no joke.
invocationimage: r.mysite.io/trogdor-burninator/cloud-communicator:0.1.0
mixins:
- exec
install:
- description: Install cloud-communicator
  exec:
    command: /usr/local/bin/talk-to-the-cloud
    arguments:
    - --setup
    - --provider
    - aws
    - -f
    - aws-setup.conf
uninstall:
- description: Uninstall cloud-communicator
  exec:
    command: /usr/local/bin/talk-to-the-cloud
    arguments:
    - --teardown
    - --provider
    - aws
    - -f
    - aws-teardown.con
```
