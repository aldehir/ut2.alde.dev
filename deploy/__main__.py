import pulumi
import pulumi_linode as linode
import pulumi_cloudflare as cloudflare

import typing

from gameservers import UT2GameServerLinode, UT2GameServerDigitalOcean

config = pulumi.Config()
stack = pulumi.get_stack()

environment: typing.Optional[str] = config.get("environment")
zone_name: typing.Optional[str] = config.get("zone_name")
domain: typing.Optional[str] = config.get("domain")
deployment: typing.Optional[dict[str, typing.Any]] = config.get_object("deployment")

assert environment
assert zone_name
assert domain
assert deployment

instances = []

for provider, provider_config in deployment.items():
    Server = UT2GameServerLinode
    if provider == "digitalocean":
        Server = UT2GameServerDigitalOcean

    base = provider_config["base"]

    for region_id, region_config in provider_config["regions"].items():
        region_name = region_config["name"]

        for i, instance_type in enumerate(region_config["instances"], 1):
            instance_name = f"ut2.{region_name}-{i}.{domain}"

            s = Server(
                instance_name,
                region=region_id,
                type=instance_type,
                zone_name=zone_name,
                server_name=instance_name,
                image=base,
                tags=[
                    "ut2-server",
                    f"ut2-server-{environment}"
                ],
            )

            instances.append(s)


pulumi.export("instances", [{"name": x.server_name, "ip": x.ip_address, "tags": x.tags} for x in instances])
