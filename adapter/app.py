import cv2
import numpy as np

from flask import Flask, request, jsonify
from nn_image_checker import NNModelChecker

from PIL import Image

from contract.download_images import download_images as contract_download_images
from contract.listen_images import listen_images as contract_listen_images


app = Flask(__name__)
image_checker = NNModelChecker()


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
    img = load_image(request.data)
    image_checker.add_image_to_storage(img, 'None')

    # build a response dict to send back to client
    response = {'message': 'image received.'}
    # encode response using jsonpickle
    response = jsonify(response)
    return response


@app.route('/image_score', methods=['POST'])
def image_score():
    img = load_image(request.data)
    scores, descriptions = image_checker.find_most_simular_images(img)

    # build a response dict to send back to client
    response = {'scores': str(scores), 'descriptions': str(descriptions)}
    # encode response using jsonpickle
    response = jsonify(response)
    return response


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
