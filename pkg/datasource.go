package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

func getQueryReqHeader(req *backend.QueryDataRequest, headerName string) string {
	headerNameCI := strings.ToLower(headerName)

	for name, value := range req.Headers {
		if strings.ToLower(name) == headerNameCI {
			return value
		}
	}

	return ""
}

func getResourceReqHeader(req *backend.CallResourceRequest, headerName string) string {
	headerNameCI := strings.ToLower(headerName)

	for name, values := range req.Headers {
		if strings.ToLower(name) == headerNameCI {
			if len(values) > 0 {
				return values[0]
			} else {
				return ""
			}
		}
	}

	return ""
}

func getCheckHealthReqHeader(req *backend.CheckHealthRequest, headerName string) string {
	headerNameCI := strings.ToLower(headerName)

	for name, value := range req.Headers {
		if strings.ToLower(name) == headerNameCI {
			return value
		}
	}

	return ""
}

func logStuff(id string, logger log.Logger, endpoint string) {
	claims := jwt.MapClaims{}
	_, _, err := jwt.NewParser(jwt.WithValidMethods([]string{"ES256"})).ParseUnverified(id, claims)
	if err != nil {
		logger.Error("error:", err)
	}
	for r, i := range claims {
		logger.Info("claim "+endpoint, fmt.Sprintf("key: %s value: %v", r, i))
	}
}

func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	grafanaIdToken := getQueryReqHeader(req, "http_X-Grafana-Id")
	logStuff(grafanaIdToken, backend.Logger.FromContext(ctx), "Query data")
	ctx = d.contextualMiddlewares(ctx)
	ctx = azusercontext.WithUserFromQueryReq(ctx, req)
	return d.Service.QueryData(ctx, req)
}

func (d *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	grafanaIdToken := getResourceReqHeader(req, "http_X-Grafana-Id")
	logStuff(grafanaIdToken, backend.Logger.FromContext(ctx), "Call resource")
	ctx = d.contextualMiddlewares(ctx)
	ctx = azusercontext.WithUserFromResourceReq(ctx, req)
	return d.Service.CallResource(ctx, req, sender)
}

func (d *Datasource) GetBuildInfo(ctx context.Context, req promlib.BuildInfoRequest) (*promlib.BuildInfoResponse, error) {
	uc, _ := azusercontext.GetCurrentUser(ctx)
	logStuff(uc.IdToken, backend.Logger.FromContext(ctx), "build info")
	ctx = d.contextualMiddlewares(ctx)

	return d.Service.GetBuildInfo(ctx, req)
}

func (d *Datasource) GetHeuristics(ctx context.Context, req promlib.HeuristicsRequest) (*promlib.Heuristics, error) {
	uc, _ := azusercontext.GetCurrentUser(ctx)
	logStuff(uc.IdToken, backend.Logger.FromContext(ctx), "heuristics")
	ctx = d.contextualMiddlewares(ctx)
	return d.Service.GetHeuristics(ctx, req)
}

func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult,
	error) {
	grafanaIdToken := getCheckHealthReqHeader(req, "http_X-Grafana-Id")
	logStuff(grafanaIdToken, backend.Logger.FromContext(ctx), "check health")
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
