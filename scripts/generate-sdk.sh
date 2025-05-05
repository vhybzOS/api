#!/bin/bash

# Install OpenAPI Generator if not already installed
if ! command -v openapi-generator-cli &> /dev/null; then
    echo "Installing OpenAPI Generator..."
    npm install @openapitools/openapi-generator-cli -g
fi

# Generate Swagger documentation
echo "Generating Swagger documentation..."
swag init -g main.go

# Generate TypeScript SDK
echo "Generating TypeScript SDK..."
openapi-generator-cli generate \
    -i docs/swagger.json \
    -g typescript-axios \
    -o sdk/typescript \
    --additional-properties=npmName=@vhybz/api-client,npmVersion=1.0.0,withInterfaces=true

# Install dependencies and build the SDK
echo "Building SDK..."
cd sdk/typescript
npm install
npm run build

echo "SDK generation complete!" 