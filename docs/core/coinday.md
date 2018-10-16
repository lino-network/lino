# CoinDay

Borrowed the idea from Bitcoin code contributors and other researchers, CoinDay, on Lino Blockchain, is a number that measures how much LINO has been held for much time in an account. CoinDay is used to compute a user's allocated bandwidth and donation power. For example, if a user received 1 LINO on day 0 and 2 LINO on day 1, this user's CoinDay on day 2 is: 1 * 2(days) + 2 * 1(day). When someone transfers a certain amount of LINO to you, the CoinDay for that amount of LINO becomes 0 and starts to grow from the time you receive it. However, CoinDay cannot grow infinitely. There is a 7 days limit for any LINO's CoinDay to grow. The specific calculation is detailed bellow:

When you receive LINO:
The blockchain maintains a CoinDay DQueue to do Lazy Evaluation of your CoinDay. When you receive LINO, a new element will be appended to the queue, which contains the amount of LINO and the start time. Then, the blockchain will go over the queue and update your CoinDay.

CoinDay = Sum(amount of LINO * min(current time - start time, 7 days))

When you transfer LINO to others:
The newest elements in the CoinDay DQueue will be poped out and your CoinDay will be updated.

When you donate LINO to content:
The oldest elements in the CoinDay DQueue will be poped out and your CoinDay will be updated.
