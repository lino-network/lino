# Post

## publish

The account on blockchain can publish a post with title and content. The post on blockchain can receive donation and get content bonus. The post can be uniquely identified by permlink, which consists of author name and post id. Post can be updated, liked, viewed, donated, reported, upvoted or deleted. Deleted post can’t accept donation and can’t be updated. The post can be deleted by post author or be censored by governance. Comment is a post with parent. Repost is a post with source.

## repost

To help distribute the content, the content creator can set redistribution split rate to encourage user redistribute the content. Lino blockchain encourages people to share and distribute the content. If people make donation to a repost, the donation will be splitted by source post’s redistribution split rate. For example, if source post’s redistribution split rate set to 5%, the donation to the repost will send 95% of the donation to the source post and repost can keep 5% of the donation. The donation to the source post still go to source post author’s account. The repost of a repost will assign the source post to the original source post. For example, If B repost A then we have the structure A -> B. If C repost B then in blockchain the struct will be changed from A -> B -> C to A -> C.

## donation

When user donate to a post the donation will be added to the post’s donation list. The donation will be divided to two parts. The 90% donation will be added to author’s balance directly and it will also be added to donation list with type direct deposit. The 10% friction will be added to daily consumption pool, which will distribute to all locked LINO holder. The donation will cost the fully charged coin day first. The coin day spent on this donation will be evaluated in reputation system. The donation power get from reputation system will then go through evaluate of content value then the result will be added to a 7 days window. After the window the evaluate result is used to share the content creator inflation pool, the shared bonus will be added to the post donation list at the end. The donation to a post will also add upvote score to the post.

## report and upvote

The report and upvote is calculated based on user’s reputation. The upvote will be added to a post when user donate to a post. Report is restrict to once a hour. Based on total upvote reputation and report reputation, a post will have a penalty score. The penalty score will affect the final bonus distribution.

## deleted post

When post is deleted, it can’t accept donation anymore, and it can’t get content bonus from reward pool neither. Deleted post can’t be updated. Currently the post can be deleted by the author or content censorship.
