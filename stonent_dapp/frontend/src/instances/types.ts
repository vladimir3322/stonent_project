import * as NSLanguages from './languages/types/redux';


export interface IReduxState {
    languages: NSLanguages.IState;
}

export type IReduxAction =
    NSLanguages.IAction;
