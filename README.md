# Standalone artifact server based on existing act code

_Forked from https://github.com/nektos/act until creating a proper PR_

This provides the integrated artifact server as a standalone binary for use in
a container for example.

This enables the use of act's `--network` option while still providing support
for artifact uploads by creating a container running this standalone artifact
server.

Instead of `podman run --network podman ...` you can also use `docker run --network default ...`.

```
make standalone
mkdir /tmp/act-artifacts
podman run --network podman \
  -v /tmp/act-artifacts:/artifacts \
  -v ./dist/local/standalone-artifact-server:/standalone-artifact-server \
  --rm \
  -ti \
  docker.io/ubuntu:latest \
  /standalone-artifact-server --artifact-server-path /artifacts
```

Instead of typing `--artifact-server-path` you can also use the short flag `-d` (for directory).

When running the artifact server, it prints the necessary environment variables
to pass to act to access the artifact server. You can either pass them via the `--env` option or by creating an
environment file and using the `--env-file` option, e.g.:

```
act \
  --network podman \
  --env ACTIONS_RUNTIME_URL=http://192.168.178.138:34567/ \
  --env ACTIONS_RESULTS_URL=http://192.168.178.138:34567/ \
  --env ACTIONS_RUNTIME_TOKEN=<a-JWT...> \
  -W .github/workflows/system-tests.yml
# OR
# create act.env
act \
  --env-file act.env \
  --network podman \
  -W .github/workflows/system-tests.yml
```

---
Original README...

![act-logo](https://raw.githubusercontent.com/wiki/nektos/act/img/logo-150.png)

# Overview [![push](https://github.com/nektos/act/workflows/push/badge.svg?branch=master&event=push)](https://github.com/nektos/act/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/nektos/act)](https://goreportcard.com/report/github.com/nektos/act) [![awesome-runners](https://img.shields.io/badge/listed%20on-awesome--runners-blue.svg)](https://github.com/jonico/awesome-runners)

> "Think globally, `act` locally"

Run your [GitHub Actions](https://developer.github.com/actions/) locally! Why would you want to do this? Two reasons:

- **Fast Feedback** - Rather than having to commit/push every time you want to test out the changes you are making to your `.github/workflows/` files (or for any changes to embedded GitHub actions), you can use `act` to run the actions locally. The [environment variables](https://help.github.com/en/actions/configuring-and-managing-workflows/using-environment-variables#default-environment-variables) and [filesystem](https://help.github.com/en/actions/reference/virtual-environments-for-github-hosted-runners#filesystems-on-github-hosted-runners) are all configured to match what GitHub provides.
- **Local Task Runner** - I love [make](<https://en.wikipedia.org/wiki/Make_(software)>). However, I also hate repeating myself. With `act`, you can use the GitHub Actions defined in your `.github/workflows/` to replace your `Makefile`!

> [!TIP]
> **Now Manage and Run Act Directly From VS Code!**<br/>
> Check out the [GitHub Local Actions](https://sanjulaganepola.github.io/github-local-actions-docs/) Visual Studio Code extension which allows you to leverage the power of `act` to run and test workflows locally without leaving your editor.

# How Does It Work?

When you run `act` it reads in your GitHub Actions from `.github/workflows/` and determines the set of actions that need to be run. It uses the Docker API to either pull or build the necessary images, as defined in your workflow files and finally determines the execution path based on the dependencies that were defined. Once it has the execution path, it then uses the Docker API to run containers for each action based on the images prepared earlier. The [environment variables](https://help.github.com/en/actions/configuring-and-managing-workflows/using-environment-variables#default-environment-variables) and [filesystem](https://docs.github.com/en/actions/using-github-hosted-runners/about-github-hosted-runners#file-systems) are all configured to match what GitHub provides.

Let's see it in action with a [sample repo](https://github.com/cplee/github-actions-demo)!

![Demo](https://raw.githubusercontent.com/wiki/nektos/act/quickstart/act-quickstart-2.gif)

# Act User Guide

Please look at the [act user guide](https://nektosact.com) for more documentation.

# Support

Need help? Ask in [discussions](https://github.com/nektos/act/discussions)!

# Contributing

Want to contribute to act? Awesome! Check out the [contributing guidelines](CONTRIBUTING.md) to get involved.

## Manually building from source

- Install Go tools 1.20+ - (<https://golang.org/doc/install>)
- Clone this repo `git clone git@github.com:nektos/act.git`
- Run unit tests with `make test`
- Build and install: `make install`
