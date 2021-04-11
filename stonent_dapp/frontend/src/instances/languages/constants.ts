import {ILanguage} from './types';


export const LS_CURRENT_LOCALE = 'LS_CURRENT_LOCALE';

export enum REDUX_ACTION {
    LANGUAGES_SET = 'LANGUAGES_SET',
}

export const DEFAULT_LOCALE = 'EN';

export const LANGUAGES: ILanguage[] = [
    {
        locale: 'EN',
        messages: {
            'landing/title': 'Check NFT picture uniqueness',
            'landing/paste-identifier': 'Past identifier here',
            'landing/supported-platforms': 'Supported collections',
            'landing/test-first': 'This is the demo project',
            'landing/test-second': 'Please, connect to the Rinkeby',
            'landing/test-third': 'test network before using',
            'check/marketplace': 'Go to the marketplace',
            'check/error-payment': 'Payment approvement is required',
            'check/error-aborted': 'Operation has been aborted',
            'check/success-updated': 'Certificate has been updated',
            'check/empty-certificate/not-certified': 'Not certified',
            'check/empty-certificate/price': 'Price',
            'check/empty-certificate/certify': 'Certify',
            'check/progress-certificate/progress': 'Certification in progress',
            'check/progress-certificate/wait': 'Please, wait',
            'check/certificate/date': 'Certification date',
            'check/certificate/version': 'Certificate version',
            'check/certificate/transaction': 'Transaction id',
            'check/certificate/oracle': 'Certification oracle',
            'check/certificate/copy': 'Copy',
            'check/certificate/score': 'Score',
            'check/certificate/recertify': 'Recertify',
            'check/certificate/plagiarism': 'Plagiarism',
            'check/certificate/original': 'Original',
            'check/certificate/copied': 'Copied',
            'not-found/not-found': 'Picture not found',
            'not-found/main-page': 'Go to the main page',
            'no-provider/no-provider': 'Для работы необходим ethereum-провайдер',
        },
    },
    {
        locale: 'RU',
        messages: {
            'landing/title': 'Проверьте уникальность NFT картины',
            'landing/paste-identifier': 'Укажите идентификатор',
            'landing/supported-platforms': 'Поддерживаемая коллекция',
            'landing/test-first': 'Это демонтрационный проект',
            'landing/test-second': 'Пожалуйста, подключитесь к тестовой',
            'landing/test-third': 'сети Rinkeby перед использованием',
            'check/marketplace': 'Отправиться в магазин',
            'check/error-payment': 'Обязательно разрешение платежей',
            'check/error-aborted': 'Операция прервана',
            'check/success-updated': 'Сертификат был обновлён',
            'check/empty-certificate/not-certified': 'Нет сертификата',
            'check/empty-certificate/price': 'Цена',
            'check/empty-certificate/certify': 'Сертифицировать',
            'check/progress-certificate/progress': 'Выполняется сертификация',
            'check/progress-certificate/wait': 'Пожалуйста, подождите',
            'check/certificate/date': 'Дата сертификации',
            'check/certificate/version': 'Версия сертификата',
            'check/certificate/transaction': 'Id транзакции',
            'check/certificate/oracle': 'Оракул сертификации',
            'check/certificate/copy': 'Скопировать',
            'check/certificate/score': 'Рейтинг',
            'check/certificate/recertify': 'Обновить сертификат',
            'check/certificate/plagiarism': 'Плагиат',
            'check/certificate/original': 'Оригинальность',
            'check/certificate/copied': 'Скопировано',
            'not-found/not-found': 'Картина не найдена',
            'not-found/main-page': 'На главную',
            'no-provider/no-provider': 'Для работы необходим ethereum-провайдер',
        },
    },
];
