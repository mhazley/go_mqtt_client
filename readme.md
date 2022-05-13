# mqttClient: a Go module

This is a go module to facilitate the use of MQTT clients by 
1. applying defaults
2. handling errors
3. printing helpful logs
4. asynchronously publishing payloads

All of this allows for neater and non-repetitive code.

## Customizing your broker

If you want to add or modify the broker options, you can access and modify `MqttClient.Options` before calling 
`StartMqttClient`. Example;

```
client := MqttClientInit(myHandlerFunction, mySubTopics, broker)
client.Options.SetClientID("myId")
client.StartMqttClient()
```

## Tests

To run the tests simply execute `go test` from the root directory. Pass the `-v` flag to see all executed tests, not 
just the failed ones.

To see `info` and `debug` logs pass the `--debug` flag.

To see logs in a human friendly way pass the `--human` flag.
