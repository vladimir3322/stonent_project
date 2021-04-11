import {ICertificate} from 'instances/stonentContract/types';


export function certificateHasBeenUpdated(first: ICertificate | null, second: ICertificate | null): boolean {
    if (!first) {
        return true;
    }
    if (!second) {
        return true;
    }
    if (first.score !== second.score) {
        return true;
    }
    if (first.oracle !== second.oracle) {
        return true;
    }
    if (first.version !== second.version) {
        return true;
    }
    if (first.date !== second.date) {
        return true;
    }
    if (first.transactionId !== second.transactionId) {
        return true;
    }

    return false;
}
