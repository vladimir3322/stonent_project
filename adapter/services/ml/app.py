import cv2
import numpy as np
import loader
from threading import Thread, Lock
from PIL import Image
from flask import Flask, request, jsonify
from nn_image_checker import NNModelChecker
from image_manager import ImageManager


app = Flask(__name__)
image_checker = NNModelChecker()
image_manager = ImageManager(image_checker)
mutex = Lock()


def register_images():
    image_manager.register_new_images()


@app.route('/check', methods=['GET'])
def call_adapter():
    def get_result(code, data=None, error=None):
        return jsonify({
            'jobRunID': f'{contract_address}_{nft_id}',
            'data': data,
            'error': error,
            'statusCode': code,
        })

    print('KEKEKEKEKEKEKEKEKEKE, 1', flush=True)
    mutex.acquire()
    contract_address = request.args.get('contract_address')
    nft_id = request.args.get('nft_id')

    print('KEKEKEKEKEKEKEKEKEKE, 2', flush=True)

    if not contract_address:
        return get_result(400, None, 'Invalid contract address')
    if not nft_id or not nft_id.isdigit():
        return get_result(400, None, 'Invalid nft id')

    print('KEKEKEKEKEKEKEKEKEKE, 3', flush=True)

    image_source, image_source_error = loader.get_image_source(contract_address, nft_id)

    print('KEKEKEKEKEKEKEKEKEKE, 4', flush=True)

    if not image_source or image_source_error:
        return get_result(500, None, image_source_error)

    print('KEKEKEKEKEKEKEKEKEKE, 5', flush=True)

    np_image_source = np.fromstring(image_source, np.uint8)
    image = cv2.imdecode(np_image_source, cv2.IMREAD_COLOR)
    image = Image.fromarray(image)

    print('KEKEKEKEKEKEKEKEKEKE, 6', flush=True)

    scores, descriptions = image_checker.find_most_similar_images(image)

    print('KEKEKEKEKEKEKEKEKEKE, 7', flush=True)

    if int(descriptions[0]) == nft_id:
        score = scores[1]
    else:
        score = scores[0]

    score = int(score * 100)
    result = {
        'score': score,
        'detailed information': {
            'scores': [*scores],
            'descriptions': [*descriptions]
        }
    }

    mutex.release()

    return get_result(200, result)


if __name__ == '__main__':
    registerer_thread = Thread(target=register_images)
    registerer_thread.start()

    app.run(debug=True, host='0.0.0.0', port=9090, threaded=False)
