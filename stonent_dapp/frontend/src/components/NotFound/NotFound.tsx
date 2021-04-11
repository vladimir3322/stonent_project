import {Button, Container, Typography} from '@material-ui/core';
import React, {FC} from 'react';
import {useIntl} from 'react-intl';
import {useHistory} from 'react-router-dom';

import Header from 'components/Helpers/Header';

import {INTL_DATA} from './intl';

import styles from './NotFound.module.scss';


const NotFound: FC = () => {
    const history = useHistory();
    const intl = useIntl();

    return (
        <>
            <Header/>
            <Container className={styles.container}>
                <div className={styles.notFound}>
                    <Typography variant={'h4'}>
                        {intl.formatMessage({id: INTL_DATA.NOT_FOUND})}
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

export default NotFound;
