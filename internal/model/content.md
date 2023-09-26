Exhaustive list of social media data model:

#### Users:
- User ID [base58 of public key]
- Username
- Email
- Password
- Profile Picture
- Bio
- Location
- Followers Count
- Following Count
- Posts Count
- Likes Count
- Comments Count
- Joined Date
- Last Active Date
- Verification Status
- Privacy Settings

#### Posts:
- Post ID
- User ID
- Caption
- Media Type (Image, Video, Text)
- Media URL
- Hashtags
- Location
- Created Date
- Updated Date
- Likes Count
- Comments Count
- Shares Count

#### Comments:
- Comment ID
- Post ID
- User ID
- Comment Text
- Created Date
- Updated Date
- Likes Count
- Replies Count

#### Likes:
- Like ID
- Post ID
- User ID
- Created Date
- Updated Date

#### Follows:
- Follow ID
- Follower User ID
- Followed User ID
- Created Date
- Updated Date

#### Notifications:
- Notification ID
- User ID
- Notification Type (Like, Comment, Mention, Follow)
- Post ID
- Comment ID
- Created Date
- Read Status

#### Messages:
- Message ID
- Sender User ID
- Receiver User ID
- Message Text
- Created Date
- Read Status

#### Groups:
- Group ID
- Group Name
- Group Description
- Group Type (Public, Private)
- Group Members Count
- Group Admins Count
- Group Rules
- Group Image URL
- Created Date
- Updated Date

#### Events:
- Event ID
- Event Name
- Event Description
- Event Type (Online, Offline)
- Event Date
- Event Time
- Event Location
- Event Image URL
- Event RSVPs Count
- Event Attendees Count
- Event Host User ID
- Created Date
- Updated Date


<!-- wont be included -->
#### Ads:
- Ad ID
- Ad Title
- Ad Description
- Ad Image URL
- Ad Target Audience
- Ad Budget
- Ad Duration
- Ad Start Date
- Ad End Date
- Ad Impressions Count
- Ad Clicks Count
- Ad Conversions Count
- Ad Cost
- Ad ROI

#### Analytics:
- User Engagement Metrics (Likes, Comments, Shares, Mentions)
- Post Performance Metrics (Reach, Impressions, Engagement Rate)
- Follower Growth Metrics (Followers Gained, Followers Lost)
- Demographic Metrics (Age, Gender, Location)
- Referral Metrics (Traffic Sources, Click-Through Rates)
- Conversion Metrics (Leads Generated, Sales Made)
- Revenue Metrics (Total Revenue, Average Order Value)
- Campaign Metrics (Ad Spend, Ad Performance)

Note: This is just an example of a social media data model and may vary depending on the specific needs of your project.
