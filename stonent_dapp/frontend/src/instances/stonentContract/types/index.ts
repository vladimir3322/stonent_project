export interface ICertificate {
    id: string;
    score: number;
    oracle: string | null;
    version: number;
    date: Date | null;
    transactionId: string;
}

export interface IEventSent {

}

export interface IEventCertified {

}
