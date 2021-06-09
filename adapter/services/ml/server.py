import http.server
import socketserver
import json
import loader
import numpy
import base64
import cv2
import config
import globals
from PIL import Image
from urllib.parse import urlparse, parse_qs


class RequestHandler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

    def do_GET(self):
        parsed_url = urlparse(self.path)
        parsed_query = parse_qs(parsed_url.query)

        status_code = 404
        response_body = {'error': 'Not found'}

        if parsed_url.path == '/check':
            contract_address = parsed_query.get('contract_address')[0]
            nft_id = parsed_query.get('nft_id')[0]

            response_body = get_adapter_result(contract_address, nft_id)
            status_code = 200
        elif parsed_url.path == '/register_new_image':
            contract_address = parsed_query.get('contract_address')[0]
            nft_id = parsed_query.get('nft_id')[0]

            status_code, response_body = register_new_image(contract_address, nft_id)
        elif parsed_url.path == '/check_registered_image':
            contract_address = parsed_query.get('contract_address')[0]
            nft_id = parsed_query.get('nft_id')[0]

            status_code, response_body = check_registered_image(contract_address, nft_id)
        elif parsed_url.path == '/register_rejected_image_by_ipfs':
            contract_address = parsed_query.get('contract_address')[0]
            nft_id = parsed_query.get('nft_id')[0]
            ipfs_host = parsed_query.get('ipfs_host')[0]
            ipfs_path = parsed_query.get('ipfs_path')[0]
            error = parsed_query.get('error')[0]

            status_code, response_body = register_rejected_image_by_ipfs(
                contract_address,
                nft_id,
                ipfs_host,
                ipfs_path,
                error
            )
        elif parsed_url.path == '/rejected_images_by_nn':
            status_code, response_body = get_rejected_images_by_nn()
        elif parsed_url.path == '/rejected_images_by_ipfs':
            status_code, response_body = get_rejected_images_by_ipfs()
        elif parsed_url.path == '/statistics':
            status_code, response_body = get_statistics()

        response = json.dumps(response_body)
        response = bytes(response, 'utf8')

        self.send_response(status_code)
        self.send_header('Content-type', 'application/json')
        self.end_headers()
        self.wfile.write(response)
        return


def get_adapter_result(contract_address, nft_id):
    def get_result(code, data=None, error=None):
        return {
            'job_run_id': f'{contract_address}_{nft_id}',
            'data': data,
            'error': error,
            'statusCode': code,
        }

    with globals.mutex:
        try:
            if not contract_address:
                return get_result(400, None, 'Invalid contract address')
            if not nft_id or not nft_id.isdigit():
                return get_result(400, None, 'Invalid nft id')

            image_source, image_source_error = loader.get_image_source(contract_address, nft_id)

            if not image_source or image_source_error:
                return get_result(500, None, image_source_error)

            np_image_source = numpy.frombuffer(base64.b64decode(image_source), numpy.uint8)
            image = cv2.imdecode(np_image_source, cv2.IMREAD_COLOR)
            image = Image.fromarray(image)

            scores, descriptions = globals.image_checker.find_most_similar_images(image)

            if not len(scores) or not len(descriptions):
                return get_result(500, None, f'Invalid NN response: scores={scores}; description={descriptions}')

            if descriptions[0] == f'{contract_address}-{nft_id}':
                score = scores[1]
            else:
                score = scores[0]

            result = {
                'score': int(score * 100),
                'detailed_information': {
                    'scores': [*scores],
                    'descriptions': [*descriptions]
                }
            }

            return get_result(200, result)
        except Exception as e:
            return get_result(200, None, e)


def register_new_image(contract_address, nft_id):
    with globals.mutex:
        try:
            if not contract_address:
                return 400, {'error': 'Invalid contract address'}
            if not nft_id or not nft_id.isdigit():
                return 400, {'error': 'Invalid nft id'}

            with open(config.registered_images_file) as file:
                data = map(lambda s: s.split(','), file.readlines())

                for registered_contract_address, registered_nft_id, _ in data:
                    if registered_contract_address == contract_address and registered_nft_id == nft_id:
                        return 400, {'error': 'Already registered'}

            image_source, image_source_error = loader.get_image_source(contract_address, nft_id)

            if not image_source or image_source_error:
                return 500, {'error': image_source_error}

            globals.image_manager.register_new_image(contract_address, nft_id, base64.b64decode(image_source))

            return 200, {'error': None}
        except Exception as e:
            return 500, {'error': e}


