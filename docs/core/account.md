# Account
Account on Lino Blockchain has a global unique username. Except genesis account, user should send register msg to create a new account. To prevent the account spam attack, each new account should have a referrer who pays for the registration fee for it. The genesis account referrer is itself. The registration fee will be added to the developer inflation pool. The register msg should also provides three public keys. After the registration blockchain recognizes user behavior authorization by checking three public keys.

Besides three public keys, each account has following parameters:

- Balance: amount of LINO.
- Coin day: fully charged coin day.
- Number of balance history: Number of balance history.
- Grant public keys: the public key this account grant permission to.
- Sequence number: the number of transaction signed by the account owner.
- Reward: statistic of reward received.
- Pending coin day list: received LINO still charging the coin day, see coin day for more details.