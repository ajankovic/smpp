package pdu

import (
	"testing"
	"time"
)

func TestParsingGoodDeliveryReceipt(t *testing.T) {
	good := "id:123123123 sub:0 dlvrd:0 submit date:1507011202 done date:1507011101 stat:DELIVRD err:0 text:Test information"
	dr, err := ParseDeliveryReceipt(good)
	if err != nil {
		t.Errorf("Error parsing good receipt %s", err)
	}
	if dr.Id != "123123123" {
		t.Errorf("Receipt id is wrong %s expected 123123123", dr.Id)
	}
	extime, _ := time.Parse(time.RFC3339, "2015-07-01T12:02:00Z")
	if dr.SubmitDate != extime {
		t.Errorf("Receipt submit date is wrong %s expected %s",
			dr.SubmitDate.Format(time.RFC3339),
			extime.Format(time.RFC3339),
		)
	}
	if dr.String() != good {
		t.Errorf("Receipt string representation is wrong %s", dr)
	}
}

func TestParsingBadDeliveryReceipt(t *testing.T) {
	keys := "id:123123123 dfdfsub:0 dlvrd:0 submit date:1507011202 done date:1507011101 stat:DELIVRD err:0 text:Test information"
	_, err := ParseDeliveryReceipt(keys)
	if err == nil {
		t.Errorf("Parsing bad receipt with wrong key name returned no error")
	}
	missingkeys := "id:123123123 sub:0 dlvrd:0 submit date:1507011202 stat:DELIVRD err:0 text:Test information"
	_, err = ParseDeliveryReceipt(missingkeys)
	if err == nil {
		t.Errorf("Parsing bad receipt with missing keys returned no error")
	}
	date := "id:123123123 sub:0 dlvrd:0 submit date:150701adsfas1202 done date:1507011101 stat:DELIVRD err:0 text:Test information"
	_, err = ParseDeliveryReceipt(date)
	if err == nil {
		t.Errorf("Parsing bad receipt with wrong date format returned no error")
	}
}

func TestParsingUUIDDeliveryReceipt(t *testing.T) {
	dlr := "id:a03ea27b-9bb4-4d5e-b87f-3f578ab46153 sub:001 dlvrd:001 submit date:161003211236 done date:161003211236 stat:DELIVRD err:000 text:-"
	r, err := ParseDeliveryReceipt(dlr)
	if err != nil {
		t.Fatalf("Error parsing UUID delivery receipt %v", err)
	}
	if r.Id != "a03ea27b-9bb4-4d5e-b87f-3f578ab46153" {
		t.Errorf("ParseDeliveryReceipt() => %s expected %s", r.Id, "a03ea27b-9bb4-4d5e-b87f-3f578ab46153")
	}
	if r.Stat != "DELIVRD" {
		t.Errorf("ParseDeliveryReceipt() => %s expected %s", r.Stat, "DELIVRD")
	}
}
