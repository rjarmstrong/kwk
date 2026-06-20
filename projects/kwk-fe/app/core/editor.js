import {linkEvent} from 'inferno';
import {info, pchPath, uri} from "../site-wide/util";
import str from '../stores/snpStore';
import UserStore from '../stores/userStore';
import BreadCrumb from "../site-wide/breadcrumb";
import CodeEditor from "./codeEditor";
import Similar from './similar';

export default class Editor extends Component {
    constructor(props) {
        info(props);
        super(props);
    }

    renderButtons() {
        if (UserStore.isUser(this.props.al.u)) {
            return (
                <div className="struct">
                    <button className="button" onClick={linkEvent(this, delSnp)}>Delete</button>
                    <button className="button" onClick={linkEvent(this, clnSnp)}>Clone</button>
                    <button className="button" onClick={linkEvent(this, updSnp)}>Save</button>
                </div>
            )
        }
    }

    render() {
        let {al} = this.props;
        info(al);
        let snp = str.getSnp(uri(al));
        if (snp === null) {
            return <div>Not found</div>
        }
        return (
            <div className="struct">
                <div className="level" style="margin-bottom:5px">
                    <BreadCrumb className="bread-crumb level-left" al={al}/>
                    v{snp.al.v}
                    <div className="level-right">
                        {this.renderButtons()}
                    </div>
                </div>
                <CodeEditor al={al}/>
                <div className="output">{ snp.prvw }</div>
                <Similar/>
            </div>
        )
    }
}

function delSnp(inst, evt) {
    if (!confirm('Are you sure you want to delete this?')) {
        evt.preventDefault();
        return;
    }
    str.deleteSnippets().then(res => {
        info(res);
        if (str.currPch.items.length === 0) {
            inst.context.router.push(pchPath(str.currPch.p));
            return
        }
        let top = str.currPch.items[0];
        inst.context.router.push(uri(top.al));
    });
    evt.preventDefault();
}

function updSnp(com, ev) {
    let isDraft = str.currSnp.al.draft;
    if (str.currSnp.al.n === 'draft') {
        str.currSnp.al.n = prompt('Please give the snippet a name:')
    }
    if (str.currSnp.al.n === null) {
        return;
    }
    str.updSnp().then(() => {
        if (isDraft) {
            com.context.router.push(uri(str.currSnp.al));
        }
    });
    ev.preventDefault();
}

function clnSnp(com, ev) {
    let p = window.prompt("Which pouch do you want to clone into?");
    if (p !== null) {
        let {al} = com.props;
        str.cln(al, p).then(res => {
            if (res.errStatus) {
                alert(res.message);
                return;
            }
            com.context.router.push(uri(res));
        });
    }
}