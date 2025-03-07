package main

import (
	"context"
	"fmt"

	"github.com/grafana/grafana-azure-sdk-go/v2/azsettings"
	"github.com/grafana/grafana-azure-sdk-go/v2/azusercontext"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	"github.com/grafana/azure-prometheus-datasource/pkg/azureauth"
	"github.com/grafana/grafana/pkg/promlib"
)

func NewDatasource(ctx context.Context, dsInstanceSettings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	plog := backend.NewLoggerWith("logger", "tsdb.azure-prometheus")
	plog.Debug("Initializing")

	return &Datasource{
		Service: promlib.NewService(sdkhttpclient.NewProvider(), plog, extendClientOpts),
	}, nil
}

type Datasource struct {
	Service *promlib.Service
}

func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	ctx = d.contextualMiddlewares(ctx)
	ctx = azusercontext.WithUserFromQueryReq(ctx, req)
	return d.Service.QueryData(ctx, req)
}

func (d *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	ctx = d.contextualMiddlewares(ctx)
	ctx = azusercontext.WithUserFromResourceReq(ctx, req)
	return d.Service.CallResource(ctx, req, sender)
}

func (d *Datasource) GetBuildInfo(ctx context.Context, req promlib.BuildInfoRequest) (*promlib.BuildInfoResponse, error) {
	ctx = d.contextualMiddlewares(ctx)
	return d.Service.GetBuildInfo(ctx, req)
}

func (d *Datasource) GetHeuristics(ctx context.Context, req promlib.HeuristicsRequest) (*promlib.Heuristics, error) {
	ctx = d.contextualMiddlewares(ctx)
	return d.Service.GetHeuristics(ctx, req)
}

func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult,
	error) {
	ctx = azusercontext.WithUserFromHealthCheckReq(ctx, req)
	ctx = d.contextualMiddlewares(ctx)
	return d.Service.CheckHealth(ctx, req)
}

func (d *Datasource) contextualMiddlewares(ctx context.Context) context.Context {
	cfg := backend.GrafanaConfigFromContext(ctx)

	middlewares := []sdkhttpclient.Middleware{
		sdkhttpclient.ResponseLimitMiddleware(cfg.ResponseLimit()),
	}

	return sdkhttpclient.WithContextualMiddleware(ctx, middlewares...)
}

func extendClientOpts(ctx context.Context, settings backend.DataSourceInstanceSettings, clientOpts *sdkhttpclient.Options, plog log.Logger) error {
	azureSettings, err := azsettings.ReadSettings(ctx)
	if err != nil {
		return fmt.Errorf("failed to read Azure settings from Grafana: %v", err)
	}

	// Set Azure authentication
	if azureSettings.AzureAuthEnabled {
		err = azureauth.ConfigureAzureAuthentication(settings, azureSettings, clientOpts, plog)
		if err != nil {
			return fmt.Errorf("error configuring Azure auth: %v", err)
		}
	}

	return nil
}
