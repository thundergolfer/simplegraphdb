import json
import tweepy
import os

this_files_path = os.path.dirname(os.path.realpath(__file__))
private_twitter_creds_path = "private_twitter_credentials.json"
output_filepath = 'your_twitter_example_db.json'

MAX_USER_PER_CALL = 100

# truncate file
with open(output_filepath, 'w') as f:
    pass


def process_group(subject, prop, group):
    print("Writing file to {}".format(output_filepath))
    with open(output_filepath, 'a') as f:
        for user in group:
            triple = {"subject": subject, "prop": prop, "object": user}
            db["triples"].append(triple)
            json.dump(triple, f)
            f.write('\n')
            print("Processed: {}".format(triple))


with open(os.path.join(this_files_path, private_twitter_creds_path), 'r') as f:
    keys = json.load(f)

consumer_key = keys["consumer_key"]
consumer_secret = keys["consumer_secret"]
access_token = keys["access_token"]
access_token_secret = keys["access_token_secret"]
twitter_handle = keys["handle"] if keys["handle"][0] != '@' else keys["handle"][1:]

auth = tweepy.auth.OAuthHandler(consumer_key, consumer_secret)
auth.set_access_token(access_token, access_token_secret)
api = tweepy.API(auth, wait_on_rate_limit=True, wait_on_rate_limit_notify=True)

if api.verify_credentials:
    print("Successfully authenticated!")
else:
    print("There was a mistake made when receiving credentials. Remove 'private_twitter_credentials.json' and try the script again")

db = {"triples": []}

# make a list because we can't re-use a tweepy.Cursor
your_friends = [user for user in tweepy.Cursor(api.friends, screen_name=twitter_handle).items()]
process_group(
    twitter_handle,
    "follows",
    [user.screen_name for user in your_friends]
)

who_follows_you = tweepy.Cursor(api.followers, screen_name=twitter_handle).items()
process_group(
    twitter_handle,
    "is_followed_by",
    [user.screen_name for user in who_follows_you]
)

# explore one level deeper into network
for user in your_friends:
    their_friends = list(tweepy.Cursor(api.friends_ids, screen_name=user.screen_name).items())
    process_group(
        user.screen_name,
        "follows",
        [u.screen_name for u in api.lookup_users(user_ids=their_friends[:MAX_USER_PER_CALL])]
    )
    who_follows_them = list(tweepy.Cursor(api.followers_ids, screen_name=user.screen_name).items())
    process_group(
        user.screen_name,
        "is_followed_by",
        [u.screen_name for u in api.lookup_users(user_ids=who_follows_them[:MAX_USER_PER_CALL])]
    )

print("done!")
