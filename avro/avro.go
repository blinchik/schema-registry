package avro

import (
	"encoding/binary"

	"github.com/linkedin/goavro"
)

// SerializeMessageConfluentAvro will serialize a Kafka message according to supplied Avro Schema
func SerializeMessageConfluentAvro(avroSchema string, schemaID int, data []byte) (confluentSchemaData []byte, err error) {

	codec, err := goavro.NewCodec(avroSchema)
	if err != nil {
		return
	}

	var binaryMsg []byte

	// Confluent serialization format version number; currently always 0.
	binaryMsg = append(binaryMsg, byte(0))

	// 4-byte schema ID as returned by Schema Registry
	binarySchemaId := make([]byte, 4)
	binary.BigEndian.PutUint32(binarySchemaId, uint32(schemaID))
	binaryMsg = append(binaryMsg, binarySchemaId...)

	native, _, err := codec.NativeFromTextual(data)
	if err != nil {
		return
	}

	binaryData, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		return
	}

	confluentSchemaData = append(binaryMsg, binaryData...)

	return confluentSchemaData, err

}

// DeserializeMessageConfluentAvro will deserialize an Avro Kafka message according to supplied Avro Schema
func DeserializeMessageConfluentAvro(avroSchema string, schemaID int, data []byte) (textual []byte, err error) {

	codec, err := goavro.NewCodec(avroSchema)
	if err != nil {
		return
	}

	// remove the magic bytes from data
	native, _, err := codec.NativeFromBinary(data[5:])
	if err != nil {
		return
	}

	textual, err = codec.TextualFromNative(nil, native)
	if err != nil {
		return
	}

	return textual, err
}
