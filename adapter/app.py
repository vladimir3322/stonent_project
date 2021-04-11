import cv2
import numpy as np
import json
import asyncio
import random
import os
from adapter import Adapter


from flask import Flask, request, jsonify
from nn_image_checker import NNModelChecker

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

from image_manager import ImageManager

app = Flask(__name__)
image_checker = NNModelChecker()
image_manager = ImageManager(image_checker)


@app.before_request
def log_request_info():
    app.logger.debug('Headers: %s', request.headers)
    app.logger.debug('Body: %s', request.get_data())



def load_image(data):
    nparr = np.fromstring(data, np.uint8)
    # decode image
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    img = Image.fromarray(img)
    return img


@app.route('/register_image', methods=['POST'])
def register_image():
    image_manager.register_new_images()

    img = load_image(request.data)
    image_checker.add_image_to_storage(img, 'None')

    # build a response dict to send back to client
    response = {'message': 'image received.'}
    # encode response using jsonpickle
    response = jsonify(response)
    return response


@app.route('/image_score', methods=['POST'])
def image_score():
    image_manager.register_new_images()

    img = load_image(request.data)
    scores, descriptions = image_checker.find_most_simular_images(img)

    # build a response dict to send back to client
    response = {'scores': str(scores), 'descriptions': str(descriptions)}
    # encode response using jsonpickle
    response = jsonify(response)
    return response


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


@app.route('/download_images', methods=['GET'])
def download_images():
    contract_download_images()
    return jsonify({'is_succeed': True})


@app.route('/listen_images', methods=['GET'])
def listen_images():
    contract_listen_images()
    return jsonify({'is_succeed': True})


@app.route('/load_image', methods=['POST'])
def load_image():
    body = json.loads(request.data)

    if 'id' not in body:
        return jsonify({'is_succeed': False})

    image_id = body['id']

    if not os.path.exists(config.DATA_PATH):
        os.mkdir(config.DATA_PATH)
    if not os.path.exists(config.SOURCE_PATH):
        os.mkdir(config.SOURCE_PATH)

    contract = get_contract()
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    data_ipfs_url = contract.functions.uri(image_id).call()
    task = asyncio.gather(download_image_data(data_ipfs_url, image_id))
    loop.run_until_complete(task)

    return jsonify({'is_succeed': True})


@app.route('/check', methods=['POST'])
def call_adapter():
    body = json.loads(request.data)
    adapter = Adapter(body, image_checker, image_manager)
    return jsonify(adapter.result)


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port='9090', threaded=True)
