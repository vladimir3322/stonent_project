import {MuiThemeProvider} from '@material-ui/core';
import React, {FC} from 'react';
import {IntlProvider} from 'react-intl';
import {BrowserRouter as Router, Route, Switch} from 'react-router-dom';

import {theme} from 'helpers/materialUI';

import {useLanguages} from 'instances/languages/hooks';

import Check from 'components/Check';
import Landing from 'components/Landing';
import NoProvider from 'components/NoProvider';
import NotFound from 'components/NotFound';


const App: FC = () => {
    const {language} = useLanguages();

    return (
        <MuiThemeProvider theme={theme}>
            <IntlProvider locale={language.locale} messages={language.messages}>
                <Router>
                    <Switch>
                        <Route exact={true} path={'/'} component={Landing}/>
                        <Route exact={true} path={'/check/:id'} component={Check}/>
                        <Route exact={true} path={'/no_provider'} component={NoProvider}/>
                        <Route exact={true} path={'*'} component={NotFound}/>
                    </Switch>
                </Router>
            </IntlProvider>
        </MuiThemeProvider>
    );
};

export default App;
