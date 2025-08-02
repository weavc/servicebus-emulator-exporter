
## servicebus-emulator-exporter

Tooling to export existing azure services bus infrastructure into a Config.json file suitable for [Azure service bus emulator](https://learn.microsoft.com/en-us/azure/service-bus-messaging/overview-emulator).

Since its likely a one time use, local tool that was built quickly it only supports full connection strings currently.

## Usage

Install via `go install`:
```bash
go install github.com/weavc/servicebus-emulator-exporter@latest
```

```bash
servicebus-emulator-exporter --cs="Endpoint=sb://<namespace>;SharedAccessKeyName=<key name>;SharedAccessKey=<key>" > Config.json
```

Download binaries from releases:
```
wget https://github.com/weavc/servicebus-emulator-exporter/releases/download/0.0.2a/servicebus-emulator-exporter
chmod u+x servicebus-emulator-exporter
./servicebus-emulator-exporter --cs="Endpoint=sb://<namespace>;SharedAccessKeyName=<key name>;SharedAccessKey=<key>" > Config.json
```

## Limitations
See here for emulator limitations: https://learn.microsoft.com/en-us/azure/service-bus-messaging/overview-emulator#known-limitations

Basic caps on durations have been implemented however caps on things like namespaces/queues/topics/subscriptions/rules have all been ignored due to this being a very user specific concern.

