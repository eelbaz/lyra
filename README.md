![Logo](https://github.com/eelbaz/lyra/blob/main/lyra-logo.png?raw=true)


# Lyra - A lightweight http Test Object Load and Performance Metrics Agent Integrated with InfluxDB


## Overview
A Simple, lightweight, concurrent agent for load simulation and monitoring http test object performance and availability -
Analyze and monitor performance metrics with configurable virtual users and InfluxDB integration for reporting & vizualization.

## ![Logo](https://github.com/eelbaz/lyra/blob/main/stats-dashboard.png)

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


### Metrics Documentation 

Lyra - is a homegrown application I cobbled together for load testing and performance measurements on multiple resources (URLs) specified in a configuration file (config.json). The application measures various metrics that contribute to an entire HTTP request for a given testObject, such as DNS lookup time, TCP connection time, TLS handshake time, server processing time, content transfer time, and total time.

The results are either displayed in the console or written to InfluxDB, a time series database deployable locally or in the cloud, depending on the value of the "use_influx_db" configuration parameter. The application uses multiple goroutines to simulate a configurable number of virtual users and launches these goroutines simultaneously. Overall, the application is an organized attempt to provide a holistic synthetic view of granular client performance metrics.

1. **DNS Lookup**: This metric measures the time it takes to resolve the hostname of the resource's URL to an IP address.
   a. DNS Lookup: This metric can be compared to the nslookup or dig command, which are used to perform DNS lookups.
   b. DNS Lookup: This metric is calculated as the difference between the time when the DNS lookup started and the time when the lookup completed. The calculation is performed using the `time.Since` function, which returns the time elapsed since a given start time.

2. **TCP Connection**: This metric measures the time it takes to establish a TCP connection to the resource's server.
   a. TCP Connection: This metric can be compared to the telnet or nc command, which are used to establish TCP connections to servers.
   b. This metric is calculated as the difference between the time when the TCP connection started and the time when the connection completed. The calculation is performed using the `time.Since` function.

3. **TLS Handshake**: This metric measures the time it takes to negotiate a secure connection using the TLS protocol.
   a. TLS Handshake: This metric can be compared to the openssl command, which is used to perform SSL/TLS handshakes.
   b. TLS Handshake: This metric is calculated as the difference between the time when the TLS handshake started and the time when the handshake completed. The calculation is performed using the `time.Since` function.

4. **Server Processing**: This metric measures the time it takes for the server to process the HTTP request and generate a response.
   a. Server Processing: This metric can be compared to the curl or wget command, which are used to make HTTP requests to servers.
   b. Server Processing: This metric is calculated as the difference between the time when the server started processing the request and the time when the response was received. The calculation is performed using the `time.Since` function.

5. **Content Transfer**: This metric measures the time it takes to transfer the response body from the server to the client.
   a. Content Transfer: This metric can be compared to the pv or pipebench command, which are used to measure the speed of data transfer.
   b. Content Transfer: This metric is calculated as the difference between the time when the content transfer started and the time when the transfer completed. The calculation is performed using the `time.Since` function.

6. **Total**: This metric measures the total time it takes to complete all the above steps for a single resource.
   a. Total: This metric is a combination of all the above metrics and can be compared to the time command, which is used to measure the execution time of a command.
   b. Total: This metric is calculated as the difference between the time when the resource measurement started and the time when the measurement completed. The calculation is performed using the `time.Since` function.

7. **Availability: This metric represents the HTTP status code of the response, indicating the availability of the resource. For example, a status code of 200 means the resource is available and a status code of 404 means the resource was not found.
a. Availability: This metric can be compared to the curl or wget command, which can be used to check the HTTP status code of a resource.



### Contributing
Pull requests and bug reports are welcome. For major changes, please open an issue first to discuss what you would like to change.

### License
MIT
