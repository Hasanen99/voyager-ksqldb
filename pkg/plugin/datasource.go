package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/pax/voyager-ksql/pkg/models"
	"github.com/thmeitz/ksqldb-go"
	ksqlnet "github.com/thmeitz/ksqldb-go/net"
	"go.uber.org/zap"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	// _ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)

	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
	_ backend.StreamHandler         = (*Datasource)(nil)
)

// NewDatasource creates a new datasource instance.
func NewDatasource(_ context.Context, _ backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &Datasource{}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct{}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
// func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
// 	// create response struct
// 	response := backend.NewQueryDataResponse()

// 	// loop over queries and execute them individually.
// 	for _, q := range req.Queries {
// 		res := d.query(ctx, req.PluginContext, q)

// 		// save the response in a hashmap
// 		// based on with RefID as identifier
// 		response.Responses[q.RefID] = res
// 	}

// 	return response, nil
// }

// type queryModel struct{}

// func (d *Datasource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
// 	var response backend.DataResponse

// 	// Unmarshal the JSON into our queryModel.
// 	var qm queryModel

// 	err := json.Unmarshal(query.JSON, &qm)
// 	if err != nil {
// 		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
// 	}

// 	// create data frame response.
// 	// For an overview on data frames and how grafana handles them:
// 	// https://grafana.com/developers/plugin-tools/introduction/data-frames
// 	frame := data.NewFrame("response")

// 	// add fields.
// 	frame.Fields = append(frame.Fields,
// 		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
// 		data.NewField("values", nil, []int64{10, 20}),
// 	)

// 	// add the frames to the response.
// 	response.Frames = append(response.Frames, frame)

// 	return response
// }

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
// func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
// 	res := &backend.CheckHealthResult{}
// 	config, err := models.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)

// 	if err != nil {
// 		res.Status = backend.HealthStatusError
// 		res.Message = "Unable to load settings"
// 		return res, nil
// 	}

// 	if config.Secrets.ApiKey == "" {
// 		res.Status = backend.HealthStatusError
// 		res.Message = "API key is missing"
// 		return res, nil
// 	}

// 	return &backend.CheckHealthResult{
// 		Status:  backend.HealthStatusOk,
// 		Message: "Data source is working",
// 	}, nil
// }

// --------------- For streaming ---------------

var Logger *zap.SugaredLogger // Declare a package-level logger

type Query struct {
	TheQuery string `json:"queryText"`
	Timeout  int    `json:"timeout"`
}

func (d *Datasource) SubscribeStream(context.Context, *backend.SubscribeStreamRequest) (*backend.SubscribeStreamResponse, error) {
	return &backend.SubscribeStreamResponse{
		Status: backend.SubscribeStreamStatusOK,
	}, nil
}

func (d *Datasource) PublishStream(context.Context, *backend.PublishStreamRequest) (*backend.PublishStreamResponse, error) {
	return &backend.PublishStreamResponse{
		Status: backend.PublishStreamStatusPermissionDenied,
	}, nil
}

func (d *Datasource) CheckHealth(_ context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	res := &backend.CheckHealthResult{}
	config, err := models.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)

	if err != nil {
		res.Status = backend.HealthStatusError
		res.Message = "Unable to load settings"
		return res, nil
	}

	if config.Ksql == "" {
		res.Status = backend.HealthStatusError
		res.Message = "KsqlDB server is a manditory parameter"
		return res, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}

func (d *Datasource) RunStream(ctx context.Context, req *backend.RunStreamRequest, sender *backend.StreamSender) error {
	config, _ := models.LoadPluginSettings(*req.PluginContext.DataSourceInstanceSettings)

	queryInput := Query{}
	json.Unmarshal(req.Data, &queryInput)

	// log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource ---------------> values are: ksql:" + config.Ksql + " Query:" + queryInput.TheQuery + " username:" + config.Username + "| pass:" + config.Secrets.Pass + "|")

	log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource ---------------> client not created yet")

	// fmt.Printf("\n\n\n\n\n conffff url: %v http: %v User: %v Pass: %v \n\n\n\n\n", q.Ksql, q.Http)
	// Create the client
	client := KsqlClient(config.Ksql, config.Http, ksqlnet.Credentials{Username: config.Username, Password: config.Secrets.Pass})
	defer client.Close()

	log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource ---------------> client created")

	rowChannel := make(chan ksqldb.Row)
	headerChannel := make(chan ksqldb.Header, 1)
	go StartPushStream(client, queryInput.TheQuery, int(queryInput.Timeout), rowChannel, headerChannel)
	log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource ---------------> channels created")

	var singleHeader ksqldb.Header
	var singleRow ksqldb.Row

	for {
		select {
		// case <-ctx.Done():
		// 	return ctx.Err()
		case singleHeader = <-headerChannel:
			fmt.Println(singleHeader)
		case singleRow = <-rowChannel:
			fmt.Println(singleRow)
			log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource -------- header & row> ")

			err := ProcessEachPush(singleRow, singleHeader, sender)
			if err != nil {
				// Logger.Error("Failed send frame", "error", err)
				log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource -------- error of process each> Failed send frame " + err.Error())
			}

			// err := sender.SendFrame(
			// 	data.NewFrame(
			// 		"response",
			// 		data.NewField("key", nil, []string{"one"}),
			// 		data.NewField("val", nil, []int32{int32(42)}),
			// 	),
			// 	data.IncludeAll,
			// )
			// we generate a random value using the intervals provided by the frontend
			// randomValue := r.Float64()*(q.UpperLimit-q.LowerLimit) + q.LowerLimit
			// err := sender.SendFrame(
			// 	data.NewFrame(
			// 		"response",
			// 		data.NewField("time", nil, []time.Time{time.Now()}),
			// 		data.NewField("value", nil, []string{q.TheQuery})),
			// 	data.IncludeAll,
			// )
		}
	}
}
