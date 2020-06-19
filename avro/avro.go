package avro

import (
	"encoding/binary"
	"fmt"

	"github.com/linkedin/goavro"
)

func AvroConfluentSchema(avroSchema string, schemaID int, data []byte) []byte {

	codec, err := goavro.NewCodec(avroSchema)
	if err != nil {
		panic(err)
	}

	var binaryMsg []byte

	// Confluent serialization format version number; currently always 0.
	binaryMsg = append(binaryMsg, byte(0))

	// 4-byte schema ID as returned by Schema Registry
	binarySchemaId := make([]byte, 4)
	binary.BigEndian.PutUint32(binarySchemaId, uint32(schemaID))
	binaryMsg = append(binaryMsg, binarySchemaId...)

	fmt.Println(string(data))

	native, _, err := codec.NativeFromTextual(data)
	if err != nil {
		panic(err)
	}

	binaryData, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		panic(err)
	}

	ConfluentSchemaData := append(binaryMsg, binaryData...)

	return ConfluentSchemaData

}
