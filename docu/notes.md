# Schema-registry operative notes

taking into consideration the following schema

```json
{
  "type": "record",
  "name": "test_topic10",
  "fields": [
    {
      "name": "user",
      "type": "string"
    },

    {
      "name": "password",
      "type": "string",
      "size": 11
    },
    {
      "name": "number",
      "type": [
        "null",
        "double"
      ],
      "default": null,
      "size": 10
    }
  ]
}
```

|      Action      |  Result ||
|-------------|------|------|
|  change `"name": "test_topic10"` to `"name": "test_topic11"` | compatible and creates new version |:heavy_check_mark:|
|  change `"size": 11` to `"size": 12` | compatible and creates new version |:heavy_check_mark:|
|  change `"type": ["null","double"]` to `"type": "double"` | incompatible with an earlier schema |:x:|
|  change `"type": ["null","double"]` to `"type": ["double"]` | incompatible with an earlier schema |:x:|
|  change `"name":"user"` to `"name":"account"` | incompatible with an earlier schema |:x:|
|  add new field `{"name":"account","type":"string"}` | incompatible with an earlier schema |:x:|
|  add `size` attribute to a field `{"name":"user","type":"string"}` | incompatible with an earlier schema |:x:|



changing type from   