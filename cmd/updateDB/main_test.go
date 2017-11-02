package main

import "github.com/HeruEwasham/CloudTecnologies-Assignment-2/exchange"
import "testing"

var testdb exchange.Storage
var normaldb exchange.Storage

func setupNormalDatabase() {
	normaldb = databaseCred(false)
}
func setupTestdatabase() {
	testdb = databaseCred(true)
}

func Test_GetTodaysCurrency(t *testing.T) {
	setupNormalDatabase() // Test this part
	setupTestdatabase()   //?
	testdb.Init()         //?
	ok := getCurrencyFromExternalDatabase(testdb, "latest")
	if !ok {
		t.Error("Function getTodaysCurrency(..) failed.")
		return
	}
	ok = getCurrencyFromExternalDatabase(testdb, "latest") // Checks for some other parts of the function
	if !ok {
		t.Error("Function getTodaysCurrency(..) failed on second time.")
		return
	}

	ok = testdb.ResetCurrency()
	if !ok {
		t.Error("Couldn't reset Currency-collection")
		return
	}
}
