import {REDUX_ACTION} from '../constants';


export interface IState {
    locale: string;
    messages: any;
}

export interface ISetAction {
    type: REDUX_ACTION.LANGUAGES_SET;
    locale: string;
    messages: any;
}

export type IAction =
    ISetAction;
