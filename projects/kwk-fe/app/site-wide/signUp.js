import {linkEvent} from 'inferno';
import UserStore from '../stores/userStore'

function signUp(inst, evt) {
    UserStore.signUp(inst.state.username, inst.state.password, inst.state.email, inst.state.inviteCode).then(x => {
        if (x.errStatus) {
            return alert(x.message);
        }
        window.location.href = '/';
    });
    evt.preventDefault();
}

function setVal(instance, event) {
    instance.setState((prev, props) => {
        prev[event.target.id] = event.target.value;
        return prev;
    });
}

function clear(inst, event) {
    inst.setState({
        username : '',
        password : '',
        email :'',
        inviteCode : ''
    });
    event.preventDefault();
}

export default class SignUp extends Component {

    constructor(props) {
        super(props);
        this.state = {};
    }

    render() {
        return (
            <form>
                <p className="help is-danger">{this.state.error}</p>

                {/*<label className="label is-small">username</label>*/}
                <div className="field has-addons">
                    <p className="control">
                        <a className="button is-medium is-static">
                            kwk.co /
                        </a>
                    </p>
                    <div className="control">
                        <input id="username" autoComplete="off"  className="input is-medium" type="text"
                               placeholder="username" value={this.state.username}
                               onChange={linkEvent(this, setVal)}/>
                    </div>
                </div>

                <div className="field">
                    <label className="label is-small">email</label>
                    <div className="control">
                        <input id="email" autoComplete="off" className="input"
                               placeholder="email" value={this.state.email} onChange={linkEvent(this, setVal)}/>
                    </div>
                </div>

                <div className="field">
                    <label className="label is-small">password</label>
                    <div className="control">
                        <input id="password" autoComplete="off" className="input" type="password"
                               placeholder="password" value={this.state.password} onChange={linkEvent(this, setVal)}/>
                    </div>
                    {/*<p className="help is-success">This username is available</p>*/}
                </div>

                {/*<div className="field">*/}
                    {/*<label className="label is-small">invite code</label>*/}
                    {/*<div className="control">*/}
                        {/*<input id="inviteCode" autoComplete="off"  className="input"*/}
                               {/*placeholder="invite code" value={this.state.inviteCode} onChange={linkEvent(this, setVal)}/>*/}
                    {/*</div>*/}
                    {/*/!*<p className="help is-success">This username is available</p>*!/*/}
                {/*</div>*/}

                <div className="field is-grouped">
                    <div className="control">
                        <button className="button is-primary" onClick={linkEvent(this, signUp)}>Sign Up</button>
                    </div>
                </div>
            </form>
        )
    }
}