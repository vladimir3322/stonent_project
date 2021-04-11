import {IConstructorData, IDefaultCallback, IHandler} from './types';
import {IResponseError} from 'tools/types';

import {FAILED_TO_PARSE_RESPONSE_BODY, UNRECOGNIZED_RESPONSE_CODE} from './constants';


class Fetch<IResult extends {error?: IResponseError}> {
    private readonly url: IConstructorData['url'];
    private readonly method: IConstructorData['method'];
    private readonly credentials: IConstructorData['credentials'] = 'include';
    private readonly headers: IConstructorData['headers'] = [['Content-type', 'application/json']];
    private readonly query: IConstructorData['query'];
    private readonly body: IConstructorData['body'];
    private readonly formData: IConstructorData['formData'];

    private statuses: number[][] = [];
    private handlers: IHandler<IResult>[] = [];
    private needToParseBody: boolean[] = [];
    private defaultCallback: IDefaultCallback | undefined;

    constructor(data: IConstructorData) {
        this.url = data.url;
        this.method = data.method;
        this.credentials = data.credentials || this.credentials;
        this.headers = data.headers || this.headers;
        this.query = data.query || this.query;
        this.body = data.body || this.body;
        this.formData = data.formData || this.formData;
    }

    public static constructQuery<IQuery = any>(query?: IQuery): string {
        if (!query) {
            return '';
        }

        const result: string[] = [];

        Object.keys(query).forEach((key) => {
            const item = (query as any)[key];

            switch (typeof item) {
                case 'object': {
                    if (item instanceof Date) {
                        result.push(`${encodeURIComponent(key)}=${encodeURIComponent(item.getTime())}`);
                    } else if (item instanceof Array) {
                        result.push(item.map((elem) => `${key}=${elem}`).join('&'));
                    } else {
                        throw Error('Unsupported query type');
                    }
                    break;
                }
                default: {
                    result.push(`${encodeURIComponent(key)}=${encodeURIComponent(item)}`);
                    break;
                }
            }
        });

        return `?${result.join('&')}`;
    }

    public on(statuses: number[], handler: IHandler<IResult>, needToParseBody = true): this {
        this.statuses.push(statuses);
        this.handlers.push(handler);
        this.needToParseBody.push(needToParseBody);

        return this;
    }

    public onDefault(callback: IDefaultCallback): this {
        this.defaultCallback = callback;

        return this;
    }

    public async exec(): Promise<IResult> {
        const response = await fetch(`${this.url}${Fetch.constructQuery(this.query)}`, {
            method: this.method,
            credentials: this.credentials,
            headers: this.formData ? [['Content-type', 'multipart/form-data']] : this.headers,
            body: this.formData ?
                this.formData :
                this.body ?
                    JSON.stringify(this.body) :
                    undefined,
        });
        const handlerIndex = this.statuses.findIndex((statuses) => statuses.includes(response.status));
        const handler = this.handlers[handlerIndex];
        const needToParseBody = this.needToParseBody[handlerIndex];

        if (handler) {
            try {
                const body = needToParseBody ? await response.json() : {error: null};

                return handler(body);
            } catch (error) {
                this.defaultCallback?.();

                return {error: FAILED_TO_PARSE_RESPONSE_BODY} as unknown as IResult;
            }
        } else {
            this.defaultCallback?.();

            return {error: UNRECOGNIZED_RESPONSE_CODE} as unknown as IResult;
        }
    }
}

export default Fetch;
