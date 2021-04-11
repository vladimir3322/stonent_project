import {ISetAction} from './types/redux';

import {useDispatch, useSelector} from 'instances/redux';

import {LANGUAGES, REDUX_ACTION} from './constants';


export function useLanguages() {
    const language = useSelector((state) => state.languages);
    const dispatch = useDispatch();

    const setLanguage = (locale: string) => {
        const language = LANGUAGES.find((language) => language.locale === locale);

        if (!language) {
            return;
        }

        dispatch<ISetAction>({
            type: REDUX_ACTION.LANGUAGES_SET,
            locale,
            messages: language.messages,
        });
    };

    return {
        language,
        setLanguage,
    };
}
