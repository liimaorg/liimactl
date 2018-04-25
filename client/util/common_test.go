package util

import "testing"

func TestBuildCommandURL(t *testing.T) {

	//Test type
	type testStruct struct {
		AppServer    []string `json:"appServer"`
		Environment  []string `json:"environment"`
		DisableMerge bool     `json:"disableMerge"`
	}

	//Test args
	type args struct {
		gType testStruct
	}

	//Tests
	tests := []struct {
		name string //Name of the test
		args args   //Arguments
		want string //Wanted testresult
	}{
		{"Test1", args{testStruct{AppServer: []string{"TestAppX", "TestAppY"}, Environment: []string{"T"}, DisableMerge: true}}, "appServer=TestAppX,TestAppY&environment=T&disableMerge=true"},
		{"Test2", args{testStruct{DisableMerge: false}}, "disableMerge=false"},
	}

	//Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildCommandURL(&tt.args.gType); got != tt.want {
				t.Errorf("BuildCommandURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagExists(t *testing.T) {

	//Test type
	type testStruct struct {
		AppServer    []string `json:"appServer"`
		Environment  []string `json:"environment"`
		DisableMerge bool     `json:"disableMerge"`
	}

	//Test args
	type args struct {
		gType   testStruct
		testTag string
	}

	//Tests
	tests := []struct {
		name string //Name of the test
		args args   //Arguments
		want bool   //Wanted testresult
	}{
		{"Test1", args{testStruct{AppServer: []string{"TestAppX", "TestAppY"}, Environment: []string{"T"}, DisableMerge: true}, "fail"}, false},     // this test must fail because "fail" not exists as json-Tag
		{"Test2", args{testStruct{AppServer: []string{"TestAppX", "TestAppY"}, Environment: []string{"T"}, DisableMerge: true}, "appServer"}, true}, // this test must be ok because "appServer" exists as json-Tag
	}

	//Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TagExists(&tt.args.gType, tt.args.testTag); got != tt.want {
				t.Errorf("TagExists(%v, %v) = %v, want %v", &tt.args.gType, tt.args.testTag, got, tt.want)
			}
		})
	}
}

func TestSetValueIfTagExists(t *testing.T) {

	//Test type
	type testStruct struct {
		AppServer    string `json:"appServer"`
		Environment  string `json:"environment"`
		DisableMerge bool   `json:"disableMerge"`
	}

	//Test args
	type args struct {
		gType     testStruct
		testTag   string
		testValue string
	}

	//Tests
	tests := []struct {
		name string //Name of the test
		args args   //Arguments
		want string //Wanted testresult
	}{
		{"Test1", args{testStruct{AppServer: "TestAppX", Environment: "T", DisableMerge: true}, "fail", "testServerA"}, "TestAppX"},         // this test must fail because "fail" not exists as json-Tag -> Want Default "TestAppX"
		{"Test2", args{testStruct{AppServer: "TestAppY", Environment: "T", DisableMerge: true}, "appServer", "testServerB"}, "testServerB"}, // this test must be ok because "appServer" exists as json-Tag
	}

	//Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetValueIfTagExists(&tt.args.gType, tt.args.testTag, tt.args.testValue)
			if got := tt.args.gType.AppServer; got != tt.want {
				t.Errorf("SetValueIfTagExists(%v, %v, %v) = %v, want %v", &tt.args.gType, tt.args.testTag, tt.args.testValue, got, tt.want)
			}
		})
	}
}
