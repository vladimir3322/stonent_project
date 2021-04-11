import {IProps} from './types';

import {Button, Card, CircularProgress, Snackbar, Typography} from '@material-ui/core';
import Alert from '@material-ui/lab/Alert';
import copy from 'copy-text-to-clipboard';
import React, {FC, useState} from 'react';
import {useIntl} from 'react-intl';

import {formatDateTime} from 'tools/functions';

import {SCORE_DIVIDER} from './constants';
import {INTL_DATA} from './intl';

import styles from './Certificate.module.scss';


const Certificate: FC<IProps> = (props) => {
    const [copyNotificationIsOpen, setCopyNotificationIsOpen] = useState(false);

    const intl = useIntl();

    const onCopy = (text: string) => {
        copy(text);
        setCopyNotificationIsOpen(true);
    };

    return (
        <>
            <div>
                <div className={styles.row}>
                    <div className={styles.solid}>
                        <div className={styles.title}>
                            <Typography>
                                <b className={styles.title}>
                                    {intl.formatMessage({id: INTL_DATA.DATE})}
                                    :
                                </b>
                                {' '}
                            </Typography>
                        </div>
                        <Typography>
                            {formatDateTime(props.certificate.date)}
                        </Typography>
                    </div>
                </div>
                <div className={styles.row}>
                    <div className={styles.solid}>
                        <div className={styles.title}>
                            <Typography>
                                <b className={styles.title}>
                                    {intl.formatMessage({id: INTL_DATA.VERSION})}
                                    :
                                </b>
                                {' '}
                            </Typography>
                        </div>
                        <Typography>
                            {props.certificate.version}
                        </Typography>
                    </div>
                </div>

                <div className={styles.row}>
                    <div className={styles.solid}>
                        <div className={styles.title}>
                            <Typography>
                                <b>
                                    {intl.formatMessage({id: INTL_DATA.TRANSACTION})}
                                    :
                                </b>
                                {' '}
                            </Typography>
                        </div>
                        <Typography className={styles.text}>
                            {props.certificate.transactionId}
                        </Typography>
                    </div>
                    <Button
                        variant={'contained'}
                        color={'primary'}
                        size={'small'}
                        onClick={() => onCopy(props.certificate.transactionId)}
                    >
                        {intl.formatMessage({id: INTL_DATA.COPY})}
                    </Button>
                </div>
                <div className={styles.row}>
                    <div className={styles.solid}>
                        <div className={styles.title}>
                            <Typography>
                                <b>
                                    {intl.formatMessage({id: INTL_DATA.ORACLE})}
                                    :
                                </b>
                                {' '}
                            </Typography>
                        </div>
                        <Typography className={styles.text}>
                            {props.certificate.oracle}
                        </Typography>
                    </div>
                    <Button
                        variant={'contained'}
                        color={'primary'}
                        size={'small'}
                        onClick={() => onCopy(props.certificate.oracle || '')}
                    >
                        {intl.formatMessage({id: INTL_DATA.COPY})}
                    </Button>
                </div>
                <div className={styles.result}>
                    <Card
                        className={props.certificate.score < SCORE_DIVIDER ? styles.danderVerdict : styles.succeedVerdict}
                        variant={'outlined'}
                    >
                        <Typography variant={'h6'}>
                            {intl.formatMessage({id: INTL_DATA.SCORE})}
                            :
                            {' '}
                            {props.certificate.score}
                        </Typography>
                    </Card>
                    <div className={styles.button}>
                        <Button
                            variant={'contained'}
                            color={'primary'}
                            size={'large'}
                            disabled={props.isFetching}
                            onClick={props.onRunCheck}
                        >
                            {
                                props.isFetching ?
                                    <CircularProgress size={25} color={'inherit'}/> :
                                    <Typography>
                                        {intl.formatMessage({id: INTL_DATA.RECERTIFY})}
                                    </Typography>
                            }
                        </Button>
                    </div>
                </div>
                <div className={styles.scoreBar}>
                    <div className={styles.title}>
                        <Typography>
                            <small>
                                {intl.formatMessage({id: INTL_DATA.PLAGIARISM})}
                            </small>
                        </Typography>
                        <Typography>
                            <small>
                                {intl.formatMessage({id: INTL_DATA.ORIGINAL})}
                            </small>
                        </Typography>
                    </div>
                    <div className={styles.bar}/>
                    <div className={styles.title}>
                        <Typography>
                            <small>
                                0
                            </small>
                        </Typography>
                        <Typography>
                            <small>
                                100
                            </small>
                        </Typography>
                    </div>
                </div>
            </div>
            <Snackbar
                open={copyNotificationIsOpen}
                autoHideDuration={3000}
                onClose={() => setCopyNotificationIsOpen(false)}
            >
                <Alert
                    variant={'filled'}
                    severity={'success'}
                    onClose={() => setCopyNotificationIsOpen(false)}
                >
                    <Typography>
                        {intl.formatMessage({id: INTL_DATA.COPIED})}
                        !
                    </Typography>
                </Alert>
            </Snackbar>
        </>
    );
};

export default Certificate;
