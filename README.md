
## servicebus-emulator-exporter

Tooling to export existing azure service bus infrastructure into a Config.json file suitable for [Azure service bus emulator](https://learn.microsoft.com/en-us/azure/service-bus-messaging/overview-emulator).

### Features
- Connect to multiple namespaces using `--cs="<connection string>"`
- Group all namespaces together using `--group="<group name>"`
- Filter queue and topic entities using regex patterns `--filter="inventory"`

### Service Bus Emulator Limitations

See here for the emulators limitations: https://learn.microsoft.com/en-us/azure/service-bus-messaging/overview-emulator#known-limitations

**Duration Limits**: Hard limits on durations have been implemented, if you find any more caps that haven't been implemented please create an issue or PR.

**Namespace Limits**: Due to there being a hard limit of 1 namespace, the `--group="<namespace name>"` parameter has been added. If provided this will merge multiple namespaces into 1.

**Queue/Topic Limits**: There is a hard limit on the emulator of 50 queues/topics. The user can filter entities by name by passing `--filter="<filter>"` parameters. While not a hard cap, it should help users filter down the queues and topics to a more managable data set. 

**Others**: No protections are provided against other caps. Most of the others seemed reasonable enough to work with most configurations, however Im sure there will be edge cases.

## Installation

Install via `go install`:
```bash
go install github.com/weavc/servicebus-emulator-exporter@latest
```

Download binaries from releases:
```
wget https://github.com/weavc/servicebus-emulator-exporter/releases/latest/download/servicebus-emulator-exporter
chmod u+x servicebus-emulator-exporter
```

## Usage

### Basic:
```
servicebus-emulator-exporter \
  --cs="Endpoint=sb://<namespace>;SharedAccessKeyName=<key name>;SharedAccessKey=<key>" \
  > Config.json
```

Pull queues and topics from 1 namespace in azure and output that to Config.json.

### Advanced:
```bash
servicebus-emulator-exporter \
  --cs="Endpoint=sb://<namespace>;SharedAccessKeyName=<key name>;SharedAccessKey=<key>" \
  --cs="Endpoint=sb://<other namespace>;SharedAccessKeyName=<other key name>;SharedAccessKey=<other key>" \
  --filter="inventory" \
  --group="e2e-testing" \
  > Config.json
```

Over 2 namespaces in Azure, find all queues and topics with "inventory" in their names, add them all to a namespace called "e2e-testing" and output that to Config.json.
