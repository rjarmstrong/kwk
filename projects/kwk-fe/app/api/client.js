import Axios from 'axios'

window.KWK_ROOT_PATH = (window.KWK_ROOT_PATH || window.location.protocol + '//' + 'localhost:8500/1') + '/';

let trimLeadingSlash = function (query) {
    if (query[0] === '/') {
        return query.substr(1, query.length - 1)
    }
    return query;
};

export function httpGet(query) {
    query = trimLeadingSlash(query);
    return Axios
        .get(KWK_ROOT_PATH + query, requestConfig())
        .then(handle);
}

export function httpPut(query, body) {
    query = trimLeadingSlash(query);
    return Axios.put(KWK_ROOT_PATH + query, body, requestConfig())
        .then(handle);
}

export function httpPost(query, body) {
    query = trimLeadingSlash(query);
    return Axios
        .post(KWK_ROOT_PATH + query, body, requestConfig())
        .then(handle);
}

function handle(res) {
    // if its a 4xx or 5xx status then copy the status to the message body
    let s = res.status - 400;
    if (s >= 0 && s < 200) {
        if (typeof res.data === 'string') {
            return {errStatus: 1, message: res.data};
        }
        res.data.errStatus = res.status;
    }
    return res.data || {};
}

function requestConfig() {
    return {
        validateStatus: false,
        headers: {
            'x-kwk-token': getToken()
        }
    }
}

export function clearToken() {
    clearCookie('x-kwk-token');
    setTimeout(function() {}, 50);
}

function clearCookie(name, domain, path){
    domain = domain || document.domain;
    path = path || "/";
    document.cookie = name + "=; expires=" + +new Date + "; domain=" + domain + "; path=" + path;
}

export function decodeToken() {
    let t = getToken();
    let segs = t.split('.');
    if (segs.length < 3) {
        return {key: 'anon'}
    }
    let base64Url = segs[1];
    let base64 = base64Url
        .replace('-', '+')
        .replace('_', '/');
    let jsonString = window.atob(base64);
    return JSON.parse(jsonString);
}

function getToken() {
    let kt = '';
    document
        .cookie
        .split(';')
        .forEach(t => {
            let kv = t.split('=');
            if (kv[0].trim() === 'x-kwk-token') {
                kt = kv[1].trim();
            }
        });
    return kt;
}