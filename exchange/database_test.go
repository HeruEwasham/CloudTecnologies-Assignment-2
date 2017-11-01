package exchange

import "testing"
import "fmt"
import "time"

var testdb Storage

func setupTestdatabase() {
	testdb = &MongoDB{
		"mongodb://CloudFullAccess:full1916@ds227045.mlab.com:27045/herus-cloud-tecnologies",
		"herus-cloud-tecnologies",
		"webhooks_test",
		"currencies_test",
	}
}

func Test_FloatToString(t *testing.T) {
	var input float32
	input = 2.35
	correctOutput := "2.35"
	output := FloatToString(input)
	if output != correctOutput {
		t.Error("Error when converting float to string, output is " + output + ", while correct output should be " + correctOutput)
	}
}

func Test_sendMsgWebhook(t *testing.T) {
	msg := MessageWebhook {"Test", time.Now().Format("2006-01-02-15:04:05"),"This is a test, if this is registered it worked", "Cloud tecnologies: Assignment 2"}
	ok := SendMessageWebhook(msg)
	if !ok {
		t.Error("Send Message didn't work")
	}
}

func Test_RegisterAndGetWebhook(t *testing.T) {
	setupTestdatabase() //?
	testdb.Init()		//?
	webhook := Webhook{WebhookURL:"http://example.com",BaseCurrency:"EUR",TargetCurrency:"NOK",MinTriggerValue: 1.5, MaxTriggerValue: 2.8}
	id, statusCode, err := testdb.RegisterWebhookToDatabase(webhook)
	if err != nil {
		t.Error("Error when registering webhook, statuscode is ", statusCode, ", and error is ", err)
		return
	}
	//if (id == "") {
	//	t.Error("ID is empty")
	//	return
	//}

	webhookGotten, statusCode, err := testdb.GetWebhook(id)
	if err != nil {
		t.Error("Error when getting webhook, id is: " + id + "statuscode is ", statusCode, ", and error is ", err)
		return
	}
	if webhookGotten.BaseCurrency != webhook.BaseCurrency || webhookGotten.MaxTriggerValue != webhook.MaxTriggerValue || webhookGotten.MinTriggerValue != webhook.MinTriggerValue || webhookGotten.TargetCurrency != webhook.TargetCurrency || webhookGotten.WebhookURL != webhook.WebhookURL {
		t.Error("base-webhook is different than the webhook gotten.")
		return
	}

	statusCode, err = SendWebhookFunc(webhookGotten, 1.8)
	if err != nil {
		t.Error("Something is wrong when sending webhook.")
	}

	webhooksGotten, statusCodeAll, errAll := testdb.GetAllWebhooks()
	if err != nil {
		t.Error("Error when getting all webhooks, statuscode is ", statusCodeAll, ", and error is ", errAll)
		return
	}
	fmt.Println(webhooksGotten)
	if webhooksGotten[0].BaseCurrency != webhook.BaseCurrency || webhooksGotten[0].MaxTriggerValue != webhook.MaxTriggerValue || webhookGotten.MinTriggerValue != webhook.MinTriggerValue || webhookGotten.TargetCurrency != webhook.TargetCurrency || webhookGotten.WebhookURL != webhook.WebhookURL {
		t.Error("base-webhook is different than the webhook gotten when calling for all.")
		return
	}

	statusCode, err = testdb.DeleteWebhook(id)
	if err != nil {
		t.Error("Error when deleting webhook, statuscode is ", statusCode, ", and error is ", err)
		return
	}
	ok := testdb.ResetWebhook()
	if !ok {
		t.Error("Couldn't reset Webhook-collection")
		return
	}
}

func Test_GetLatest(t *testing.T) {
	setupTestdatabase() //?
	testdb.Init()		//?
	rateMap := make(map[string]float32) // Make a map with rates
	rateMap["NOK"] = 1.56
	currency := Currency{"EUR","2100-01-01",rateMap}
	statusCode, err := testdb.RegisterCurrencyToDatabase(currency)
	if err != nil {
		t.Error("Error when registering currency, statuscode is ", statusCode, ", and error is ", err)
		return
	}
	latestCurrency, _, statusCode, err := testdb.GetLatest("NOK")		// Testing by getting latest Currency for Norwegian Kroner
	if err != nil {
		t.Error("Error when getting latest currency, statuscode is ", statusCode, ", and error is ", err)
		return
	}

	if latestCurrency != rateMap["NOK"] {
		t.Error("Latest currency gotten is not the one inserted just before, latest currency inserted is ", rateMap["NOK"], ", while we got ", latestCurrency)
		return
	}
	ok := testdb.ResetCurrency()
	if !ok {
		t.Error("Couldn't reset Currency-collection")
		return
	}
}

func Test_GetAverage(t *testing.T) {
	setupTestdatabase() //?
	testdb.Init()		//?
	rateMap := make(map[string]float32) // Make a map with rates
	var rateAverage float32
	for i:=1;i<=3;i++ {
		rateMap["NOK"] = 1.56 + float32(i)
		rateAverage += rateMap["NOK"]
		currency := Currency{"EUR","2100-01-0" + string(i), rateMap}
		statusCode, err := testdb.RegisterCurrencyToDatabase(currency)
		if err != nil {
			t.Error("Error when registering currency number ", i, ", statuscode is ", statusCode, ", and error is ", err)
			return
		}
	}
	averageCurrrency, statusCode, err := testdb.GetAverage("NOK")		// Testing by getting latest Currency for Norwegian Kroner
	if err != nil {
		t.Error("Error when getting average currency, statuscode is ", statusCode, ", and error is ", err)
		return
	}

	if averageCurrrency != (rateAverage/3) {
		t.Error("Average currency gotten is not the one inserted just before, latest currency inserted is ", rateAverage, ", while we got ", averageCurrrency)
		return
	}
	
	ok := testdb.ResetCurrency()

	if !ok {
		t.Error("Error when resetting currency")
	}
}