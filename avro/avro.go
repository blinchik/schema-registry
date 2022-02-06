package avro

import (
	"encoding/binary"
	"log"

	"github.com/linkedin/goavro"
)

// SerializeMessageConfluentAvro will serialize a Kafka message according to supplied Avro Schema
func SerializeMessageConfluentAvro(avroSchema string, schemaID int, data []byte) []byte {

	codec, err := goavro.NewCodec(avroSchema)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	binaryData, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		log.Fatal(err)
	}

	ConfluentSchemaData := append(binaryMsg, binaryData...)

	return ConfluentSchemaData

}

// DeserializeMessageConfluentAvro will deserialize an Avro Kafka message according to supplied Avro Schema
func DeserializeMessageConfluentAvro(avroSchema string, schemaID int, data []byte) []byte {

	codec, err := goavro.NewCodec(avroSchema)
	if err != nil {
		log.Fatal(err)
	}

	// remove the magic bytes from data
	native, _, err := codec.NativeFromBinary(data[5:])
	if err != nil {
		log.Fatal(err)
	}

	textual, err := codec.TextualFromNative(nil, native)
	if err != nil {
		log.Fatal(err)
	}

	return textual
}
