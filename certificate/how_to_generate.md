Brug openssl:
```bash
choco install openssl
```


lav key og certificate med command:
(credit https://raymii.org/s/tutorials/OpenSSL_generate_self_signed_cert_with_Subject_Alternative_name_oneliner.html)
```bash
openssl req -nodes -x509 -sha256 -newkey rsa:4096 -keyout priv.key -out server.crt -days 356 -subj "/C=DK/ST=Copenhagen/L=Copenhagen/O=Me/OU=mpc/CN=localhost" -addext "subjectAltName = DNS:localhost,IP:0.0.0.0"
```