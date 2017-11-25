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

func TestPresentableResults(t *testing.T) {
	var hexastore *Hexastore = newHexastore()
	var testTriple = Triple{Subject: 0, Prop: 0, Object: 1}
	var testTripleTwo = Triple{Subject: 0, Prop: 0, Object: 2}

	hexastore.props.Put("property")
	hexastore.entities.Put("A")
	hexastore.entities.Put("B")
	hexastore.entities.Put("C")

	var testResults = []Triple{testTriple, testTripleTwo}

	got := *(PresentableResults(&testResults, hexastore))

	if len(got) != 2 {
		t.Error("Expected to put in 2 results and get 2 presentables back. got ", len(got))
	}

	if got[0] != "A -> property -> B" {
		t.Error("Incorrect presentation string returned")
	}

	if got[1] != "A -> property -> C" {
		t.Error("Incorrect presentation string returned")
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

func TestHexastoreRemove(t *testing.T) {
	var hexastore *Hexastore = newHexastore()
	triple := Triple{Subject: 1, Prop: 2, Object: 3, Value: "hello world"}

	hexastore.add(&triple)

	if hexastore.SPO[1][2][3] != "hello world" {
		t.Error("Minimal test of whether triple was correctly added failed!")
	}

	hexastore.remove(&triple)

	_, ok := hexastore.SPO[1][2][3]
	if ok != false {
		t.Error("Failed to remove triple from SPO")
	}
	if _, ok := hexastore.SOP[1][3][2]; ok {
		t.Error("Failed to remove triple from SOP")
	}
	if _, ok := hexastore.OPS[3][2][1]; ok {
		t.Error("Failed to remove triple from OPS")
	}
	if _, ok := hexastore.OSP[3][1][2]; ok {
		t.Error("Failed to remove triple from OSP")
	}
	if _, ok := hexastore.POS[2][3][1]; ok {
		t.Error("Failed to remove triple from POS")
	}
	if _, ok := hexastore.PSO[2][1][3]; ok {
		t.Error("Failed to remove triple from PSO")
	}
}

func TestQuerySXX(t *testing.T) {
	var hexastore *Hexastore = newHexastore()
	triple1 := Triple{Subject: 1, Prop: 2, Object: 3, Value: "hello world"}
	triple2 := Triple{Subject: 2, Prop: 2, Object: 3, Value: "hello world"}
	triple3 := Triple{Subject: 1, Prop: 3, Object: 3, Value: "hello world"}
	triple4 := Triple{Subject: 1, Prop: 4, Object: 3, Value: "hello world"}

	hexastore.add(&triple1)
	hexastore.add(&triple2)
	hexastore.add(&triple3)
	hexastore.add(&triple4)

	results := hexastore.QuerySXX(1)

	if len(*results) != 3 {
		t.Error("Subject-oriented query returned incorrect number of records. Expected 3, got ", len(*results))
	}
}

func TestQuerySPX(t *testing.T) {
	var hexastore *Hexastore = newHexastore()
	triple1 := Triple{Subject: 1, Prop: 2, Object: 3, Value: "hello world"}
	triple2 := Triple{Subject: 2, Prop: 2, Object: 3, Value: "hello world"}
	triple3 := Triple{Subject: 1, Prop: 3, Object: 3, Value: "hello world"}
	triple4 := Triple{Subject: 1, Prop: 4, Object: 3, Value: "hello world"}

	hexastore.add(&triple1)
	hexastore.add(&triple2)
	hexastore.add(&triple3)
	hexastore.add(&triple4)

	results := hexastore.QuerySPX(1, 3)

	if len(*results) != 1 {
		t.Error("Subject+Property oriented query returned incorrect num of records. Expected 1, got ", len(*results))
	}
}

func TestQueryXPO(t *testing.T) {
	var hexastore *Hexastore = newHexastore()
	triple1 := Triple{Subject: 1, Prop: 2, Object: 3, Value: "hello world"}
	triple2 := Triple{Subject: 2, Prop: 2, Object: 3, Value: "hello world"}
	triple3 := Triple{Subject: 4, Prop: 3, Object: 3, Value: "hello world"}
	triple4 := Triple{Subject: 4, Prop: 4, Object: 3, Value: "hello world"}

	hexastore.add(&triple1)
	hexastore.add(&triple2)
	hexastore.add(&triple3)
	hexastore.add(&triple4)

	results := hexastore.QueryXPO(2, 3)

	if len(*results) != 2 {
		t.Error("Subject-oriented query returned incorrect number of records. Expected 3, got ", len(*results))
	}

	for _, triple := range *results {
		if !(triple.Subject == 1 || triple.Subject == 2) {
			t.Error("Query returned incorrect triple. Expected Subject ID to be 1, got ", triple.Subject)
		}
	}
}

func TestQueryXXX(t *testing.T) {
	var hexastore *Hexastore = newHexastore()
	triple1 := Triple{Subject: 1, Prop: 2, Object: 3, Value: "hello world"}
	triple2 := Triple{Subject: 2, Prop: 2, Object: 3, Value: "hello world"}
	triple3 := Triple{Subject: 4, Prop: 3, Object: 3, Value: "hello world"}
	triple4 := Triple{Subject: 4, Prop: 4, Object: 3, Value: "hello world"}

	hexastore.add(&triple1)
	hexastore.add(&triple2)
	hexastore.add(&triple3)
	hexastore.add(&triple4)

	results := hexastore.QueryXXX()

	if len(*results) != 4 {
		t.Error("*-type query did not return all records in hexastore")
	}
}

func TestQuerySPO(t *testing.T) {
	var hexastore *Hexastore = newHexastore()
	triple1 := Triple{Subject: 1, Prop: 2, Object: 3, Value: "hello mars"}
	hexastore.add(&triple1)

	result := hexastore.QuerySPO(1, 2, 3)
	if len(*result) != 1 {
		t.Error("Full subject-property-object query did not return expected value. Wanted 'hello mars', got ", *result)
	}

	first := (*result)[0]
	if first.Value != "hello mars" {
		t.Error("Something went wrong")
	}
}
