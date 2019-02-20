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