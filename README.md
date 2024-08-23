# Go Auth

#### Generate key:
Note  cd đến thư mục `/configs/keys` để và chạy lệnh bên dưới

```bash
  # Generate private key
  openssl genrsa -out private.key 2048

  # Generate public key from private key
  openssl rsa -in private.key -pubout -out public.key
```
