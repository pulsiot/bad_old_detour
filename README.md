# detour

deTour is a small proxy server which redirects.

## Config File.

deTour uses `config.yaml` file to read the routes.

## Dependencies

Please note that the program assumes you have the necessary config.yaml, cert.pem, and key.pem files in the same directory as the program.

## Create SSL certificates

```
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 3650 -nodes -subj "/C=IN/ST=Kerala/L=Kochi/O=Detour/OU=DetourProxy/CN=localhost"
```
