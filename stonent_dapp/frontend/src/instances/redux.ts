import {IReduxAction, IReduxState} from './types';

import {useDispatch as reactReduxUseDispatch, useSelector as reactReduxUseSelector} from 'react-redux';
import {combineReducers, Dispatch} from 'redux';

import languages from './languages/redux';


export const reducer = combineReducers<IReduxState, IReduxAction>({
    languages,
});

export const useDispatch = () => reactReduxUseDispatch<Dispatch<IReduxAction>>();
export const useSelector = <IResult>(selector: (state: IReduxState) => IResult) => reactReduxUseSelector<IReduxState, IResult>(selector);
