import {IState} from './types/redux';
import {IReduxAction} from 'instances/types';

import {LANGUAGES, REDUX_ACTION} from './constants';
import {getLSLocale, setLSLocale} from './functions';


const initState: IState = {
    locale: getLSLocale(),
    messages: LANGUAGES.find((language) => language.locale === getLSLocale())?.messages,
};

export default function(state: IState = initState, action: IReduxAction) {
    switch (action.type) {
        case REDUX_ACTION.LANGUAGES_SET: {
            const {locale, messages} = action;

            setLSLocale(locale);

            return {
                ...state,
                locale,
                messages,
            };
        }
        default: {
            return state;
        }
    }
}
