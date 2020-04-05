import json
import requests
import os
import urllib
from urllib.parse import urlparse
import base64


def get_top_images_from_subreddit(subreddit, count=20):
    url = f"https://www.reddit.com/r/{subreddit}/top/.json"

    headers = {"User-Agent": "memecan-scraper"}
    response = requests.get(url, params={"limit": count, "t": "day"}, headers=headers)

    json_data = response.json()
    if json_data.get("error"):
        raise Exception(json_data["message"])


    for post in json_data.get("data", {}).get("children", {}):
        image_url = post["data"]["url"]
        print(image_url)

        image_name = os.path.basename(urlparse(image_url).path)

        r = requests.get(image_url)
        if r.status_code == 200:
            b64_image = base64.b64encode(r.content)

            body = {
                "tags": ["reddit", "me_irl"],
                "image": {
                    "base64": b64_image.decode(),
                    "filename": image_name
                }
            }

            response = requests.post("http://localhost:3000/memes", json=body)


get_top_images_from_subreddit("me_irl", count=10)
