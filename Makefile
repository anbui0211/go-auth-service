# Run application
run-air:
	@air

# Generate the private key
gen-key:
	mkdir -p ./configs/keys
	@openssl genrsa -out ./configs/keys/private.key 2048
	@openssl rsa -in ./configs/keys/private.key -pubout -out ./configs/keys/public.key

# Migrate
migrate:
	@curl -X POST http://localhost:8000/migrate
