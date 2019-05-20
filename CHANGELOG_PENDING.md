# Changelog

## v0.3.0
---

*date*

IMPROVEMENTS

  * [tendermint] upgrade to v0.31.5 (shall improve mempool performance and fix a leak issue)
  * [cosmos-sdk] upgrade to v0.34.4
  * [build] remove cleveldb related patches as tendermint/iavl are upgraded, cosmos's patch is required.
  * [build] remove cosmos clelveldb patch as they now support it through build tags.

BUG FIXES

  * [module] xxxxxx