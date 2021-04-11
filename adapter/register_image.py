"""
Script to register image via http.
"""

import argparse
import json

import cv2
import requests

parser = argparse.ArgumentParser()
parser.add_argument("image_path")

if __name__ == '__main__':
    args = parser.parse_args()
    img = cv2.imread(args.image_path)
    # encode image as jpeg
    _, img_encoded = cv2.imencode('.jpg', img)
    # send http request with image and receive response
    content_type = 'image/jpeg'
    headers = {'content-type': content_type}
    response = requests.post('http://0.0.0.0:9090/register_image', data=img_encoded.tostring(), headers=headers)
    # decode response
    print(json.loads(response.text))
