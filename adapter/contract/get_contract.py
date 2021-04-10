from web3 import Web3
import json

from . import config


def get_contract():
    w3 = Web3(Web3.HTTPProvider(config.config.get_web3_provider()))
    with open(config.config.CONTRACT_ABI_FILENAME) as json_file:
        abi = json.load(json_file)
    return w3.eth.contract(address=config.config.CONTRACT_ADDRESS, abi=abi)
