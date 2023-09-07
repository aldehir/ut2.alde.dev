import pulumi
import pulumi_linode as linode
import pulumi_cloudflare as cloudflare

from gameservers import UT2GameServerLinode

config = pulumi.Config()
stack = pulumi.get_stack()

linode_image = "private/21473290"

chi01 = UT2GameServerLinode(
    f"ut2-01-chi-staging",
    region="us-ord",
    type="g6-dedicated-2",
    zone_name="kokuei.dev",
    server_name=f"ut2-01.chi.staging.kokuei.dev",
    image=linode_image,
)

pulumi.export(chi01.server_name, {"ip": chi01.ip_address, "tags": chi01.tags})

dfw01 = UT2GameServerLinode(
    f"ut2-01-dfw-staging",
    region="us-central",
    type="g6-dedicated-2",
    zone_name="kokuei.dev",
    server_name=f"ut2-01.dfw.staging.kokuei.dev",
    image=linode_image,
)

pulumi.export(dfw01.server_name, {"ip": dfw01.ip_address, "tags": dfw01.tags})
