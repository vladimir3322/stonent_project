import os
from os.path import join
from py_dotenv import read_dotenv


class Config:
    DATA_PATH = './data'
    SOURCE_PATH = './source'
    WEB_3_PROVIDER_URL = 'https://mainnet.infura.io/v3/'
    CONTRACT_ADDRESS = '0xd07dc4262BCDbf85190C01c996b4C06a461d2430'
    CONTRACT_ABI_FILENAME = './contract/abi.json'
    IPFS_IMAGE_SOURCE_URL = 'https://ipfs.daonomic.com'
    FIRST_BLOCK_NUMBER = 10147631
    SAFE_BLOCK_STEP = 142

    def __init__(self):
        pass

    def get_web3_provider(self):
        dotenv_path = join('.env')
        read_dotenv(dotenv_path)
        provider_key = os.getenv('WEB_3_PROVIDER_KEY')

        return f'{self.WEB_3_PROVIDER_URL}{provider_key}'


config = Config()
