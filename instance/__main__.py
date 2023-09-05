import pulumi
import pulumi_linode as linode
import pulumi_cloudflare as cloudflare

config = pulumi.Config()
zone_id = config.require("cfZoneId")
dns_suffix = config.get("dnsSuffix", "kokuei.dev")

filtered_images = linode.get_images(filters=[
    linode.GetImagesFilterArgs(name="label", values=["rocky8"]),
    linode.GetImagesFilterArgs(name="is_public", values=["false"]),
])

base_image_id = filtered_images.images[0].id

region_map = {
    "us-central": "dfw",
    "us-ord": "chi",
}

instances = []

for count, region in [(1, "us-central"), (1, "us-ord")]:
    for i in range(1, count+1):
        fqdn = f"ut2-{i:02d}.{region_map[region]}.{dns_suffix}"

        instance = linode.Instance(
            fqdn,
            label=fqdn,
            type='g6-nanode-1',
            region=region,
            image=base_image_id,
            tags=["ut2-server", "game-server"],
        )

        record = cloudflare.Record(
            fqdn,
            zone_id=zone_id,
            name=fqdn,
            value=instance.ip_address,
            type="A",
            ttl=3600,
            proxied=False,
        )

        instances.append({
            "label": instance.label,
            "region": instance.region,
            "type": instance.type,
            "record": {
                "type": record.type,
                "name": record.name,
                "value": record.value,
            },
        })

pulumi.export('instances', instances)
