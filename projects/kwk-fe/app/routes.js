import createBrowserHistory from 'history/createBrowserHistory';
import {doAllAsyncBefore, IndexRoute, match, Route, Router} from "inferno-router";

import MasterTemplate from './site-wide/masterTemplate';
import str from './stores/snpStore';
import Settings from './support/settings';
import Home from './core/home';
import EditLs from './core/editLs';
import PchLs from './core/pchLs';
import SnpLs from './core/snpLs';
import {info} from './site-wide/util';

const browserHistory = createBrowserHistory();

function NoMatch({children}) {
    info('route not found', children);
}

const routes = (
    <Router history={browserHistory} asyncBefore={url => {
        return doAllAsyncBefore(match(routes, url));
    }
    }>
        <Route component={MasterTemplate}>
            <IndexRoute component={Home}/>
            <Route path="/settings" component={Settings}/>

            {/* Ensure that all other non-wild card routes are longer than
        1 char since the router sorts best matches by longest first */}
            <Route path="/:u" component={PchLs} asyncBefore={PchLs.load}>
                <Route path="/:p" component={SnpLs} asyncBefore={PchLs.load}>
                    <Route path="/:n" component={EditLs} />
                </Route>
            </Route>
            <Route path="*" component={NoMatch}/>
        </Route>
    </Router>
);

export default routes