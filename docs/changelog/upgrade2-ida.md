# In-app Digital Assets

## Concept
---
**MiniDollar**: Introduce types.MiniDollar as the unit of consensus
price of Lino. One MiniDollar is 10^(-8) USD. Internally it is a sdk.Int.

**MiniIDA**: Introduce types.MiniIDA as the unit of IDA. One MiniIDA = 10^5 IDA.

## Changes to Post
---

1. Simplify post struct, remove unused fields.
2. Remove view/upvote/report tx.
3. Storage use MustUnmarshalBinaryLengthPrefixed, same as cosmos modules.

### Minors

1. experimental refactor on post structure.

## Changes to Donation
---

1. Reputation use MiniDollar as the unit. For users before upgrade2,
   reputation scores are converted to minidollar.
2. Impact Factor(also called donation power) use MiniDollar as unit as well.
3. Consumption window use MiniDollar.
4. still support donating using LINO.

## BREAKING
---

**CreatePostMsg**: two new fields: signer, createdBy. Both must be provided.
**Codec**: post signbytes are now sorted by sdk.MustSortJSON.
**links**: now are in content.
**CreateAffiliateAccounts**: app need to create affilicate accounts to create posts.
