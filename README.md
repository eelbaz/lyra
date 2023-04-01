![Logo](https://github.com/eelbaz/lyra/blob/main/lyra-logo.png?raw=true)

# Lyra - A lightweight http Test Object Load and Performance Metrics Agent Integrated with InfluxDB

## Overview
A Simple, lightweight, concurrent agent for load simulation and monitoring http test object performance and availability -
Analyze and monitor performance metrics with configurable virtual users and InfluxDB integration for reporting & vizualization.


## Table of Contents
- [Overview](#overview)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

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
      "cdn": "example_cdn_test_object",
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
