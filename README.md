# KNX Prometheus Exporter

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fchr-fritz%2Fknx-exporter.svg?type=small)](https://app.fossa.com/projects/git%2Bgithub.com%2Fchr-fritz%2Fknx-exporter?ref=badge_small)
![Go build](https://github.com/chr-fritz/knx-exporter/workflows/Go%20build/badge.svg)

The KNX Prometheus Exporter is a small bridge to export values measured by KNX sensors to
Prometheus. It takes the values either from cyclic sent `GroupValueWrite` telegrams and can request
values itself using `GroupValueRead` telegrams.

[TOC]: # "## Table of Contents"

## Table of Contents
- [Usage](#usage)
  - [Converting the ETS 5 Group Export to a configuration](#converting-the-ets-5-group-export-to-a-configuration)
  - [Preparing the configuration](#preparing-the-configuration)
  - [Running the exporter](#running-the-exporter)
- [Exported metrics](#exported-metrics)
- [Health Check Endpoints](#health-check-endpoints)
- [Contributing](#contributing)
- [License](#license)
- [Maintainer](#maintainer)


## Usage

### Converting the ETS 5 Group Export to a configuration

The KNX Prometheus Exporter will only export the values from configured group addresses. A good
starting point is to
[export the group addresses](https://support.knx.org/hc/en-us/articles/115001825324-Group-Address-Export)
from ETS 5 into the XML format and convert them.

Please refer the KNX Documentation
"[Group Address & Export](https://support.knx.org/hc/en-us/articles/115001825324-Group-Address-Export)"
for more information about how to export the group addresses to XML.

With the exported group addresses you can call the `knx-exporter` and convert them:

```shell script
knx-exporter convertGA [SOURCE] [TARGET]
```

You must replace `[SOURCE]` with the path to your group address export file. `[TARGET]` is the path
where the converted configuration should be stored.

### Preparing the configuration

Converting the group addresses is a good starting point for preparing the actual configuration. The
simplest fully configuration will look like this:

```yaml
Connection:
  Type: "Tunnel"
  Endpoint: "192.168.1.15:3671"
  PhysicalAddress: 2.0.1
MetricsPrefix: knx_
AddressConfigs:
  0/0/1:
    Name: dummy_metric
    DPT: 1.*
    Export: true
    MetricType: "counter"
    ReadActive: true
    MaxAge: 10m
    Comment: dummy comment
```

In the next sections we will describe the meanings for every statement.

#### The `Connection` section

The `Connection` section contains all settings about how to connect to your KNX system and how the
KNX Prometheus Exporter will identify itself within it. It has three properties:

- `Type` This defines the connection type to your KNX system. It can be either `Router` or `Tunnel`.
- `Endpoint` This defines the ip address or hostname including the port to where the KNX Prometheus
  Exporter should open the connection. In case of you are using `Router` in `Type` the default might
  be `224.0.23.12:3671`.
- `PhysicalAddress` This defines the physical address how the KNX Prometheus Exporter will identify
  itself within your KNX address.

#### The `MetricsPrefix`

The `MetricsPrefix` defines a single string that will be added to all your exported metrics. The
format must be compliant with the
[prometheus metrics names](https://prometheus.io/docs/practices/naming/#metric-names).

#### The `AddressConfigs` section

The `AddressConfigs` section defines all the information about the group addresses which should be
exported to prometheus. It contains the following structure for every exported group address.

```yaml
  0/0/1:
    Name: dummy_metric
    DPT: 1.*
    Export: true
    MetricType: "counter"
    ReadActive: true
    MaxAge: 10m
    Comment: dummy comment
```

The first line `0/0/1` defines the group address which should be exported. It can be written in all
three valid forms:
- Two layers: `0/1`
- Three layers: `0/0/1`
- Free structure `123`

You can find more information about the formatting of group addresses here:
[cemi@v0.0 NewGroupAddrString](https://pkg.go.dev/github.com/vapourismo/knx-go/knx/cemi@v0.0.0-20201122213738-75fe09ace330#NewGroupAddrString)

Next it defines the actual information for a single group address:

- `Name` is the short name of the exported metric. The format must be compliant with the
  [prometheus metrics names](https://prometheus.io/docs/practices/naming/#metric-names).
- `DPT` The knx data point type. This defines how the KNX Prometheus Exporter will interpret the
  payload of received telegrams. All available data point types can be found here:
  [knx dpt](https://pkg.go.dev/github.com/vapourismo/knx-go@v0.0.0-20201122213738-75fe09ace330/knx/dpt)
- `Export` Can be either `true` or `false`. Allows to disable exporting the group address as value.
- `MetricType` defines the type of the exported metric. Can be either `counter` or `gauge`. See
  [Prometheus documentation counter vs. gauge](https://prometheus.io/docs/practices/instrumentation/#counter-vs-gauge-summary-vs-histogram)
  for more information about it.
- `ReadActive` can either be `true` or `false`. If set to `true` the KNX Prometheus Exporter will
  send a `GroupValueRead` telegram to the group address to active ask for a new value if the last
  received value is older than `MaxAge`.
- `MaxAge` defines the maximum age of a value until the KNX Prometheus Exporter will send a
  `GroupValueRead` telegram to active request a new value for the group address. This setting will
  be ignored if `ReadActive` is set to `false`.
- `Comment` a short comment for the group address. Will be also exported as comment within the
  prometheus metrics.

### Running the exporter

To run the metrics export just run the following command:

```shell script
knx-exporter run -f [CONFIG-FIlE]
```

You must replace `[CONFIG-FILE]` with the path to your configuration file that you prepared in the
previous step. After starting the exporter you can open
[`http://localhost:8080/metrics`](http://localhost:8080/metrics) to view the exported metrics.

## Exported metrics

Beside exported metrics from KNX group addresses it exports some additional metrics. This metrics
can be grouped into three groups:

1. **KNX Message Metrics:** Three counter metrics which counts the number of
   - received messages `knx_messages{direction="received",processed="false"}`
   - processed received messages `knx_messages{direction="received",processed="true"}` and
   - sent messages `knx_messages{direction="sent",processed="true"}`.
2. **HTTP Metrics:** Counts the processed number of successfully and failed http requests. All
   metrics starts with `promhttp_`.
3. **GoLang Metrics:** These are metrics which indicates some health information about memory, cpu
   usage, threads and other GoLang internal information.

## Health Check Endpoints

The KNX Prometheus Exporter provides endpoints for liveness and readiness checks like they were
recommended in enviroments like Kubernetes. They are reachable under `/live` and `/ready`. With
`/live?full=1` and `/ready?full=1` they also print the status of every single check.

## Contributing

1. Fork it
2. Download your fork to your PC (`git clone https://github.com/your_username/knx-exporter && cd
   knx-exporter`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Make changes and add them (`git add .`)
5. Commit your changes (`git commit -m 'Add some feature'`)
6. Push to the branch (`git push origin my-new-feature`)
7. Create new pull request

## License

The KNX Exporter is released under the Apache 2.0 license. See
[LICENSE](https://github.com/chr-fritz/knx-exporter/blob/master/LICENSE)

## Maintainer

Christian Fritz (@chrfritz)
