import {IProps} from './types';

import {Card, Typography} from '@material-ui/core';
import React, {FC} from 'react';
import {useIntl} from 'react-intl';

import {INTL_DATA} from './intl';

import styles from './ProgressCertificate.module.scss';


const ProgressCertificate: FC<IProps> = () => {
    const intl = useIntl();

    return (
        <div className={styles.progressCertificate}>
            <Card
                className={styles.verdict}
                variant={'outlined'}
            >
                <Typography variant={'h6'}>
                    {intl.formatMessage({id: INTL_DATA.PROGRESS})}
                </Typography>
            </Card>
            <div className={styles.wait}>
                <Typography variant={'h6'}>
                    {intl.formatMessage({id: INTL_DATA.WAIT})}
                </Typography>
            </div>
        </div>
    );
};

export default ProgressCertificate;
