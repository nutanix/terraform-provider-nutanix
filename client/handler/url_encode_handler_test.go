package handler

import (
	"testing"
)

type InputService3TestShapeStructType struct {
	_ struct{} `type:"structure"`

	InstanceIds []*string `locationName:"InstanceId" locationNameList:"InstanceId" type:"list"`
}

func TestURLEncodeMarshalHander(t *testing.T) {
	expOut := "Action=DescribeInstances&InstanceId.1=i-76536489&Version=2017-12-15"

	ID1 := "i-76536489"
	input := &InputService3TestShapeStructType{
		InstanceIds: []*string{&ID1},
	}

	res, err := URLEncodeMarshalHander(input, "DescribeInstances", "2017-12-15")
	if err != nil {
		t.Fatalf("Got error: %s", err)
	}
	if res != expOut {
		t.Fatalf("Error Marshal: Got:(%s), Have(%s)", expOut, res)
	}
}
