import {AppBar, Button, Input, Toolbar, Typography} from '@material-ui/core';
import SearchIcon from '@material-ui/icons/Search';
import React, {FC, useState} from 'react';
import {useIntl} from 'react-intl';
import {useHistory} from 'react-router-dom';

import Languages from './Languages';

import {INTL_DATA} from './intl';

import styles from './Header.module.scss';


const Header: FC = () => {
    const [value, setValue] = useState('');

    const history = useHistory();
    const intl = useIntl();

    const onSearch = () => {
        if (!value) {
            return;
        }

        window.location.replace(`/check/${value}`);
    };

    return (
        <AppBar position={'sticky'}>
            <Toolbar className={styles.header}>
                <div>
                    <Button color={'inherit'} onClick={() => history.push('/')}>
                        <Typography variant={'h6'}>
                            STONENT
                        </Typography>
                    </Button>
                </div>
                <div className={styles.right}>
                    <div>
                        <Input
                            className={styles.input}
                            value={value}
                            onChange={(e) => setValue(e.target.value)}
                            onKeyPress={(e) => e.key === 'Enter' && onSearch()}
                            placeholder={intl.formatMessage({id: INTL_DATA.PASTE_IDENTIFIER})}
                            endAdornment={
                                <SearchIcon
                                    className={styles.searchIcon}
                                    onClick={onSearch}
                                />
                            }
                        />
                    </div>
                    <div>
                        <Languages/>
                    </div>
                </div>
            </Toolbar>
        </AppBar>
    );
};

export default Header;
