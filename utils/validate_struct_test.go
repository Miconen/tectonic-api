package utils

import (
	"testing"
)

type Interfc interface {
	DoIt()
}

type Impl struct {
	implField map[string]string
}

func (i Impl) DoIt() {}

type Data struct {
	X            int
	Y            int
	Str          string
	ptr          *string
	structMember Impl
	ifMember     Interfc
	mapMember    map[string]interface{}
	sliceMember  []string
}

func TestValidateStruct(t *testing.T) {
	var ptrString string
	impl := Impl{implField: make(map[string]string)}

	testCases := []struct {
		name string
		data Data
		pass bool
	}{
		{
			name: "AllFieldsSet",
			data: Data{
				X:            1,
				Y:            2,
				Str:          "test",
				ptr:          &ptrString,
				structMember: impl,
				ifMember:     Impl{implField: map[string]string{"key": "value"}},
				mapMember:    map[string]interface{}{"key": "value"},
				sliceMember:  []string{"element"},
			},
			pass: true,
		},
		{
			name: "AllFieldsSetZeroValue",
			data: Data{
				X:            0,
				Y:            2,
				Str:          "test",
				ptr:          &ptrString,
				structMember: impl,
				ifMember:     Impl{implField: map[string]string{"key": "value"}},
				mapMember:    map[string]interface{}{"key": "value"},
				sliceMember:  []string{"element"},
			},
			pass: true,
		},
		{
			name: "AllFieldsSetEmptyString",
			data: Data{
				X:            1,
				Y:            2,
				Str:          "",
				ptr:          &ptrString,
				structMember: impl,
				ifMember:     Impl{implField: map[string]string{"key": "value"}},
				mapMember:    map[string]interface{}{"key": "value"},
				sliceMember:  []string{"element"},
			},
			pass: true,
		},
		{
			name: "MissingFields1",
			data: Data{X: 1}, // Only 'X' field set, others are unset
			pass: false,
		},
		{
			name: "AllFieldsSet2",
			data: Data{
				X:           1,
				Str:         "",
				ptr:         &ptrString,
				mapMember:   map[string]interface{}{"key": "value"},
				sliceMember: []string{"element"},
			},
			pass: false,
		},
		{
			name: "AllFieldsSet3",
			data: Data{
				ifMember:    Impl{implField: map[string]string{"key": "value"}},
				mapMember:   map[string]interface{}{"key": "value"},
				sliceMember: []string{"element"},
			},
			pass: false,
		},
		// Add more test cases here as needed
	}

	for _, tc := range testCases {
		err := ValidateStruct(tc.data)

		if tc.pass && err != nil {
			t.Errorf("%s: Expected no error, got: %v", tc.name, err)
		} else if !tc.pass && err == nil {
			t.Errorf("%s: Expected an error, but received none", tc.name)
		}
	}
}
