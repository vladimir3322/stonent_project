export interface IConstructorData {
    url: string;
    method: 'GET' | 'POST' | 'PATCH' | 'DELETE';

    credentials?: 'omit' | 'same-origin' | 'include';
    headers?: string[][];
    query?: any;
    body?: any;
    formData?: FormData;
}

export type IHandler<IResult> = (body: any) => IResult;
export type IDefaultCallback = () => void;
