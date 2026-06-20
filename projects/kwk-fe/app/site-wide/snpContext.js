import str from '../stores/snpStore'
import {fileName, info, uri} from "./util";
import {linkEvent} from 'inferno';

export default class SnpContext extends Component {
    constructor() {
        super();
        this.state = {
            visible: false
        };
        str.subscribe(['invoke-snp-context'], 'snpContext', (chan, val) => {
            this.setState(val);
            info(val);
        })
    }

    render() {
        if (this.state.visible) {
            return (
                <div className="snp-context" style={{
                    left: this.state.x,
                    top: this.state.y
                }}>
                    <div className="menu">
                        <div class="dropdown is-active">
                            <div className="dropdown-menu" role="menu">
                                <div onClick={linkEvent(this, rnm)} className="dropdown-content">
                                    <a href="#" className="dropdown-item">
                                        Rename..
                                    </a>
                                </div>
                                <div onClick={linkEvent(this, mv)} className="dropdown-content">
                                    <a href="#" className="dropdown-item">
                                        Move..
                                    </a>
                                </div>
                                <div onClick={linkEvent(this, rnm)} className="dropdown-content">
                                    <a href="#" className="dropdown-item">
                                        Share..
                                    </a>
                                </div>
                                <div onClick={linkEvent(this, cln)} className="dropdown-content">
                                    <a href="#" className="dropdown-item">
                                        Clone
                                    </a>
                                </div>
                                <div onClick={linkEvent(this, rnm)} className="dropdown-content">
                                    <a href="#" className="dropdown-item">
                                        Delete
                                    </a>
                                </div>
                                <div onClick={linkEvent(this, rnm)} className="dropdown-content">
                                    <a href="#" className="dropdown-item">
                                        Find Similar
                                    </a>
                                </div>
                                <div onClick={linkEvent(this, rnm)} className="dropdown-content">
                                    <a href="#" className="dropdown-item">
                                       Merge
                                    </a>
                                </div>
                            </div>
                        </div>
                    </div>

                </div>)
        }
    }
}

function rnm(com, ev) {
    info("rename", com.state);
    let {al} = com.state;
    com.setState({visible: false}, () => {
        setTimeout(() => {
            let nn = prompt("Rename", fileName(al));
            if (nn !== null && nn !== fileName(al)) {
                str.rnm(al, nn);
            }
        }, 10);
    });
}

function mv(com, ev) {
    info("mv", com.state);
    let {al} = com.state;
    com.setState({visible: false}, () => {
        setTimeout(() => {
            let p = prompt("Choose a pouch");
            if (p !== null) {
                str.mv(al, p).then(() => {
                    let top = str.currPch.items[0];
                    com.context.router.push(uri(top.al));
                });
            }
        }, 10);
    });
}

function cln(com, ev) {
    let {al} = com.state;
    com.setState({visible: false}, () => {
        setTimeout(() => {
            let p = window.prompt("Which pouch do you want to clone into?");
            if (p !== null) {
                str.cln(al, p).then(res => {
                    if(res.errStatus) {
                        alert(res.message);
                        return;
                    }
                    com.context.router.push(uri(res));
                });
            }
        }, 10);
    });
}