# Reputation

In Lino Blockchain, the content creator can get content bonus from content creator inflation pool, which is huge incentive for content creator. However, someone may make fake donation to his own content to get bonus from the inflation pool. To prevent bad content stealing bonus, the account on blockchain will have a reputation score with it. The reputation score consists of two parts: free score and customer score.

The free score comes from stake in. A certain percentage of free score will be added to user’s reputation if user stake in and removed once user stake out.

The default customer score is 1 LINO for all accounts. To increase the customer score, one should be in order front of donors to daily top N content. The total donation power a user can spend on a post is limited by his reputation.

For each donation, based on user’s reputation, it will add sum of the reputation of the post. For each report, it will minus sum of the reputation of the post. The penalty score of a post is negative of sum of reputation divided by predefined hard cap. The penalty score will affect the content bonus get from inflation pool.  
