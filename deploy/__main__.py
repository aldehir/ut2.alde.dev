import pulumi
import pulumi_linode as linode
import pulumi_cloudflare as cloudflare

from helpers import get_base_image, create_server, generate_export

stack = pulumi.get_stack()

image = get_base_image("rocky8")

servers = []

if stack == "production":
    servers += [
        create_server(image.id, 1, "us-ord", "g6-dedicated-2")
    ]

elif stack == "dev":
    servers += [
        create_server(image.id, 1, "us-ord", "g6-dedicated-2")
    ]

pulumi.export("servers", [generate_export(x) for x in servers])
