package main

import "testing"

func TestMakeTriple(t *testing.T) {
	var triple *Triple = MakeTriple(1, 2, 3, "1234")

	if triple.Subject != 1 {
		t.Error("Expected 1, got ", triple.Subject)
	}
	if triple.Prop != 2 {
		t.Error("Expected 2, got ", triple.Prop)
	}
	if triple.Object != 3 {
		t.Error("Expected 3, got ", triple.Object)
	}
	if triple.Value != "1234" {
		t.Error("Expected '1234', got ", triple.Value)
	}
}

func TestNewHexastore(t *testing.T) {
	var hexastore *Hexastore = newHexastore()

	if hexastore.SPO == nil || hexastore.SOP == nil || hexastore.POS == nil || hexastore.PSO == nil || hexastore.OSP == nil || hexastore.OPS == nil {
		t.Error("Didn't initialise Hexastore correctly")
	}
}

func TestGetKey(t *testing.T) {
	var dict Dictionary

	dict.m = make(map[int]string)

	dict.m[dict.NextKey] = "hello world"
	dict.m[dict.NextKey+1] = "goodbye world"

	hello_world_key, ok := dict.GetKey("hello world")
	if hello_world_key != 0 {
		t.Error("Expected 0, got ", hello_world_key)
	}
	if ok != true {
		t.Error("Expected success variable to be true")
	}
	goodbye_world_key, ok := dict.GetKey("goodbye world")
	if goodbye_world_key != 1 {
		t.Error("Expected 1, got ", goodbye_world_key)
	}
	if ok != true {
		t.Error("Expected success variable to be true")
	}
}

func TestHexastoreAdd(t *testing.T) {
	var hexastore *Hexastore = newHexastore()
	triple := Triple{Subject: 1, Prop: 2, Object: 3, Value: "hello world"}

	hexastore.add(&triple)

	if hexastore.SPO[1][2][3] != "hello world" {
		t.Error("Failed to add triple into SPO index")
	}
	if hexastore.SOP[1][3][2] != "hello world" {
		t.Error("Failed to add triple into SOP index")
	}
	if hexastore.OPS[3][2][1] != "hello world" {
		t.Error("Failed to add triple into OPS index: ", hexastore.OPS[3][2][1])
	}
	if hexastore.OSP[3][1][2] != "hello world" {
		t.Error("Failed to add triple into OSP index")
	}
	if hexastore.POS[2][3][1] != "hello world" {
		t.Error("Failed to add triple into POS index")
	}
	if hexastore.PSO[2][1][3] != "hello world" {
		t.Error("Failed to add triple into PSO index", hexastore.PSO[2][1][3])
	}
}