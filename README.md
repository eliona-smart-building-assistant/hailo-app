# Hailo App
The [Hailo app](https://github.com/eliona-smart-building-assistant/hailo-app) enables the [Hailo Digital Hub](https://www.hailodigitalhub.de/) for an [Eliona](https://www.eliona.io/) enviroment.

This app collects the data from configurable Hailo FDS endpoints. For each endpoint the app read the data for Hailo smart devices (smart waste stations and smart waste bins). Each device corresponds with an Eliona asset, which are created automatically, and writes various Eliona data for these assets.


## Configuration

The app needs environment variables and database tables for configuration. To edit the database tables the app provides an own API access.


### Registration in Eliona ###

To start and initialize an app in an Eliona environment, the app have to registered in Eliona. For this, an entry in the database table `public.eliona_app` is necessary.


### Environment variables ###

- `APPNAME`: must be set to `hailo`. Some resources use this name to identify the app inside an Eliona environment.

- `CONNECTION_STRING`: configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). Otherwise, the app can't be initialized and started. (e.g. `postgres://user:pass@localhost:5432/iot`)

- `API_ENDPOINT`:  configures the endpoint to access the [Eliona API v2](https://github.com/eliona-smart-building-assistant/eliona-api). Otherwise, the app can't be initialized and started. (e.g. `http://api-v2:3000/v2`)

- `API_TOKEN`: defines the secret to authenticate the app and access the API.  

- `API_SERVER_PORT`(optional): define the port the API server listens. The default value is Port `3000`.

- `DEBUG_LEVEL`(optional): defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/log). Not defined the default level is `info`.


### Database tables ###

The app requires configuration data that remains in the database. To do this, the app creates its own database schema `hailo` during initialization. To modify and handle the configuration data the Hailo app provides an API access. Have a look at the [API specification](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/hailo-app/develop/openapi.yaml) how the configuration tables should be used.

- `hailo.config`: contains Hailo FDS endpoints. Each row stands for one endpoint with configurable timeouts and polling intervals.

- `hailo.asset`: maps each Hailo smart device to an Eliona asset. For different Eliona projects different assets are used. The app collect and writes data separate for each configured project. The mapping is created automatically by the app.

**Generation**: to generate access method to database see Generation section below.

**Migration**: Versions of this app prior 2.0 use different mapping of assets and Hailo smart devices. So the mapping have to migrate manually with the following command. You must ensure that you have read permissions to the table `public.asset`.

    insert into hailo.asset
        (config_id, asset_id, proj_id, device_id)
    select
        app_id, asset.asset_id, proj_id, device_pkey
    from
        hailo.config, public.asset
    where
        asset_type like 'Hailo%'

## References

### Hailo app API ###

The Hailo app provides its own API to access configuration data and other functions. The full description of the API is defined in the `openapi.yaml` OpenAPI definition file.

- [API Reference](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/hailo-app/develop/openapi.yaml) shows Details of the API

**Generation**: to generate api server stub see Generation section below.

### Eliona assets ###

The app creates corresponding Eliona asset types and attribute sets during initialization. See [eliona/init.go](eliona/init.go) for details.

The hailo app writes data for each Hailo smart device. The data is structured into different subtypes of Eliona assets. See [eliona/heaps.go](eliona/heaps.go) for details. The following subtypes are defined:

- `Input`: Data like current volume in percent or count of openings for bins and stations (`stationDataPayload` and `binDataPayload`)
- `Status`: Statistic data like expected filling level at next service (`statusDataPayload`)
- `Info`: Static data which specifies a Hailo smart device like total volume and registration date (`deviceDataPayload`)


## Tools

### Test mode ###
The app can start in test mode to print out all information from a Hailo FDS endpoint. For this no further configuration (database or environment variables) are necessary. The endpoint is defined by command line arguments. After printing out the information the app quits.

```bash
./hailo --help
    Usage:
    -auth string
        Authentication endpoint for Hailo FDS (used for -t)
    -fds string
        FDS endpoint for for Hailo FDS (used for -t)
    -password string
        Password for Hailo FDS authentication endpoint (used for -t)
    -t    Test get data from Hailo FDS API without Eliona
    -user string
        Username for Hailo FDS authentication endpoint (used for -t)
./hailo -t --user username --password password --fds FDS-endpoint --auth auth-endpoint
```

### Generate API server stub ###

For the API server the [OpenAPI Generator](https://openapi-generator.tech/docs/generators/openapi-yaml) for go-server is used to generate a server stub. The easiest way to generate the server files is to use one of the predefined generation script which use the OpenAPI Generator Docker image.

```
.\generate-api-server.cmd # Windows
./generate-api-server.sh # Linux
```

### Generate Database access ###

For the database access [SQLBoiler](https://github.com/volatiletech/sqlboiler) is used. The easiest way to generate the database files is to use one of the predefined generation script which use the SQLBoiler implementation. Please note that the database connection in the `sqlboiler.toml` file have to be configured.

```
.\generate-db.cmd # Windows
./generate-db.sh # Linux
```

