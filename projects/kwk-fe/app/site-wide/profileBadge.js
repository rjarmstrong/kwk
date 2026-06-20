import {clearToken, decodeToken} from "../api/client";
import {linkEvent} from 'inferno'
import SignIn from './signIn'

function signOut() {
    clearToken();
    window.location.href = '/';
}

export default function ProfileBadge({props}) {
    let jwt = decodeToken();
    if (jwt.key !== 'anon') {
        return (
            <div className="navbar-item has-dropdown is-hoverable is-hidden-touch">
                <div className="navbar-link">
                    {jwt.key}
                </div>
                <div id="moreDropdown" className="navbar-dropdown ">
                    <div className="level">
                        <div className="level-item" style="cursor: pointer" onClick={linkEvent(props, signOut)}>
                            <p>
                                <strong>Sign Out</strong>
                            </p>
                        </div>
                    </div>
                    <div className="level">
                        <div className="level-item" style="cursor: pointer" onClick={linkEvent(props, signOut)}>
                            <p>
                                <strong>Settings</strong>
                            </p>
                        </div>
                    </div>
                    <div className="level">
                        <div className="level-item" style="cursor: pointer" onClick={linkEvent(props, signOut)}>
                            <p>
                                <strong>Import</strong>
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        )
    }
    return (
        <SignIn/>
    )
}