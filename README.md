# HP throttling Middleware

[Homing pigeon](https://github.com/softonic/homing-pigeon) middleware that throttles messages to limit writes.

### Overview

This middleware is a simple implementation of a throttling middleware. It will limit the number of messages per second that can be sent to the next middleware, allowing for
burst usage.

### How it works

The middleware uses the `rate` library of Go to implement throttling and burst functionality. This means that you can configure the number of messages allowed per second, as well as the ability to
handle occasional bursts of traffic. The configuration is done according to the documentation provided by the [`rate` library](https://pkg.go.dev/golang.org/x/time/rate).

- **Limit**: This defines the rate at which events are allowed. It's expressed in events per second. If you set a `Limit` of 5, you are allowing 5 events to occur per second. If the value is set to 0,
  throttling is disabled.

- **Burst**: This is the maximum number of events that are allowed to burst beyond the `Limit` in a short period. Even if the rate is controlled by the `Limit`, a burst allows occasional bursts of
  event traffic greater than the normal rate.

### Usage

If you are using the [homing-pigeon-chart](https://github.com/softonic/homing-pigeon-chart/) Helm chart, you can enable the throttling middleware by setting the following values in the `values.yaml` file:

```yaml
requestMiddlewares:
 - name: throttling
   repository: softonic/hp-throttling
   tag: latest
   pullPolicy: IfNotPresent
   env:
   - name: THROTTLE_LIMIT
     value: "100"
   - name: THROTTLE_BURST
     value: "100"
```
