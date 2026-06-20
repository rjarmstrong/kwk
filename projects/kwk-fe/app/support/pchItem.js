import str from '../stores/snpStore'
import {info, pchPath,} from '../site-wide/util';
import {Link} from 'inferno-router'

export default class PchItem extends Component {
    constructor(props) {
        super(props);
        let that = this;
        str.subscribe(['pouch-hover'], pchPath(this.props.p), (chan, val) => {
            that.setState({});
        });
    }

    render() {
        info();
        let pch = this.props;
        let cn = 'pouch ';
        if (str.getPouch() && str.getPouch() === this.props.p) {
            cn += ' selected';
        }
        return (<Link to={pchPath(pch)}>
            <div id={pchPath(pch)} className={cn}>
                <div className={ getClass(pch) }/>
                {pch.p} <span className="subdued">{pch.stats.snips}</span>
            </div></Link>
        )
    }
}

function getClass(p) {
    let c = 'pch-icon';
    if(p.mkPvt) {
        c += ' pvt';
    }
    if(str.stats.lastP === p.p) {
        c += ' last';
    }
    return c;
}