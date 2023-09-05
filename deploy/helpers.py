import pulumi
import pulumi_linode as linode
import pulumi_cloudflare as cloudflare


config = pulumi.Config()
stack = pulumi.get_stack()

tags = [
    "ut2-server",
    "ut2-server-" + stack,
]

zone_id = config.require("cfZoneId")
dns_suffix = config.get("dnsSuffix", "kokuei.dev")

_REGION_MAP = {
    "us-central": "dfw",
    "us-ord": "chi",
}


class Server:
    inst: linode.Instance
    record: cloudflare.Record


def get_base_image(label: str):
    images = linode.get_images(filters=[
        linode.GetImagesFilterArgs(name="label", values=[label]),
        linode.GetImagesFilterArgs(name="is_public", values=["false"]),
    ])

    return images.images[0]


def create_server(
    image: pulumi.Input[str],
    id: int,
    region: str,
    type: str = 'g6-nanode-1'
):
    fqdn = f"ut2-{id:02d}.{_REGION_MAP[region]}.{dns_suffix}"

    server = Server()

    inst = server.inst = linode.Instance(
        fqdn,
        label=fqdn,
        type=type,
        region=region,
        image=image,
        tags=tags,
    )

    server.record = cloudflare.Record(
        fqdn,
        zone_id=zone_id,
        name=fqdn,
        value=inst.ip_address,
        type="A",
        ttl=3600,
        proxied=False,
    )

    return server


def generate_export(server: Server):
    return {
        "label": server.inst.label,
        "region": server.inst.region,
        "type": server.inst.type,
        "record": {
            "type": server.record.type,
            "name": server.record.name,
            "value": server.record.value,
        },
    }
