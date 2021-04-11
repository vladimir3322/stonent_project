import {ICertificate} from 'instances/stonentContract/types';


export interface IProps {
    certificate: ICertificate;
    isFetching: boolean;

    onRunCheck(): void;
}
