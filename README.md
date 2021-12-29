# git commit -m "first commit"

## Documentation
https://www.keycloak.org/getting-started

## Starting with docker
Run: `docker run -p 8080:8080 -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin quay.io/keycloak/keycloak:15.1.0`
Then enter in: `http://localhost:8080/auth`
- User: 'admin'
- Password: 'admin'

## Starting with docker compose
Run: `cd theme-docker && docker-compose up -d`

## Changing the Theme
- Enter in the running docker container
  > docker exec -it keycloak bash

- Checking the themes
  > cd /opt/jboss/keycloak/themes/