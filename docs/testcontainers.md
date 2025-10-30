# Testcontainers

TRex uses https://github.com/testcontainers/testcontainers-go/ for integration tests to spin up ephemeral containers for tests.

The containers used by the tests are initialized/destroyed in the  `integration_testing` environment.


## Compatibility with podman

testcontainers project only supports Docker officially and some errors can appear with podman.

If you encounter the following error:

```
Failed to start PostgreSQL testcontainer: create container: container create: Error response from daemon: container create: unable to find network with name or ID bridge: network not found: creating reaper failed
```
It can happen because testcontainers spin up an additional [testcontainers/ryuk](https://github.com/testcontainers/moby-ryuk) container that manages the lifecycle of the containers used in the tests and performs cleanup in case there are fails.


One way to bypass this problem is not to use ryuk setting the environment variable

```bash
TESTCONTAINERS_RYUK_DISABLED=true
```

Or setting a property in `~/.testcontainers.properties`

```
ryuk.disabled=true
```

Ryuk needs to execute with root permissions in the podman machine to manage other containers. This [issue](https://github.com/testcontainers/testcontainers-go/issues/2781#issuecomment-2619626043) in testcontainer's repository offers an alternative solution. Be mindful of the elevated permissions required:

```bash
# verify socket path inside podman machine
$ podman machine ssh
Connecting to vm podman-machine-default. To close connection, use `~.` or `exit`
Fedora CoreOS 40.20240808.2.0

root@localhost:~# ls -al /var/run/podman/podman.sock 
srw-rw----. 1 root root 0 Dec 20 14:32 /var/run/podman/podman.sock
exit

# On the host machine
$ sudo mkdir /var/run/podman
$ sudo ln -s /Users/Your.User/.local/share/containers/podman/machine/podman.sock /var/run/podman/podman.sock

export DOCKER_HOST="unix:///var/run/podman/podman.sock"
export TESTCONTAINERS_RYUK_CONTAINER_PRIVILEGED=true

# if it still fails, give permissions to /var/run/podman/podman.sock within the podman machine
$ sudo chmod a+xrw /var/run/podman
$ sudo chmod a+xrw /var/run/podman/podman.sock
```



