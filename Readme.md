Reverse API Gateway
======

Expose external HTTP(S) APIs as insecure local services.

Free client code from the burden of checking the SSL certificates and handling authentication details.

Features
---
- Configurable from yaml file
- Multiple (outgoing) endpoints and credentials
- Support for OAuth 1a
- Support for OAuth 2 (TODO)
- Live updating security tokens (TODO)

Limitations
---
- No listening on HTTPS
- Single listen address/port
- No credentials checking for incoming connections

Maybe:
---
- Transparent proxy
- Work in serverless environment
