import pulumi
from pulumi.resource import ResourceOptions
import pulumi_linode as linode
import pulumi_cloudflare as cloudflare


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
    ):
        super().__init__("pkg:index:UT2GameServerLinode", name, None, opts)

        self.server_name = server_name

        zone = cloudflare.get_zone(name=zone_name)
        tags = self.tags = [
            "ut-server",
            f"ut-server-{pulumi.get_stack()}",
            f"ut-server-{pulumi.get_stack()}-{name}",
        ]

        inst = linode.Instance(
            f"{name}-instance",
            label=server_name,
            type=type,
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
            value=inst.ip_address,
            opts=ResourceOptions(parent=self),
        )

        self.ip_address = inst.ip_address

        self.register_outputs({"ipAddress": inst.ip_address})
