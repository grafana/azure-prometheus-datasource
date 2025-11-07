## Compatibility

Azure Monitor Managed Service for Prometheus is not compatible with Grafana<11.5.0 If you are running version 11.4.x and lower, please use plugin the core Prometheus plugin instead.

# Azure Monitor Managed Service for Prometheus Data Source

Azure Monitor Managed Service for Prometheus is a Prometheus-compatible service that monitors and provides alerts on containerized applications and infrastructure at scale.

Read more about it here: [Azure Monitor and Prometheus](https://learn.microsoft.com/en-us/azure/azure-monitor/metrics/prometheus-metrics-overview)

## Add the data source

To configure Azure authentication see [Configure Azure Active Directory (AD) authentication](https://grafana.com/docs/grafana/latest/datasources/azure-monitor/#configure-azure-active-directory-ad-authentication).

In Grafana Enterprise, update the .ini configuration file: [Configure Grafana](https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana/). Depending on your setup, the .ini file is located [here](https://grafana.com/docs/grafana/latest/setup-grafana/configure-grafana/#configuration-file-location).
Add the following setting in the **[auth]** section :

```bash
[auth]
azure_auth_enabled = true
```

{{% admonition type="note" %}}
If you are using Azure authentication settings do not enable `Forward OAuth identity`. Both use the same HTTP authorization headers. Azure settings will get overwritten by the OAuth token.
{{% /admonition %}}

Read more about connecting to Prometheus using Azure authenitication here: [Connect to Azure Monitor Managed Service for Prometheus](https://grafana.com/docs/grafana/latest/datasources/prometheus/configure/azure-authentication/).