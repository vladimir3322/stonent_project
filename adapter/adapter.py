from bridge import Bridge
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




class Adapter:
    id_params = ['id', 'nft', '_id']

    def __init__(self, input):
        self.id = input.get('id', '1')
        self.result = ""   #kostil
        self.request_data = input.get('data')
        if self.validate_request_data():
            self.bridge = Bridge()
            self.set_params()
            self.create_request()
        else:
            self.result_error('No data provided')

    def validate_request_data(self):
        if self.request_data is None:
            return False
        if self.request_data == {}:
            return False
        return True

    def set_params(self):
        for param in self.id_params:
            self.id_params = self.request_data.get(param)
            if self.id_params is not None:
                break


    def create_request(self):
        try:
            #params = {
            #    'fsym': self.from_param,
            #    'tsyms': self.to_param,
            #}
            #response = self.bridge.request(self.base_url, params)

            nftID = int(self.id_params)
            contract = get_contract()

            loop = asyncio.new_event_loop()
            asyncio.set_event_loop(loop)
            task = asyncio.gather(download_image_data(contract, nftID))
            loop.run_until_complete(task)

            task_result = task.result()


            rndmScore =  1000
            if nftID % 2 == 0:
                rndmScore = random.randint(80, 100)
            else:
                rndmScore = random.randint(0, 40)


            #data = response.json()
            #self.result = data[self.to_param]
            data= {'score': rndmScore}
            self.result_success(data)
        except Exception as e:
            self.result_error(e)
        finally:
            self.bridge.close()

    def result_success(self, data):
        self.result = {
            'jobRunID': self.id,
            'data': data,
            'result': self.result,
            'statusCode': 200,
        }

    def result_error(self, error):
        self.result = {
            'jobRunID': self.id,
            'status': 'errored',
            'error': f'There was an error: {error}',
            'statusCode': 500,
        }
