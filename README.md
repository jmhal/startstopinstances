# startstopinstances

A small tool for wake up my development instances and stop them.

```bash
$ ./startstopinstances (start|stop|status)
```

1. Start: reads the instances IDs from ~/.aws/instances and starts them.
2. Stop: reads the instances IDs from ~/.aws/instances and stops them.
3. Status: shows the state and assigned public IPs.

~/.aws/instances holds the IDs you want to control. One ID per line. If you want to control a different region, configure ~/.aws/config properly.