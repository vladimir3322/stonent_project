import {IProps} from './types';

import {Snackbar, Typography} from '@material-ui/core';
import Alert from '@material-ui/lab/Alert';
import React, {FC} from 'react';


const SuccessMessage: FC<IProps> = (props) => {
    return (
        <Snackbar open={props.isOpen} autoHideDuration={5000} onClose={props.onClose}>
            <Alert variant={'filled'} severity={'success'} onClose={props.onClose}>
                <Typography>
                    {props.message}
                </Typography>
            </Alert>
        </Snackbar>
    );
};

export default SuccessMessage;
