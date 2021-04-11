import {DEFAULT_LOCALE, LS_CURRENT_LOCALE} from './constants';


export function getLSLocale() {
    return localStorage.getItem(LS_CURRENT_LOCALE) || DEFAULT_LOCALE;
}

export function setLSLocale(locale: string) {
    localStorage.setItem(LS_CURRENT_LOCALE, locale);
}
