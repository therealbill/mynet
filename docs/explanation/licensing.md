---
title: "Licensing"
description: "License model for the Mynet marketplace and its plugins"
weight: 3
---

# Licensing

The Mynet marketplace uses a custom license that addresses a specific tension: how to keep individual plugins freely usable while preventing the collection from being copied wholesale into a competing marketplace. Standard open-source licenses (MIT, Apache, GPL) do not distinguish between "use this plugin" and "clone this entire marketplace," so a custom license was necessary.

## The core tension

Open-source licensing typically treats software as a single unit. A user either has the right to redistribute or they do not. But a plugin marketplace is structurally different -- it is both a collection of independent components and a curated product in its own right. The individual plugins have value on their own, but the curation, organization, and completeness of the marketplace is also valuable work that deserves protection.

The Mynet license resolves this by granting broad rights at the individual plugin level while restricting aggregation at the collection level.

## How the license works

The license defines four key provisions, each addressing a different use case.

**Aggregation prohibited (Section 1).** The software cannot be bundled into multi-plugin collections, aggregations, or repositories that gather components from multiple sources. This directly prevents the scenario where a third party copies all 19 plugins into their own repository and presents them as their own marketplace. The prohibition is broad -- it covers archives, packages, distributions, and any other form of collection, whether or not the software is modified.

**Marketplace reference permitted (Section 2).** External marketplaces and plugin registries can list Mynet plugins in their catalogs. The requirement is that listings reference the original repository (https://github.com/therealbill/mynet) and fetch the software from the original source at install time. No source code, documentation, or other files may be stored in the external marketplace itself. This provision enables ecosystem integration -- a Claude Code plugin directory could list Mynet plugins and link users to the source -- without allowing redistribution.

**Individual redistribution (Section 3).** Any single plugin can be redistributed on its own, in source or binary form, provided that the license and copyright notice are included, the original repository URL is displayed prominently, and any modifications are clearly marked. A user who takes `sql-pro` and adapts it for their team is fully permitted to share their version, as long as it is shared individually and not bundled into a competing collection.

**Attribution (Section 4).** All uses require the copyright notice, full license text, repository link, and identification of the original author. Attribution ensures that even permitted uses maintain a trail back to the source.

## Why this model

The license balances four interests:

| Interest | How the license addresses it |
|----------|------------------------------|
| Individual users | Full freedom to use, modify, and share individual plugins |
| Organizations | Can adopt and adapt plugins internally without restriction |
| Marketplace platforms | Can reference and link to plugins without storing source |
| Original author | Retains control of the collection; aggregation is prohibited |

The alternative approaches each sacrifice something important. A permissive license (MIT) would allow wholesale cloning. A copyleft license (GPL) would restrict how organizations use plugins internally. A proprietary license would prevent individual sharing entirely. The custom license threads the needle by distinguishing between individual plugin use (permissive) and collection aggregation (restricted).

## Individual plugin licenses

Plugins may include their own `LICENSE` files with additional terms. The root license applies to the marketplace structure and to any plugins that do not carry their own license. When a plugin has its own license, both the root license and the plugin-specific license apply -- the root license governs the plugin's inclusion in the marketplace, while the plugin license may add or refine terms for the plugin's own content.

## Related

The full license text is in the root [`LICENSE`](https://github.com/therealbill/mynet/blob/master/LICENSE) file of the repository.
