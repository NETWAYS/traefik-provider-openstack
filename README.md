# traefik-provider-openstack

This is a still experimental HTTP provider for [Traefik](https://traefik.io/) to add routers and services for
[OpenStack](https://www.openstack.org/) servers (virtual machines).

## Goals

We want to avoid using many external IPs or expose too many VMs directly to the internet, while being able to manage
TLS/HTTPs termination at a central point.

## Usage

Start up a new instance in Openstack and install [Docker](https://docs.docker.com/engine/install/)
and [Docker-Compose](https://docs.docker.com/compose/install/). Clone the Repo and use the docker-compose command to
build the `traefik-` and `openstack-`container.

> Note: The openstack-container is the GO-Code itself. See the Dockerfile

Configure your credentials for the Openstack-API in the `docker-compose.yml`

Place your certificates inside the `./config/certs` directory

```
$ docker-compose up --build
```

Now you just have to connect to the address for the
Traefik-Dashboard: `https://<address of openstack instance>.<NDD>:8081`

### How it works

This`openstack-traefik-provider`uses the [http-provider](https://doc.traefik.io/traefik/providers/http/) from Traefik
itself in order to discover the specified VM's/Instances inside an Openstack-Project. To receive data in the
JSON-Format, the `openstack-container`(
the actual GO-Code) uses an openstack-GO-client to poll and discover information about the instances inside the
Project-ID from Openstack. This information will be published via HTTP-Response from`http://openstack:8080/traefik` from
the `openstack-container`. Traefik polls every 10sec from this address to update its routes, middlewares, etc.

## License

Copyright (C) 2021 [NETWAYS GmbH](mailto:info@netways.de)

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public
License as published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not,
see <http://www.gnu.org/licenses/>.
