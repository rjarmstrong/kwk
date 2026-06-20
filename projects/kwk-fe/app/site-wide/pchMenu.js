import str from '../stores/snpStore';
import {linkEvent} from 'inferno';
import {Link} from "inferno-router";
import {info, parseAl, pchPath} from "./util";

export default class PchMenu extends Component {

    constructor(props) {
        super(props)
    }

    render() {
        let cp = str.currPch;
        info(cp.p);
        return (
            <div style="margin-bottom:5px">
                <Link to={pchPath(cp.p)}><h2><span className="pch-icon"/> {cp.p.p} <span className="subdued">{cp.tot}</span></h2></Link>
                <div style="margin: 11px 0;">
                    <button className="button is-small" id={pchPath(cp.p)} onClick={linkEvent(this, delPch)}>
                        Delete Pouch</button>
                    {cp.p.mkPvt ? <button className="button is-small" id={pchPath(cp.p)} onClick={linkEvent(this, mkPvt)}>Make Public</button> : <button className="button is-small" id={pchPath(cp.p)} onClick={linkEvent(this, mkPvt)}>Make Pouch Private</button>}


                </div>
            </div>
        )
    }
}

function mkPvt(com, evt) {
    str.mkPvt(parseAl(evt.target.id)).then(() => {
        com.setState({});
    });
    evt.preventDefault();
}

function delPch(inst, evt) {
    info(evt.target.id);
    if(!confirm('Are you sure you want to delete this?')) {
        evt.preventDefault();
        return;
    }
    str.delPch(parseAl(evt.target.id)).then(() => {
        inst.context.router.push('/');
    });
    evt.preventDefault();
}