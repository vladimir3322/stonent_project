import {Button, Container, Typography} from '@material-ui/core';
import React, {FC} from 'react';
import {useIntl} from 'react-intl';
import {useHistory} from 'react-router-dom';

import Header from 'components/Helpers/Header';

import {INTL_DATA} from './intl';

import styles from './NoProvider.module.scss';


const NoProvider: FC = () => {
    const history = useHistory();
    const intl = useIntl();

    return (
        <>
            <Header/>
            <Container className={styles.container}>
                <div className={styles.noProvider}>
                    <Typography variant={'h4'}>
                        {intl.formatMessage({id: INTL_DATA.NO_PROVIDER})}
                    </Typography>
                </div>
                <div className={styles.mainPage}>
                    <Button
                        variant={'contained'}
                        color={'primary'}
                        onClick={() => history.push('/')}
                        size={'large'}
                    >
                        {intl.formatMessage({id: INTL_DATA.MAIN_PAGE})}
                    </Button>
                </div>
            </Container>
        </>
    );
};

export default NoProvider;
