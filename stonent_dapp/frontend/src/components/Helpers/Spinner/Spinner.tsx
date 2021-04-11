import {IProps} from './types';

import {CircularProgress} from '@material-ui/core';
import React from 'react';

import styles from './Spinner.module.scss';


const Spinner = (props: IProps) => {
    return (
        <div className={props.isPage ? styles.spinnerPage : props.className}>
            <CircularProgress size={props.size || 100}/>
        </div>
    );
};

export default Spinner;
