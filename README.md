# Hailo App
The [Hailo app](https://github.com/eliona-smart-building-assistant/hailo-app) enables the [Hailo Digital Hub](https://www.hailodigitalhub.de/) for an [Eliona](https://www.eliona.io/) enviroment.

This app collects the data from configurable Hailo FDS endpoints. For each endpoint the app read the data for Hailo smart devices (smart waste stations and smart waste bins). Each device corresponds with an Eliona asset, which are created automatically, and writes various Eliona data for these assets.

## Configuration

The app needs environment variables and database tables for configuration. To edit the database tables the app provides an own API access.

### Eliona ###

To start and initialize an app in an Eliona environment, the app have to registered in Eliona. For this, an entry in the database table `public.eliona_app` is necessary.

### Environment variables ###

#### APPNAME

The `APPNAME` MUST be set to `hailo`. Some resources use this name to identify the app inside an Eliona environment. For running as a Docker container inside an Eliona environment, the `APPNAME` have to set in the [Dockerfile](Dockerfile). If the app runs outside an Eliona environment the `APPNAME` must be set explicitly.

```bash
export APPNAME=hailo
```

#### CONNECTION_STRING

The `CONNECTION_STRING` variable configures the [Eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). Otherwise, the app can't be initialized and started.

```bash
export CONNECTION_STRING=postgres://user:pass@localhost:5432/iot
```

#### API_ENDPOINT

The `API_ENDPOINT` variable configures the endpoint to access the [Eliona API v2](https://github.com/eliona-smart-building-assistant/eliona-api). Otherwise, the app can't be initialized and started.

```bash
export API_ENDPOINT=http://api-v2:80/v2
```

#### API_SERVER_PORT (optional)

You can set `API_SERVER_PORT` to define the port the API server listens. The default value is Port `80`.


```bash
export API_SERVER_PORT=8082 # optionally, default is '80'
```

#### DEBUG_LEVEL (optional)

The `DEBUG_LEVEL` variable defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/log). Not defined the default level is `info`.

```bash
export LOG_LEVEL=debug # optionally, default is 'info'
```

### Database tables ###

The app requires configuration data that remains in the database. To do this, the app creates its own database schema `hailo` during initialization. Some data in this schema should be made editable by Eliona frontend. This allows the app to be configured by the user without direct database access.

To modify and handle the configuration data the Hailo app provides an API access. Have a look at the [API specification](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/hailo-app/develop/openapi.yaml) ([openapi.yaml](https://raw.githubusercontent.com/eliona-smart-building-assistant/hailo-app/develop/openapi.yaml)) how the configuration tables should be used. The API is available at `http://servername:80/v1`. The port can be configured with the `API_SERVER_PORT` environment variable.

#### hailo.config

The database table `hailo.config` contains Hailo FDS endpoints. Each row stands for one endpoint with configurable timeouts and polling intervals. This table should be editable by the Eliona frontend. The table is filled during the initialization to demonstrate the configuration. This demo data contains no real endpoint and must be change for proper working.

For a detailed description have a look at the [API specification](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/hailo-app/develop/openapi.yaml) ([openapi.yaml](https://raw.githubusercontent.com/eliona-smart-building-assistant/hailo-app/develop/openapi.yaml)).

Insert in the `proj_ids` column all Eliona project ids for which the endpoint should create assets and collect data from all Hailo smart devices. For example, if the column contains `{1, 99}` the app does the following:

1. Collect specification and data for all devices
2. Create assets for each device in both projects `1` and `99`, if not already exists. Remember created asset ids in table `hailo.asset`.
3. Writes data for each device for the corresponding assets in both projects `1` and `99`

#### hailo.asset

The table `hailo.asset` maps each Hailo smart device to an Eliona asset. For different Eliona projects different assets are used. The app collect and writes data separate for each configured project (see column `proj_ids` above).  The mapping is created automatically by the app. So this table does not necessarily have to be editable via the frontend or displayed in the frontend.

For a detailed description have a look at the [API specification](https://eliona-smart-building-assistant.github.io/open-api-docs/?https://raw.githubusercontent.com/eliona-smart-building-assistant/hailo-app/develop/openapi.yaml) ([openapi.yaml](https://raw.githubusercontent.com/eliona-smart-building-assistant/hailo-app/develop/openapi.yaml)).

If specification of a Hailo smart device is read (at the first time or if new devices are added later) the app looks if there is already a mapping. If so, the app uses the mapped asset id for writing the data. If not, the app creates a new Eliona asset or updates an existing one and inserts the mapping in this table for further use.

You can move or delete automatically created assets via the Eliona frontend. The app doesn't recreate this asset as long as the mapping is present in the table `hailo.asset`. In this case of deleted asset, the app just skips the writing of the data.

### Migrate configuration prior version 2.0

Versions of this app prior 2.0 use different mapping of assets and Hailo smart devices. So the mapping have to migrate manually with the following command. You must ensure that you have read permissions to the table `public.asset`.

    insert into hailo.asset
        (config_id, asset_id, proj_id, device_id)
    select
        app_id, asset.asset_id, proj_id, device_pkey
    from
        hailo.config, public.asset
    where
        asset_type like 'Hailo%'

## Generation and Implementation ##

The Hailo app uses generators to create API server capabilities and database access.

For the API server the [OpenAPI Generator](https://openapi-generator.tech/docs/generators/openapi-yaml) for go-server is used. The easiest way to generate the server files is to use one of the predefined generation script which use the OpenAPI Generator Docker image.

```
.\generate-api-server.cmd # Windows
./generate-api-server.sh # Linux
```

For the database access [SQLBoiler](https://github.com/volatiletech/sqlboiler) is used. The easiest way to generate the database files is to use one of the predefined generation script which use the SQLBoiler implementation. Please note that the database connection in the `sqlboiler.toml` file have to be configured.

```
.\generate-db.cmd # Windows
./generate-db.sh # Linux
```

## API Reference

The hailo app writes data for each Hailo smart device. The data is structured into different subtypes of Eliona assets. See [eliona/heaps.go](eliona/heaps.go) for details. The following subtypes are defined:

- `Input`: Data like current volume in percent or count of openings for bins and stations ( `struct stationDataPayload {}` and `struct binDataPayload {}`)
- `Status`: Statistic data like expected filling level at next service (`struct statusDataPayload {}`)
- `Info`: Static data which specifies a Hailo smart device like total volume and registration date (`struct deviceDataPayload {}`)

The app creates corresponding Eliona asset types and attribute sets during initialization. See [eliona/init.go](eliona/init.go) for details.

## Tools

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
