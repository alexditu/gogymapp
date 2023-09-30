#!/bin/bash

COUNTRY_NAME="RO"
STATE="Romania"
CITY="Bucharest"
COMPANY_NAME="GymApp Ltd."
OU=""
FQDN="localhost"
EMAIL="ditu.alexandru@gmail.com"

openssl req -newkey rsa:2048 -x509 -days 3650 -nodes -out "./cert.pem" -keyout "./key.pem" >/dev/null 2>&1 <<EOF
${COUNTRY_NAME}
${STATE}
${CITY}
${COMPANY_NAME}
${OU}
${FQDN}
${EMAIL}
EOF