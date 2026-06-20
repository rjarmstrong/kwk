import {linkEvent} from 'inferno';
import {httpGet} from '../api/client'
import {info, relPath, snpNav, uri} from "./util";
import createElement from 'inferno-create-element';
import snpStore from "../stores/snpStore";

export default class FindBox extends Component {
    constructor(props) {
        super(props);
        this.state = {
            fnd: [],
            ss: [],
            ps: [],
            ix: 0,
        };
    }

    render() {
        info();
        let dropDown = null;
        if (this.state.fnd.length > 0) {
            let cn = 'item';
            let i = 0;
            dropDown = <div className="drop-down">
                <div>
                {
                    this.state.fnd.map(x => {
                        if (i === this.state.ix) {
                            cn = 'item selected';
                        } else {
                            cn = 'item';
                        }
                        if(!x.al.n) {
                            cn += ' pouch';
                        }
                        i++;
                        let ht = createElement('div', {}, []);
                        if (x.hlts['cnt']) {
                            for (let i = 0; i < x.hlts['cnt'].seg.length; i++) {
                                let className = i % 2 === 1 ? 'highlight' : 'plain';
                                ht.children.push(createElement('span', {className: className}, x.hlts['cnt'].seg[i]));
                            }
                        }
                        let head = createElement('div', {className: 'header'}, []);
                        if (x.hlts['relative_name']) {
                            for (let i = 0; i < x.hlts['relative_name'].seg.length; i++) {
                                let className = i % 2 === 1 ? 'highlight' : 'plain';
                                head.children.push(createElement('span', {className: className}, x.hlts['relative_name'].seg[i]));
                            }
                        } else {
                            head.children.push(createElement('span', {}, relPath(x.al)));
                        }
                        return <div className={cn} id={uri(x.al)} onMouseDown={linkEvent(this, handleMouseDown)}>
                                {head}
                            <pre className="prvw">{ht}</pre>
                        </div>
                    })
                }
                </div>
            </div>
        }

        return ([
            <div>Hola</div>,
            <div>
                <form>
                    <div className="search-box">
                        <span className="prompt">&gt;</span><input placeholder={'kwk search'} autoComplete="off"
                                                                   onKeyDown={linkEvent(this, handleKeyDown)}
                                                                   onKeyUp={linkEvent(this, handleKeyUp)} name="term"/>
                    </div>
                </form>
                {dropDown}
            </div>]
        )
    }
}

function handleMouseDown(com, ev) {
    info(ev.target.id);
    com.context.router.push(ev.target.id);
    com.setState({fnd: [] });
    ev.preventDefault();
}

// TASK: add parallel upgrade KeyboardEvent.code e.g. ev.code === 'Enter'
function handleKeyDown(com, ev) {
    // escape
    if(ev.keyCode === 27) {
        com.setState({fnd: [], ix : 0});
        return;
    }
    if (ev.keyCode === 13) {
        com.context.router.push(uri(com.state.fnd[com.state.ix].al));
        com.setState({fnd: [] });
        ev.preventDefault();
    }
    if (ev.keyCode === 38 && com.state.ix > -1) {
        ev.preventDefault();
    }
}

function handleKeyUp(com, ev) {
    // escape
    if(ev.keyCode === 27) {
        com.setState({fnd: [], ix : 0});
        return;
    }
    if (ev.keyCode === 13) {
        return;
    }
    info(ev.keyCode);
    if (ev.keyCode === 39 || ev.keyCode === 40) {
        info('down', com.state.ix);
        let j = ev.keyCode === 39 ? 1 : 3;
        let ix = (com.state.ix + j);
        if (ix > com.state.fnd.length) {
            ix = 0;
        }
        com.setState({ix: ix});
        ev.preventDefault();
        return;
    }

    if (ev.keyCode === 38 || ev.keyCode === 37) {
        info('up', com.state.ix);
        let j = ev.keyCode === 37 ? 1 : 3;
        let ix = (com.state.ix - j);
        if (ix < 0) {
            ix = (com.state.fnd.length -1)
        }
        com.setState({ix: ix});
        ev.preventDefault();
        return;
    }
    let v = ev.target.value;
    if (v.length < 1) {
        com.setState({fnd: [], ix : 0});
        return;
    }

    let u = snpStore.getUsername();
    httpGet('find?u=' + u + '&t=' + v).then(res => {
        com.setState(res);
    })
}