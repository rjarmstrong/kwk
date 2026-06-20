import {decodeToken, httpPost} from '../api/client'

/**
 * @class UserStore
 */

class UserStore {

    constructor(request) {
        this.request = request;
    }

    signUp(username, password, email, inviteCode) {
        return httpPost('sign-up', {
            username: username,
            password: password,
            email: email,
            inviteCode: inviteCode,
        })
    }

    signIn(username, password) {
        return httpPost('sign-in', {username: username, password: password});
    }

    isSignedIn() {
        let t = decodeToken();
        return t.key !== 'anon';
    }

    isUser(u) {
        return u === this.username();
    }

    username() {
        let t = decodeToken();
        return t.key;
    }

}

let userStore = new UserStore();

export default userStore