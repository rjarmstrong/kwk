import {linkEvent} from 'inferno';
import UserStore from '../stores/userStore'

function signIn(inst, evt) {
    UserStore.signIn(inst.state.username, inst.state.password).then(x => {
        if (x.errStatus) {
            return alert(x.message);
        }
        window.location.href = '/' + inst.state.username;
    });
    evt.preventDefault();
}

function setVal(instance, event) {
    // instance.setState((prev, props) => {
    //     prev[event.target.id] = event.target.value;
    //     return prev;
    // });
}

export default class SignIn extends Component {

    constructor(props) {
        super(props);
        this.state = {
            newClass: 'modal',
            buttonClass: 'button'
        };
    }

    render() {
        return (
            <div>
                <button className={this.state.buttonClass} onClick={linkEvent(this, openForm)}>Sign in</button>
                <div className={this.state.newClass}>
                    <div className="modal-background"/>
                    <div className="modal-content">
                        <h1 className="is-size-5">Sign In</h1>
                        <p className="help is-danger">{this.state.error}</p>
                        <form>
                            <div className="field">
                                <label className="label is-small">Username / Email</label>
                                <div className="control">
                                    <input id="username" className="input" type="text"
                                           placeholder="username/email" value={this.state.username} onChange={linkEvent(this, setVal)}/>
                                </div>
                            </div>
                            {/*<p className="help is-success">This username is available</p>*/}


                            <div className="field">
                                <label className="label is-small">Password</label>
                                <div className="control">
                                    <input id="password" className="input" type="password"
                                           placeholder="password" value={this.state.password} onChange={linkEvent(this, setVal)}/>
                                </div>
                                {/*<p className="help is-success">This username is available</p>*/}
                            </div>

                            <div className="field is-grouped">
                                <div className="control">
                                    <button className="button is-primary" onClick={linkEvent(this, signIn)}>Submit</button>
                                </div>
                                <div className="control">
                                    <button className="button is-link" onClick={linkEvent(this, closeForm)}>Cancel</button>
                                </div>
                            </div>
                        </form>


                    </div>
                    <button
    onClick={linkEvent(this, closeForm)}
    className="modal-close is-large"
    aria-label="close"/>
                </div>
            </div>
        )
    }
}

function closeForm(instance, event) {
    instance.setState((prev, props) => {
        prev.newClass = 'modal';
        prev.buttonClass = 'button';
        return prev;
    });
    event.preventDefault();
}

function openForm(instance, event) {
    instance.setState((prev, props) => {
        prev.newClass = 'modal is-active';
        prev.buttonClass = 'hide';
        return prev;
    })
}