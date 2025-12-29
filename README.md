# oapit

OpenAPI 3.0 CLI toolkit for validating OpenAPI schemas and payloads against them.

## Features

- **Schema Validation**: Validate OpenAPI 3.0 schema files for correctness
- **Payload Validation**: Validate JSON payloads against schema components
- **Flexible Options**: Control validation behavior with various flags
- **External References**: Support for resolving external schema references
- **Request/Response Validation**: Validate payloads as either requests or responses

## Install

### Pre-built Binaries

Pre-built binaries are available on the [releases page](https://github.com/karupanerura/oapit/releases).

```bash
VERSION=0.0.3
curl -sfLO https://github.com/karupanerura/oapit/releases/download/v${VERSION}/oapit_${VERSION}_$(go env GOOS)_$(go env GOARCH).tar.gz
tar zxf oapit_${VERSION}_$(go env GOOS)_$(go env GOARCH).tar.gz
install -m 0755 oapit /usr/local/bin/
rm oapit oapit_${VERSION}_$(go env GOOS)_$(go env GOARCH).tar.gz
```

### From Source

```bash
go install github.com/karupanerura/oapit@latest
```

## Usage

### Basic Commands

#### Validate OpenAPI Schema

Validate an OpenAPI 3.0 schema file:

```bash
oapit validate schema -f schema.yml
```

#### Validate JSON Payload

Validate a JSON payload against a schema component:

```bash
oapit validate payload -f schema.yml -s MySchema payload.json
```

### Command Reference

#### Global Flags

- `-f, --schema-file`: Path to OpenAPI 3.0 schema file (required)
- `--[no-]resolve-reference`: Enable/disable resolution of external references (default: enabled)

#### `validate schema` Command

Validates an OpenAPI 3.0 schema file for correctness.

**Flags:**
- `--allow-extra-sibling-fields`: Allow extra sibling fields for specified schema properties
- `--[no-]validate-examples`: Enable/disable validation of examples in the schema (default: enabled)
- `--[no-]validate-defaults`: Enable/disable validation of default values (default: enabled)
- `--[no-]validate-patterns`: Enable/disable validation of regex patterns (default: enabled)

**Example:**

```bash
# Validate schema with all checks enabled
oapit validate schema -f openapi.yml

# Validate schema without example validation
oapit validate schema -f openapi.yml --no-validate-examples

# Allow extra sibling fields for specific properties
oapit validate schema -f openapi.yml --allow-extra-sibling-fields=x-custom,x-internal
```

#### `validate payload` Command

Validates a JSON payload against a schema component.

**Arguments:**
- `<payload>`: Path to JSON payload file (use `-` for stdin)

**Flags:**
- `-s, --schema`: Schema component name from the OpenAPI spec (required)
- `--as`: Validate as `request` or `response` (default: `request`)

**Examples:**

```bash
# Validate a request payload
oapit validate payload -f openapi.yml -s UserCreateRequest user.json

# Validate a response payload
oapit validate payload -f openapi.yml -s UserResponse --as response user-response.json

# Validate from stdin
cat payload.json | oapit validate payload -f openapi.yml -s MySchema -
```

## Examples

### Example OpenAPI Schema

```yaml
openapi: 3.0.0
info:
  title: Example API
  version: 1.0.0
components:
  schemas:
    User:
      type: object
      required:
        - id
        - name
      properties:
        id:
          type: integer
        name:
          type: string
        email:
          type: string
          format: email
```

### Validate the Schema

```bash
oapit validate schema -f api-spec.yml
```

### Validate a User Payload

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

```bash
oapit validate payload -f api-spec.yml -s User user.json
```

## Development

### Build

```bash
go build -o oapit main.go
```

### Run Tests

```bash
go test ./...
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
