// Copyright © 2023 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opencdc

import (
	"fmt"

	opencdcv1 "github.com/conduitio/conduit-commons/proto/opencdc/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	var cTypes [1]struct{}
	_ = cTypes[int(OperationCreate)-int(opencdcv1.Operation_OPERATION_CREATE)]
	_ = cTypes[int(OperationUpdate)-int(opencdcv1.Operation_OPERATION_UPDATE)]
	_ = cTypes[int(OperationDelete)-int(opencdcv1.Operation_OPERATION_DELETE)]
	_ = cTypes[int(OperationSnapshot)-int(opencdcv1.Operation_OPERATION_SNAPSHOT)]
}

// -- From Proto To OpenCDC ----------------------------------------------------

// FromProto takes data from the supplied proto object and populates the
// receiver. If the proto object is nil, the receiver is set to its zero value.
// If the function returns an error, the receiver could be partially populated.
func (r *Record) FromProto(proto *opencdcv1.Record) error {
	if proto == nil {
		*r = Record{}
		return nil
	}

	var err error
	r.Key, err = dataFromProto(proto.Key)
	if err != nil {
		return fmt.Errorf("error converting key: %w", err)
	}

	if proto.Payload != nil {
		err := r.Payload.FromProto(proto.Payload)
		if err != nil {
			return fmt.Errorf("error converting payload: %w", err)
		}
	} else {
		r.Payload = Change{}
	}

	r.Position = proto.Position
	r.Metadata = proto.Metadata
	r.Operation = Operation(proto.Operation)
	return nil
}

// FromProto takes data from the supplied proto object and populates the
// receiver. If the proto object is nil, the receiver is set to its zero value.
// If the function returns an error, the receiver could be partially populated.
func (c *Change) FromProto(proto *opencdcv1.Change) error {
	if proto == nil {
		*c = Change{}
		return nil
	}

	var err error
	c.Before, err = dataFromProto(proto.Before)
	if err != nil {
		return fmt.Errorf("error converting before: %w", err)
	}

	c.After, err = dataFromProto(proto.After)
	if err != nil {
		return fmt.Errorf("error converting after: %w", err)
	}

	return nil
}

func dataFromProto(proto *opencdcv1.Data) (Data, error) {
	if proto == nil {
		return nil, nil //nolint:nilnil // This is the expected behavior.
	}

	switch v := proto.Data.(type) {
	case *opencdcv1.Data_RawData:
		return RawData(v.RawData), nil
	case *opencdcv1.Data_StructuredData:
		return StructuredData(v.StructuredData.AsMap()), nil
	case nil:
		return nil, nil //nolint:nilnil // This is the expected behavior.
	default:
		return nil, ErrInvalidProtoDataType
	}
}

// -- From OpenCDC To Proto ----------------------------------------------------

// ToProto takes data from the receiver and populates the supplied proto object.
// If the function returns an error, the proto object could be partially
// populated.
func (r Record) ToProto(proto *opencdcv1.Record) error {
	if r.Key != nil {
		if proto.Key == nil {
			proto.Key = &opencdcv1.Data{}
		}
		err := r.Key.ToProto(proto.Key)
		if err != nil {
			return fmt.Errorf("error converting key: %w", err)
		}
	} else {
		proto.Key = nil
	}

	if proto.Payload == nil {
		proto.Payload = &opencdcv1.Change{}
	}
	err := r.Payload.ToProto(proto.Payload)
	if err != nil {
		return fmt.Errorf("error converting payload: %w", err)
	}

	proto.Position = r.Position
	proto.Metadata = r.Metadata
	proto.Operation = opencdcv1.Operation(r.Operation)
	return nil
}

// ToProto takes data from the receiver and populates the supplied proto object.
// If the function returns an error, the proto object could be partially
// populated.
func (c Change) ToProto(proto *opencdcv1.Change) error {
	if c.Before != nil {
		if proto.Before == nil {
			proto.Before = &opencdcv1.Data{}
		}
		err := c.Before.ToProto(proto.Before)
		if err != nil {
			return fmt.Errorf("error converting before: %w", err)
		}
	} else {
		proto.Before = nil
	}

	if c.After != nil {
		if proto.After == nil {
			proto.After = &opencdcv1.Data{}
		}
		err := c.After.ToProto(proto.After)
		if err != nil {
			return fmt.Errorf("error converting after: %w", err)
		}
	} else {
		proto.After = nil
	}

	return nil
}

// ToProto takes data from the receiver and populates the supplied proto object.
func (d RawData) ToProto(proto *opencdcv1.Data) error {
	protoRawData, ok := proto.Data.(*opencdcv1.Data_RawData)
	if !ok {
		protoRawData = &opencdcv1.Data_RawData{}
		proto.Data = protoRawData
	}
	protoRawData.RawData = d
	return nil
}

// ToProto takes data from the receiver and populates the supplied proto object.
func (d StructuredData) ToProto(proto *opencdcv1.Data) error {
	protoStructuredData, ok := proto.Data.(*opencdcv1.Data_StructuredData)
	if !ok {
		protoStructuredData = &opencdcv1.Data_StructuredData{}
		proto.Data = protoStructuredData
	}
	data, err := structpb.NewStruct(d)
	if err != nil {
		return fmt.Errorf("could not convert structured data to proto: %w", err)
	}
	protoStructuredData.StructuredData = data
	return nil
}
