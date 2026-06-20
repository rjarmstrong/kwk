import {httpGet, httpPost, httpPut} from '../api/client';
import {fileName, info, parseAl, parseMal, parseName, pathname, pchPath, uri} from '../site-wide/util';
import Publisher from '../core/publisher'
import UserStore from "./userStore";

function Pch(pal) {
    pal = pal || {};
    this.p = pal.p;
    this.u = pal.u;
    this.stats = {};
    this.crtd = Date.now();
    this.mkPvt = false; // TASK: THIS SHOULD BE DEFAULT = true, but not currently as there is a server-side bug in Linux
}

function Snp(u, p) {
    this.al = {u: u, p: p, n: 'draft', x: 'txt', draft: true};
    this.cnt = '';
    this.crtd = Date.now();
}

/**
 * @class SnpStore
 */
class SnpStore extends Publisher {

    constructor() {
        super();
        this.currSnp = {};
        this.currPch = { p : new Pch(), items : []};
        this.ps = [];
        this.latest = [];
        this.stats = {};
        this.trySetFromStatic();
    }

    trySetFromStatic() {
        info('try set static');
        if (window.KWK_STATIC_CONTENT && window.KWK_STATIC_CONTENT.currPch) {
            this.currPch = (KWK_STATIC_CONTENT.currPch);
            let d = this.getDraft(this.currPch.p.u, this.currPch.p.p);
            if (d !== null) {
                this.currPch.items = this.currPch.items || [];
                this.currPch.items.unshift(d)
            }
        }
        if (window.KWK_STATIC_CONTENT && KWK_STATIC_CONTENT.root && KWK_STATIC_CONTENT.root.ps) {
            info('try set static root with', KWK_STATIC_CONTENT.root.ps.length);
            this.ps = KWK_STATIC_CONTENT.root.ps || [];
            this.stats = KWK_STATIC_CONTENT.root.stats || {};
            this.latest = KWK_STATIC_CONTENT.root.ss || [];
        }
    }

    /**
     * Get the currently selected snippets.
     * @returns {*} An array of aliases.
     */
    getSelected() {
        return parseMal(pathname());
    }

    getPouch() {
        return parseAl(pathname()).p;
    }

    getUsername() {
        return parseAl(pathname()).u;
    }

    getSnp(_uri) {
        info(_uri);
        let snp = this.currPch.items.find(x => {
            return uri(x.al) === _uri
        });
        let al = parseAl(_uri);
        if (!snp) {
            return this.getDraft(al.u, al.p);
        }
        return snp;
    }

    removeFromCurrPch(als) {
        let t = this;
        als.forEach(i => {
            let ix = t.currPch.items.findIndex(x => {
                return uri(i) === uri(x.al);
            });
            t.currPch.items.splice(ix, 1);
        });
    }

    deleteSnippets() {
        let als = snpStore.getSelected();
        let nonDraft = [];
        let d = this.getDraft(this.currSnp.al.u, this.currSnp.al.p);
        if (d !== null) {
            nonDraft = als.filter(x => {
                return uri(x) !== uri(d.al);
            });
            this.setDraft(d.al.u, d.al.p, null);
        } else {
            nonDraft = als;
        }
        if (nonDraft.length < 1) {
            this.removeFromCurrPch([d.al]);
            return new Promise(res => {
                res({});
            })
        }
        return httpPut('snippets-delete', {
            als: nonDraft
        }).then(res => {
            if (res.errStatus) {
                return res;
            }
            this.removeFromCurrPch(als);
            return res
        })
    }

    getDraft(u, p) {
        let d = localStorage.getItem(u + ':' + p + ':draft');
        if (!d) {
            return null;
        }
        return JSON.parse(d)
    }

    setDraft(u, p, draft) {
        let key = u + ':' + p + ':draft';
        if (draft === null) {
            localStorage.removeItem(key);
            return;
        }
        localStorage.setItem(key, JSON.stringify(draft))
    }

    setCurrPch(u, p) {
        let pch = this.ps.find(x => { return pchPath(x) === pchPath({u:u, p:p});});
        this.currPch =  { p : pch, items : [] }
    }

    crePch(pal) {
        let t = this;
        let p = new Pch(pal);
        return httpPost('pouches', p).then(res => {
            if (res && res.errStatus) {
                alert(res.message);
                return;
            }
            t.currPch = p;
            t.ps.unshift(t.currPch);
            t.stats.lastP = pal.p;
        });
    }

    delPch(pal) {
        let t = this;
        return httpPut('pouches-delete', pal).then(res => {
            if (res.errStatus) {
                alert(res.message);
                return;
            }
            t.currPch = { p : new Pch(), items : []};
            // TASK: need to delete all snippets too
            let ix = t.ps.findIndex(x => { return pchPath(x) === pchPath(pal);});
            t.ps.splice(ix, 1);
            t.stats.lastP = null;
        })
    }

    mkPvt(pal) {
        pal.mkPvt = !this.currPch.p.mkPvt;
        return httpPut("pouches-mk-pvt", pal).then(res => {
                // TO ALL THE UPDATES
            this.currPch.p.mkPvt = !this.currPch.p.mkPvt;
        });
    }