def check_registered_image(contract_address, nft_id):
    with globals.mutex:
        try:
            if not contract_address:
                return 400, {'error': 'Invalid contract address'}
            if not nft_id or not nft_id.isdigit():
                return 400, {'error': 'Invalid nft id'}

            with open(config.registered_images_file) as file:
                data = map(lambda s: s.split(','), file.readlines())

                for registered_contract_address, registered_nft_id, _ in data:
                    if registered_contract_address == contract_address and registered_nft_id == nft_id:
                        return 200, {'is_registered': True}
            return 200, {'is_registered': False}
        except Exception as e:
            return 500, {'error': e}


def register_rejected_image_by_ipfs(contract_address, nft_id, ipfs_host, ipfs_path, error):
    with globals.mutex:
        if not contract_address:
            return 400, {'error': 'contract_address is required'}
        if not nft_id:
            return 400, {'error': 'nft_id is required'}
        if not ipfs_host:
            return 400, {'error': 'ipfs_host is required'}
        if not ipfs_path:
            return 400, {'error': 'ipfs_path is required'}
        if not error:
            return 400, {'error': 'error is required'}

        try:
            with open(config.rejected_images_by_IPFS_file, 'a') as file:
                error = error.replace(',', '')
                print(f'{contract_address},{nft_id},{ipfs_host},{ipfs_path},{error}', file=file)
            return 200, {'error': None}
        except Exception as e:
            return 500, {'error': e}


def get_rejected_images_by_ipfs():
    def map_list_to_dict(item):
        return {
            'contract_address': item[0],
            'nft_id': item[1],
            'ipfs_host': item[2],
            'ipfs_path': item[3],
            'description': item[4],
        }

    with globals.mutex:
        try:
            with open(config.rejected_images_by_IPFS_file) as file:
                data = map(lambda s: s.rstrip().split(','), file.readlines())
                data = map(map_list_to_dict, data)

            return 200, {'rejected_images_by_nn': list(data)}
        except Exception as e:
            return 500, {'error': e}


def get_rejected_images_by_nn():
    with globals.mutex:
        try:
            with open(config.rejected_images_by_NN_file) as file:
                data = map(lambda s: s.rstrip().split(','), file.readlines())
                data = map(lambda a: {'contract_address': a[0], 'nft_id': a[1], 'description': a[2]}, data)

            return 200, {'rejected_images_by_nn': list(data)}
        except Exception as e:
            return 500, {'error': e}


def get_statistics():
    with globals.mutex:
        try:
            with open(config.registered_images_file) as file:
                registered_images_count = len(file.readlines()) or 0
            with open(config.rejected_images_by_IPFS_file) as file:
                rejected_by_ipfs_images_count = len(file.readlines()) or 0
            with open(config.rejected_images_by_NN_file) as file:
                rejected_by_nn_images_count = len(file.readlines()) or 0

            loader_response, loader_error = loader.get_statistics()

            if not loader_response or loader_error:
                raise loader_error
            else:
                found_images_count = loader_response.get('CountOfFound')
                precessed_images_count = loader_response.get('CountOfDownloaded')

            return 200, {
                'statistics': {
                    'found_images_count': found_images_count,
                    'precessed_images_count': precessed_images_count,
                    'registered_images_count': registered_images_count,
                    'rejected_by_ipfs_images_count': rejected_by_ipfs_images_count,
                    'rejected_by_nn_images_count': rejected_by_nn_images_count,
                    'is_completed': globals.all_images_has_been_downloaded,
                }
            }
        except Exception as e:
            return 500, {'error': e}


def run_server():
    with socketserver.TCPServer(('', config.server_port), RequestHandler) as httpd:
        print(f'Server started at port: {str(config.server_port)}', flush=True)
        httpd.serve_forever()
