import pulumi
from pulumi.resource import ResourceOptions
import pulumi_linode as linode
import pulumi_cloudflare as cloudflare
import pulumi_digitalocean as digitalocean


class UT2GameServerLinode(pulumi.ComponentResource):
    def __init__(
        self,
        name,
        opts=None,
        region: str = "us-ord",
        type: str = "g6-nanode-1",
        zone_name: str = "",
        server_name: str = "",
        image: str = "",
        tags: list[str] = [],
    ):
        super().__init__("pkg:index:UT2GameServerLinode", name, None, opts)

        self.server_name = server_name
        self.tags = tags

        zone = cloudflare.get_zone(name=zone_name)

        inst = linode.Instance(
            f"{name}-instance",
            label=server_name,
            type=type,
            region=region,
            image=image,
            tags=tags,
            swap_size=2048,
            opts=pulumi.ResourceOptions(parent=self),
        )

        cloudflare.Record(
            f"{name}-record",
            zone_id=zone.id,
            name=server_name,
            type="A",
            value=inst.ip_address,
            opts=ResourceOptions(parent=self),
        )

        self.ip_address = inst.ip_address

        self.register_outputs({"ipAddress": inst.ip_address})


class UT2GameServerDigitalOcean(pulumi.ComponentResource):
    def __init__(
        self,
        name,
        opts=None,
        region: str = "nyc1",
        type: str = "s-vcpu-1gb",
        zone_name: str = "",
        server_name: str = "",
        image: str = "",
        tags: list[str] = [],
    ):
        super().__init__("pkg:index:UT2GameServerDigitalOcean", name, None, opts)

        self.server_name = server_name
        self.tags = tags

        zone = cloudflare.get_zone(name=zone_name)

        inst = digitalocean.Droplet(
            f"{name}-instance",
            name=server_name,
            size=type,
            region=region,
            image=image,
            tags=tags,
            opts=pulumi.ResourceOptions(parent=self),
        )

        cloudflare.Record(
            f"{name}-record",
            zone_id=zone.id,
            name=server_name,
            type="A",
            value=inst.ipv4_address,
            opts=ResourceOptions(parent=self),
        )

        self.ip_address = inst.ipv4_address

        self.register_outputs({"ipAddress": inst.ipv4_address})
