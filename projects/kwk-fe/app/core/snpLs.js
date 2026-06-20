import SnpItem from '../support/snpItem'
import {info} from '../site-wide/util';
import str from '../stores/snpStore';
import Grid from "./grid";
import List from "./list";
import BreadCrumb from "../site-wide/breadcrumb";
import PchMenu from '../site-wide/pchMenu';

export default class SnpLs extends Component {
    constructor({props}) {
        super(props);
        info();
        str.subscribe(['snp-changed'], null, () => {
           this.setState({})
        });
    }

    componentWillUnmount() {
        info();
    }

    render({params, children}) {
        info();
        let cp = str.currPch;
        let items = str.currPch.items;

        // No snippets
        if (!cp.p || !cp.items) {
            return (
                <div>
                    <BreadCrumb className="bread-crumb" al={cp.p}/>
                    <PchMenu/>
                    No snippets
                </div>)
        }
        // Grid view
        if (!params.n) {
            return (
                <div className="struct">
                    <BreadCrumb className="bread-crumb" al={cp.p}/>
                    <PchMenu/>
                    <Grid items={items}>
                        <SnpItem showPrvw="true" gridView="true" />
                    </Grid>
                </div>)
        }
        // Snippet View
        return (
            <div className="struct">
                <div className="columns is-12 struct">
                    <div className="column is-narrow struct">
                        <PchMenu/>
                        <List items={items}>
                            <SnpItem/>
                        </List>
                    </div>
                    <div className="column struct" style="padding: 0 0 0 15px;">
                        {children}
                    </div>
                </div>
            </div>
        )
    }
}


{/*<h2 className="editable is-size-4" onClick={linkEvent(this, rename)}>{str.currentPouch.p.p}</h2>*/
}
{/*{ std.length>0 ? <DragHandle num={std.length}/> : '' }*/
}