import requests
import config


def get_image_source(address, image_id):
    res = requests.get(f'{config.loader_url}/image_source?address={address}&id={image_id}')

    if res.status_code != 200:
        res.encoding = 'utf-8'

        return None, f'Invalid response code from loader server: {res.status_code} with error: {res.text}'

    return res.text, None


def get_statistics():
    res = requests.get(f'{config.loader_url}/statistics')

    if res.status_code != 200:
        res.encoding = 'utf-8'

        return None, f'Invalid response code from loader server: {res.status_code} with error: {res.text}'

    return res.json(), None
