# Evaluate of Content Value

To get bonus from content creator inflation pool, the content creator need to publish post then get donation. The donation will subtract the coin with most charged coin day first. The coin day consumed will be the input of reputation system, after that a donation power will be calculated and as the input to the Evaluate of Content value. Evaluated result then will be added to a seven-days window. After 7 days, get content bonus and send to content author. Evaluate of Content Value result:
```
 r = A * K * κ * τ * η 
``` 
where A is the donation power; K is the adjustment of cumulative consumption amount on the content; κ is the consumption amount adjustment; τ is the consumption time adjustment, and η is the adjustment according to the number of times of consumption of the same content creator.
```
 K = 1 + 1/(1+e ^(c/1000−5)) 
``` 
Where c is the total consumption amount on the content.
```
 κ = (A)^-0.2 
``` 

```
 τ = 1/(1+e^(Δt/3153600−5)) 
``` 

Where Δt (in seconds) is the consumption time minus content creation time.

```
 η = 1 + 1/(1+e^(n−7)) 
``` 

Where n is the number of times of consumptions of the same content creator.

After the evaluate result added to the window, after 7 days, this evaluate result will be pop from the window and get the portion of inflation from the content creator inflation pool. The equation:
```
 reward = (r/w)*p 
``` 
Where w is total amount in seven-days window and p is the penalty score. The penalty score is calculated as:
```
 p = 1 - min(report/upvote, 1) 
```