    cln(al, p) {
        let s = this.getSnp(uri(al));
        let clnd = JSON.parse(JSON.stringify(s));
        clnd.al.u = UserStore.username();
        clnd.al.p = p;
        return httpPost("snippets-cln", {al: al, newAl: clnd.al }).then((res) => {
            if (res.errStatus) {
                return res;
            }
            this.currPch.items.unshift(clnd);
            return clnd.al;
        });
    }

    similar() {
        if(!this.currSnp) {
            return
        }
        return httpGet("similar?uri=" + uri(this.currSnp.al)).then(res => {
            if (res.errStatus) {
                return res;
            }
            res.cnt = res.cnt || [];
            res.time = res.time || [];
            this.currSnp.similar = res;
            this.notify("similar-fetched", null)
        });
    }

    // al is only used if getSelected returns no aliases.
    mv(al, p) {
        let t = this;
        let sel = this.getSelected();
        if (sel.length === 0) {
            sel.push(al);
        }
        return httpPut("snippets-mv", {als: sel, p:p}).then(res => {
            let ix = t.ps.findIndex(x => { return x.p === p; });
            // increment target
            t.ps[ix].stats.snips += t.getSelected().length;
            // remove from current pouch items
            t.removeFromCurrPch(t.getSelected());
        })
    }

    rnm(al, nn) {
        let sn = parseName(nn);
        let nal = {
            u : al.u,
            p : al.p,
            n : sn.n,
            x : sn.x
        };
        return httpPut("snippets-rnm", {al:al, newN: {n: nal.n, x : nal.x }}).then(res => {
            // find snp and rename it
            let ix = this.currPch.items.findIndex(x => { return uri(x.al) === uri(al); });
            if (ix > -1) {
                this.currPch.items[ix].al = nal;
            }
            if (this.currSnp.al !== undefined && uri(al) === uri(this.currSnp.al)) {
                this.currSnp.al = nal;
            }
            this.notify('snp-changed', null, { al : nal })
        });
    }

    createDraft(u) {
        if(!this.currPch.p) {
            return new Promise(res => { res() });
        }
        return new Promise((res) => {
            let {p} = this.currPch.p;
            let d = this.getDraft(u, p);
            if (d !== null) {
                // if there is already one just return;
                this.currSnp = d;
                res();
                return;
            }
            this.currSnp = new Snp(u, p);
            this.currPch.items = this.currPch.items || [];
            this.currPch.items.unshift(this.currSnp);
            this.setDraft(u, p, this.currSnp);
            res();
        });
    }

    // TODO: make private
    createFromDraft() {
        let t = this;
        if (!this.currSnp.al.draft) {
            throw new Error('Cannot create from non-draft');
        }
        return httpPost('snippets', this.currSnp).then(res => {
            if (res.errStatus) {
                return res;
            }
            t.setDraft(t.currSnp.al.u, t.currSnp.al.p, null);
            delete t.currSnp.al.draft;
        })
    }

    updSnp() {
        this.currSnp.al.n = this.currSnp.al.n.toLowerCase();
        if (this.currSnp.al.draft) {
            if (this.currSnp.al.n === 'draft') {
                return new Promise(res => {
                    return res({errStatus: 400, message: 'Please give the snippet a meaningful name'});
                });
            }
            return this.createFromDraft()
        }
        return httpPut('snippets-patch', this.currSnp).then(res => {
            if (res.errStatus) {
                return res;
            }
            snpStore.latest = [this.currSnp];
            return res;
        })
    }

    /**
     * mAl is an array of multiple aliases, used only in the browser to show multiple snippets in the same pouch.
     * @param al The alias which is being toggled.
     * @param als (existing aliases in the mal.
     * @returns {*}
     */
    recreateMal(al, als) {
        let selected = als.findIndex(x => {
            return uri(x) === uri(al)
        });
        if (selected > -1) {
            als.splice(selected, 1);
        } else {
            als.push(al);
        }
        return als;
    }

    loadPouch(params) {
        let t = this;
        // if is root clear current pouch
        if (!params.p && !params.n) {
            this.currPch = { p : new Pch(), items : [] };
            return;
        }
        // to prevent double loading when changing state and pushing history we exit if route is the same.
        if (!params.p || t.gettingPouch || pchPath(t.currPch.p) === pchPath(params)) {
            return;
        }
        t.gettingPouch = true;
        return httpGet(pchPath(params)).then(res => {
            if (res) {
                t.currPch = res || new Pch();
            }
            let d = t.getDraft(params.u, params.p);
            if (d !== null) {
                t.currPch.items = t.currPch.items || [];
                t.currPch.items.unshift(d)
            }
            t.gettingPouch = false;
        });
    }

    isSelected(al) {
        let selected = false;
        let s = this.getSelected();
        for (let i = 0; i < s.length; i++) {
            let uri2 = uri(s[i]);
            let uri3 = uri(al);
            if (uri2 === uri3) {
                selected = true;
                break;
            }
        }
        return selected;
    }

    loadRoot(username) {
        info('load root');
        return httpGet(username).then(res => {
            if (res && res.root) {
                info('replacing pouches with:', res.root.ps.length);
                this.ps = res.root.ps;
            }
            this.loaded = true
        });
    }
}


function findSnp(pouch, name) {
    let snp = null;
    for (let i = 0; i < pouch.items.length; i++) {
        if (fileName(pouch.items[i].al) === name) {
            snp = pouch.items[i];
            break;
        }
    }
    return snp;
}

let snpStore = new SnpStore();

export default snpStore