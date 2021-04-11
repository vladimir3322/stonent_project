import {IImageMetadata} from './types';

import config from 'config';
import {useEffect, useState} from 'react';
import Web3 from 'web3';
import {Contract} from 'web3-eth-contract';

import abi from './abi.json';
import {ADDRESS} from './constants';


export function useRaribleContract() {
    const [web3, setWeb3] = useState<Web3 | null>(null);
    const [raribleContract, setRaribleContract] = useState<Contract | null>(null);

    const getImageMetadataUrl = async (id: string): Promise<string | null> => {
        if (!raribleContract) {
            return null;
        }

        const imageMetadataURI: string = await raribleContract.methods.uri(id).call();

        return /\/ipfs\/.*$/.test(imageMetadataURI) ? imageMetadataURI : null;
    };
    const getImageMetadata = async (metadataURL: string): Promise<IImageMetadata | null> => {
        const res = await fetch(metadataURL);

        try {
            if (res.status === 200) {
                const data = await res.json();

                return {
                    name: data.name,
                    description: data.description,
                    imageIPFSURL: data.image,
                    externalUrl: data.external_url,
                    attributes: data.attributes,
                };
            } else {
                return null;
            }
        } catch (error) {
            console.log(error);

            return null;
        }
    };

    useEffect(() => {
        (async () => {
            let provider: any = null;

            switch (config.mode) {
                case 'DEV': {
                    provider = 'wss://mainnet.infura.io/ws/v3/54f99e303d714f1899660b512775e88a';
                    break;
                }
                case 'TEST': {
                    provider = 'wss://mainnet.infura.io/ws/v3/54f99e303d714f1899660b512775e88a';
                    break;
                }
                case 'PROD': {
                    try {
                        // @ts-ignore
                        await window.web3.currentProvider.enable();
                    } catch {
                        window.location.replace('/no_provider');

                        return;
                    }

                    // @ts-ignore
                    provider = window.web3.currentProvider;
                    break;
                }
            }

            const web3 = new Web3(provider);
            // @ts-ignore
            const contract = new web3.eth.Contract(abi, ADDRESS);

            setWeb3(web3);
            setRaribleContract(contract);
        })();
    }, []);

    return {
        raribleContract,
        getImageMetadataUrl,
        getImageMetadata,
    };
}
