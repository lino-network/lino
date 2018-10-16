# Coin day

In Lino Blockchain, each account will have a coin day. It is used to evaluate the donation and allocate bandwidth to user. If someone transfers a certain amount of LINO to you, the coin day for that amount of LINO is 0 when you received them. During 7 days refill window the coin day for that LINO will gradually grow till full. The detail is as following:

When new coin transferred or added to your balance, the coin will be added to your balance immediately, and your pending coin day queue will generate a new pending coin day about this record, which records the start time, end time and amount of LINO. Then your coin day will be refreshed.

When refresh the coin day, it consists of two parts, 1) the coin day is still growing, 2) and the coin day finish grew. First, the blockchain will go through your pending coin day queue. For each pending coin day, if 1) the pending coin day doesn't reach the end time, the blockchain adds the (pending coin day amount) * (now - last evaluate time) / (end time - start time) to your coin day is still growing. 2) the pending coin day reach the end time, we remove the amount of (pending coin day amount) * (last evaluate time - start time) / end time - start time) from your coin day is still growing then add the amount of this pending coin day to your coin day finish grew. After the refresh, the last update time set to now.

When you spend money from your balance, the coin will be removed from your balance immediately, then the blockchain will do following steps:

Refresh the coin day first, follow the above evaluate process, then go through the pending coin day queue from the newest one and check if 1) the LINO to cost is larger than this pending coin day, this pending coin day will be removed from the queue and minus the amount of this pending coin day from your coin day is growing, then minus the amount of this pending coin day and go to next newest pending stake in the queue. 2) the LINO to cost is less than this pending coin day, the blockchain will subtract a portion of the LINO from this pending coin day, then minus the cost * (now - pending stake start time) / (end time - start time) from the coin day is still growing. 3) there is no pending coin day in queue, just minus the LINO from the coin day is grew.


