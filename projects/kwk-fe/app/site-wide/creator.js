import {linkEvent} from 'inferno';
import UserStore from '../stores/userStore'
import str from '../stores/snpStore'
import {snpNav} from '../site-wide/util'
import {pchPath} from "./util";

export default class Creator extends Component {

    constructor(props) {
        super(props);
        this.state = {
            newClass: 'modal'
        };
    }

    render() {
        if (!UserStore.isSignedIn()) {
            return
        }
        return (
            <div>
                <div style="cursor: pointer; display:inline-block;color:#8a8a8a" onClick={linkEvent(this, createPch)}>
                <span className="icon is-large">
                    <i data-feather="plus"/>
                </span>
                    Pouch
                </div>
                <div style="cursor: pointer; display:inline-block;" onClick={linkEvent(this, createDraft)}>
                <span className="icon is-large">
                    <i data-feather="plus-circle"/>
                </span>
                    Snippet
                </div>
                <div style="display:inline-block; color: #626f55; margin: 0 0 0 32px">
                    <div className="navbar-item has-dropdown is-hoverable is-hidden-touch">
                        <div className="navbar-link">
                            <span className="icon is-large">
                                        <i data-feather="sun"/>
                                </span>
                        </div>
                        <div id="moreDropdown" className="navbar-dropdown ">
                            <div className="level">
                                <div className="level-item" style="cursor: pointer">
                                    <p>
                                        <strong>Find dupes</strong>
                                    </p>
                                </div>
                            </div>
                            <div className="level">
                                <div className="level-item" style="cursor: pointer">
                                    <p>
                                        <strong>Find stale</strong>
                                    </p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
}

function createDraft(com, ev) {
    if(!str.currPch.p.p) {
        let p = prompt("Choose a pouch");
        if (!p) {
            return;
        }
        str.setCurrPch(UserStore.username(), p);
    }
    str.createDraft(UserStore.username()).then(() => {
        if(!str.currSnp || !str.currSnp.al) return;
        snpNav(com.context.router, [str.currSnp.al]);
    });
}

function createPch(com, ev) {
    let p = prompt("Name");
    if (!p) {
        return;
    }
    let pal =  { u : UserStore.username(), p : p };
    str.crePch(pal).then(() => {
        com.context.router.push(pchPath(pal));
    });
}

function find(instance, event) {
    let v = event.target.value;
    for (let i = 0; i < pouches.length; i++) {
        if (pouches[i].startsWith(event.target.value)) {
            instance.setState((prev, props) => {
                prev.fullName = pouches[i];
            });
            break
        }
    }
}

// function setVal(instance, event) {
//     instance.setState((prev, props) => {
//         prev[event.target.id] = event.target.value;
//         return prev;
//     })
// }
//
// function createSnippet(instance, event) {
//     str.create(instance.state.relName, instance.state.content).then(res => {
//         if (res.errStatus) {
//             alert(res.message);
//             return
//         }
//         closeForm(instance, event);
//         instance.context.router.push('/' + UserStore.username() + '/' + instance.state.relName);
//     });
// }