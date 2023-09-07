import pulumi
import pulumi_linode as linode
import pulumi_cloudflare as cloudflare

from gameservers import UT2GameServerLinode

config = pulumi.Config()
stack = pulumi.get_stack()

linode_image = "private/21473290"

s = UT2GameServerLinode(
    f"ut2-01-chi-staging",
    region="us-ord",
    type="g6-dedicated-2",
    zone_name="kokuei.dev",
    server_name=f"ut2-01.chi.staging.kokuei.dev",
    image=linode_image,
)

pulumi.export(s.server_name, {"ip": s.ip_address, "tags": s.tags})
