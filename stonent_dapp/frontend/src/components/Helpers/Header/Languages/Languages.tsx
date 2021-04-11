import {Button, Menu, MenuItem} from '@material-ui/core';
import React, {FC, useState} from 'react';

import {LANGUAGES} from 'instances/languages/constants';
import {useLanguages} from 'instances/languages/hooks';

import {LANGUAGE_MENU_ID} from './constants';


const Languages: FC = () => {
    const [anchor, setAnchor] = useState<HTMLButtonElement | null>(null);

    const {language, setLanguage} = useLanguages();

    const onClose = () => {
        setAnchor(null);
    };
    const onChange = (locale: string) => {
        setLanguage(locale);
        onClose();
    };

    return (
        <>
            <Button
                color={'inherit'}
                aria-controls={LANGUAGE_MENU_ID}
                aria-haspopup={'true'}
                onClick={(e) => setAnchor(e.currentTarget)}
            >
                {language.locale}
            </Button>
            <Menu
                id={LANGUAGE_MENU_ID}
                anchorEl={anchor}
                keepMounted={true}
                open={!!anchor}
                onClose={onClose}
            >
                {
                    LANGUAGES.map(
                        (language) =>
                            <MenuItem key={language.locale} onClick={() => onChange(language.locale)}>
                                {language.locale}
                            </MenuItem>,
                    )
                }
            </Menu>
        </>
    );
};

export default Languages;
