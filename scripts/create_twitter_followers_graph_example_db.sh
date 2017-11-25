#!/usr/bin/env bash


function setup_creds {
  echo "You will need to create a Twitter App to get consumer keys."
  echo "Redirecting you to app creation webpage. Make sure you're logged in to Twitter"
  echo "Important: keep your new application open in the browser. You'll be asked by this script for the keys."
  sleep 1

  $(xdg-open "https://apps.twitter.com/app/new")

  read -p "Enter 'Consumer Key' (from 'Keys and Access Tokens' tab): " CONSUMER_KEY
  read -p "Enter 'Consumer Secret': " CONSUMER_SECRET

  echo "Now, you haven't yet authorised this app for your account. Go down to the bottom of the tab you're on and "
  echo "click the button to create an access token"

  sleep 3

  echo "After clicking, you should see 'Access Token' and 'Access Token Secret' appear"
  read -p "Enter 'Access Token': " ACCESS_TOKEN
  read -p "Enter 'Access Token Secret': " ACCESS_TOKEN_SECRET

  echo "Great! One last thing, we need your twitter handle. It starts with '@'"
  read -p "Enter Twitter handle: " TWITTER_HANDLE

  KEY_FILE_CONTENTS="{\"consumer_key\": \"$CONSUMER_KEY\", \"consumer_secret\": \"$CONSUMER_SECRET\", \"access_token\": \"$ACCESS_TOKEN\", \"access_token_secret\": \"$ACCESS_TOKEN_SECRET\", \"handle\": \"$TWITTER_HANDLE\"}"
  echo "****************************************************"
  echo "Creating .gitignore-d private credentials .json file"
  echo "****************************************************"
  echo "$KEY_FILE_CONTENTS" > private_twitter_credentials.json
}

if [ -e "private_twitter_credentials.json" ]
then
  echo "creds file already exists!. moving on to scraping"
else
  setup_creds
fi
