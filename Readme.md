# Go Analysis Tool

## Overview
This Go Analysis Tool is a powerful and versatile application designed to analyze various aspects of your Go applications. It leverages popular monitoring tools such as DataDog, Prometheus, and Grafana to provide real-time insights into the performance, resource utilization, and behavior of your Go applications.

## Features
- **Real-time Monitoring:** Utilize DataDog, Prometheus, and Grafana to monitor your Go applications in real-time, gaining valuable insights into their performance and resource usage.
- **Customizable Dashboards:** Create custom dashboards in Grafana to visualize the metrics that matter most to you and your team.
- **Alerting:** Set up alerts in DataDog or Prometheus to receive notifications when critical thresholds are breached, ensuring proactive issue resolution.
- **Historical Data Analysis:** Analyze historical data to identify trends and patterns in your application's performance, helping you make data-driven decisions.
- **Easy Integration:** Seamlessly integrate this tool into your existing Go applications with minimal configuration.

## Runtime images
![Prometheus](image.png)
![Grafana](image-1.png)

## Installation
To get started, follow these steps:

1. Clone this repository to your local machine:
```bash
    git clone https://github.com/swarajkumarsingh/go-analysis-tools.git
```
2. Install the required dependencies:
```bash
    go mod tidy
```
3. Configure DataDog, Prometheus, and Grafana to collect and visualize metrics from your Go application.
4. Build and run the Go Analysis Tool:
```bash
    go build -o analysis-tool
    ./analysis-tool
```

## Configure
1. Run Prometheus
```bash
    docker run -d -p 9090:9090 -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```
2. Run Grafana
```bash
    docker run -d --name=grafana -p 3000:3000 grafana/grafana-enterprise
```

## Contributing
Contributions are welcome! If you have ideas for improvements or new features, please open an issue or submit a pull request. Make sure to follow our code of conduct.

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
Thanks to the Go community for their support and contributions.

## Contact
For questions or support, please contact sswaraj169@gmail.com.