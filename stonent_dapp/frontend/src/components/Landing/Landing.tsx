import {Container, TextField, Typography} from '@material-ui/core';
import SearchIcon from '@material-ui/icons/Search';
import React, {FC, useState} from 'react';
import {useIntl} from 'react-intl';
import {useHistory} from 'react-router-dom';

import Header from 'components/Helpers/Header';

import {INTL_DATA} from './intl';

import styles from './Landing.module.scss';


const Landing: FC = () => {
    const [searchValue, setSearchValue] = useState('');

    const history = useHistory();
    const intl = useIntl();

    const onSearch = async () => {
        if (!searchValue) {
            return;
        }

        // @ts-ignore
        await window.web3.currentProvider.enable();

        history.push(`/check/${searchValue}`);
    };

    return (
        <>
            <Header/>
            <Container className={styles.landing}>
                <div className={styles.content}>
                    <div className={styles.title}>
                        <Typography variant={'h3'}>
                            {intl.formatMessage({id: INTL_DATA.TITLE})}
                        </Typography>
                    </div>
                    <div
                        className={styles.search}
                        onKeyPress={(e) => e.key === 'Enter' && onSearch()}
                    >
                        <TextField
                            className={styles.searchInput}
                            label={intl.formatMessage({id: INTL_DATA.PASTE_IDENTIFIER})}
                            value={searchValue}
                            onChange={(e) => setSearchValue(e.target.value)}
                            InputProps={{
                                endAdornment:
                                    <SearchIcon
                                        className={styles.searchIcon}
                                        onClick={onSearch}
                                    />,
                            }}
                        />
                    </div>
                    <div className={styles.supportedCollections}>
                        <Typography>
                            {intl.formatMessage({id: INTL_DATA.SUPPORTED_PLATFORMS})}
                            :
                            {' '}
                            <span
                                className={styles.link}
                                onClick={() => window.open('https://etherscan.io/address/0xd07dc4262bcdbf85190c01c996b4c06a461d2430')}
                            >
                                Rarible Collection
                            </span>
                        </Typography>
                    </div>
                    <div className={styles.testAlert}>
                        <div>
                            <Typography>
                                <b>
                                    {intl.formatMessage({id: INTL_DATA.TEST_FIRST})}
                                </b>
                            </Typography>
                        </div>
                        <div>
                            <Typography>
                                <b>
                                    {intl.formatMessage({id: INTL_DATA.TEST_SECOND})}
                                </b>
                            </Typography>
                        </div>
                        <div>
                            <Typography>
                                <b>
                                    {intl.formatMessage({id: INTL_DATA.TEST_THIRD})}
                                </b>
                            </Typography>
                        </div>
                    </div>
                </div>
            </Container>
        </>
    );
};

export default Landing;
