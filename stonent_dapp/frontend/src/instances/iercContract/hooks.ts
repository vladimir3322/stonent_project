import {useEffect, useState} from 'react';
import Web3 from 'web3';
import {Contract} from 'web3-eth-contract';

import {ADDRESS as STONET_ADDRESS} from 'instances/stonentContract/constants';

import abi from './abi.json';
import {ADDRESS} from './constants';


export function useIercContract() {
    const [web3, setWeb3] = useState<Web3 | null>(null);
    const [iercContract, setIercContract] = useState<Contract | null>(null);

    const approveMoneyUsage = async (): Promise<boolean> => {
        if (!web3) {
            return false;
        }
        if (!iercContract) {
            return false;
        }

        const accounts = await web3.eth.getAccounts();

        if (!accounts) {
            return false;
        }

        const defaultAccount = accounts[0];

        if (!defaultAccount) {
            return false;
        }

        return !!await iercContract.methods.approve(STONET_ADDRESS, '99000000000000000000000000000').send({from: defaultAccount});
    };
    const checkAccessMoneyUsage = async () => {
        if (!web3) {
            return false;
        }
        if (!iercContract) {
            return false;
        }

        const accounts = await web3.eth.getAccounts();

        if (!accounts) {
            return false;
        }

        const defaultAccount = accounts[0];

        if (!defaultAccount) {
            return false;
        }

        return !!await iercContract.methods.allowance(defaultAccount, STONET_ADDRESS).call();
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
            setIercContract(contract);
        })();
    }, []);

    return {
        iercContract,
        approveMoneyUsage,
        checkAccessMoneyUsage,
    };
}
