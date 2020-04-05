import base64
import requests
import sys


for path in sys.argv[1:]:
    with open(path, "rb") as image_file:
        data = image_file.read()

        b64data = base64.b64encode(data)

        body = {
            "tags": ["asdf", "123"],
            "image": {
                "base64": b64data.decode(),
                "filename": "my-image.jpg"
            }
        }

        response = requests.post("http://localhost:3000/memes", json=body)

        print(response.text)


