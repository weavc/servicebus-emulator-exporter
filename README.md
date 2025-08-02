
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

I have added some protection against known limitation in regards to messaging entity properties (Max TTL, duplicate duration etc), however I do not cap things like namespace/queue/topic limits. These should be filtered out by the user after generating the config file.

