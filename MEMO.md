
## 自己証明書
openssl req -x509 -newkey rsa:2048 -nodes `
  -keyout server.key `
  -out server.crt `
  -days 365 `
  -subj "/CN=localhost" `
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

