import {Link} from 'inferno-router';
import {linkEvent} from 'inferno';
import {info} from './util';

export default class BreadCrumb extends Component {

    constructor({props,}) {
        super(props);
    }

    render() {
        if(!this.props.n) {
            return;
        }
        let {u,p,n,x} = this.props.al;
        let kwk = <Link to="/">kwk.co</Link>;
        if (!n && !p) {
            return <h1 className="uri">{kwk} / <span className="uri-highlight">{u}</span></h1>
        }

        if (!n) {
            return (<h1 className="uri">{kwk} / <Link to={'/' + u}>{u}</Link> / <span className="uri-highlight" onClick={linkEvent(this, rename)}>{p}</span></h1>)
        }
        let nms = n.split(',');
        return (<h1 className="uri">{kwk} / <Link to={'/' + u}>{u}</Link> / <Link
            to={'/' + u + '/' + p}>{p}</Link> / <span className="uri-highlight">{nms.shift()}.{x}</span>{
            nms.map(x => {
                return <small>, {x}</small>
            })
        }</h1>)
    }
}

function rename(com, ev) {
    let {u,p,n} = this.props.al;
    let res = window.prompt("Rename", p);
    if (res !== p) {
        info('saving...')
    }
}
