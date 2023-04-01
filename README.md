
# Lyra - A lightweight HTTP Load and Performance Metrics Agent

## Overview
A lightweight, concurrent agent for load simulation and monitoring http test object performance and availability with InfluxDB integration.
Analyze and monitor performance metrics with configurable virtual users and InfluxDB integration.


## Table of Contents
- <link>Installation</link>
- Configuration
- Usage
- Contributing
- License

### Installation
To install the Lyra Performance Metrics Agent, simply clone the repository:
``` 
git clone https://github.com/user/repo.git 
```


### Configuration
Before running the agent, you need to configure it using a config.json file. Here is a sample configuration:


```
{
  "tag_prefix": "lyra",
  "num_users": 5,
  "debug": true,
  "use_influx_db": true,
  "influx_db_uri": "http://localhost:8086",
  "influx_db_api_key": "your_api_key",
  "influx_db_org": "your_org",
  "influx_db_bucket": "your_bucket",
  "influx_point_measurement_name": "performance_metrics",
  "resources": [
    {
      "url": "https://example.com",
      "cdn": "example_cdn",
      "workflow": "example_workflow"
    }
  ]
}
```



### Usage
To run the Lyra Performance Metrics Agent, execute the following command:

```
go run main.go
```

The agent will start sending performance metrics to InfluxDB according to the specified configuration.

### Contributing
Pull requests and bug reports are welcome. For major changes, please open an issue first to discuss what you would like to change.

### License
MIT
