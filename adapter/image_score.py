"""
Script to score image via http.
"""

import argparse
import json

from skimage import io
import cv2
import requests

parser = argparse.ArgumentParser()
parser.add_argument("image_path")

if __name__ == '__main__':
    args = parser.parse_args()
    img = io.imread(args.image_path)
    # encode image as jpeg
    _, img_encoded = cv2.imencode('.jpg', img)
    # send http request with image and receive response
    content_type = 'image/jpeg'
    headers = {'content-type': content_type}
    response = requests.post('http://0.0.0.0:9090/image_score', data=img_encoded.tostring(), headers=headers)
    # decode response
    print(json.loads(response.text))
