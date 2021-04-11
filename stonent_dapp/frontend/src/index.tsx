import {IReduxAction, IReduxState} from 'instances/types';

import React from 'react';
import ReactDOM from 'react-dom';
import {Provider} from 'react-redux';

import {createStore} from 'tools/redux';

import {reducer} from 'instances/redux';

import App from './App';

import './index.scss';


ReactDOM.render(
    <React.StrictMode>
        <Provider store={createStore<IReduxState, IReduxAction>(reducer)}>
            <App/>
        </Provider>
    </React.StrictMode>,
    document.getElementById('root'),
);
