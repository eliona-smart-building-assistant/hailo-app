# Hailo App
The [Hailo app](https://github.com/eliona-smart-building-assistant/hailo-app) enables the [Hailo Digital Hub](https://www.hailodigitalhub.de/) for an [Eliona](https://www.eliona.io/) enviroment.

This app collects the data from configurable Hailo FDS endpoints. For each endpoint the app read the data for Hailo smart devices (smart waste stations and smart waste bins). Each device corresponds with an Eliona asset, which are created automatically, and writes various Eliona heap data for these assets.

## Configuration

The app needs environment variables and database tables for configuration.

### Environment variables ###

#### APPNAME

The `APPNAME` MUST be set to `hailo`. Some resources use this name to identify the app inside an Eliona environment. For running as a Docker container inside an Eliona environment, the `APPNAME` have to set in the [Dockerfile](Dockerfile). If the app runs outside an Eliona environment the `APPNAME` must be set explicitly.

```bash
export APPNAME=hailo # For running in Eliona environment set app name in Dockerfile
```

#### CONNECTION_STRING

The `CONNECTION_STRING` variable configures the [eliona database](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/db). If the app runs as a Docker container inside an Eliona environment it is ensured that the variable is already set by the environment. If you run the app standalone you must provide this variable. Otherwise the app can't be initialized and started.

```bash
export CONNECTION_STRING=postgres://user:pass@localhost::5432/iot # only if run standalone
```

#### DEBUG_LEVEL (optional)

The `DEBUG_LEVEL` variable defines the minimum level that should be [logged](https://github.com/eliona-smart-building-assistant/go-eliona/tree/main/log). Not defined the default level is `info`.

```bash
export LOG_LEVEL=debug # This is optionally, default is info
```

### Database tables ###

The app requires configuration data that remains in the database. To do this, the app creates its own database schema `hailo` during initialization. Some data in this schema should be made editable by Eliona frontend. This allows the app to be configured by the user without direct database access.

#### hailo.config

The database table `hailo.config` contains Hailo FDS endpoints. Each row stands for one endpoint with configurable timeouts and polling intervals. This table should be editable by the Eliona frontend. The table is filled during the initialization to demonstrate the configuration. This demo data contains no real endpoint and must be change for proper working.  

| Column          | Description                                                                        |
|-----------------|:-----------------------------------------------------------------------------------|
| app_id          | Id to identify the configured endpoint (automatically by sequence)                 |
| config          | JSON string to configure the endpoint (see below)                                  |
| enable          | Flag to enable or disable the endpoint                                             |
| description     | Description of the endpoint (optional)                                             |
| asset_id        | Id of an parent asset with groups all device assets (optional)                     |
| interval_sec    | Interval in seconds for collecting data from endpoint                              |
| auth_timeout    | Timeout in seconds for authentication server (default `5`s)                        |
| request_timeout | Timeout in seconds for FDS server (default `120`s)                                 |
| active          | Set to `true` by the app when running and to `false` when app is stopped           |
| proj_ids        | List of Eliona project ids for which this endpoint should collect data (see below) |

The `config` column have to contain a JSON to configure the Hailo FDS endpoint:

    {
      "username":    "Login for auth endpoint"
      "password":    "Password for auth endpoint"
      "auth_server": "Url to Hailo auth endpoint"
      "fds_server":  "Url to Hailo FDS endpoint"
    }

Insert in the `proj_ids` column all Eliona project ids for which the endpoint should create assets and collect data from all Hailo smart devices. For example, if the column contains `{1, 99}` the app does the following:

1. Collect specification and data for all devices
2. Create assets for each device in both projects `1` and `99`, if not already exists. Remember created asset ids in table `hailo.asset`.
3. Writes heap data for each device for the corresponding assets in both projects `1` and `99`

#### hailo.asset

The table `hailo.asset` maps each Hailo smart device to an Eliona asset. For different Eliona projects different assets are used. The app collect and writes data separate for each configured project (see column `proj_ids` above).

The mapping is created automatically by the app. So this table does not necessarily have to be editable via the frontend or displayed in the frontend.

If specification of a Hailo smart device is read (at the first time or if new devices are added later) the app looks if there is already a mapping. If so, the app uses the mapped asset id for writing the heap data. If not, the app creates a new Eliona asset or updates an existing one and inserts the mapping in this table for further use.

You can move or delete automatically created assets via the Eliona frontend. The app doesn't recreate this asset as long as the mapping is present in the table `hailo.asset`. In this case of deleted asset, the app just skips the writing of the heap data.

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

## API Reference

The hailo app writes heap data for each Hailo smart device. The data is structured into different subtypes of Eliona assets. See [eliona/heaps.go](eliona/heaps.go) for details. The following subtypes are defined:

- `Input`: Data like current volume in percent or count of openings for bins and stations ( `struct stationHeapData {}` and `struct binHeapData {}`)
- `Status`: Statistic data like expected filling level at next service (`struct statusHeapData {}`)
- `Info`: Static data which specifies a Hailo smart device like total volume and registration date (`struct deviceHeapData {}`)

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
