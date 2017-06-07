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
