import {ICertificate, IEventCertified, IEventSent} from './types';

import {useEffect, useState} from 'react';
import Web3 from 'web3';
import {Contract} from 'web3-eth-contract';

import {formatPrice} from 'instances/stonentContract/functions';

import abi from './abi.json';
import {ADDRESS} from './constants';


export function useStonentContract() {
    const [web3, setWeb3] = useState<Web3 | null>(null);
    const [stonentContract, setStonentContract] = useState<Contract | null>(null);
    const [certificationPrice, setCertificationPrice] = useState<number | null>(null);

    const getCertificate = async (imageId: string): Promise<ICertificate | null> => {
        if (!stonentContract) {
            return null;
        }

        const transactionId = await stonentContract.methods.lastCertification(imageId).call();

        if (!transactionId) {
            return null;
        }

        const certificate = await stonentContract.methods.certificates(transactionId).call();

        if (!certificate) {
            return null;
        }

        const parsedDate  = parseInt(certificate.Date, 10);

        return {
            id: certificate.ID,
            score: parseInt(certificate.Score, 10),
            oracle: certificate.Oracle === '0x0000000000000000000000000000000000000000' ? null : certificate.Oracle,
            version: certificate.Version,
            date: parsedDate ? new Date(parsedDate * 1000) : null,
            transactionId,
        };
    };
    const checkImage = async (imageId: string) => {
        if (!web3) {
            return;
        }
        if (!stonentContract) {
            return;
        }

        const accounts = await web3.eth.getAccounts();

        if (!accounts) {
            return false;
        }

        const defaultAccount = accounts[0];

        if (!defaultAccount) {
            return false;
        }

        await stonentContract.methods.check(imageId).send({from: defaultAccount});
    };
    const getPrice = async () => {
        if (!stonentContract) {
            return;
        }

        const price = await stonentContract.methods.price().call();
        const parsedPrice = parseInt(price, 10);

        if (isNaN(parsedPrice)) {
            return null;
        }

        const formattedPrice = formatPrice(parsedPrice);

        setCertificationPrice(formattedPrice);
    };
    const subscribeToCertificateSent = (cb: (data: IEventSent) => void) => {
        if (!stonentContract) {
            return;
        }

        stonentContract.events.RequestSended({fromBlock: 'latest'}, (error: any, data: IEventSent) => {
            cb(data);
        });
    };
    const subscribeToCertified = (cb: (data: IEventCertified) => void) => {
        if (!stonentContract) {
            return;
        }

        stonentContract.events.RequestCertified({fromBlock: 'latest'}, (error: any, data: IEventCertified) => {
            cb(data);
        });
    };

    useEffect(() => {
        (async () => {
            try {
                // @ts-ignore
                await window.web3.currentProvider.enable();
            } catch {
                window.location.replace('/no_provider');

                return;
            }

            // @ts-ignore
            const web3 = new Web3(window.web3.currentProvider);
            // @ts-ignore
            const contract = new web3.eth.Contract(abi, ADDRESS);

            setWeb3(web3);
            setStonentContract(contract);
        })();
    }, []);

    return {
        stonentContract,
        certificationPrice,
        getCertificate,
        checkImage,
        getPrice,
        subscribeToCertificateSent,
        subscribeToCertified,
    };
}
