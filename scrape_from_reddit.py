import json
import requests
import os
import urllib
from urllib.parse import urlparse


def get_top_images_from_subreddit(subreddit, count=20):
    url = f"https://www.reddit.com/r/{subreddit}/top/.json"

    response = requests.get(url, params={"limit": count, "t": "day"})

    json_data = response.json()
    if json_data.get("error"):
        raise Exception(json_data["message"])


    for post in json_data.get("data", {}).get("children", {}):
        image_url = post["data"]["url"]
        print(image_url)

        image_name = os.path.basename(urlparse(image_url).path)
        save_path = os.path.join("reddit", image_name)

        r = requests.get(image_url, stream=True)
        if r.status_code == 200:
            with open(save_path, 'wb') as f:
                for chunk in r.iter_content(1024):
                    f.write(chunk)

get_top_images_from_subreddit("me_irl", count=20)
