
export function snpNav(router, als) {
    if(als.length === 0) {
        return
    }
    let fns = als.map(x => { return fileName(x); });
    router.push(uri({u:als[0].u, p:als[0].p, n:fns.join(',')}));
}

export function parseAl(uri) {
    uri = uri.substr(1, uri.length-1);
    let ts = uri.split('/');
    if(ts.length === 1) {
        return {u: ts[0]};
    }
    if (ts.length === 2 ) {
        return {u:ts[0], p:ts[1]};
    }
    if (ts.length === 3 ) {
        let sn = parseName(ts[2]);
        return { u : ts[0], p : ts[1], n : sn.n, x : sn.x }
    }
    return {};
}

export function parseName(txt) {
    let ix = txt.lastIndexOf('.');
    let n = txt.substring(0, ix);
    let x = txt.substring(ix+1);
    return { n : n, x : x }
}

export function parseMal(uri) {
    uri = uri.substr(1, uri.length-1);
    let ts = uri.split('/');
    if(ts.length!==3) {
        return [];
    }
    let mAls = [];
    let ft = ts[2].split(',');
    ft.forEach(x => {
        let sn = parseName(x);
        mAls.push({ u : ts[0], p : ts[1], n : sn.n, x : sn.x });
    });
    return mAls;
}

export function uri(al) {
    if(!al.p) return '/' + al.u;
    if(!al.n) return '/' + al.u + '/' + al.p;
    return '/' + al.u + '/' + relPath(al);
}

export function pchPath(p) {
    return '/' + p.u + '/' + p.p;
}

export function relPath(al) {
    return al.p + '/' + fileName(al);
}

export function fileName(al) {
    if(!al.n) {
        info(al);
    }
    let snps = al.n.split(',');
    if (snps.length>1) {
        // multi
    }
    if (!al.x) {
        return  al.n;
    }
    return al.n + '.' + al.x;
}

export const urlRegex = '(https?:\\/\\/(.*\\.)?[-a-zA-Z0-9@:%._\\+~#=]{2,256}\\.*[a-z]{2,6}\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*))';

window.KWK_LOG = false;
window.KWK_LOG_FILTER = null;
window.KWK_LOG_STACK = false;

/**
 @param {...*} args
 */
export function info() {
    let msg = Array.prototype.slice.call(arguments);
    let strMsg = JSON.stringify(msg) || '';
    if (window.KWK_LOG) {
        let l = line();
        if(window.KWK_LOG_FILTER !== null) {
            if (strMsg.indexOf(window.KWK_LOG_FILTER) === -1 &&
                l.indexOf(window.KWK_LOG_FILTER) === -1) {
                return;
            }
        }
        if (msg.length>0) {
            console.log('KWK:', l, Array(80 - l.length).join(' '), strMsg);
            return;
        }
        console.log('KWK:', l, Array(80-l.length).join(' '));
    }
    if (window.KWK_LOG_STACK) {
        let err = new Error();
        msg.push(err.stack);
        console.log('KWK:', strMsg);
    }
}

function line() {
    const e = new Error();
    if (!e.stack) try {
        // IE requires the Error to actually be throw or else the Error's 'stack'
        // property is undefined.
        throw e;
    } catch (e) {
        if (!e.stack) {
            return 0; // IE < 10, likely
        }
    }
    const frame = e.stack.toString().split(/\r\n|\n/)[3];
    let url = new RegExp(urlRegex).exec(frame);
    let trimmed = frame.replace('    at ', '').split(/\(|\[/)[0];
    let prefix = trimmed.substr(0, 30);
    return '-> ' + prefix + Array(31-prefix.length+1).join(' ') +  url[0];
}

/*
The MIT License (MIT)

Copyright (c) 2015 Yahoo Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
 */
export const elementsFromPoint = (document.elementsFromPoint || function (x,y) {
    let elements = [], previousPointerEvents = [], current, i, d;

    // get all elements via elementFromPoint, and remove them from hit-testing in order
    while ((current = document.elementFromPoint(x,y)) && elements.indexOf(current) === -1 && current !== null) {

        // push the element and its current style
        elements.push(current);
        previousPointerEvents.push({
            value: current.style.getPropertyValue('pointer-events'),
            priority: current.style.getPropertyPriority('pointer-events')
        });

        // add "pointer-events: none", to get to the underlying element
        current.style.setProperty('pointer-events', 'none', 'important');
    }

    // restore the previous pointer-events values
    for(i = previousPointerEvents.length; d=previousPointerEvents[--i]; ) {
        elements[i].style.setProperty('pointer-events', d.value ? d.value: '', d.priority);
    }

    // return our results
    return elements;

}).bind(document);

export function getPos(ele){
    let x=0;
    let y=0;
    while(true){
        x += ele.offsetLeft;
        y += ele.offsetTop;
        if(ele.offsetParent === null){
            break;
        }
        ele = ele.offsetParent;
    }
    return [x, y];
}

export function humanDate(kwkTime) {
    const delta = Math.round((Date.now() - kwkTime) / 1000);

    const minute = 60,
        hour = minute * 60,
        day = hour * 24,
        week = day * 7;

    let fuzzy;

    if (delta < minute) {
        fuzzy = delta + ' seconds ago.';
    } else if (delta < 2 * minute) {
        fuzzy = 'a minute ago.'
    } else if (delta < hour) {
        fuzzy = Math.floor(delta / minute) + ' mins ago.';
    } else if (Math.floor(delta / hour) === 1) {
        fuzzy = '1 hour ago.'
    } else if (delta < day) {
        fuzzy = Math.floor(delta / hour) + ' hrs ago.';
    } else if (delta < day * 2) {
        fuzzy = 'yesterday at ' + new Date(kwkTime);
    } else {
        fuzzy = new Date(kwkTime).toDateString();
    }
    return fuzzy;
}

export function pathname() {
    return window.location.pathname
}

// returns k, v pairs
export function queryString() {
    let query = (window.location.search || '?').substr(1),
        map   = {};
    query.replace(/([^&=]+)=?([^&]*)(?:&+|$)/g, function(match, key, value) {
        (map[key] = map[key] || []).push(value);
    });
    return map;
}
