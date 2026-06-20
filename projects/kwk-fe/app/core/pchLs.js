import {info, pchPath} from '../site-wide/util';
import PchItem from '../support/pchItem'
import str from '../stores/snpStore';
import Grid from "./grid";
import List from "./list";
import BreadCrumb from "../site-wide/breadcrumb";
import SnpItem from "../support/snpItem";

export default class PchLs extends Component {
    constructor({props}) {
        super(props);
        info();
    }

    static load(params) {
        return str.loadPouch(params)
    }

    render() {
        info();
        let {params} = this.props;

        if (!params.p) {
            return (
                <div className="root">
                    <BreadCrumb className="bread-crumb level-left" al={params}/>
                    {
                        str.latest.map(x => {
                            return <SnpItem showPrvw="true" showDate="true" al={x.al} cnt={x.cnt} crtd={x.crtd}/>
                        })
                    }
                    <div className="scroll"><Grid items={str.ps}><PchItem/></Grid>
                    </div>
                </div>);
        }
        return (
            <div className="struct">
                <div className="columns struct">
                    <div className="column is-narrow struct">
                        <div className="pch-header">
                            <h2>{str.ps.length} Pouches</h2>
                        </div>
                        <List items={str.ps}>
                            <PchItem/>
                        </List>
                    </div>
                    <div className="column struct">
                                    {this.props.children}
                    </div>
                </div>
            </div>
        )
    }
}
