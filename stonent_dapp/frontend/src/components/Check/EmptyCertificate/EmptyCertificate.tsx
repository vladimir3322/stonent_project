import {IProps} from './types';

import {Button, Card, CircularProgress, Typography} from '@material-ui/core';
import React, {FC} from 'react';
import {useIntl} from 'react-intl';

import {INTL_DATA} from './intl';

import styles from './EmptyCertificate.module.scss';


const EmptyCertificate: FC<IProps> = (props) => {
    const intl = useIntl();

    return (
        <div className={styles.emptyCertificate}>
            <Card
                className={styles.verdict}
                variant={'outlined'}
            >
                <Typography variant={'h6'} color={'error'}>
                    {intl.formatMessage({id: INTL_DATA.NOT_CERTIFIED})}
                </Typography>
            </Card>
            <div className={styles.certify}>
                <Typography variant={'h6'}>
                    {intl.formatMessage({id: INTL_DATA.PRICE})}
                    :
                    {' '}
                    {props.certificationPrice}
                    {' '}
                    USDT
                </Typography>
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
                                    {intl.formatMessage({id: INTL_DATA.CERTIFY})}
                                </Typography>
                        }
                    </Button>
                </div>
            </div>
        </div>
    );
};

export default EmptyCertificate;
