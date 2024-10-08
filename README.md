# Rill

The HTTP/WS/Signaling Server for the [Rill](https://rill.one) Project. Rill is a platform that allows users to create and share interactive live streams.

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

> Once the configuration file is ready, you can execute any command listed below.

### Quick note on Secrets

Rill uses 1Password to store its secrets. The secrets are fetched using the `op` CLI tool. For more information, please refer to the [1Password CLI documentation](https://developer.1password.com/docs/cli).

You can also see how we make use of `1Password` referring to `.github/workflows/deploy.yaml`.

## Commands :hammer:

Rill uses a Makefile to manage its commands. You can run the following commands:

> [!Warning]
> The commands must be prefixed with `make`.

| Command | Arguments | Action |
|:--------------:|:---------------------------------------------:|:-----------------------------------------------:|
| `run`          | -                                             | Run the server                                  |
| `build`        | -                                             | Build the server into a binary                  |
| `dev`          | -                                             | Run the server with hot-reloading               |
| `test`         | -                                             | Run tests                                       |
| `docker_build` | [env(local,staging,production)\|pkl_version\] | Build a Docker image of the server              |
| `docker_run`   | inherited from `docker_build`                 | Build and run a Docker container of the server  |
| `gen_config`   | -                                             | Generate the types from the configuration file  |

## Project Structure :open_file_folder:

Rill tries to follow the [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) principles. The project is divided into the following directories:

```txt
.
├── cmd
│   └── rill
├── config
│   ├── local
│   ├── production
│   └── staging
├── internal
│   ├── app
│   ├── broadcasts
│   │   ├── application
│   │   ├── domain
│   │   └── infrastructure
│   │       ├── repository
│   │       └── ui
│   └── pkg
│       └── publicid
└── scripts
```

The three main directories are:

- `cmd`: Contains the entry point to the application.
- `internal`: Contains the application's business logic.
- `config`: Contains the configuration files.

Inside the `internal` directory, are define Rill's modules. Each module is divided into three subdirectories:

- `application`: Contains the use cases.
- `domain`: Contains the business logic.
- `infrastructure`: Contains the implementation details. Here, the repository and the UI are defined.

> [!Note]
> The `app` module is a special module that contains the main application logic.
> Here is defined the server definition and the application's configuration.
>
> The `pkg` module contains the public ID package. This package is used to generate unique public IDs.
