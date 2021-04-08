# import cv2
import numpy as np
import json
import asyncio
import random

from flask import Flask, request, jsonify
from image_checker import ImageChecker

from PIL import Image

from contract.config import config
from contract.download_images import download_images as contract_download_images
from contract.download_images import metadata as download_images_metadata
from contract.listen_images import listen_images as contract_listen_images
from contract.listen_images import metadata as listen_images_metadata
from contract.get_contract import get_contract
from contract.download_image_data import download_image_data
from contract.download_image_data import metadata as download_images_data_metadata
from contract.download_image_data import errors as download_images_data_errors


app = Flask(__name__)
image_checker = ImageChecker()


@app.before_request
def log_request_info():
    app.logger.debug('Headers: %s', request.headers)
    app.logger.debug('Body: %s', request.get_data())


# def load_image(data):
#     nparr = np.fromstring(data, np.uint8)
#     # decode image
#     img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
#     img = Image.fromarray(img)
#     return img
#
#
# @app.route('/register_image', methods=['POST'])
# def register_image():
#     img = load_image(request.data)
#     image_checker.register_new_image(img, 'None')
#
#     # build a response dict to send back to client
#     response = {'message': 'image received.'}
#     # encode response using jsonpickle
#     response = jsonify(response)
#     return response
#
#
# @app.route('/image_score', methods=['POST'])
# def image_score():
#     img = load_image(request.data)
#     score = image_checker.get_image_score(img)
#
#     # build a response dict to send back to client
#     response = {'score': str(score)}
#     # encode response using jsonpickle
#     response = jsonify(response)
#     return response


@app.route('/info', methods=['GET'])
def info():
    return jsonify({
        'lauding_contract': config.CONTRACT_ADDRESS,
        'analysis_network': 'ethereum mainnet',
        'adapter_version': '0.0.1',
        'found_images_count_while_downloading': download_images_metadata['found_images'],
        'downloaded_images_count': download_images_metadata['downloaded_images'],
        'found_images_count_while_watching': listen_images_metadata['found_images'],
        'watched_images_count': listen_images_metadata['downloaded_images'],
        'last_downloaded_urls': download_images_data_metadata['last_urls'],
    })


@app.route('/check', methods=['POST'])
def check():
    body = json.loads(request.data)

    if 'data' not in body:
        return jsonify({'is_succeed': False})

    data = body['data']

    if not isinstance(data, dict):
        return jsonify({'is_succeed': False})

    if 'id' not in data:
        return jsonify({'is_succeed': False})

    image_id = data['id']

    contract = get_contract()

    nftID = int(image_id)

    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    task = asyncio.gather(download_image_data(contract, nftID ))
    loop.run_until_complete(task)

    task_result = task.result()

    if task_result == download_images_data_errors['not_found_by_contract']:
        return {'score': 404}

    if nftID % 2 == 0:
        return jsonify({'score': random.randint(80, 100)})
    else:
        return jsonify({'score': random.randint(0, 40)})


@app.route('/download_images', methods=['GET'])
def download_images():
    contract_download_images()

    return jsonify({'is_succeed': True})


@app.route('/listen_images', methods=['GET'])
def listen_images():
    contract_listen_images()

    return jsonify({'is_succeed': True})


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port='8080', threaded=True)
