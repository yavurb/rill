# Rill

The HTTP/WS/Signaling Server for the [Rill](https://rill.one) Project.

## Configuration :wrench:

Rill uses [Pkl](https://pkl-lang.org) as its configuration language. The primary configuration template file is located at `config/Config.pkl`. This file should be used to amend the final configuration file.

> [!Note]
> The final config file should be placed in one of the following locations based on the environment:
>
> - `config/local/config.pkl`
> - `config/staging/config.pkl`
> - `config/production/config.pkl`

Example configuration:

```pkl
// config/local/config.pkl

amends "../Config.pkl"

host = read?("env:HOST") ?? "0.0.0.0"
port = read?("env:PORT")?.toInt() ?? 8910
cors = new Cors {
  allowOrigins = new Listing { "http://localhost:4321" }
  allowMethods = new Listing { "GET" "POST" "PUT" "DELETE" }
}
webRTC = new WebRTC {
  iceServers = new Listing {
    new ICEServer {
      urls = new Listing {
        "stun:stun.l.google.com:19302"
        "stun:stun1.l.google.com:19302"
        "stun:stun2.l.google.com:19302"
      }
    }
  }
}

logLevel = read?("env:LOG_LEVEL")?.trim()?.toLowerCase() ?? "debug"
```

[^Once the configuration file is ready, you can execute any command listed below.]

## Commands :hammer:

Rill uses a Makefile to manage its commands. You can run the following commands:

> [!Note]
> The commands must be prefixed with `make`.

| Command | Arguments | Action |
|:--------------:|:---------------------------------------------:|:-----------------------------------------------:|
| `run`          | -                                             | Run the server                                  |
| `build`        | -                                             | Build the server                                |
| `dev`          | -                                             | Run the server with hot-reloading               |
| `test`         | -                                             | Run tests                                       |
| `docker_build` | [env(local,staging,production)\|pkl_version\] | Build a Docker image                            |
| `docker_run`   | inherited from `docker_build`                 | Run a Docker container                          |
| `gen_config`   | -                                             | Generate the types from the configuration file  |



