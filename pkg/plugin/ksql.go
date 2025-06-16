package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/thmeitz/ksqldb-go"
	ksqlnet "github.com/thmeitz/ksqldb-go/net"
)

func KsqlClient(url string, httpAllow bool, credentials ksqlnet.Credentials) ksqldb.KsqldbClient {
	cons := ksqlnet.Options{
		Credentials: credentials,
		BaseUrl:     url,
		AllowHTTP:   httpAllow,
	}
	client, err := ksqldb.NewClientWithOptions(cons)
	if err != nil {
		log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource ---------------> ERROR: can't connect to server !!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	}

	// ctx, _ := context.WithTimeout(context.TODO(), time.Duration(700)*time.Second)
	// info, _ := client.GetServerInfo(ctx)
	// log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource --------------->" + info.KafkaClusterID)
	return client
}

func ProcessPushData(rowChannel chan ksqldb.Row) {
	// var PUBLICSUFFIX_TLD1 string
	// var querycount float64

	for row := range rowChannel {
		if row != nil {
			// Handle the timestamp

			// t := int64(row[0].(string))
			// ts := time.Unix(t/1000, 0).Format(time.RFC822)
			parsedTime, err := time.Parse(time.RFC3339Nano, row[0].(string))
			if err == nil {
				row = append(row[1:], parsedTime)
			}
			// log.Infof("🐾 New dog at %v: '%v' is %v and %v (id %v)\n", ts, name, dogSize, age, id)
			// P.Println(row...)
			fmt.Println(row...)
		}
	}
}

func ProcessEachPush(row ksqldb.Row, head ksqldb.Header, sender *backend.StreamSender) error {
	tempFrame := data.NewFrame("response")
	// frames := []data.Frame

	// ["abc" true "nnc" 767575351]
	for i := range row {
		switch row[i].(type) {
		case string:
			if head.Columns[i].Name != "DATETIME" {
				tempFrame.Fields = append(tempFrame.Fields, data.NewField(head.Columns[i].Name, nil, []string{row[i].(string)}))
			} else {
				val, _ := time.Parse(time.RFC3339Nano, row[i].(string))
				tempFrame.Fields = append(tempFrame.Fields, data.NewField("DATETIME", nil, []time.Time{val}))
			}
		case int8:
			tempFrame.Fields = append(tempFrame.Fields, data.NewField(head.Columns[i].Name, nil, []int8{row[i].(int8)}))
		case int32:
			tempFrame.Fields = append(tempFrame.Fields, data.NewField(head.Columns[i].Name, nil, []int32{row[i].(int32)}))
		case int64:
			tempFrame.Fields = append(tempFrame.Fields, data.NewField(head.Columns[i].Name, nil, []int64{row[i].(int64)}))
		case float32:
			tempFrame.Fields = append(tempFrame.Fields, data.NewField(head.Columns[i].Name, nil, []float32{row[i].(float32)}))
		case float64:
			tempFrame.Fields = append(tempFrame.Fields, data.NewField(head.Columns[i].Name, nil, []float64{row[i].(float64)}))
		case bool:
			tempFrame.Fields = append(tempFrame.Fields, data.NewField(head.Columns[i].Name, nil, []bool{row[i].(bool)}))
		}
	}

	err := sender.SendFrame(
		tempFrame,
		data.IncludeAll,
	)
	log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource ---------------> We are inside")

	return err
}

func StartPushStream(client ksqldb.KsqldbClient, query string, timeout int, rowChannel chan ksqldb.Row, headerChannel chan ksqldb.Header) {
	// rowChannel := make(chan ksqldb.Row)
	// headerChannel := make(chan ksqldb.Header, 1)

	fmt.Println(query, timeout)
	// query = "select * from DNSSTRAEM EMIT CHANGES;"
	// This Go routine will handle rows as and when they are sent to the channel
	// go ProcessPushData(rowChannel)

	// go showHeader(headerChannel)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(timeout)*time.Second)
	defer cancel()

	log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource ---------------> Before the push")
	e := client.Push(ctx, ksqldb.QueryOptions{Sql: query}, rowChannel, headerChannel)
	log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource ---------------> After the push")
	log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource --------------->" + e.Error())
	if e != nil {
		log.DefaultLogger.Error("plugin.pax-voyagerksql-datasource --------------->" + e.Error())
	}

	// return rowChannel, headerChannel
}

// client := KsqlClient("http://10.4.4.195:3005", true, ksqlnet.Credentials{Username: "", Password: ""})
// defer client.Close()
// query := "select * from DNSSTRAEM EMIT CHANGES;"
// StartPushStream(client, query, 60) --> returns: 2 channels
// from cahnnel header, create the
