---
aliases:
  - ../data-sources/azure-prometheus/alerting/
  - ../data-sources/azure-monitor-prometheus/alerting/
description: Set up alerting with the Azure Monitor Managed Service for Prometheus data source, including Grafana-managed and workspace-managed rules.
keywords:
  - grafana
  - prometheus
  - azure
  - azure monitor
  - alerting
  - recording rules
labels:
  products:
    - cloud
    - enterprise
    - oss
menuTitle: Alerting
title: Azure Monitor Managed Service for Prometheus alerting
weight: 450
review_date: 2026-08-20
---

# Azure Monitor Managed Service for Prometheus alerting

The Azure Monitor Managed Service for Prometheus data source works with Grafana Alerting. You can create Grafana-managed alert rules that query your workspace, and you can manage alerting and recording rules stored in the workspace from the Grafana Alerting UI.

## Before you begin

Before you set up alerting, ensure you have:

- [Configured the Azure Monitor Managed Service for Prometheus data source](https://grafana.com/docs/plugins/grafana-azureprometheus-datasource/latest/configure/).
- Azure credentials with permission to read and, if you manage rules from Grafana, write Prometheus rule groups in the workspace.

## Alert approaches

The data source supports two complementary alerting approaches.

### Grafana-managed alert rules

Grafana-managed alert rules run in Grafana and can query your Azure Monitor workspace as their data source. Grafana evaluates the rule, manages its state, and routes notifications. Use this approach when you want a single alerting experience across all your data sources, including expressions, images in notifications, and rules that query multiple data sources.

To create a Grafana-managed alert rule, refer to [Configure Grafana-managed alert rules](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/alerting-rules/create-grafana-managed-rule/).

For example, to alert when the rate of HTTP 5xx errors exceeds five percent over the last five minutes, use a query such as:

```promql
sum(rate(http_requests_total{status=~"5.."}[5m]))
/
sum(rate(http_requests_total[5m]))
> 0.05
```

Grafana-managed alert rules require credentials that work without a signed-in user. If the data source uses **Current User** authentication, configure [fallback service credentials](https://grafana.com/docs/plugins/grafana-azureprometheus-datasource/latest/configure/#current-user).

### Workspace-managed rules

The data source can also read and write the recording and alerting rules stored in your Azure Monitor workspace. These rules are evaluated by Azure Monitor, not by Grafana. Use this approach when you want your rules to run close to the data and remain available independently of Grafana.

The data source manages two kinds of rules:

- **Recording rules:** Pre-compute frequently used or expensive expressions and store the results as new time series for faster queries.
- **Alerting rules:** Evaluate a condition and fire when the condition holds.

To manage these rules from Grafana, enable **Manage alerts via Alerting UI** on the data source [configuration page](https://grafana.com/docs/plugins/grafana-azureprometheus-datasource/latest/configure/). When enabled, the rules appear in the Grafana Alerting UI, where you can view and edit them.

The data source reaches the workspace ruler through the `/rules` and `/config/v1/rules` paths under the **Prometheus server URL** you configure, so you don't need to set a separate ruler URL. If the rules fail to load with an **Unable to fetch alert rules** error, refer to [Troubleshooting](https://grafana.com/docs/plugins/grafana-azureprometheus-datasource/latest/troubleshooting/).

Managing rules in the workspace requires an identity with permission to read and write Prometheus rule groups on the Azure Monitor workspace, in addition to query access.

{{< admonition type="note" >}}
Workspace-managed rules require Prometheus rule groups in your Azure Monitor workspace. For more information, refer to the [Azure Monitor Prometheus rule groups documentation](https://learn.microsoft.com/en-us/azure/azure-monitor/metrics/prometheus-rule-groups).
{{< /admonition >}}

## Next steps

- [Configure the data source](https://grafana.com/docs/plugins/grafana-azureprometheus-datasource/latest/configure/)
- [Grafana Alerting documentation](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/)
