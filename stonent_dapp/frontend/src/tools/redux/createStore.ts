import {createBrowserHistory} from 'history';
import {routerMiddleware} from 'react-router-redux';
import {Action, applyMiddleware, createStore, Reducer} from 'redux';


export default function<IReduxState, IReduxAction extends Action>(reducer: Reducer<IReduxState, IReduxAction>) {
    const history = createBrowserHistory();
    const middleware = [
        routerMiddleware(history),
    ];

    return createStore<IReduxState, IReduxAction, null, null>(reducer);
}
