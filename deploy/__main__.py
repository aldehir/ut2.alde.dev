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
region_instance_counts = {}

for provider_config in deployment:
    provider = provider_config["provider"]
    base = provider_config["base"]

    Server = UT2GameServerLinode
    if provider == "digitalocean":
        Server = UT2GameServerDigitalOcean

    for region_config in provider_config["regions"]:
        region_id = region_config["id"]
        region_name = region_config["name"]

        for instance_type in region_config["instances"]:
            # Ensure that regions across providers don't overlay
            count = region_instance_counts.get(region_name, 0)
            count = region_instance_counts[region_name] = count + 1

            instance_name = f"ut2.{region_name}-{count}.{domain}"
            dns_name = f"{region_name}-{count}.{domain}"

            s = Server(
                instance_name,
                region=region_id,
                type=instance_type,
                zone_name=zone_name,
                server_name=dns_name,
                image=base,
                tags=[
                    "ut2-server",
                    f"ut2-server-{environment}",
                    "monitoring",
                ],
            )

            instances.append(s)


pulumi.export("instances", [{"name": x.server_name, "ip": x.ip_address, "tags": x.tags} for x in instances])
