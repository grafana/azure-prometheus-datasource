## Compatibility

Azure Managed Service for Prometheus is not compatible with Grafana<11.5.0 If you are running version 11.4.x and lower, please use plugin the core Prometheus plugin instead.

# Azure Managed Service for Prometheus Data Source

Azure Managed Service for Prometheus is a Prometheus-compatible service that monitors and provides alerts on containerized applications and infrastructure at scale.

Read more about it here:

[Azure Monitor and Prometheus](https://learn.microsoft.com/en-us/azure/azure-monitor/metrics/prometheus-metrics-overview)

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

## Migrate from core Prometheus to Azure Managed Service for Prometheus plugin

Azure AD authentication is deprecated in the core Prometheus plugin. You will need to do the following to migrate your Prometheus plugin using Azure AD to Azure Managed Prometheus plugin.

1. Get the `UID` for Prometheus using Azure Authentication and get the `UID` for your new Azure Managed Prometheus.
    - Navigate to the configuration page for your new data source.
    - Find the `UID` in the url.
      - Example: `“connections/datasources/edit/<DATA SOURCE UID>”`
    - Copy the UID for both the old and new data sources.
      - `"<Prom AD UID>"`
      - `"<AMP UID>"`

2. Update dashboards with the new datasource `UID`
    - Navigate to the dashboard that contains your old Prometheus data source.
    - Click on Dashboard settings.
    - Select the JSON model in the tabs.
    - Search for the `"<Prom AD UID>"`.
      - Example: `"uid": "<Prom AD UID>"`
      - Change the `UID` to `<AMP UID>`
      - Change the type to `grafana-azureprometheus-datasource` (the old type is `prometheus`)
      - Example for changing `type` and `UID`:
      ```
        "type": "prometheus",
        "uid": "<Prom AD UID>"

        // can be change to

        "type": "grafana-azureprometheus-datasource",
        "uid": "<AMP UID>"
      ```
    - Confirm the "datasource" change for all of the following categories in the JSON model.
      - `annotations`
      - `panels`
      - `targets`
      - `templating`
    - Click “Save changes” in the bottom left side of the JSON model UI.
  Save your dashboard.

2. Update alert rules by exporting provisioning files and updating the data source in the model or create new alert rules.
    - Alert rule data sources cannot be changed without wiping out the query.
    - There are two ways to migrate alert rules.
      - First, edit the rule by exporting to a provisioned file [documentation here](https://grafana.com/docs/grafana/latest/alerting/set-up/provision-alerting-resources/export-alerting-resources/#export-alerting-resources).
        - Navigate to the Alert rules page.
        - Identify the alert rule that uses the Prom AD data source.
        - Select “Export rule.”
        - Export the rule in your choice of JSON, YAML or Terraform.
        - Search the exported rule for the `"<Prom AD UID>"`.
        - Change the `UID` and the `type` in the exported rule.
      ```
        "type": "prometheus",
        "uid": "<Prom AD UID>"

        // can be change to

        "type": "grafana-azureprometheus-datasource",
        "uid": "<AMP UID>"
      ```
        - Rename the rule and update the `name` field in the exported rule.
        - Delete the old rule.
      - OR second, Create a new alert rule [documentation here](https://grafana.com/tutorials/alerting-get-started/).
        - Copy the fields from the Prom AD rule.
        - Create a new alert rule and select your AMP data source.
        - Add the fields from the old rule to the new rule.
        - Delete the old rule.

3. Recreate correlations that are using your old Prometheus data source using your new Prometheus data source.
    - Create a new correlation, see [documentation here](https://grafana.com/docs/grafana/latest/administration/correlations/create-a-new-correlation/).
    - Identify the correlation that uses the old Prom AD data source.
    - Copy the fields from your correlation.
    - Create a brand new correlation using the new AMP data source.
    - Add the fields from the old correlation to the new correlation.
    - Delete the old correlation.

4. Recreate recorded queries that are using your old Prometheus data source using your new Prometheus data source.
    - Only available in Grafana Enterprise and Grafana Cloud.
    - Create a new recorded query, see [documentation here](https://grafana.com/docs/grafana/latest/administration/correlations/create-a-new-correlation/).
    - Identify the recorded query that uses the old Prom AD data source.
    - Copy the fields from your recorded query.
    - Create a brand new recorded query using the new AMP data source.
    - Add the fields from the old recorded query to the new recorded query.
    - Delete the old recorded query.
